package store

import "sync"

type Client struct {
	Store
}

func NewClient(opt *Options) *Client {
	opt.init()

	return &Client{
		Store{
			data: sync.Map{},
		},
	}
}
