package gorip

import (
	"io/ioutil"
	"net/http"
	"strings"
)

type Gorip struct {
	UrlHandlers map[string]*ResponseHandler
}

type ResponseHandler struct {
	HandleFunc  func() string
	Status      int
	ContentType string
}

func Double(urlHandlers map[string]*ResponseHandler) *Gorip {
	self := &Gorip{urlHandlers}
	http.DefaultClient.Transport = self

	return self
}

func (g *Gorip) RoundTrip(req *http.Request) (*http.Response, error) {
	if handler := g.findHandler(req.URL.String()); handler != nil {
		body := ioutil.NopCloser(strings.NewReader(handler.HandleFunc()))
		resp := &http.Response{
			Header:     make(http.Header),
			Body:       body,
			StatusCode: handler.Status,
		}
		resp.Header.Set("Content-Type", handler.ContentType)

		return resp, nil
	} else {
		return http.DefaultTransport.RoundTrip(req)
	}
}

func (g *Gorip) findHandler(targetUrl string) *ResponseHandler {
	for url, handler := range g.UrlHandlers {
		if strings.Contains(targetUrl, url) {
			return handler
		}
	}
	return nil
}

func (g *Gorip) Close() {
	http.DefaultClient.Transport = nil
}
