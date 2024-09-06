package node

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net"
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"
)

// Node represents a node in the IBC network
type Node struct {
	id         string
	chainKey   *ecdsa.PrivateKey
	listener   net.Listener
	logger     log.Logger
	packetChan chan *Packet
}

// NewNode creates a new Node instance
func NewNode(id string, chainKey *ecdsa.PrivateKey, logger log.Logger) *Node {
	return &Node{
		id:         id,
		chainKey:   chainKey,
		logger:     logger,
		packetChan: make(chan *Packet, 10),
	}
}

// Start starts the node
func (n *Node) Start() error {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		return err
	}
	n.listener = listener
	go n.handleIncomingPackets()
	return nil
}

// handleIncomingPackets handles incoming packets from other nodes
func (n *Node) handleIncomingPackets() {
	for {
		conn, err := n.listener.Accept()
		if err != nil {
			n.logger.Error(err)
			continue
		}
		go n.handlePacket(conn)
	}
}

// handlePacket handles a single incoming packet
func (n *Node) handlePacket(conn net.Conn) {
	packet := &Packet{}
	err := json.NewDecoder(conn).Decode(packet)
	if err != nil {
		n.logger.Error(err)
		return
	}
	n.packetChan <- packet
}

// SendPacket sends a packet to another node
func (n *Node) SendPacket(packet *Packet, destNodeID string) error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:8080", destNodeID))
	if err != nil {
		return err
	}
	defer conn.Close()
	return json.NewEncoder(conn).Encode(packet)
}

// GetPacketChan returns the packet channel
func (n *Node) GetPacketChan() <-chan *Packet {
	return n.packetChan
}
