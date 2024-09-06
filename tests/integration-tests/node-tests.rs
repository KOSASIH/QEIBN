use std::collections::{HashMap, VecDeque};
use std::sync::{Arc, Mutex};

use serde::{Serialize, Deserialize};
use serde_json::{json, Value};

use crypto::{Hash, PublicKey, SecretKey};
use crypto::hashes::{sha256, Sha256};

use tokio::prelude::*;
use tokio::sync::mpsc;

use node_architecture::node::Node;
use node_architecture::ibc::Connection;

#[tokio::test]
async fn test_ibc_protocol() {
    let chain_key = Arc::new(SecretKey::generate());
    let node1 = Node::new("node1".to_string(), chain_key.clone(), slog::Logger::root(slog::Discard, o!()));
    let node2 = Node::new("node2".to_string(), chain_key.clone(), slog::Logger::root(slog::Discard, o!()));

    // Create a new IBC connection between node1 and node2
    let conn = Connection::new(node1.clone(), node2.clone(), "conn_id".to_string()).await.unwrap();

    // Create a new packet to be sent from node1 to node2
    let packet = node1.create_packet("src_chain_id".to_string(), "src_port_id".to_string(), "src_channel_id".to_string(), "dest_chain_id".to_string(), "dest_port_id".to_string(), "dest_channel_id".to_string(), vec![1, 2, 3]).await.unwrap();

    // Send the packet from node1 to node2 via the IBC connection
    conn.send_packet(packet.clone()).await.unwrap();

    // Verify that the packet was received by node2
    let received_packet = conn.receive_packet().await.unwrap();
    assert_eq!(received_packet, packet);

    // Verify that the packet was successfully relayed from node2 to node1
    conn.relay_packet(received_packet).await.unwrap();
}

#[tokio::test]
async fn test_ibc_protocol_timeout() {
    let chain_key = Arc::new(SecretKey::generate());
    let node1 = Node::new("node1".to_string(), chain_key.clone(), slog::Logger::root(slog::Discard, o!()));
    let node2 = Node::new("node2".to_string(), chain_key.clone(), slog::Logger::root(slog::Discard, o!()));

    // Create a new IBC connection between node1 and node2
    let conn = Connection::new(node1.clone(), node2.clone(), "conn_id".to_string()).await.unwrap();

    // Create a new packet to be sent from node1 to node2
    let packet = node1.create_packet("src_chain_id".to_string(), "src_port_id".to_string(), "src_channel_id".to_string(), "dest_chain_id".to_string(), "dest_port_id".to_string(), "dest_channel_id".to_string(), vec![1, 2, 3]).await.unwrap();

    // Set a timeout for the packet relay
    let timeout = tokio::time::sleep(std::time::Duration::from_secs(5));

    // Send the packet from node1 to node2 via the IBC connection
    conn.send_packet(packet.clone()).await.unwrap();

    // Wait for the packet to be relayed or the timeout to expire
    tokio::select! {
        _ = timeout => {
            assert!(false, "packet relay timed out");
        }
        received_packet = conn.receive_packet() => {
            assert_eq!(received_packet, packet);
        }
    }
}
