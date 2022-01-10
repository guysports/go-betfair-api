package transport

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/guysports/go-betfair-api/pkg/types"
	"github.com/hashicorp/go-retryablehttp"
)

type (
	JsonRPCClient struct {
		AuthData types.Authenticate
		Client   *retryablehttp.Client
		Config   *types.Config
		Ctx      context.Context
	}
)

const (
	authenticateUrl = "https://identitysso-cert.betfair.com/api/certlogin"
	jsonRPCUrl      = "https://api.betfair.com/exchange/betting/json-rpc/v1"
	defaultRootCA   = "certs/rootca.pem"
)

func NewJsonRPCClient(ctx context.Context, config *types.Config) (*JsonRPCClient, error) {
	var cacert []byte
	var err error
	if len(config.RootCAPath) > 0 {
		cacert, err = ioutil.ReadFile(config.RootCAPath)
	} else {
		cacert, err = ioutil.ReadFile(defaultRootCA)
	}
	if err != nil {
		return nil, err
	}

	// CA certificate pool
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(cacert)

	// Keypair
	certificate, err := tls.LoadX509KeyPair(config.CertPath, config.KeyPath)
	if err != nil {
		return nil, err
	}

	client := retryablehttp.NewClient()
	client.HTTPClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: caCertPool,
			Certificates: []tls.Certificate{
				certificate,
			},
		},
	}
	client.CheckRetry = retryablehttp.CheckRetry(func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		// do not retry on context.Canceled or context.DeadlineExceeded
		if ctx.Err() != nil {
			return false, ctx.Err()
		}

		// retry on any other connection errors
		if err != nil {
			return true, err
		}

		// do not retry if we received any http response code
		return false, nil
	})

	client.Logger = nil

	return &JsonRPCClient{
		Client: client,
		Config: config,
		Ctx:    ctx,
	}, nil
}

func (r *JsonRPCClient) SetSessionKey(key string) {
	r.AuthData.SessionToken = key
}

func (r *JsonRPCClient) Authenticate() (*types.Authenticate, error) {
	// Load a session key if it hasn't expired yet

	body := []byte(fmt.Sprintf("username=%s&password=%s", r.Config.User, r.Config.Password))
	req, err := retryablehttp.NewRequest(http.MethodPost, authenticateUrl, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Application", r.Config.AppKey)
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	req = req.WithContext(r.Ctx)

	resp, err := r.Client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unable to authenticate with error %s [%d]", resp.Status, resp.StatusCode)
	}
	buf, _ := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	var authenticate types.Authenticate
	err = json.Unmarshal(buf, &authenticate)
	if err != nil {
		return nil, err
	}

	r.AuthData = authenticate

	return &authenticate, nil
}

func (r *JsonRPCClient) Do(id int, method string, filter *types.MarketFilter, additionalParams interface{}) ([]byte, error) {

	params := types.Params{}
	if additionalParams != nil {
		if marketParams, ok := additionalParams.(*types.MarketFilterParams); ok {
			params = createParams(filter, marketParams)
		}
		if instructionParams, ok := additionalParams.(*types.PlaceInstructionParams); ok {
			params = createPlaceParams(instructionParams)
		}
	} else {
		params.Filter = filter
		params.Locale = "en"
	}
	query := types.JsonRPC{
		JsonRPC:   "2.0",
		RPCParams: params,
		Method:    fmt.Sprintf("SportsAPING/v1.0/%s", method),
		ID:        id,
	}
	body, err := json.Marshal(&query)
	if err != nil {
		return nil, err
	}
	req, err := retryablehttp.NewRequest(http.MethodPost, jsonRPCUrl, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Application", r.Config.AppKey)
	req.Header.Set("X-Authentication", r.AuthData.SessionToken)
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Accept", "application/json")
	req = req.WithContext(r.Ctx)

	resp, err := r.Client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unable to authenticate with error %s [%d]", resp.Status, resp.StatusCode)
	}
	buf, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var rpcresp types.JsonRPCResponse
	_ = json.Unmarshal(buf, &rpcresp)
	if rpcresp.Error != nil {
		fmt.Printf("Error returned from API %d [%s]", rpcresp.Error.Code, rpcresp.Error.Message)
		return nil, fmt.Errorf("Error returned from API %d [%s]", rpcresp.Error.Code, rpcresp.Error.Message)
	}

	payload, err := json.Marshal(rpcresp.Result)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func createParams(filter *types.MarketFilter, marketParams *types.MarketFilterParams) types.Params {
	params := types.Params{
		Filter: filter,
		Locale: "en",
	}
	if marketParams != nil {
		if marketParams.Granularity != "" {
			params.Granularity = &marketParams.Granularity
		}
		if marketParams.MarketId != "" {
			params.MarketId = marketParams.MarketId
		}
		if marketParams.MarketIds != nil {
			params.MarketIds = marketParams.MarketIds
		}
		if marketParams.SelectionId != 0 {
			params.SelectionId = marketParams.SelectionId
		}
		if marketParams.MarketProjection != nil {
			params.MarketProjection = marketParams.MarketProjection
		}
		if marketParams.MaxResults != 0 {
			params.MaxResults = marketParams.MaxResults
		}
		if marketParams.MatchProjection != "" {
			params.MatchProjection = marketParams.MatchProjection
		}
		if marketParams.OrderProjection != "" {
			params.OrderProjection = marketParams.OrderProjection
		}
		if marketParams.PriceProjection != nil {
			params.PriceProjection = marketParams.PriceProjection
		}
		if marketParams.DateRange != nil {
			params.DateRange = *marketParams.DateRange
		}
	}

	return params
}

func createPlaceParams(instructionParams *types.PlaceInstructionParams) types.Params {
	params := types.Params{
		Locale: "en",
	}

	if instructionParams.CustomerRef != "" {
		params.CustomerRef = instructionParams.CustomerRef
	}
	if instructionParams.CustomerStrategyRef != "" {
		params.CustomerStrategyRef = instructionParams.CustomerStrategyRef
	}
	if instructionParams.MarketID != "" {
		params.MarketId = instructionParams.MarketID
	}
	if instructionParams.Instructions != nil {
		params.Instructions = instructionParams.Instructions
	}

	return params
}
