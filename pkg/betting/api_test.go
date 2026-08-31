package betting

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/guysports/go-betfair-api/pkg/types"
)

type stubTransport struct {
	response       []byte
	err            error
	lastID         int
	lastMethod     string
	lastFilter     *types.MarketFilter
	lastAdditional interface{}
}

func (s *stubTransport) Authenticate() (*types.Authenticate, error) {
	return nil, nil
}

func (s *stubTransport) SetSessionKey(key string) {}

func (s *stubTransport) Do(id int, method string, filter *types.MarketFilter, additionalParams interface{}) ([]byte, error) {
	s.lastID = id
	s.lastMethod = method
	s.lastFilter = filter
	s.lastAdditional = additionalParams
	return s.response, s.err
}

func TestAPI_ListEventTypes_UsesTransportAndUnmarshalsResponse(t *testing.T) {
	filter := &types.MarketFilter{EventIds: []string{"1"}}
	response := []types.EventTypeWrapper{{EventType: &types.Detail{ID: "1", Name: "Football"}, MarketCount: 2}}
	payload, _ := json.Marshal(response)
	transport := &stubTransport{response: payload}
	api := &API{Client: transport}

	got, err := api.ListEventTypes(filter)
	if err != nil {
		t.Fatalf("ListEventTypes returned error: %v", err)
	}

	if transport.lastID != betfairId {
		t.Fatalf("expected betfair id %d, got %d", betfairId, transport.lastID)
	}
	if transport.lastMethod != "listEventTypes" {
		t.Fatalf("expected method listEventTypes, got %s", transport.lastMethod)
	}
	if !reflect.DeepEqual(transport.lastFilter, filter) {
		t.Fatalf("expected filter %#v, got %#v", filter, transport.lastFilter)
	}
	if transport.lastAdditional != nil {
		t.Fatalf("expected nil additional params, got %#v", transport.lastAdditional)
	}
	if !reflect.DeepEqual(got, response) {
		t.Fatalf("expected %#v, got %#v", response, got)
	}
}

func TestAPI_ListTimeRanges_AddsTimeRangeAndForwardsParams(t *testing.T) {
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	filter := &types.MarketFilter{EventIds: []string{"1"}}
	response := []types.RangeWrapper{{Range: &types.TimeRange{From: from.Format(time.RFC3339), To: to.Format(time.RFC3339)}, MarketCount: 3}}
	payload, _ := json.Marshal(response)
	transport := &stubTransport{response: payload}
	api := &API{Client: transport}

	got, err := api.ListTimeRanges(&from, &to, filter, "DAY")
	if err != nil {
		t.Fatalf("ListTimeRanges returned error: %v", err)
	}

	if filter.MarketStartTime == nil {
		t.Fatal("expected MarketStartTime to be set")
	}
	if filter.MarketStartTime.From != from.Format(time.RFC3339) {
		t.Fatalf("expected From %q, got %q", from.Format(time.RFC3339), filter.MarketStartTime.From)
	}
	if filter.MarketStartTime.To != to.Format(time.RFC3339) {
		t.Fatalf("expected To %q, got %q", to.Format(time.RFC3339), filter.MarketStartTime.To)
	}

	params, ok := transport.lastAdditional.(*types.MarketFilterParams)
	if !ok {
		t.Fatalf("expected MarketFilterParams, got %T", transport.lastAdditional)
	}
	if params.Granularity != "DAY" {
		t.Fatalf("expected granularity DAY, got %s", params.Granularity)
	}
	if !reflect.DeepEqual(got, response) {
		t.Fatalf("expected %#v, got %#v", response, got)
	}
}

func TestAPI_ListCurrentOrders_AndPlaceOrders(t *testing.T) {
	currentOrdersResponse := &types.CurrentOrdersWrapper{Orders: []types.CurrentOrder{{BetId: "bet-1"}}}
	currentPayload, _ := json.Marshal(currentOrdersResponse)
	currentTransport := &stubTransport{response: currentPayload}
	currentAPI := &API{Client: currentTransport}

	gotCurrent, err := currentAPI.ListCurrentOrders()
	if err != nil {
		t.Fatalf("ListCurrentOrders returned error: %v", err)
	}

	params, ok := currentTransport.lastAdditional.(*types.MarketFilterParams)
	if !ok {
		t.Fatalf("expected MarketFilterParams, got %T", currentTransport.lastAdditional)
	}
	if params.DateRange == nil {
		t.Fatal("expected DateRange to be present")
	}
	if !reflect.DeepEqual(*gotCurrent, *currentOrdersResponse) {
		t.Fatalf("expected %#v, got %#v", currentOrdersResponse, gotCurrent)
	}

	placeReport := &types.PlaceExecutionReport{Status: "SUCCESS", OrderStatus: "EXECUTION_COMPLETE"}
	placePayload, _ := json.Marshal(placeReport)
	placeTransport := &stubTransport{response: placePayload}
	placeAPI := &API{Client: placeTransport}
	instructions := []types.PlaceInstruction{{SelectionId: 10, Side: "BACK"}}
	placeParams := &types.PlaceInstructionParams{MarketID: "1.234", Instructions: instructions, CustomerRef: "ref-1"}

	gotPlace, err := placeAPI.PlaceOrders(placeParams)
	if err != nil {
		t.Fatalf("PlaceOrders returned error: %v", err)
	}

	if placeTransport.lastMethod != "placeOrders" {
		t.Fatalf("expected method placeOrders, got %s", placeTransport.lastMethod)
	}
	if placeTransport.lastFilter != nil {
		t.Fatalf("expected nil filter, got %#v", placeTransport.lastFilter)
	}
	if !reflect.DeepEqual(placeTransport.lastAdditional, placeParams) {
		t.Fatalf("expected additional params %#v, got %#v", placeParams, placeTransport.lastAdditional)
	}
	if !reflect.DeepEqual(gotPlace, placeReport) {
		t.Fatalf("expected %#v, got %#v", placeReport, gotPlace)
	}
}

