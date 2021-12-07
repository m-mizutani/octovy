package opa

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
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

type client struct {
	opaClient *opaclient.Client
	config    Config
	basePath  string
}

type Config struct {
	BaseURL      string
	Path         string
	UseGoogleIAP bool
}

type googleIAPClient struct{}

func (x *googleIAPClient) Do(req *http.Request) (*http.Response, error) {
	logger.With("req", req).Debug("Creating IAP HTTP Client")

	client, err := idtoken.NewClient(req.Context(), req.URL.String())
	if err != nil {
		logger.Err(err).Error("failed idtoken.NewClient")
		return nil, goerr.Wrap(err, "failed idtoken.NewClient for GCP IAP").With("req", req)
	}

	logger.With("req", req).Debug("Created IAP HTTP client")

	return client.Do(req)

}

func New(cfg *Config) (*client, error) {
	logger.With("cfg", cfg).Debug("New opa.Client")

	var options []opaclient.Option
	/*
		if cfg.UseGoogleIAP {
			logger.Debug("Add HTTP Client")
			options = append(options, opaclient.WithHTTPClient(&googleIAPClient{}))
		}
	*/
	opaClient, err := opaclient.New(cfg.BaseURL, options...)
	if err != nil {
		return nil, goerr.Wrap(err)
	}

	return &client{
		opaClient: opaClient,
		config:    *cfg,
		basePath:  strings.TrimRight(cfg.Path, "/"),
	}, nil
}

func (x *client) Data(ctx context.Context, pkg RegoPkg, input interface{}, result interface{}) error {
	type dataInput struct {
		Input interface{} `json:"input"`
	}

	type dataOutput struct {
		Result interface{} `json:"result"`
	}

	logger.With("input", input).Debug("querying to OPA server")
	url := strings.TrimRight(x.config.BaseURL, "/") + "/v1/data/" + x.basePath
	client, err := idtoken.NewClient(ctx, url)
	if err != nil {
		logger.Err(err).Error("failed idtoken.NewClient")
		return goerr.Wrap(err, "failed idtoken.NewClient for GCP IAP")
	}

	logger.With("input", input).Debug("marshal input")
	raw, err := json.Marshal(dataInput{Input: input})
	if err != nil {
		logger.Err(err).Error("json.Marshal(dataInput{Input: input})")
		return goerr.Wrap(err)
	}

	logger.With("url", url).Debug("creating request")
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(raw))
	if err != nil {
		logger.Err(err).Error("failed to create a new request")
		return goerr.Wrap(err)
	}

	logger.With("req", req).Debug("sending request")
	resp, err := client.Do(req)
	if err != nil {
		return goerr.Wrap(err)
	}

	defer resp.Body.Close()
	logger.With("resp", resp).Debug("checking response")
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return goerr.New("status code is not OK").
			With("code", resp.StatusCode).
			With("body", string(body))
	}

	logger.Debug("reading body")
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return goerr.Wrap(err).With("body", string(raw))
	}
	var output dataOutput
	logger.Debug("unmarshal body")
	if err := json.Unmarshal(body, &output); err != nil {
		return goerr.Wrap(err)
	}

	logger.Debug("marshal result")
	rawResult, err := json.Marshal(output.Result)
	if err != nil {
		return goerr.Wrap(err)
	}

	logger.Debug("unmarshal result")
	if err := json.Unmarshal(rawResult, result); err != nil {
		return goerr.Wrap(err).With("body", string(raw))
	}

	return nil
}
