# gg [![Circle CI](https://circleci.com/gh/yoppi/gg.svg?style=svg)](https://circleci.com/gh/yoppi/gg)

Test Double library for controller test with net/http.
A web application often use external api in action.
But cannot use that in always, each on testing.
`gg` help such cases.

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

type ExternalApi struct {
  Response string
}

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

  "github.com/yoppi/gg"
)

func response() string {
  return `{"response":"test"}`
}

func TestExample(t *testing.T) {
  double := gg.Double(map[string]func() string{
    "http://example.com/api/v2/test": response,
  })
  defer double.Close()

  ts := httptest.NewServer(http.HandlerFunc(TestAction))
  defer ts.Close()

  res, err := http.Get(ts.URL)
  if err != nil {
    t.Error("unexpected", err)
  }
  defer res.Body.Close()

  if res.StatusCode != 200 {
    t.Error("unexpected", err)
  }

  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    t.Error("unexpected", err)
  }

  var ret &ExternalApi
  err := json.Unmarshal(body, &ret)
  if err != nil {
    t.Error("unexpected", err)
  }

  if ret.Response != "test" {
    t.Error("`Response` should be test")
  }
}
```

## license

This software is released under the MIT License, see [LICENSE](https://github.com/yoppi/gg/blob/master/LICENSE)
