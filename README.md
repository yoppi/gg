# gg [![Circle CI](https://circleci.com/gh/yoppi/gg.svg?style=svg)](https://circleci.com/gh/yoppi/gg)

Test Double library for controller test with net/http.
A web application often use external api in action.
But cannot use that in always, each on testing.
`gg` (pronounced `double go`) help such cases.

## install

```
$ go get github.com/yoppi/gg
```

## usage

A example application code.

```go
import (
  "net/http"
  "encoding/json"
)

var externalApi = "https://example.com/api/v2/test"

func TestAction(w http.ResponseWriter, r *http.Request) {
  res, _ := http.Get(externalApi)
  defer res.Body.Close()

  body, _ := ioutil.ReadAll(res.Body)

  w.Header().Add("Content-Type", "application/json")
  w.Write(body)
}

func main() {
  http.HandleFunc("/api/test", TestAction)
  http.ListenAndServe(":8080", nil)
}
```

corresponding test.

```go
import (
  "net/http"
  "testing"
  "ioutil"
  "strings"

  "github.com/yoppi/gg"
)

func response() string {
  return `{"response":"test"}`
}

func TestExample(t *testing.T) {
	double := gg.Double(map[string]*gg.ResponseHandler{
		"http://example.com/api/test": &gg.ResponseHandler{
			HandleFunc:  apiResponseHandler,
			Status:      http.StatusOK,
			ContentType: "application/json",
		},
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
```

## license

This software is released under the MIT License, see [LICENSE](https://github.com/yoppi/gg/blob/master/LICENSE)