func TestAPI_SimpleListMethods_ForwardFilterAndUnmarshal(t *testing.T) {
	filter := &types.MarketFilter{EventIds: []string{"1"}}

	competitionsResponse := []types.CompetitionWrapper{{Competition: &types.Detail{ID: "1", Name: "League"}, MarketCount: 4}}
	competitionsPayload, _ := json.Marshal(competitionsResponse)
	competitionsTransport := &stubTransport{response: competitionsPayload}
	competitionsAPI := &API{Client: competitionsTransport}

	gotCompetitions, err := competitionsAPI.ListCompetitions(filter)
	if err != nil {
		t.Fatalf("ListCompetitions returned error: %v", err)
	}
	if competitionsTransport.lastMethod != "listCompetitions" {
		t.Fatalf("expected method listCompetitions, got %s", competitionsTransport.lastMethod)
	}
	if !reflect.DeepEqual(gotCompetitions, competitionsResponse) {
		t.Fatalf("expected %#v, got %#v", competitionsResponse, gotCompetitions)
	}

	eventsResponse := []types.EventWrapper{{Event: &types.Detail{ID: "2", Name: "Match"}, MarketCount: 7}}
	eventsPayload, _ := json.Marshal(eventsResponse)
	eventsTransport := &stubTransport{response: eventsPayload}
	eventsAPI := &API{Client: eventsTransport}

	gotEvents, err := eventsAPI.ListEvents(filter)
	if err != nil {
		t.Fatalf("ListEvents returned error: %v", err)
	}
	if eventsTransport.lastMethod != "listEvents" {
		t.Fatalf("expected method listEvents, got %s", eventsTransport.lastMethod)
	}
	if !reflect.DeepEqual(gotEvents, eventsResponse) {
		t.Fatalf("expected %#v, got %#v", eventsResponse, gotEvents)
	}

	marketTypesResponse := []types.MarketTypeWrapper{{MarketType: "MATCH_ODDS", MarketCount: 5}}
	marketTypesPayload, _ := json.Marshal(marketTypesResponse)
	marketTypesTransport := &stubTransport{response: marketTypesPayload}
	marketTypesAPI := &API{Client: marketTypesTransport}

	gotMarketTypes, err := marketTypesAPI.ListMarketTypes(filter)
	if err != nil {
		t.Fatalf("ListMarketTypes returned error: %v", err)
	}
	if marketTypesTransport.lastMethod != "listMarketTypes" {
		t.Fatalf("expected method listMarketTypes, got %s", marketTypesTransport.lastMethod)
	}
	if !reflect.DeepEqual(gotMarketTypes, marketTypesResponse) {
		t.Fatalf("expected %#v, got %#v", marketTypesResponse, gotMarketTypes)
	}

	countriesResponse := []types.CountryWrapper{{Country: "GB", MarketCount: 2}}
	countriesPayload, _ := json.Marshal(countriesResponse)
	countriesTransport := &stubTransport{response: countriesPayload}
	countriesAPI := &API{Client: countriesTransport}

	gotCountries, err := countriesAPI.ListCountries(filter)
	if err != nil {
		t.Fatalf("ListCountries returned error: %v", err)
	}
	if countriesTransport.lastMethod != "listCountries" {
		t.Fatalf("expected method listCountries, got %s", countriesTransport.lastMethod)
	}
	if !reflect.DeepEqual(gotCountries, countriesResponse) {
		t.Fatalf("expected %#v, got %#v", countriesResponse, gotCountries)
	}

	venuesResponse := []types.VenueWrapper{{Venue: "Wembley", MarketCount: 3}}
	venuesPayload, _ := json.Marshal(venuesResponse)
	venuesTransport := &stubTransport{response: venuesPayload}
	venuesAPI := &API{Client: venuesTransport}

	gotVenues, err := venuesAPI.ListVenues(filter)
	if err != nil {
		t.Fatalf("ListVenues returned error: %v", err)
	}
	if venuesTransport.lastMethod != "listVenues" {
		t.Fatalf("expected method listVenues, got %s", venuesTransport.lastMethod)
	}
	if !reflect.DeepEqual(gotVenues, venuesResponse) {
		t.Fatalf("expected %#v, got %#v", venuesResponse, gotVenues)
	}
}
