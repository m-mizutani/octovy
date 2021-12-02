package opa

import (
	"context"

	"github.com/m-mizutani/goerr"
	opaclient "github.com/m-mizutani/opa-go-client"
)

type Interface interface {
	Data(ctx context.Context, path string, input interface{}, result interface{}) error
}

type Client struct {
	client *opaclient.Client
}

type Config struct {
	BaseURL      string
	UseGoogleIAP bool
}

func New(cfg *Config) (*Client, error) {
	var options []opaclient.Option
	if cfg.UseGoogleIAP {
		options = append(options, opaclient.OptEnableGoogleIAP())
	}
	client, err := opaclient.New(cfg.BaseURL, options...)
	if err != nil {
		return nil, goerr.Wrap(err)
	}

	return &Client{
		client: client,
	}, nil
}

func (x *Client) Data(ctx context.Context, path string, input interface{}, result interface{}) error {
	req := &opaclient.DataRequest{
		Path:  path,
		Input: input,
	}
	if err := x.client.GetData(ctx, req, result); err != nil {
		return goerr.Wrap(err)
	}

	return nil
}
