package entanglement_protocol_tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"

	"node-architecture/node"
)

func TestCreatePacket(t *testing.T) {
	chainKey, _ := crypto.GenerateKey()
	node := node.NewNode("node1", chainKey, log.NewNopLogger())
	packet, err := node.CreatePacket("src_chain_id", "src_port_id", "src_channel_id", "dest_chain_id", "dest_port_id", "dest_channel_id", []byte("data"))
	assert.NoError(t, err)
	assert.NotNil(t, packet)
}

func TestVerifyPacket(t *testing.T) {
	chainKey, _ := crypto.GenerateKey()
	node := node.NewNode("node1", chainKey, log.NewNopLogger())
	packet, _ := node.CreatePacket("src_chain_id", "src_port_id", "src_channel_id", "dest_chain_id", "dest_port_id", "dest_channel_id", []byte("data"))
	err := node.VerifyPacket(packet)
	assert.NoError(t, err)
}

func TestSignPacket(t *testing.T) {
	chainKey, _ := crypto.GenerateKey()
	node := node.NewNode("node1", chainKey, log.NewNopLogger())
	packet, _ := node.CreatePacket("src_chain_id", "src_port_id", "src_channel_id", "dest_chain_id", "dest_port_id", "dest_channel_id", []byte("data"))
	sig, err := node.SignPacket(packet)
	assert.NoError(t, err)
	assert.NotNil(t, sig)
}

func TestVerifySignature(t *testing.T) {
	chainKey, _ := crypto.GenerateKey()
	node := node.NewNode("node1", chainKey, log.NewNopLogger())
	packet, _ := node.CreatePacket("src_chain_id", "src_port_id", "src_channel_id", "dest_chain_id", "dest_port_id", "dest_channel_id", []byte("data"))
	sig, _ := node.SignPacket(packet)
	err := node.VerifySignature(packet, sig)
	assert.NoError(t, err)
}
