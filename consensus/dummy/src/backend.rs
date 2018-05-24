//! Ekiden dummy consensus backend.
use std::collections::HashMap;
use std::sync::{Arc, Mutex};

use ekiden_common::bytes::{B256, H256};
use ekiden_common::error::{Error, Result};
use ekiden_common::futures::sync::{mpsc, oneshot};
use ekiden_common::futures::{future, BoxFuture, BoxStream, Executor, Future, FutureExt, Stream,
                             StreamExt};
use ekiden_common::hash::empty_hash;
use ekiden_common::signature::Signed;
use ekiden_common::subscribers::StreamSubscribers;
use ekiden_common::uint::U256;
use ekiden_consensus_base::*;
use ekiden_scheduler_base::{Committee, CommitteeNode, CommitteeType, Role, Scheduler};
use ekiden_storage_base::StorageBackend;

/// Round state.
#[derive(Eq, PartialEq)]
enum State {
    WaitingCommitments,
    WaitingRevealsAndBlock,
}

/// Try finalize result.
#[derive(Eq, PartialEq)]
enum FinalizationResult {
    StillWaiting,
    NotifyReveals,
    Finalized(Block),
}

/// State needed for managing a protocol round.
struct Round {
    /// Storage backend.
    storage: Arc<StorageBackend>,
    /// Computation committee.
    committee: Committee,
    /// Computation group, mapped by public key hashes.
    computation_group: HashMap<B256, CommitteeNode>,
    /// Commitments from computation group nodes.
    commitments: HashMap<B256, Commitment>,
    /// Reveals from computation group nodes.
    reveals: HashMap<B256, Reveal<Header>>,
    /// Current block.
    current_block: Block,
    /// Next block.
    next_block: Option<Block>,
    /// Round state.
    state: State,
}

impl Round {
    /// Create new round descriptor.
    fn new(storage: Arc<StorageBackend>, committee: Committee, block: Block) -> Self {
        // Index computation group members by their public key hash.
        let mut computation_group = HashMap::new();
        for node in &committee.members {
            computation_group.insert(node.public_key.clone(), node.clone());
        }

        Self {
            storage,
            committee,
            computation_group,
            commitments: HashMap::new(),
            reveals: HashMap::new(),
            current_block: block,
            next_block: None,
            state: State::WaitingCommitments,
        }
    }

    /// Reset round.
    fn reset(&mut self) {
        self.commitments.clear();
        self.reveals.clear();
        self.next_block = None;
        self.state = State::WaitingCommitments;
    }

    /// Add new commitment from a node in this round.
    fn add_commitment(&mut self, commitment: Commitment) -> Result<()> {
        if self.state != State::WaitingCommitments {
            return Err(Error::new("commitment cannot be sent at this point"));
        }

        // Ensure commitment is from a valid compute node.
        let node_id = commitment.signature.public_key.clone();
        if !self.computation_group.contains_key(&node_id) {
            return Err(Error::new("node not part of computation group"));
        };

        if !commitment.verify() {
            return Err(Error::new("commitment has invalid signature"));
        }

        // Ensure node did not already submit a commitment.
        if self.commitments.contains_key(&node_id) {
            return Err(Error::new("node already sent commitment"));
        }

        self.commitments.insert(node_id, commitment);

        Ok(())
    }

    /// Add new reveal from a node in this round.
    fn add_reveal(&mut self, reveal: Reveal<Header>) -> Result<()> {
        if self.state != State::WaitingRevealsAndBlock {
            return Err(Error::new("reveal cannot be sent at this point"));
        }

        // Ensure commitment is from a valid compute node.
        let node_id = reveal.signature.public_key.clone();
        if !self.computation_group.contains_key(&node_id) {
            return Err(Error::new("node not part of computation group"));
        };

        if !reveal.verify() {
            return Err(Error::new("reveal has invalid signature"));
        }

        // Ensure node submitted a commitment.
        if !self.commitments.contains_key(&node_id) {
            return Err(Error::new("node did not send commitment"));
        }

        // Ensure node did not already submit a reveal.
        if self.reveals.contains_key(&node_id) {
            return Err(Error::new("node already sent reveal"));
        }

        self.reveals.insert(node_id, reveal);

        Ok(())
    }

    /// Add new block submission from a leader in this round.
    fn add_submit(&mut self, block: Signed<Block>) -> Result<()> {
        if self.state != State::WaitingRevealsAndBlock {
            return Err(Error::new("block cannot be sent at this point"));
        }

        // Ensure commitment is from a valid compute node and that the node is a leader.
        let node_id = block.signature.public_key.clone();
        let node = match self.computation_group.get(&node_id) {
            Some(node) => node,
            None => return Err(Error::new("node not part of computation group")),
        };

        if node.role != Role::Leader {
            return Err(Error::new("node is not a leader"));
        }

        // Ensure block has a valid signature.
        let block = block.open(&BLOCK_SUBMIT_SIGNATURE_CONTEXT)?;

        // Ensure node did not already submit a block.
        if self.next_block.is_some() {
            return Err(Error::new("node already sent block"));
        }

        self.next_block = Some(block);

        Ok(())
    }

