package proxy

import (
	"accompany-sdk/ai_struct"
	"accompany-sdk/pkg/ternary"
	"context"
	"fmt"
	"github.com/openimsdk/tools/log"
	"golang.org/x/net/proxy"
	"net/http"
	"net/url"
)

func NewProxy(conf *ai_struct.ProxyConfig) *Proxy {
	pp := &Proxy{}
	if conf.ProxyURL == "" {
		if conf.Socks5Proxy != "" {
			var err error
			pp.Socks5, err = proxy.SOCKS5("tcp", conf.Socks5Proxy, nil, proxy.Direct)
			if err != nil {
				log.ZError(context.Background(), fmt.Sprintf("invalid socks5 proxy url: %s", conf.Socks5Proxy), err)
				return nil
			}
		}

		return pp
	}

	p, err := url.Parse(conf.ProxyURL)
	if err != nil {
		log.ZError(context.Background(), fmt.Sprintf("invalid proxy url: %s", conf.ProxyURL), err)
		return nil
	}

	pp.HttpProxy = http.ProxyURL(p)

	return pp
}

type Proxy struct {
	Socks5    proxy.Dialer
	HttpProxy func(*http.Request) (*url.URL, error)
}

func (pp *Proxy) BuildTransport() *http.Transport {
	return ternary.IfLazy(
		pp.HttpProxy != nil,
		func() *http.Transport {
			return &http.Transport{Proxy: pp.HttpProxy}
		},
		func() *http.Transport {
			return &http.Transport{Dial: pp.Socks5.Dial}
		},
	)
}

func ShouldLoad(conf *ai_struct.ProxyConfig) bool {
	return conf.Socks5Proxy != "" || conf.ProxyURL != ""
}
