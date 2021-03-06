//
//   date  : 2016-05-13
//   author: xjdrew
//

package proxy

import (
	"errors"
	"net"
	"net/url"
)

// A Dialer is a means to establish a connection.
type Dialer interface {
	// Dial connects to the given address via the proxy.
	Dial(network, addr string) (c net.Conn, err error)
}

// proxySchemes is a map from URL schemes to a function that creates a Dialer
// from a URL with such a scheme.
var proxySchemes = make(map[string]func(*url.URL, Dialer) (Dialer, error))

// RegisterDialerType takes a URL scheme and a function to generate Dialers from
// a URL with that scheme and a forwarding Dialer. Registered schemes are used
// by FromURL.
func registerDialerType(scheme string, f func(*url.URL, Dialer) (Dialer, error)) {
	proxySchemes[scheme] = f
}

func GetDialerByURL(u *url.URL, forward Dialer) (Dialer, error) {
	if f, ok := proxySchemes[u.Scheme]; ok {
		return f(u, forward)
	}
	return nil, errors.New("proxy: unknown scheme: " + u.Scheme)
}

type Proxy struct {
	Url *url.URL
	D   Dialer
}

func (p *Proxy) Dial(network, addr string) (net.Conn, error) {
	return p.D.Dial(network, addr)
}

func FromUrl(rawURL string) (*Proxy, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	d, err := GetDialerByURL(u, DirectInstance)
	if err != nil {
		return nil, err
	}

	proxy := &Proxy{
		Url: u,
		D:   d,
	}

	return proxy, nil
}