    /// Try to finalize the round.
    fn try_finalize(&mut self) -> BoxFuture<FinalizationResult> {
        // Check if all nodes sent commitments.
        if self.commitments.len() != self.computation_group.len() {
            info!("Still waiting for other round participants to commit");
            return Box::new(future::ok(FinalizationResult::StillWaiting));
        }

        if self.state == State::WaitingCommitments {
            info!("Commitments received, now waiting for reveals");
            self.state = State::WaitingRevealsAndBlock;
            return Box::new(future::ok(FinalizationResult::NotifyReveals));
        }

        // Check if all nodes sent reveals.
        if self.reveals.len() != self.computation_group.len() {
            info!("Still waiting for other round participants to reveal");
            return Box::new(future::ok(FinalizationResult::StillWaiting));
        }

        // Check if leader sent the block.
        let block = match self.next_block.take() {
            Some(block) => block,
            None => return Box::new(future::ok(FinalizationResult::StillWaiting)),
        };

        // Everything is ready, try to finalize round.
        info!("Attempting to finalize round");
        for node_id in self.computation_group.keys() {
            let reveal = self.reveals.get(node_id).unwrap();
            let commitment = self.commitments.get(node_id).unwrap();

            if !reveal.verify_commitment(&commitment) {
                return Box::new(future::err(Error::new(format!(
                    "commitment from node {} does not match reveal",
                    node_id
                ))));
            }

            if !reveal.verify_value(&block.header) {
                return Box::new(future::err(Error::new(format!(
                    "reveal from node {} does not match block",
                    node_id
                ))));
            }
        }

        // Check if block was internally consistent.
        if !block.is_internally_consistent() {
            return Box::new(future::err(Error::new(
                "submitted block is not internally consistent",
            )));
        }

        // Check if block is based on the previous block.
        if !block.header.is_parent_of(&self.current_block.header) {
            return Box::new(future::err(Error::new(
                "submitted block is not based on previous block",
            )));
        }

        // Check if storage backend contains correct state root.
        // TODO: Currently we just check a single key, we would need to check against a log.
        if block.header.state_root != empty_hash() {
            self.storage
                .get(block.header.state_root)
                .and_then(move |_| {
                    info!("Round has been finalized");
                    Ok(FinalizationResult::Finalized(block))
                })
                .map_err(|_error| Error::new("state root not found in storage"))
                .into_box()
        } else {
            // There is no state root which means the state is empty.
            info!("Round has been finalized");
            future::ok(FinalizationResult::Finalized(block)).into_box()
        }
    }
}

#[derive(Debug)]
enum Command {
    Commit(B256, Commitment, oneshot::Sender<Result<()>>),
    Reveal(B256, Reveal<Header>, oneshot::Sender<Result<()>>),
    Submit(B256, Signed<Block>, oneshot::Sender<Result<()>>),
}

struct Inner {
    /// Scheduler.
    scheduler: Arc<Scheduler>,
    /// Storage backend.
    storage: Arc<StorageBackend>,
    /// In-memory blockchain.
    blocks: Mutex<HashMap<B256, Vec<Block>>>,
    /// Current rounds.
    rounds: Mutex<HashMap<B256, Arc<Mutex<Round>>>>,
    /// Block subscribers.
    block_subscribers: StreamSubscribers<Block>,
    /// Event subscribers.
    event_subscribers: StreamSubscribers<(B256, Event)>,
    /// Shutdown signal sender (until used).
    shutdown_sender: Mutex<Option<oneshot::Sender<()>>>,
    /// Shutdown signal receiver (until initialized).
    shutdown_receiver: Mutex<Option<oneshot::Receiver<()>>>,
    /// Command sender.
    command_sender: mpsc::UnboundedSender<Command>,
    /// Command receiver (until initialized).
    command_receiver: Mutex<Option<mpsc::UnboundedReceiver<Command>>>,
}

/// A dummy consensus backend which simulates consensus in memory.
///
/// **This backend should only be used to test implementations that use the consensus
/// interface but it only simulates a consensus backend.***
pub struct DummyConsensusBackend {
    inner: Arc<Inner>,
}

