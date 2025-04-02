package store

import (
	"context"
	"errors"
	"sync"
	"time"
)

type Client struct {
	Store
}

type ClientOptions struct {
	Addr string
}

func (opt *ClientOptions) init() {
	if opt.Addr == "" {
		opt.Addr = ":8080"
	}
}

func NewClient(opt *ClientOptions) *Client {
	opt.init()

	return &Client{
		Store{
			data: sync.Map{},
		},
	}
}

// Get retrieves a value by key
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.Store.Get(ctx, key)
}

// Set stores a value with an optional TTL
func (c *Client) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	return c.Store.Set(ctx, key, value, ttl)
}

// Delete removes a key-value pair
func (c *Client) Delete(ctx context.Context, key string) error {
	return c.Store.Delete(ctx, key)
}

var ErrKeyNotFound = errors.New("key not found")
