package es

import (
	"time"
)

type ElasticsearchConfig struct {
	Addresses       []string `mapstructure:"addresses" json:"addresses" yaml:"addresses"`
	Username        string   `mapstructure:"username" json:"username" yaml:"username"`
	Password        string   `mapstructure:"password" json:"password" yaml:"password"`
	TransportConfig `mapstructure:"transport" json:"transport" yaml:"transport"`
}

type TransportConfig struct {
	DialerConfig struct {
		Timeout   time.Duration `mapstructure:"timeout" json:"timeout" yaml:"timeout"`
		KeepAlive time.Duration `mapstructure:"keep_alive" json:"keep_alive" yaml:"keep_alive"`
	} `mapstructure:"dialer" json:"dialer" yaml:"dialer"`
	ExpectContinueTimeout time.Duration `mapstructure:"expect_continue_timeout" json:"expect_continue_timeout" yaml:"expect_continue_timeout"`
	DisableCompression    bool          `mapstructure:"disable_compression" json:"disable_compression" yaml:"disable_compression"`
	MaxIdleConnsPerHost   int           `mapstructure:"max_idle_conns_per_host" json:"max_idle_conns_per_host" yaml:"max_idle_conns_per_host"`
	IdleConnTimeout       time.Duration `mapstructure:"idle_conn_timeout" json:"idle_conn_timeout" yaml:"idle_conn_timeout"`
	MaxConnsPerHost       int           `mapstructure:"max_conns_per_host" json:"max_conns_per_host" yaml:"max_conns_per_host"`
	MaxIdleConns          int           `mapstructure:"max_idle_conns" json:"max_idle_conns" yaml:"max_idle_conns"`
	ResponseHeaderTimeout time.Duration `mapstructure:"response_header_timeout" json:"response_header_timeout" yaml:"response_header_timeout"`
	TLSHandshakeTimeout   time.Duration `mapstructure:"tls_handshake_timeout" json:"tls_handshake_timeout" yaml:"tls_handshake_timeout"`
}

var DefaultTransportConfig = &TransportConfig{
	ExpectContinueTimeout: 1 * time.Second,
	DisableCompression:    false,
	MaxIdleConnsPerHost:   64,
	IdleConnTimeout:       120 * time.Second,
	MaxConnsPerHost:       64,
	MaxIdleConns:          128,
	ResponseHeaderTimeout: 2 * time.Second,
	TLSHandshakeTimeout:   3 * time.Second,
}
