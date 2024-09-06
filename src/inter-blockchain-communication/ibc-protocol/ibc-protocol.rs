use std::collections::{HashMap, VecDeque};
use std::sync::{Arc, Mutex};
use std::time::{SystemTime, UNIX_EPOCH};

use serde::{Serialize, Deserialize};
use serde_json::{json, Value};

use crypto::{Hash, PublicKey, SecretKey};
use crypto::hashes::{sha256, Sha256};

use tokio::prelude::*;
use tokio::sync::mpsc;

// IBC packet struct
#[derive(Debug, Serialize, Deserialize, Clone)]
struct Packet {
    source_chain_id: String,
    source_port_id: String,
    source_channel_id: String,
    destination_chain_id: String,
    destination_port_id: String,
    destination_channel_id: String,
    data: Vec<u8>,
}

// IBC protocol struct
pub struct IBCProtocol {
    chain_id: String,
    chain_key: Arc<SecretKey>,
    logger: slog::Logger,
}

impl IBCProtocol {
    // Create a new IBC protocol instance
    pub fn new(chain_id: String, chain_key: Arc<SecretKey>, logger: slog::Logger) -> Self {
        IBCProtocol {
            chain_id,
            chain_key,
            logger,
        }
    }

    // Create a new IBC packet
    pub async fn create_packet(&self, src_chain_id: String, src_port_id: String, src_channel_id: String, dest_chain_id: String, dest_port_id: String, dest_channel_id: String, data: Vec<u8>) -> Result<Vec<u8>, String> {
        let packet = Packet {
            source_chain_id: src_chain_id,
            source_port_id: src_port_id,
            source_channel_id: src_channel_id,
            destination_chain_id: dest_chain_id,
            destination_port_id: dest_port_id,
            destination_channel_id: dest_channel_id,
            data,
        };
        let packet_json = json!(packet);
        let packet_bytes = packet_json.to_string().as_bytes();
        Ok(packet_bytes.to_vec())
    }

    // Verify an IBC packet
    pub async fn verify_packet(&self, packet: Vec<u8>) -> Result<(), String> {
        let packet_json = json!(packet);
        let packet: Packet = serde_json::from_value(packet_json).unwrap();
        if packet.destination_chain_id != self.chain_id {
            return Err("Packet is not destined for this chain".to_string());
        }
        Ok(())
    }

    // Sign an IBC packet
    pub async fn sign_packet(&self, packet: Vec<u8>) -> Result<Vec<u8>, String> {
        let hash = sha256(packet);
        let sig = self.chain_key.sign(hash);
        Ok(sig)
    }

    // Verify the signature of an IBC packet
    pub async fn verify_signature(&self, packet: Vec<u8>, sig: Vec<u8>) -> Result<(), String> {
        let hash = sha256(packet);
        self.chain_key.verify(hash, sig)
    }
}

// IBC packet channel struct
pub struct PacketChannel {
    channel_id: String,
    packet_queue: Arc<Mutex<VecDeque<Packet>>>,
}

impl PacketChannel {
    // Create a new packet channel
    pub fn new(channel_id: String) -> Self {
        PacketChannel {
            channel_id,
            packet_queue: Arc::new(Mutex::new(VecDeque::new())),
        }
    }

    // Send a packet to the channel
    pub async fn send_packet(&self, packet: Packet) -> Result<(), String> {
        self.packet_queue.lock().unwrap().push_back(packet);
        Ok(())
    }

    // Receive a packet from the channel
    pub async fn recv_packet(&self) -> Result<Packet, String> {
        let mut packet_queue = self.packet_queue.lock().unwrap();
        if let Some(packet) = packet_queue.pop_front() {
            Ok(packet)
        } else {
            Err("No packets available".to_string())
        }
    }
}
