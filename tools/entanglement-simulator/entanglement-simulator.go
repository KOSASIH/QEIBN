package entanglement_simulator

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"

	"node-architecture/node"
	"node-architecture/ibc"
)

type EntanglementSimulator struct {
	nodes        []*node.Node
	connections []*ibc.Connection
	rand         *rand.Rand
	logger       log.Logger
}

func NewEntanglementSimulator(numNodes int, logger log.Logger) *EntanglementSimulator {
	es := &EntanglementSimulator{
		nodes:        make([]*node.Node, numNodes),
		connections: make([]*ibc.Connection, 0),
		rand:         rand.New(rand.NewSource(time.Now().UnixNano())),
		logger:       logger,
	}

	for i := 0; i < numNodes; i++ {
		es.nodes[i] = node.NewNode(fmt.Sprintf("node-%d", i), crypto.GenerateKey(), logger)
	}

	return es
}

func (es *EntanglementSimulator) CreateConnection(node1, node2 *node.Node) (*ibc.Connection, error) {
	conn, err := ibc.NewConnection(node1, node2, fmt.Sprintf("conn-%d-%d", node1.ID(), node2.ID()))
	if err != nil {
		return nil, err
	}
	es.connections = append(es.connections, conn)
	return conn, nil
}

func (es *EntanglementSimulator) SimulateEntanglement() {
	wg := &sync.WaitGroup{}

	for _, conn := range es.connections {
		wg.Add(1)
		go func(conn *ibc.Connection) {
			defer wg.Done()
			es.simulateEntanglementOnConnection(conn)
		}(conn)
	}

	wg.Wait()
}

func (es *EntanglementSimulator) simulateEntanglementOnConnection(conn *ibc.Connection) {
	packet, err := conn.Node1.CreatePacket("src_chain_id", "src_port_id", "src_channel_id", "dest_chain_id", "dest_port_id", "dest_channel_id", []byte("entangled_data"))
	if err != nil {
		es.logger.Error(err)
		return
	}

	err = conn.SendPacket(packet)
	if err != nil {
		es.logger.Error(err)
		return
	}

	receivedPacket, err := conn.ReceivePacket()
	if err != nil {
		es.logger.Error(err)
		return
	}

	es.logger.Info("Entanglement simulated successfully!")
}

func (es *EntanglementSimulator) Close() {
	for _, node := range es.nodes {
		node.Close()
	}
}
