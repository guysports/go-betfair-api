package fake

import (
	"fmt"

	"github.com/guysports/go-betfair-api/pkg/types"
)

type (
	FakeTransportClient struct {
		SessionExpiredError bool
	}
)

var (
	validAuthentication = types.Authenticate{
		SessionToken: "validsessionkey",
		LoginStatus:  "active",
	}
)

func (f *FakeTransportClient) Authenticate() (*types.Authenticate, error) {
	if f.SessionExpiredError {
		return nil, fmt.Errorf("session expired [ANGX-0003]")
	}
	return &validAuthentication, nil
}

func (f *FakeTransportClient) SetSessionKey(key string) {
	validAuthentication.SessionToken = key
}

func (f *FakeTransportClient) Do(id int, method string, filter *types.MarketFilter, additionalParams interface{}) ([]byte, error) {
	return nil, nil
}
