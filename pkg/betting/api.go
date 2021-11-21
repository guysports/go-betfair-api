package betting

import (
	"context"
	"encoding/json"
	"reflect"
	"time"

	"github.com/guysports/go-betfair-api/pkg/transport"
	"github.com/guysports/go-betfair-api/pkg/types"
)

const (
	betfairId = 1
)

type (
	API struct {
		Client types.TransportInterface
	}

	APIInterface interface {
		ListEventTypes(filter *types.MarketFilter) ([]types.EventTypeWrapper, error)
		ListCompetitions(filter *types.MarketFilter) ([]types.CompetitionWrapper, error)
		ListTimeRanges(from, to *time.Time, filter *types.MarketFilter, granularity string) ([]types.RangeWrapper, error)
		ListEvents(filter *types.MarketFilter) ([]types.EventWrapper, error)
		ListMarketTypes(filter *types.MarketFilter) ([]types.MarketTypeWrapper, error)
		ListCountries(filter *types.MarketFilter) ([]types.CountryWrapper, error)
		ListVenues(filter *types.MarketFilter) ([]types.VenueWrapper, error)
		ListMarketCatalogue(filter *types.MarketFilter, maxResults int, marketProjection []string) ([]types.MarketCatalogueWrapper, error)
		ListMarketBook(marketIds []string, priceProjection *types.PriceProjection, orderProjection string, matchProjection string) ([]types.MarketBookWrapper, error)
		ListRunnerBook(marketId string, selectionId int, priceProjection *types.PriceProjection, orderProjection string, matchProjection string) ([]types.MarketBookWrapper, error)
		ListCurrentOrders() (*types.CurrentOrdersWrapper, error)
	}
)

func NewAPI(ctx context.Context, config *types.Config) (*API, error) {
	client, err := transport.NewJsonRPCClient(ctx, config)
	if err != nil {
		return nil, err
	}

	return &API{
		Client: client,
	}, nil
}

func (a *API) ListEventTypes(filter *types.MarketFilter) ([]types.EventTypeWrapper, error) {
	buf, err := a.Client.Do(betfairId, "listEventTypes", filter, nil)
	if err != nil {
		return nil, err
	}

	var result []types.EventTypeWrapper
	_ = json.Unmarshal(buf, &result)

	return result, nil
}

func (a *API) ListCompetitions(filter *types.MarketFilter) ([]types.CompetitionWrapper, error) {
	buf, err := a.Client.Do(betfairId, "listCompetitions", filter, nil)
	if err != nil {
		return nil, err
	}

	var result []types.CompetitionWrapper
	_ = json.Unmarshal(buf, &result)

	return result, nil
}

func (a *API) ListTimeRanges(from, to *time.Time, filter *types.MarketFilter, granularity string) ([]types.RangeWrapper, error) {
	mfParams := types.MarketFilterParams{
		Granularity: granularity,
	}

	marketRange := types.TimeRange{}
	if from != nil {
		marketRange.From = from.Format(time.RFC3339)
	}
	if to != nil {
		marketRange.To = to.Format(time.RFC3339)
	}
	if !reflect.DeepEqual(marketRange, types.TimeRange{}) {
		filter.MarketStartTime = &marketRange
	}

	buf, err := a.Client.Do(betfairId, "listTimeRanges", filter, &mfParams)
	if err != nil {
		return nil, err
	}

	var result []types.RangeWrapper
	_ = json.Unmarshal(buf, &result)

	return result, nil
}

func (a *API) ListEvents(filter *types.MarketFilter) ([]types.EventWrapper, error) {
	buf, err := a.Client.Do(betfairId, "listEvents", filter, nil)
	if err != nil {
		return nil, err
	}

	var result []types.EventWrapper
	_ = json.Unmarshal(buf, &result)

	return result, nil
}

func (a *API) ListMarketTypes(filter *types.MarketFilter) ([]types.MarketTypeWrapper, error) {
	buf, err := a.Client.Do(betfairId, "listMarketTypes", filter, nil)
	if err != nil {
		return nil, err
	}

	var result []types.MarketTypeWrapper
	_ = json.Unmarshal(buf, &result)

	return result, nil
}

func (a *API) ListCountries(filter *types.MarketFilter) ([]types.CountryWrapper, error) {
	buf, err := a.Client.Do(betfairId, "listCountries", filter, nil)
	if err != nil {
		return nil, err
	}

	var result []types.CountryWrapper
	_ = json.Unmarshal(buf, &result)

	return result, nil
}

func (a *API) ListVenues(filter *types.MarketFilter) ([]types.VenueWrapper, error) {
	buf, err := a.Client.Do(betfairId, "listVenues", filter, nil)
	if err != nil {
		return nil, err
	}

	var result []types.VenueWrapper
	_ = json.Unmarshal(buf, &result)

	return result, nil
}

func (a *API) ListMarketCatalogue(filter *types.MarketFilter, maxResults int, marketProjection []string) ([]types.MarketCatalogueWrapper, error) {
	mfParams := types.MarketFilterParams{
		MaxResults: maxResults,
	}
	if marketProjection != nil {
		mfParams.MarketProjection = marketProjection
	}

	buf, err := a.Client.Do(betfairId, "listMarketCatalogue", filter, &mfParams)
	if err != nil {
		return nil, err
	}

	var result []types.MarketCatalogueWrapper
	_ = json.Unmarshal(buf, &result)
	return result, nil
}

func (a *API) ListMarketBook(marketIds []string, priceProjection *types.PriceProjection, orderProjection string, matchProjection string) ([]types.MarketBookWrapper, error) {
	params := types.MarketFilterParams{
		MarketIds:       marketIds,
		PriceProjection: priceProjection,
		OrderProjection: orderProjection,
		MatchProjection: matchProjection,
	}

	buf, err := a.Client.Do(betfairId, "listMarketBook", nil, &params)
	if err != nil {
		return nil, err
	}
	var result []types.MarketBookWrapper
	_ = json.Unmarshal(buf, &result)
	return result, nil
}

func (a *API) ListRunnerBook(marketId string, selectionId int, priceProjection *types.PriceProjection, orderProjection string, matchProjection string) ([]types.MarketBookWrapper, error) {
	params := types.MarketFilterParams{
		MarketId:        marketId,
		SelectionId:     selectionId,
		PriceProjection: priceProjection,
		OrderProjection: orderProjection,
		MatchProjection: matchProjection,
	}

	buf, err := a.Client.Do(betfairId, "listRunnerBook", nil, &params)
	if err != nil {
		return nil, err
	}
	var result []types.MarketBookWrapper
	_ = json.Unmarshal(buf, &result)
	return result, nil
}

func (a *API) ListCurrentOrders() (*types.CurrentOrdersWrapper, error) {
	buf, err := a.Client.Do(betfairId, "listCurrentOrders", nil, &types.MarketFilterParams{
		DateRange: &types.TimeRange{},
	})
	if err != nil {
		return nil, err
	}
	var result *types.CurrentOrdersWrapper
	_ = json.Unmarshal(buf, &result)
	return result, nil
}

func (a *API) PlaceOrder(params *types.PlaceInstructionParams) (*types.PlaceExecutionReport, error) {
	buf, err := a.Client.Do(betfairId, "placeOrder", nil, &params)
	if err != nil {
		return nil, err
	}
	var result *types.PlaceExecutionReport
	_ = json.Unmarshal(buf, result)
	return result, nil
}
