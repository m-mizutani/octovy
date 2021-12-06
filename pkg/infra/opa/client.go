package opa

import (
	"context"
	"io"
	"net/http"

	"github.com/m-mizutani/goerr"
	opaclient "github.com/m-mizutani/opa-go-client"
	"google.golang.org/api/idtoken"
)

type Interface interface {
	Data(ctx context.Context, input interface{}, result interface{}) error
}

type Client struct {
	client *opaclient.Client
	config Config
}

type Config struct {
	BaseURL      string
	Path         string
	UseGoogleIAP bool
}

func googleIAPRequest(ctx context.Context, method, url string, data io.Reader) (*http.Response, error) {
	client, err := idtoken.NewClient(ctx, url)
	if err != nil {
		return nil, goerr.Wrap(err, "failed idtoken.NewClient for GCP IAP").With("url", url)
	}
	httpReq, err := http.NewRequestWithContext(ctx, method, url, data)
	if err != nil {
		return nil, err
	}

	if data != nil {
		httpReq.Header.Add("Content-Type", "application/json")
	}

	return client.Do(httpReq)
}

func New(cfg *Config) (*Client, error) {
	var options []opaclient.Option
	if cfg.UseGoogleIAP {
		options = append(options, opaclient.OptHTTPRequest(googleIAPRequest))
	}
	client, err := opaclient.New(cfg.BaseURL, options...)
	if err != nil {
		return nil, goerr.Wrap(err)
	}

	return &Client{
		client: client,
		config: *cfg,
	}, nil
}

func (x *Client) Data(ctx context.Context, input interface{}, result interface{}) error {
	req := &opaclient.DataRequest{
		Path:  x.config.Path,
		Input: input,
	}
	if err := x.client.GetData(ctx, req, result); err != nil {
		return goerr.Wrap(err)
	}

	return nil
}