impl DummyConsensusBackend {
    /// Create new dummy consensus backend.
    pub fn new(scheduler: Arc<Scheduler>, storage: Arc<StorageBackend>) -> Self {
        // Create channels.
        let (command_sender, command_receiver) = mpsc::unbounded();
        let (shutdown_sender, shutdown_receiver) = oneshot::channel();

        Self {
            inner: Arc::new(Inner {
                scheduler,
                storage,
                blocks: Mutex::new(HashMap::new()),
                rounds: Mutex::new(HashMap::new()),
                block_subscribers: StreamSubscribers::new(),
                event_subscribers: StreamSubscribers::new(),
                shutdown_sender: Mutex::new(Some(shutdown_sender)),
                shutdown_receiver: Mutex::new(Some(shutdown_receiver)),
                command_sender,
                command_receiver: Mutex::new(Some(command_receiver)),
            }),
        }
    }

    fn get_genesis_block(contract_id: B256) -> Block {
        let mut block = Block {
            header: Header {
                version: 0,
                namespace: contract_id,
                round: U256::from(0),
                previous_hash: H256::zero(),
                group_hash: H256::zero(),
                transaction_hash: H256::zero(),
                state_root: empty_hash(),
                commitments_hash: H256::zero(),
            },
            computation_group: vec![],
            transactions: vec![],
            commitments: vec![],
        };

        block.update();
        block
    }

    /// Send a command to the backend task.
    fn send_command(
        &self,
        command: Command,
        receiver: oneshot::Receiver<Result<()>>,
    ) -> BoxFuture<()> {
        if let Err(_) = self.inner.command_sender.unbounded_send(command) {
            return Box::new(future::err(Error::new("command channel closed")));
        }

        Box::new(receiver.then(|result| match result {
            Ok(result) => result,
            Err(_) => Err(Error::new("response channel closed")),
        }))
    }

    /// Get or create round for specified contract.
    fn get_round(inner: Arc<Inner>, contract_id: B256) -> BoxFuture<Arc<Mutex<Round>>> {
        Box::new(
            inner
                .scheduler
                .get_committees(contract_id)
                .and_then(move |mut committees| {
                    // Get the computation committee.
                    let committee = committees
                        .drain(..)
                        .filter(|c| c.kind == CommitteeType::Compute)
                        .next();
                    if let Some(committee) = committee {
                        let block = {
                            let mut blocks = inner.blocks.lock().unwrap();

                            if blocks.contains_key(&contract_id) {
                                // Get last block.
                                let blocks = blocks.get(&contract_id).unwrap();
                                blocks.last().unwrap().clone()
                            } else {
                                // No blockchain yet for this contract. Create a new one.
                                let block = Self::get_genesis_block(contract_id);
                                blocks.insert(contract_id.clone(), vec![block.clone()]);

                                block
                            }
                        };

                        // Check if we already have a round and if the round is for the same committee/block.
                        let mut rounds = inner.rounds.lock().unwrap();

                        let existing_round = if rounds.contains_key(&contract_id) {
                            // Round already exists for this contract.
                            let shared_round = rounds.get(&contract_id).unwrap();
                            let round = shared_round.lock().unwrap();

                            if round.current_block == block && round.committee == committee {
                                // Existing round is the same.
                                Some(shared_round.clone())
                            } else {
                                // New round needed as either block or committee has changed.
                                None
                            }
                        } else {
                            // No round exists for this contract.
                            None
                        };

                        match existing_round {
                            Some(round) => Ok(round),
                            None => {
                                let new_round = Arc::new(Mutex::new(Round::new(
                                    inner.storage.clone(),
                                    committee,
                                    block,
                                )));
                                rounds.insert(contract_id.clone(), new_round.clone());

                                Ok(new_round)
                            }
                        }
                    } else {
                        // No compute committee, this is an error.
                        error!("No compute committee received for current round");
                        panic!("scheduler gave us no compute committee");
                    }
                }),
        )
    }

    /// Attempt to finalize the current round.
    fn try_finalize(inner: Arc<Inner>, round: Arc<Mutex<Round>>) -> BoxFuture<()> {
        let round_clone = round.clone();
        let mut round_guard = round_clone.lock().unwrap();
        let inner = inner.clone();
        let contract_id = round_guard.current_block.header.namespace.clone();

        Box::new(round_guard.try_finalize().then(move |result| {
            match result {
                Ok(FinalizationResult::Finalized(block)) => {
                    // Round has been finalized, block is ready.
                    {
                        let mut blocks = inner.blocks.lock().unwrap();
                        let mut blockchain = blocks.get_mut(&contract_id).unwrap();
                        blockchain.push(block.clone());
                    }

                    inner.block_subscribers.notify(&block);
                }
                Ok(FinalizationResult::StillWaiting) => {
                    // Still waiting for some round participants.
                }
                Ok(FinalizationResult::NotifyReveals) => {
                    // Notify round participants that they should reveal.
                    inner
                        .event_subscribers
                        .notify(&(contract_id.clone(), Event::CommitmentsReceived));
                }
                Err(error) => {
                    // Round has failed.
                    error!("Round has failed: {:?}", error);

                    {
                        let mut round = round.lock().unwrap();
                        round.reset();
                    }

                    inner
                        .event_subscribers
                        .notify(&(contract_id.clone(), Event::RoundFailed(error)));
                }
            }

            Ok(())
        }))
    }
}

