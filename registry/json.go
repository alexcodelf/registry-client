package registry

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
)

var (
	ErrNoMorePages = errors.New("no more pages")
)

func (registry *Registry) getPaginatedJSON(url string, response interface{}) (string, error) {
	resp, err := registry.Client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(response)
	if err != nil {
		return "", err
	}
	return getNextLink(resp)
}

var nextLinkRE = regexp.MustCompile(`^ *<?([^;>]+)>? *(?:;[^;]*)*; *rel="?next"?(?:;.*)?`)

func getNextLink(resp *http.Response) (string, error) {
	for _, link := range resp.Header[http.CanonicalHeaderKey("Link")] {
		parts := nextLinkRE.FindStringSubmatch(link)
		if parts != nil {
			return parts[1], nil
		}
	}
	return "", ErrNoMorePages
}
