package store

type Options struct {
	// host:port address
	Addr string
}

func (o *Options) init() {
	if o.Addr == "" {
		o.Addr = "127.0.0.1:6379"
	}
}
