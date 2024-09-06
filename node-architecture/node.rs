use std::collections::{HashMap, VecDeque};
use std::net::{TcpListener, TcpStream};
use std::sync::{Arc, Mutex};
use std::time::{SystemTime, UNIX_EPOCH};

use serde::{Serialize, Deserialize};
use serde_json::{json, Value};

use crypto::{Hash, PublicKey, SecretKey};
use crypto::hashes::{sha256, Sha256};

use tokio::prelude::*;
use tokio::sync::mpsc;

// Node struct
pub struct Node {
    id: String,
    chain_key: Arc<SecretKey>,
    listener: TcpListener,
    logger: slog::Logger,
    packet_chan: Arc<Mutex<VecDeque<Packet>>>,
}

impl Node {
    // Create a new Node instance
    pub fn new(id: String, chain_key: Arc<SecretKey>, logger: slog::Logger) -> Self {
        Node {
            id,
            chain_key,
            listener: TcpListener::bind("0.0.0.0:8080").unwrap(),
            logger,
            packet_chan: Arc::new(Mutex::new(VecDeque::new())),
        }
    }

    // Start the node
    pub async fn start(&self) -> Result<(), String> {
        self.listener.listen(10).unwrap();
        let mut incoming = self.listener.incoming();
        while let Some(stream) = incoming.next().await {
            let stream = stream?;
            self.handle_incoming_packet(stream).await?;
        }
        Ok(())
    }

    // Handle an incoming packet
    async fn handle_incoming_packet(&self, stream: TcpStream) -> Result<(), String> {
        let packet: Packet = serde_json::from_reader(stream).unwrap();
        self.packet_chan.lock().unwrap().push_back(packet);
        Ok(())
    }

    // Send a packet to another node
    pub async fn send_packet(&self, packet: Packet, dest_node_id: String) -> Result<(), String> {
        let mut stream = TcpStream::connect(format!("{}:8080", dest_node_id)).await?;
        serde_json::to_writer(&mut stream, &packet)?;
        Ok(())
    }

    // Get the packet channel
    pub fn get_packet_chan(&self) -> Arc<Mutex<VecDeque<Packet>>> {
        self.packet_chan.clone()
    }
}
