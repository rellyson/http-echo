package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/rellyson/http-echo/pkg/netutils"
	"github.com/rellyson/http-echo/pkg/version"
)

const (
	httpHeaderAppName    = "X-App-Name"
	httpHeaderAppVersion = "X-App-Version"
	httpHeaderAppBuild   = "X-App-Build"
)

// EchoResponse is the response for the echo request.
type EchoResponse struct {
	HostInfo HostInfoResponse `json:"host"`
	HttpInfo HttpInfoResponse `json:"http"`
}

// HostInfo contains the host information.
type HostInfoResponse struct {
	Hostname string `json:"hostname"` // application hostname
	IP       string `json:"ip"`       // application IP address
}

// HttpInfo contains the HTTP information.
type HttpInfoResponse struct {
	Headers map[string]string `json:"headers"`           // request headers
	Queries map[string]string `json:"queries,omitempty"` // query parameters
	Params  []string          `json:"params,omitempty"`  // path parameters
	Body    interface{}       `json:"body,omitempty"`    // request body
}

// Handles the echo request.
func Echo(w http.ResponseWriter, r *http.Request) {
	hostinfo, err := netutils.GetHostInfo()

	if err != nil {
		panic(err)
	}

	setHeaders(w)
	w.WriteHeader(http.StatusOK)

	response := &EchoResponse{
		HostInfo: HostInfoResponse{
			Hostname: hostinfo.Hostname,
			IP:       hostinfo.IP.String(),
		},
		HttpInfo: HttpInfoResponse{
			Headers: mapHeaders(r.Header),
			Queries: mapQuery(r.URL.Query()),
			Params:  mapPathParams(r.URL.Path),
			Body:    mapBody(r.Body),
		},
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

// setHeaders sets the headers for the response.
func setHeaders(w http.ResponseWriter) {
	v, err := version.GetVersion()

	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set(httpHeaderAppName, "http-echo")
	w.Header().Set(httpHeaderAppVersion, v.Version)
	w.Header().Set(httpHeaderAppBuild, v.Build)
}

// mapHeaders maps the headers to a map.
func mapHeaders(h http.Header) map[string]string {
	headers := make(map[string]string)

	for k, v := range h {
		headers[k] = v[0]
	}

	return headers
}

// mapQuery maps the query to a map.
func mapQuery(q url.Values) map[string]string {
	queries := make(map[string]string)

	for k, v := range q {
		queries[k] = v[0]
	}

	return queries
}

// mapPathParams maps the path params to a map.
func mapPathParams(p string) []string {
	var params []string

	for _, param := range strings.Split(p, "/") {
		if param == "" {
			continue
		}

		params = append(params, param)
	}

	return params
}

// mapBody maps the body to a map.
func mapBody(b io.ReadCloser) interface{} {
	var body interface{}

	if err := json.NewDecoder(b).Decode(&body); err != nil {
		return nil
	}

	return body
}
