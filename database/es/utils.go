package es

import (
	"net"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

func GetTransport(v *viper.Viper) *http.Transport {
	var config ElasticsearchConfig
	v.Sub("elastic").Unmarshal(&config)

	transport := &http.Transport{}

	if config.TransportConfig.ExpectContinueTimeout > 0 {
		transport.ExpectContinueTimeout = config.TransportConfig.ExpectContinueTimeout
	} else {
		transport.ExpectContinueTimeout = DefaultTransportConfig.ExpectContinueTimeout
	}
	//
	if config.TransportConfig.DisableCompression {
		transport.DisableCompression = config.TransportConfig.DisableCompression
	} else {
		transport.DisableCompression = DefaultTransportConfig.DisableCompression
	}
	if config.TransportConfig.MaxIdleConnsPerHost > 0 {
		transport.MaxIdleConnsPerHost = config.TransportConfig.MaxIdleConnsPerHost
	} else {
		transport.MaxIdleConnsPerHost = DefaultTransportConfig.MaxIdleConnsPerHost
	}
	if config.TransportConfig.MaxIdleConns > 0 {
		transport.MaxIdleConns = config.TransportConfig.MaxIdleConns
	} else {
		transport.MaxIdleConns = DefaultTransportConfig.MaxIdleConns
	}
	if config.TransportConfig.IdleConnTimeout > 0 {
		transport.IdleConnTimeout = config.TransportConfig.IdleConnTimeout
	} else {
		transport.IdleConnTimeout = DefaultTransportConfig.IdleConnTimeout
	}
	if config.TransportConfig.IdleConnTimeout > 0 {
		transport.IdleConnTimeout = config.TransportConfig.IdleConnTimeout
	} else {
		transport.IdleConnTimeout = DefaultTransportConfig.IdleConnTimeout
	}
	if config.TransportConfig.MaxConnsPerHost > 0 {
		transport.MaxConnsPerHost = config.TransportConfig.MaxConnsPerHost
	} else {
		transport.MaxConnsPerHost = DefaultTransportConfig.MaxConnsPerHost
	}
	if config.TransportConfig.ResponseHeaderTimeout > 0 {
		transport.ResponseHeaderTimeout = config.TransportConfig.ResponseHeaderTimeout
	} else {
		transport.ResponseHeaderTimeout = DefaultTransportConfig.ResponseHeaderTimeout
	}
	if config.TransportConfig.TLSHandshakeTimeout > 0 {
		transport.TLSHandshakeTimeout = config.TransportConfig.TLSHandshakeTimeout
	} else {
		transport.TLSHandshakeTimeout = DefaultTransportConfig.TLSHandshakeTimeout
	}
	var dialTimeout, dialKeepAlive time.Duration

	if config.TransportConfig.DialerConfig.Timeout > 0 {
		dialTimeout = config.TransportConfig.DialerConfig.Timeout
	} else {
		dialTimeout = DefaultTransportConfig.DialerConfig.Timeout
	}
	if config.TransportConfig.DialerConfig.KeepAlive > 0 {
		dialKeepAlive = config.TransportConfig.DialerConfig.KeepAlive
	} else {
		dialKeepAlive = DefaultTransportConfig.DialerConfig.KeepAlive
	}

	transport.DialContext = (&net.Dialer{
		Timeout:   dialTimeout,   // 连接超时
		KeepAlive: dialKeepAlive, // 保持连接时间
	}).DialContext

	return transport
}
