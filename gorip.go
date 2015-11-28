package gorip

import (
	"io/ioutil"
	"net/http"
	"strings"
)

type Gorip struct {
	DoubleUrlHandlers map[string]func() string
}

func Double(doubleUrlHandlers map[string]func() string) *Gorip {
	self := &Gorip{doubleUrlHandlers}
	http.DefaultClient.Transport = self

	return self
}

func (g *Gorip) RoundTrip(req *http.Request) (*http.Response, error) {
	if handler := g.findDoubleHandler(req.URL.String()); handler != nil {
		body := ioutil.NopCloser(strings.NewReader(handler()))
		resp := &http.Response{
			Header:     make(http.Header),
			Body:       body,
			StatusCode: http.StatusOK,
		}
		resp.Header.Set("Content-Type", "application/json")

		return resp, nil
	} else {
		return http.DefaultTransport.RoundTrip(req)
	}
}

func (g *Gorip) findDoubleHandler(url string) func() string {
	for doubleUrl, handler := range g.DoubleUrlHandlers {
		if strings.Contains(url, doubleUrl) {
			return handler
		}
	}
	return nil
}

func (g *Gorip) Close() {
	http.DefaultClient.Transport = nil
}
