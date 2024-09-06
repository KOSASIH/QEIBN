use std::collections::{HashMap, VecDeque};
use std::sync::{Arc, Mutex};

use serde::{Serialize, Deserialize};
use serde_json::{json, Value};

use crypto::{Hash, PublicKey, SecretKey};
use crypto::hashes::{sha256, Sha256};

use tokio::prelude::*;
use tokio::sync::mpsc;

use node_architecture::node::Node;

#[tokio::test]
async fn test_create_packet() {
    let chain_key = Arc::new(SecretKey::generate());
    let node = Node::new("node1".to_string(), chain_key, slog::Logger::root(slog::Discard, o!()));
    let packet = node.create_packet("src_chain_id".to_string(), "src_port_id".to_string(), "src_channel_id".to_string(), "dest_chain_id".to_string(), "dest_port_id".to_string(), "dest_channel_id".to_string(), vec![1, 2, 3]).await;
    assert!(packet.is_ok());
}

#[tokio::test]
async fn test_verify_packet() {
    let chain_key = Arc::new(SecretKey::generate());
    let node = Node::new("node1".to_string(), chain_key, slog::Logger::root(slog::Discard, o!()));
    let packet = node.create_packet("src_chain_id".to_string(), "src_port_id".to_string(), "src_channel_id".to_string(), "dest_chain_id".to_string(), "dest_port_id".to_string(), "dest_channel_id".to_string(), vec![1, 2, 3]).await.unwrap();
    let result = node.verify_packet(packet.clone()).await;
    assert!(result.is_ok());
}

#[tokio::test]
async fn test_sign_packet() {
    let chain_key = Arc::new(SecretKey::generate());
    let node = Node::new("node1".to_string(), chain_key, slog::Logger::root(slog::Discard, o!()));
    let packet = node.create_packet("src_chain_id".to_string(), "src_port_id".to_string(), "src_channel_id".to_string(), "dest_chain_id".to_string(), "dest_port_id".to_string(), "dest_channel_id".to_string(), vec![1, 2, 3]).await.unwrap();
    let sig = node.sign_packet(packet.clone()).await.unwrap();
    assert!(sig.len() > 0
}

#[tokio::test]
async fn test_verify_signature() {
    let chain_key = Arc::new(SecretKey::generate());
    let node = Node::new("node1".to_string(), chain_key, slog::Logger::root(slog::Discard, o!()));
    let packet = node.create_packet("src_chain_id".to_string(), "src_port_id".to_string(), "src_channel_id".to_string(), "dest_chain_id".to_string(), "dest_port_id".to_string(), "dest_channel_id".to_string(), vec![1, 2, 3]).await.unwrap();
    let sig = node.sign_packet(packet.clone()).await.unwrap();
    let result = node.verify_signature(packet, sig).await;
    assert!(result.is_ok());
}
