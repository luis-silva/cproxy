package cproxy

import (
	"time"

	"github.com/magisterquis/connectproxy"
	"golang.org/x/net/proxy"
)

type connectDialer struct {
	timeout time.Duration
	logger  logger
}

func NewConnectDialer(config *configuration) *connectDialer {
	return &connectDialer{timeout: config.DialTimeout, logger: config.Logger}
}

func (this *connectDialer) Dial(address string) Socket {

	if socket, err := connectproxy.New(address, proxy.Direct); err == nil {
		return socket
	} else {
		this.logger.Printf("[INFO] Unable to establish connection to [%s]: %s", address, err)
	}

	return nil
}
