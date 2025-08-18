package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

const (
	maxIdleConns        = 100
	maxConnsPerHost     = 100
	maxIdleConnsPerHost = 100
)

// OptionFunc defines the function to update *http.Request
type OptionFunc func(*http.Request, ...context.Context)

type transport struct {
	*http.Transport
	options []OptionFunc
}

type Response[T any] struct {
	Code     int
	Response T
}

func NewHTTPClient(options ...OptionFunc) *http.Client {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = maxIdleConns
	t.MaxConnsPerHost = maxConnsPerHost
	t.TLSClientConfig.InsecureSkipVerify = true
	t.MaxIdleConnsPerHost = maxIdleConnsPerHost
	options = append(options, jsonOption)
	return &http.Client{
		Timeout: 30 * time.Second,
		Transport: transport{
			Transport: t,
			options:   options,
		},
	}
}

func Get[RES any](ctx context.Context, client *http.Client, url string, apiKey, token *string) (response Response[RES], err error) {
	return do[any, RES](ctx, client, http.MethodGet, url, nil, apiKey, token)
}

func Post[REQ, RES any](ctx context.Context, client *http.Client, url string, payload REQ, apiKey, token *string) (response Response[RES], err error) {
	return do[REQ, RES](ctx, client, http.MethodPost, url, payload, apiKey, token)
}

func do[REQ, RES any](ctx context.Context, client *http.Client, method, url string, payload REQ, apiKey, token *string) (response Response[RES], err error) {
	req, err := newRequest(ctx, client, method, url, payload)

	if apiKey != nil {
		req.Header.Add("api-key", *apiKey)
	}

	if token != nil {
		req.Header.Add("Authorization", "Bearer "+*token)
	}

	if err != nil {
		return response, err
	}

	return doRequest[RES](client, req)
}

func doRequest[RES any](client *http.Client, req *http.Request) (response Response[RES], err error) {

	bodyByte, err := io.ReadAll(req.Body)
	if err != nil {

	}

	req.Body = io.NopCloser(bytes.NewBuffer(bodyByte))
	var reqBody any

	if json.Unmarshal(bodyByte, &reqBody) != nil {
		reqBody = string(bodyByte)
	}

	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var respBody any
	if json.Unmarshal(bodyBytes, &respBody) != nil {
		respBody = string(bodyBytes)
	}

	response = Response[RES]{
		Code: resp.StatusCode,
	}

	var v RES
	if err = json.Unmarshal(bodyBytes, &v); err != nil {
		err = errors.New(string(bodyBytes))
		return
	}

	response.Response = v

	return
}

func newRequest(ctx context.Context, client *http.Client, method, url string, payload any) (*http.Request, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(&payload); err != nil {
		return nil, err
	}

	var req *http.Request
	req, err := http.NewRequest(method, url, &buf)
	if err != nil {
		return nil, err
	}

	if t, ok := client.Transport.(transport); ok {
		for _, option := range t.options {
			option(req, ctx)
		}
	}

	return req.WithContext(ctx), nil
}

func NewSetClient() *http.Client {
	return NewHTTPClient()
}
