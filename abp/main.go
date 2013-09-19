// abp is an ad blocking proxy designed to stress the gorilla/http client library.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	stdhttp "net/http"
	"sync/atomic"

	"github.com/gorilla/http"
	"github.com/gorilla/http/client"
)

var counter uint32
var listen = flag.String("-l", "localhost:8080", "listen address")

func proxy(w stdhttp.ResponseWriter, req *stdhttp.Request) {
	reqn := atomic.AddUint32(&counter, 1)
	status, h, body, err := send(reqn, req.Method, req.URL.String(), map[string][]string(req.Header), req.Body)
	if err != nil {
		stdhttp.Error(w, err.Error(), 504)
		return
	}
	fmt.Println(reqn, "resp:", status, h)
	w.WriteHeader(status.Code)
	header := w.Header()
	for k, v := range h {
		header[k] = v
	}
	if body != nil {
		defer body.Close()
		n, err := io.Copy(w, body)
		fmt.Println(reqn, "resp body:", n, err)
	}
}

func send(reqn uint32, method, url string, headers map[string][]string, body io.Reader) (client.Status, map[string][]string, io.ReadCloser, error) {
	var buf bytes.Buffer
	n, err := io.Copy(&buf, body)
	if err != nil {
		return client.Status{}, nil, nil, err
	}
	if n == 0 {
		body = nil
	} else {
		body = &buf
	}
	fmt.Println(reqn, "req:", method, url, headers)
	return http.DefaultClient.Do(method, url, headers, body)
}

func main() {
	flag.Parse()
	stdhttp.HandleFunc("/", proxy)
	stdhttp.ListenAndServe(*listen, nil)
}
