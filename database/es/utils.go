package es

import (
	"net"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

func GetTransportFromViper(viper *viper.Viper, key string) *http.Transport {
	if key == "" {
		key = "elastic.transport"
	}
	var transportConfig TransportConfig
	err := viper.Sub(key).Unmarshal(&transportConfig)

	if err != nil {
		panic(err)
	}
	if transportConfig.DialerConfig.Timeout == 0 {
		transportConfig.DialerConfig.Timeout = DefaultTransportConfig.DialerConfig.Timeout
	}

	if transportConfig.DialerConfig.KeepAlive == 0 {
		transportConfig.DialerConfig.KeepAlive = DefaultTransportConfig.DialerConfig.KeepAlive
	}
	if transportConfig.ExpectContinueTimeout == 0 {
		transportConfig.ExpectContinueTimeout = DefaultTransportConfig.ExpectContinueTimeout
	}
	if transportConfig.DisableCompression {
		transportConfig.DisableCompression = DefaultTransportConfig.DisableCompression
	}
	if transportConfig.MaxIdleConnsPerHost == 0 {
		transportConfig.MaxIdleConnsPerHost = DefaultTransportConfig.MaxIdleConnsPerHost
	}
	if transportConfig.IdleConnTimeout == 0 {
		transportConfig.IdleConnTimeout = DefaultTransportConfig.IdleConnTimeout
	}
	if transportConfig.MaxConnsPerHost == 0 {
		transportConfig.MaxConnsPerHost = DefaultTransportConfig.MaxConnsPerHost
	}
	if transportConfig.MaxIdleConns == 0 {
		transportConfig.MaxIdleConns = DefaultTransportConfig.MaxIdleConns
	}
	if transportConfig.ResponseHeaderTimeout == 0 {
		transportConfig.ResponseHeaderTimeout = DefaultTransportConfig.ResponseHeaderTimeout
	}
	if transportConfig.TLSHandshakeTimeout == 0 {
		transportConfig.TLSHandshakeTimeout = DefaultTransportConfig.TLSHandshakeTimeout
	}
	if transportConfig.DialerConfig.Timeout == 0 {
		transportConfig.DialerConfig.Timeout = DefaultTransportConfig.DialerConfig.Timeout
	}
	if transportConfig.DialerConfig.KeepAlive == 0 {
		transportConfig.DialerConfig.KeepAlive = DefaultTransportConfig.DialerConfig.KeepAlive
	}

	return GetTransport(&transportConfig)
}

func GetDefaultTransport() *http.Transport {
	return GetTransport(DefaultTransportConfig)
}

func GetTransport(transportConfig *TransportConfig) *http.Transport {
	transport := &http.Transport{}

	if transportConfig.ExpectContinueTimeout > 0 {
		transport.ExpectContinueTimeout = transportConfig.ExpectContinueTimeout
	} else {
		transport.ExpectContinueTimeout = DefaultTransportConfig.ExpectContinueTimeout
	}
	//
	if transportConfig.DisableCompression {
		transport.DisableCompression = transportConfig.DisableCompression
	} else {
		transport.DisableCompression = DefaultTransportConfig.DisableCompression
	}
	if transportConfig.MaxIdleConnsPerHost > 0 {
		transport.MaxIdleConnsPerHost = transportConfig.MaxIdleConnsPerHost
	} else {
		transport.MaxIdleConnsPerHost = DefaultTransportConfig.MaxIdleConnsPerHost
	}
	if transportConfig.MaxIdleConns > 0 {
		transport.MaxIdleConns = transportConfig.MaxIdleConns
	} else {
		transport.MaxIdleConns = DefaultTransportConfig.MaxIdleConns
	}
	if transportConfig.IdleConnTimeout > 0 {
		transport.IdleConnTimeout = transportConfig.IdleConnTimeout
	} else {
		transport.IdleConnTimeout = DefaultTransportConfig.IdleConnTimeout
	}
	if transportConfig.IdleConnTimeout > 0 {
		transport.IdleConnTimeout = transportConfig.IdleConnTimeout
	} else {
		transport.IdleConnTimeout = DefaultTransportConfig.IdleConnTimeout
	}
	if transportConfig.MaxConnsPerHost > 0 {
		transport.MaxConnsPerHost = transportConfig.MaxConnsPerHost
	} else {
		transport.MaxConnsPerHost = DefaultTransportConfig.MaxConnsPerHost
	}
	if transportConfig.ResponseHeaderTimeout > 0 {
		transport.ResponseHeaderTimeout = transportConfig.ResponseHeaderTimeout
	} else {
		transport.ResponseHeaderTimeout = DefaultTransportConfig.ResponseHeaderTimeout
	}
	if transportConfig.TLSHandshakeTimeout > 0 {
		transport.TLSHandshakeTimeout = transportConfig.TLSHandshakeTimeout
	} else {
		transport.TLSHandshakeTimeout = DefaultTransportConfig.TLSHandshakeTimeout
	}
	var dialTimeout, dialKeepAlive time.Duration

	if transportConfig.DialerConfig.Timeout > 0 {
		dialTimeout = transportConfig.DialerConfig.Timeout
	} else {
		dialTimeout = DefaultTransportConfig.DialerConfig.Timeout
	}
	if transportConfig.DialerConfig.KeepAlive > 0 {
		dialKeepAlive = transportConfig.DialerConfig.KeepAlive
	} else {
		dialKeepAlive = DefaultTransportConfig.DialerConfig.KeepAlive
	}

	transport.DialContext = (&net.Dialer{
		Timeout:   dialTimeout,   // 连接超时
		KeepAlive: dialKeepAlive, // 保持连接时间
	}).DialContext

	return transport
}
