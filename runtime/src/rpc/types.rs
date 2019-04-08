//! RPC protocol types.
use ring::rand::{SecureRandom, SystemRandom};
use serde_cbor::Value;
use serde_derive::{Deserialize, Serialize};

impl_bytes!(
    SessionID,
    32,
    "Session identifier for multiplexing multiple sessions over the \
     same transport"
);

impl SessionID {
    /// Generate a random session identifier.
    pub fn random() -> Self {
        let rng = SystemRandom::new();
        let mut session_id = [0u8; 32];
        rng.fill(&mut session_id)
            .expect("random session id generation must succeed");

        SessionID(session_id)
    }
}

/// Frame.
#[derive(Debug, Serialize, Deserialize)]
pub struct Frame {
    pub session: SessionID,
    #[serde(with = "serde_bytes")]
    pub payload: Vec<u8>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct Request {
    pub method: String,
    pub args: Value,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct Error {
    pub message: String,
}

#[derive(Debug, Serialize, Deserialize)]
pub enum Body {
    Success(Value),
    Error(String),
}

#[derive(Debug, Serialize, Deserialize)]
pub struct Response {
    pub body: Body,
}

/// Protocol message.
#[derive(Debug, Serialize, Deserialize)]
pub enum Message {
    Request(Request),
    Response(Response),
    Close,
}
