package httpclient

import (
	"context"
	"fmt"
	"net/http"
)

func AuthorizationOption(token string) func(r *http.Request) {
	return func(r *http.Request) {
		r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	}
}

func ApiKeyOption(apiKey string) func(r *http.Request) {
	return func(r *http.Request) {
		r.Header.Add("api-key", apiKey)
	}
}

func jsonOption(r *http.Request, ctx ...context.Context) {
	r.Header.Set("Content-Type", "application/json")
}
