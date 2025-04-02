package raft

import (
	"context"
	"sync"
	"time"
)

// NodeState represents the current state of a Raft node
type NodeState int

const (
	Follower NodeState = iota
	Candidate
	Leader
)

// Node represents a Raft node in the cluster
type Node struct {
	mu sync.RWMutex

	// Persistent state
	currentTerm int64
	votedFor    string

	// Volatile state
	state       NodeState
	commitIndex int64
	lastApplied int64

	// Leader state
	nextIndex  map[string]int64
	matchIndex map[string]int64

	// Configuration
	id        string
	peers     map[string]string // peer ID -> address
	election  *time.Timer
	heartbeat *time.Timer

	// Channels for communication
	stopCh chan struct{}
}

// NewNode creates a new Raft node
func NewNode(id string, peers map[string]string) *Node {
	return &Node{
		id:         id,
		peers:      peers,
		state:      Follower,
		nextIndex:  make(map[string]int64),
		matchIndex: make(map[string]int64),
		stopCh:     make(chan struct{}),
	}
}

// Start begins the Raft node's operation
func (n *Node) Start(ctx context.Context) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	// Initialize timers
	n.election = time.NewTimer(randomElectionTimeout())
	n.heartbeat = time.NewTimer(50 * time.Millisecond) // Heartbeat interval
	n.heartbeat.Stop()                                 // Stop heartbeat timer initially

	// Start background tasks
	go n.run(ctx)

	return nil
}

// Stop gracefully stops the Raft node
func (n *Node) Stop() {
	close(n.stopCh)
}

func (n *Node) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-n.stopCh:
			return
		case <-n.election.C:
			n.handleElectionTimeout()
		case <-n.heartbeat.C:
			n.handleHeartbeatTimeout()
		}
	}
}

func (n *Node) handleElectionTimeout() {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.state == Follower {
		n.startElection()
	}
}

func (n *Node) startElection() {
	n.state = Candidate
	n.currentTerm++
	n.votedFor = n.id

	// Reset election timer
	n.election.Reset(randomElectionTimeout())

	// TODO: Implement request votes RPC
}

func (n *Node) handleHeartbeatTimeout() {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.state == Leader {
		// TODO: Implement append entries RPC
	}
}

// Helper function to generate random election timeout
func randomElectionTimeout() time.Duration {
	// Random timeout between 150ms and 300ms
	return time.Duration(150+time.Now().UnixNano()%150) * time.Millisecond
}
