package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func ValuesToString(values map[string][]string) string {
	result := ""
	for k, v := range values {
		result += fmt.Sprint(k, ": ", strings.Join(v, ", "), "\n")
	}
	return result
}

func HelloServer(c http.ResponseWriter, req *http.Request) {
	var protocol string // URL.Scheme does not really work
	if req.TLS == nil {
		protocol = "http"
	} else {
		protocol = "https"
	}
	result := "hello, world!\n"
	result += "\n- Stuff:\n"
	result += "Method: " + req.Method + "\n"
	result += "UA: " + req.UserAgent() + "\n"
	result += "Protocol: " + protocol + "\n"
	result += "Host: " + req.Host + "\n"
	result += "Path: " + req.URL.Path + "\n"
	result += "URL: " + req.URL.String() + "\n"
	result += "URI: " + req.RequestURI + "\n"
	result += "Hash: " + req.URL.Fragment + "\n" // Always empty, browser doesn't send this to server
	result += "Remote: " + req.RemoteAddr + "\n"

	result += "\n- Headers:\n"
	result += ValuesToString(req.Header)

	query, _ := url.ParseQuery(req.URL.RawQuery)
	result += "\n- Query:\n"
	result += ValuesToString(query)

	bodyBytes, _ := ioutil.ReadAll(req.Body) // Note: this consumes the body and will not be available in ParseForm()
	bodyString := string(bodyBytes)
	body, _ := url.ParseQuery(bodyString)
	result += "\n- Body:\n"
	result += ValuesToString(body)

	result += "\n- Simple Query & Body:\n"
	req.ParseForm()
	result += ValuesToString(req.Form)

	c.Header().Set("Content-Type", "text/plain")
	c.Header().Set("Content-Length", strconv.Itoa(len(result)))
	io.WriteString(c, result)
}

func main() {
	http.Handle("/", http.HandlerFunc(HelloServer))
	http.ListenAndServe("0.0.0.0:8080", nil)
}
