package ibc_protocol_tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"

	"node-architecture/node"
	"node-architecture/ibc"
)

func TestIBCProtocol(t *testing.T) {
	chainKey, _ := crypto.GenerateKey()
	node1 := node.NewNode("node1", chainKey, log.NewNopLogger())
	node2 := node.NewNode("node2", chainKey, log.NewNopLogger())

	// Create a new IBC connection between node1 and node2
	conn, err := ibc.NewConnection(node1, node2, "conn_id")
	assert.NoError(t, err)

	// Create a new packet to be sent from node1 to node2
	packet, err := node1.CreatePacket("src_chain_id", "src_port_id", "src_channel_id", "dest_chain_id", "dest_port_id", "dest_channel_id", []byte("data"))
	assert.NoError(t, err)

	// Send the packet from node1 to node2 via the IBC connection
	err = conn.SendPacket(packet)
	assert.NoError(t, err)

	// Verify that the packet was received by node2
	receivedPacket, err := conn.ReceivePacket()
	assert.NoError(t, err)
	assert.NotNil(t, receivedPacket)

	// Verify that the packet was successfully relayed from node2 to node1
	err = conn.RelayPacket(receivedPacket)
	assert.NoError(t, err)
}

func TestIBCProtocolTimeout(t *testing.T) {
	chainKey, _ := crypto.GenerateKey()
	node1 := node.NewNode("node1", chainKey, log.NewNopLogger())
	node2 := node.NewNode("node2", chainKey, log.NewNopLogger())

	// Create a new IBC connection between node1 and node2
	conn, err := ibc.NewConnection(node1, node2, "conn_id")
	assert.NoError(t, err)

	// Create a new packet to be sent from node1 to node2
	packet, err := node1.CreatePacket("src_chain_id", "src_port_id", "src_channel_id", "dest_chain_id", "dest_port_id", "dest_channel_id", []byte("data"))
	assert.NoError(t, err)

	// Set a timeout for the packet relay
	timeout := time.After(5 * time.Second)

	// Send the packet from node1 to node2 via the IBC connection
	err = conn.SendPacket(packet)
	assert.NoError(t, err)

	// Wait for the packet to be relayed or the timeout to expire
	select {
	case <-timeout:
		assert.Fail(t, "packet relay timed out")
	case receivedPacket, err := conn.ReceivePacket():
		assert.NoError(t, err)
		assert.NotNil(t, receivedPacket)
	}
}
