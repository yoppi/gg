package gorip

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var testExternalApi = "http://example.com/api/test"

func api(w http.ResponseWriter, r *http.Request) {
	res, _ := http.Get(testExternalApi)
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	w.Header().Add("Content-Type", "application/json")
	w.Write(body)
}

func apiResponseHandler() string {
	return `{"test":"ok"}`
}

func TestGorip(t *testing.T) {
	double := Double(map[string]func() string{
		"http://example.com/api/test": apiResponseHandler,
	})
	defer double.Close()

	ts := httptest.NewServer(http.HandlerFunc(api))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Error("unexpected", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error("unexpected")
	}

	if string(body) == "" {
		t.Error("should have response")
	}
	if !strings.Contains(string(body), "\"test\":\"ok\"") {
		t.Error("should have api reponse")
	}
}
