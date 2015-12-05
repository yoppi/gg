package gg

import (
	"io/ioutil"
	"net/http"
	"strings"
)

type GG struct {
	UrlHandlers map[string]*ResponseHandler
}

type ResponseHandler struct {
	HandleFunc  func() string
	Status      int
	ContentType string
}

func Double(urlHandlers map[string]*ResponseHandler) *GG {
	self := &GG{urlHandlers}
	http.DefaultClient.Transport = self

	return self
}

func (g *GG) RoundTrip(req *http.Request) (*http.Response, error) {
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

func (g *GG) findHandler(targetUrl string) *ResponseHandler {
	for url, handler := range g.UrlHandlers {
		if strings.Contains(targetUrl, url) {
			return handler
		}
	}
	return nil
}

func (g *GG) Close() {
	http.DefaultClient.Transport = nil
}
