package opa

import (
	"context"
	"net/http"
	"strings"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/octovy/pkg/utils"
	opaclient "github.com/m-mizutani/opa-go-client"
	"google.golang.org/api/idtoken"
)

var logger = utils.Logger

type Interface interface {
	Data(ctx context.Context, pkg RegoPkg, input interface{}, result interface{}) error
}

type RegoPkg string

const (
	Check RegoPkg = "check"
)

type Client struct {
	client   *opaclient.Client
	config   Config
	basePath string
}

type Config struct {
	BaseURL      string
	Path         string
	UseGoogleIAP bool
}

type googleIAPClient struct{}

func (x *googleIAPClient) Do(req *http.Request) (*http.Response, error) {
	client, err := idtoken.NewClient(req.Context(), req.URL.String())
	if err != nil {
		return nil, goerr.Wrap(err, "failed idtoken.NewClient for GCP IAP").With("req", req)
	}

	logger.With("req", req).Debug("Created IAP HTTP request")

	return client.Do(req)

}

func New(cfg *Config) (*Client, error) {
	var options []opaclient.Option
	if cfg.UseGoogleIAP {
		options = append(options, opaclient.WithHTTPClient(&googleIAPClient{}))
	}
	client, err := opaclient.New(cfg.BaseURL, options...)
	if err != nil {
		return nil, goerr.Wrap(err)
	}

	return &Client{
		client:   client,
		config:   *cfg,
		basePath: strings.TrimRight(cfg.Path, "/"),
	}, nil
}

func (x *Client) Data(ctx context.Context, pkg RegoPkg, input interface{}, result interface{}) error {
	req := &opaclient.DataRequest{
		Path:  x.basePath + "/" + string(pkg),
		Input: input,
	}
	if err := x.client.GetData(ctx, req, result); err != nil {
		return goerr.Wrap(err)
	}

	return nil
}
