package registry

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
)

type Registry struct {
	URL    string
	Client *http.Client
}

func New(registryURL, username, password string) (*Registry, error) {
	transport := http.DefaultTransport

	return newFromTransport(registryURL, username, password, transport)
}

func NewInsecure(registryURL, username, password string) (*Registry, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, //nolint:gosec
		},
	}

	return newFromTransport(registryURL, username, password, transport)
}

func WrapTransport(transport http.RoundTripper, url, username, password string) http.RoundTripper {
	tokenTransport := &TokenTransport{
		Transport: transport,
		Username:  username,
		Password:  password,
	}
	basicAuthTransport := &BasicTransport{
		Transport: tokenTransport,
		URL:       url,
		Username:  username,
		Password:  password,
	}
	errorTransport := &ErrorTransport{
		Transport: basicAuthTransport,
	}
	return errorTransport
}

func newFromTransport(registryURL, username, password string, transport http.RoundTripper) (*Registry, error) {
	url := strings.TrimSuffix(registryURL, "/")
	transport = WrapTransport(transport, url, username, password)
	registry := &Registry{
		URL: url,
		Client: &http.Client{
			Transport: transport,
		},
	}

	if err := registry.Ping(); err != nil {
		return nil, err
	}

	return registry, nil
}

func (r *Registry) url(pathTemplate string, args ...interface{}) string {
	pathSuffix := fmt.Sprintf(pathTemplate, args...)
	url := fmt.Sprintf("%s%s", r.URL, pathSuffix)
	return url
}

func (r *Registry) Ping() error {
	url := r.url("/v2/")
	resp, err := r.Client.Get(url)
	if resp != nil {
		defer resp.Body.Close()
	}
	return err
}
