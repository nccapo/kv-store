package raft

import (
	"context"
	"sync"
)

// RequestVoteArgs represents the arguments for RequestVote RPC
type RequestVoteArgs struct {
	Term         int64
	CandidateID  string
	LastLogIndex int64
	LastLogTerm  int64
}

// RequestVoteReply represents the reply for RequestVote RPC
type RequestVoteReply struct {
	Term        int64
	VoteGranted bool
}

// AppendEntriesArgs represents the arguments for AppendEntries RPC
type AppendEntriesArgs struct {
	Term         int64
	LeaderID     string
	PrevLogIndex int64
	PrevLogTerm  int64
	Entries      []LogEntry
	LeaderCommit int64
}

// AppendEntriesReply represents the reply for AppendEntries RPC
type AppendEntriesReply struct {
	Term    int64
	Success bool
}

// LogEntry represents a single entry in the Raft log
type LogEntry struct {
	Term    int64
	Command interface{}
}

// RPCClient represents a client for making RPC calls to other nodes
type RPCClient struct {
	mu    sync.RWMutex
	peers map[string]RPCClientInterface
}

// RPCClientInterface defines the interface for RPC clients
type RPCClientInterface interface {
	RequestVote(ctx context.Context, args *RequestVoteArgs, reply *RequestVoteReply) error
	AppendEntries(ctx context.Context, args *AppendEntriesArgs, reply *AppendEntriesReply) error
}

// NewRPCClient creates a new RPC client
func NewRPCClient() *RPCClient {
	return &RPCClient{
		peers: make(map[string]RPCClientInterface),
	}
}

// AddPeer adds a new peer to the RPC client
func (c *RPCClient) AddPeer(id string, client RPCClientInterface) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.peers[id] = client
}

// RemovePeer removes a peer from the RPC client
func (c *RPCClient) RemovePeer(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.peers, id)
}

// RequestVote sends a RequestVote RPC to all peers
func (c *RPCClient) RequestVote(ctx context.Context, args *RequestVoteArgs) map[string]*RequestVoteReply {
	c.mu.RLock()
	defer c.mu.RUnlock()

	replies := make(map[string]*RequestVoteReply)
	for id, peer := range c.peers {
		reply := &RequestVoteReply{}
		if err := peer.RequestVote(ctx, args, reply); err != nil {
			// TODO: Handle error
			continue
		}
		replies[id] = reply
	}
	return replies
}

// AppendEntries sends an AppendEntries RPC to all peers
func (c *RPCClient) AppendEntries(ctx context.Context, args *AppendEntriesArgs) map[string]*AppendEntriesReply {
	c.mu.RLock()
	defer c.mu.RUnlock()

	replies := make(map[string]*AppendEntriesReply)
	for id, peer := range c.peers {
		reply := &AppendEntriesReply{}
		if err := peer.AppendEntries(ctx, args, reply); err != nil {
			// TODO: Handle error
			continue
		}
		replies[id] = reply
	}
	return replies
}
