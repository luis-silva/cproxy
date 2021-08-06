package cproxy

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/magisterquis/connectproxy"
	netproxy "golang.org/x/net/proxy"
)

type connectDialer struct {
	timeout      time.Duration
	logger       logger
	proxyAddress string
	proxyAuth    string
}

func NewConnectDialer(config *configuration) *connectDialer {
	return &connectDialer{timeout: config.DialTimeout, logger: config.Logger, proxyAddress: config.ProxyAddress, proxyAuth: config.ProxyAuth}
}

func (this *connectDialer) Dial(targetAddress string) Socket {
	this.logger.Printf("address %s", targetAddress)
	var proxyScheme, proxyHost string

	if strings.HasPrefix(this.proxyAddress, "http://") {
		proxyScheme = "http"
		proxyHost = strings.TrimPrefix(this.proxyAddress, "http://")
	} else if strings.HasPrefix(this.proxyAddress, "https://") {
		proxyScheme = "https"
		proxyHost = strings.TrimPrefix(this.proxyAddress, "https://")
	} else {
		this.logger.Printf("Invalid scheme in proxy address provided to Dialer: %s", this.proxyAddress)
		return nil
	}

	proxyDetails := url.URL{
		Scheme: proxyScheme,
		Host:   proxyHost,
	}

	this.logger.Printf("%+v", targetAddress)
	this.logger.Printf("%+v", this.proxyAddress)
	proxyConfig := connectproxy.Config{
		InsecureSkipVerify: true,
		Header: http.Header{
			"Proxy-Authorization": []string{this.proxyAuth},
		},
		DialTimeout: this.timeout,
	}
	cdialer, err := connectproxy.NewWithConfig(&proxyDetails, netproxy.Direct, &proxyConfig)
	if err != nil {
		this.logger.Printf("Error instantiating a new connect proxy dialer: %v", err)
	}

	// Connect string will be set to this
	if socket, err := cdialer.Dial("tcp", targetAddress); err == nil {
		return socket
	} else {
		this.logger.Printf("[INFO] Unable to establish connection to [%s]: %s", this.proxyAddress, err)
	}

	return nil
}
