package ibcprotocol

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/golang/protobuf/proto"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"
)

// IBCProtocol represents the Inter-Blockchain Communication protocol
type IBCProtocol struct {
	chainID  string
	chainKey *ecdsa.PrivateKey
	logger   log.Logger
}

// NewIBCProtocol creates a new IBCProtocol instance
func NewIBCProtocol(chainID string, chainKey *ecdsa.PrivateKey, logger log.Logger) *IBCProtocol {
	return &IBCProtocol{
		chainID:  chainID,
		chainKey: chainKey,
		logger:   logger,
	}
}

// CreatePacket creates a new IBC packet
func (p *IBCProtocol) CreatePacket(srcChainID string, srcPortID string, srcChannelID string, destChainID string, destPortID string, destChannelID string, data []byte) ([]byte, error) {
	packet := &Packet{
		SourceChainID: srcChainID,
		SourcePortID:  srcPortID,
		SourceChannelID: srcChannelID,
		DestinationChainID: destChainID,
		DestinationPortID:  destPortID,
		DestinationChannelID: destChannelID,
		Data: data,
	}
	return proto.Marshal(packet)
}

// VerifyPacket verifies an IBC packet
func (p *IBCProtocol) VerifyPacket(packet []byte) error {
	packetProto := &Packet{}
	err := proto.Unmarshal(packet, packetProto)
	if err != nil {
		return err
	}
	if packetProto.DestinationChainID != p.chainID {
		return fmt.Errorf("packet is not destined for this chain")
	}
	return nil
}

// SignPacket signs an IBC packet
func (p *IBCProtocol) SignPacket(packet []byte) ([]byte, error) {
	hash := crypto.Sha256(packet)
	sig, err := ecdsa.Sign(rand.Reader, p.chainKey, hash[:])
	if err != nil {
		return nil, err
	}
	return sig, nil
}

// VerifySignature verifies the signature of an IBC packet
func (p *IBCProtocol) VerifySignature(packet []byte, sig []byte) error {
	hash := crypto.Sha256(packet)
	pubKey := p.chainKey.Public()
	return ecdsa.Verify(pubKey, hash[:], sig)
}

type Packet struct {
	SourceChainID string
	SourcePortID  string
	SourceChannelID string
	DestinationChainID string
	DestinationPortID  string
	DestinationChannelID string
	Data []byte
}

func (p *Packet) Marshal() ([]byte, error) {
	return proto.Marshal(p)
}

func (p *Packet) Unmarshal(data []byte) error {
	return proto.Unmarshal(data, p)
}
