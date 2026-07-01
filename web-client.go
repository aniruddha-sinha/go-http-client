package goHttpClient

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	ErrReadResponse = errors.New("errors reading response")
)

type ClientAPI interface {
	Do(ctx context.Context, url, method string, headers map[string]string, payload io.Reader) (*http.Response, error)
}

type Client struct {
	httpClient *http.Client
}

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Client{
		httpClient: httpClient,
	}
}

func (c *Client) requestBuilder(ctx context.Context, url, method string, headers map[string]string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	//set the default headers
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	//custom headers
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	return req, nil
}

func statusCodeAnalyzer(res *http.Response) (*http.Response, error) {
	if res.StatusCode >= 300 {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, ErrReadResponse
		}
		_ = res.Body.Close()
		return nil, fmt.Errorf("Client API error (status : %d): %s", res.StatusCode, string(bodyBytes))
	}

	return res, nil
}

func (c *Client) Do(ctx context.Context, url, method string, headers map[string]string, body io.Reader) (*http.Response, error) {
	req, err := c.requestBuilder(ctx, method, url, headers, body)
	if err != nil {
		return nil, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return statusCodeAnalyzer(res)
}