impl ConsensusBackend for DummyConsensusBackend {
    fn start(&self, executor: &mut Executor) {
        info!("Starting dummy consensus backend");

        // Create command processing channel.
        let command_receiver = self.inner
            .command_receiver
            .lock()
            .unwrap()
            .take()
            .expect("start already called");
        let command_processor: BoxFuture<()> = {
            let shared_inner = self.inner.clone();

            Box::new(
                command_receiver
                    .map_err(|_| Error::new("command channel closed"))
                    .for_each(move |command| -> BoxFuture<()> {
                        let shared_inner = shared_inner.clone();

                        // Decode command.
                        let (contract_id, sender, command): (
                            _,
                            _,
                            Box<Fn(&mut Round) -> _ + Send>,
                        ) = match command {
                            Command::Commit(contract_id, commitment, sender) => (
                                contract_id,
                                sender,
                                Box::new(move |round| round.add_commitment(commitment.clone())),
                            ),
                            Command::Reveal(contract_id, reveal, sender) => (
                                contract_id,
                                sender,
                                Box::new(move |round| round.add_reveal(reveal.clone())),
                            ),
                            Command::Submit(contract_id, block, sender) => (
                                contract_id,
                                sender,
                                Box::new(move |round| round.add_submit(block.clone())),
                            ),
                        };

                        // Fetch the current round and process command.
                        Self::get_round(shared_inner.clone(), contract_id)
                            .and_then(move |round| {
                                let result = command(&mut round.lock().unwrap());

                                // Try to finalize the round.
                                Self::try_finalize(shared_inner.clone(), round)
                                    .and_then(move |_| Ok(result))
                            })
                            .then(move |result| {
                                let result = match result {
                                    Ok(result) => result,
                                    Err(error) => Err(error),
                                };
                                drop(sender.send(result));

                                Ok(())
                            })
                            .into_box()
                    }),
            )
        };

        // Create shutdown signal handler.
        let shutdown_receiver = self.inner
            .shutdown_receiver
            .lock()
            .unwrap()
            .take()
            .expect("start already called");
        let shutdown = Box::new(shutdown_receiver.then(|_| Err(Error::new("shutdown"))));

        executor.spawn(Box::new(
            future::join_all(vec![command_processor, shutdown]).then(|_| future::ok(())),
        ));
    }

    fn shutdown(&self) {
        info!("Shutting down dummy consensus backend");

        if let Some(shutdown_sender) = self.inner.shutdown_sender.lock().unwrap().take() {
            drop(shutdown_sender.send(()));
        }
    }

    fn get_blocks(&self, contract_id: B256) -> BoxStream<Block> {
        let (sender, receiver) = self.inner.block_subscribers.subscribe();
        {
            let mut blocks = self.inner.blocks.lock().unwrap();
            let block = if blocks.contains_key(&contract_id) {
                let blockchain = blocks.get(&contract_id).unwrap();
                blockchain.last().expect("empty blockchain").clone()
            } else {
                // No blockchain yet for this contract. Create a new one.
                let block = Self::get_genesis_block(contract_id);
                blocks.insert(contract_id.clone(), vec![block.clone()]);

                block
            };

            drop(sender.unbounded_send(block));
        }

        receiver
            .filter(move |block| block.header.namespace == contract_id)
            .into_box()
    }

    fn get_events(&self, contract_id: B256) -> BoxStream<Event> {
        self.inner
            .event_subscribers
            .subscribe()
            .1
            .filter(move |&(cid, _)| cid == contract_id)
            .map(|(_, event)| event)
            .into_box()
    }

    fn commit(&self, contract_id: B256, commitment: Commitment) -> BoxFuture<()> {
        let (sender, receiver) = oneshot::channel();
        self.send_command(Command::Commit(contract_id, commitment, sender), receiver)
    }

    fn reveal(&self, contract_id: B256, reveal: Reveal<Header>) -> BoxFuture<()> {
        let (sender, receiver) = oneshot::channel();
        self.send_command(Command::Reveal(contract_id, reveal, sender), receiver)
    }

    fn submit(&self, contract_id: B256, block: Signed<Block>) -> BoxFuture<()> {
        let (sender, receiver) = oneshot::channel();
        self.send_command(Command::Submit(contract_id, block, sender), receiver)
    }
}
