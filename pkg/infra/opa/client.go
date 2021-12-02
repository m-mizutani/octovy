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

func New(url string) (*Client, error) {
	client, err := opaclient.New(url)
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
