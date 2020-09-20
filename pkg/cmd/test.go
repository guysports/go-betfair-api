package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/guysports/go-betfair-api/pkg/betting"
	"github.com/guysports/go-betfair-api/pkg/types"
	"github.com/jedib0t/go-pretty/v6/table"
)

type (
	Test struct {
		RootCAPath string `help:"Path to the RootCA certificate that Betfair signs their certificate with"`
		CertPath   string `help:"Path to the User certificate created and added in the Betfair account"`
		KeyPath    string `help:"Path to the User private key for the Betfair account"`
		User       string `help:"Username for the Betfair account"`
		Password   string `help:"Password for the Betfair account"`
		Operation  string `help:"Betfair operation to query"`
	}
)

func (t *Test) Run(globals *types.Globals) error {
	cfg := types.Config{
		CertPath: t.CertPath,
		KeyPath:  t.KeyPath,
		User:     t.User,
		Password: t.Password,
		AppKey:   globals.AppKey,
	}

	ctx, cancel := context.WithTimeout(context.Background(), types.DefaultTimeout)
	defer cancel()
	client, err := betting.NewAPI(ctx, &cfg)
	if err != nil {
		return err
	}

	_, err = client.Client.Authenticate()
	if err != nil {
		return err
	}

	switch t.Operation {
	case "listeventtypes":
		filter := types.MarketFilter{
			EventTypeIds: []string{"1", "2"},
		}
		eventTypes, err := client.ListEventTypes(&filter)
		if err != nil {
			return err
		}
		fmt.Println("Event Types returned...")
		for _, eventTypeW := range eventTypes {
			fmt.Printf("Event Type ID %s, Market %s, Number of Markets %d\n", eventTypeW.EventType.ID, eventTypeW.EventType.Name, eventTypeW.MarketCount)
		}
	case "listcompetitions":
		filter := types.MarketFilter{
			TextQuery: "Premier League",
		}
		competitions, err := client.ListCompetitions(&filter)
		if err != nil {
			return err
		}
		fmt.Println("Competitions returned...")
		for _, competitionW := range competitions {
			fmt.Printf("Competition ID %s, Market %s, Number of Markets %d, Region %s\n", competitionW.Competition.ID, competitionW.Competition.Name, competitionW.MarketCount, competitionW.Region)
		}
	case "listtimeranges":
		from := time.Now()
		to := from.Add(72 * time.Hour)
		marketsInRange, err := client.ListTimeRanges(&from, &to, &types.MarketFilter{}, "DAYS")
		if err != nil {
			return err
		}
		fmt.Println("Number of Markets for each day returned...")
		for _, market := range marketsInRange {
			fmt.Printf("From %s, To %s, Number of Markets %d\n", market.Range.From, market.Range.To, market.MarketCount)
		}

	case "listevents":
		to := time.Now().Add(7 * 24 * time.Hour).Format(time.RFC3339)
		eventsBefore := types.TimeRange{
			To: to,
		}
		filter := types.MarketFilter{
			EventTypeIds:    []string{"1"},
			CompetitionIds:  []string{"10932509"},
			MarketStartTime: &eventsBefore,
		}
		events, err := client.ListEvents(&filter)
		if err != nil {
			return err
		}
		fmt.Println("Premier League Fixtures in next 7 days")
		for _, event := range events {
			fmt.Printf("Event Type ID %s, Fixture %s, Start Time %s, Number of Markets %d\n", event.Event.ID, event.Event.Name, event.Event.OpenDate, event.MarketCount)
		}

	case "listmarkettypes":
		event, err := getEventID(client)
		if err != nil {
			return err
		}

		filter := types.MarketFilter{
			EventIds: []string{event.ID},
		}
		markets, err := client.ListMarketTypes(&filter)
		if err != nil {
			return err
		}
		fmt.Printf("Markets for fixtures %s\n", event.Name)
		for _, market := range markets {
			fmt.Printf("Market %s, Number of Markets %d\n", market.MarketType, market.MarketCount)
		}

	case "listcountries":
		countries, err := client.ListCountries(&types.MarketFilter{})
		if err != nil {
			return err
		}
		fmt.Println("Countries and number of Markets")
		for _, country := range countries {
			fmt.Printf("Country Code %s, Number of Markets %d\n", country.Country, country.MarketCount)
		}
	case "listvenues":
		venues, err := client.ListVenues(&types.MarketFilter{MarketCountries: []string{"GB"}})
		if err != nil {
			return err
		}
		fmt.Println("Racing Venues today")
		for _, venue := range venues {
			fmt.Printf("Venue %s, Number of Markets %d\n", venue.Venue, venue.MarketCount)
		}
	case "listmarketcatalogue":
		event, err := getEventID(client)
		if err != nil {
			return err
		}
		filter := types.MarketFilter{
			EventIds:        []string{event.ID},
			MarketTypeCodes: []string{"MATCH_ODDS"},
		}

		catalogue, err := client.ListMarketCatalogue(&filter, 1, []string{"RUNNER_METADATA"})
		if err != nil {
			return err
		}
		fmt.Printf("Selections for %s market in fixture %s (%s)\n", catalogue[0].MarketName, event.Name, catalogue[0].MarketId)
		for _, entry := range catalogue[0].Selections {
			fmt.Printf("Selection ID %d, Selection %s, \n", entry.SelectionId, entry.Name)
		}

	case "listmarketbook":
		event, err := getEventID(client)
		if err != nil {
			return err
		}
		filter := types.MarketFilter{
			EventIds:        []string{event.ID},
			MarketTypeCodes: []string{"MATCH_ODDS"},
		}
		catalogue, err := client.ListMarketCatalogue(&filter, 1, []string{"RUNNER_METADATA"})
		if err != nil {
			return err
		}
		marketBook, err := client.ListMarketBook([]string{catalogue[0].MarketId}, &types.PriceProjection{PriceData: []string{"EX_BEST_OFFERS"}}, "EXECUTABLE", "ROLLED_UP_BY_AVG_PRICE")
		if err != nil {
			return err
		}
		fmt.Printf("Market Book for match odds for fixture %s\n", catalogue[0].MarketName)
		for _, runner := range marketBook[0].Runners {
			tw := table.NewWriter()
			name := findSelectionName(catalogue[0].Selections, runner.SelectionID)
			tw.SetTitle(fmt.Sprintf("Prices for selection %s", name))
			tw.AppendHeader(table.Row{"AvailableToBack", "", "AvailableToLay", ""})
			tw.AppendHeader(table.Row{"Odds", "Amount", "Odds", "Amount"})
			for i := range runner.Exchange.AvailableToBack {
				tw.AppendRow([]interface{}{runner.Exchange.AvailableToBack[i].Price, runner.Exchange.AvailableToBack[i].Size, runner.Exchange.AvailableToLay[i].Price, runner.Exchange.AvailableToLay[i].Size})
			}
			fmt.Println(tw.Render())
		}
	case "listrunnerbook":
		event, err := getEventID(client)
		if err != nil {
			return err
		}
		filter := types.MarketFilter{
			EventIds:        []string{event.ID},
			MarketTypeCodes: []string{"MATCH_ODDS"},
		}
		catalogue, err := client.ListMarketCatalogue(&filter, 1, []string{"RUNNER_METADATA"})
		if err != nil {
			return err
		}
		runnerBook, err := client.ListRunnerBook(catalogue[0].MarketId, catalogue[0].Selections[0].SelectionId, &types.PriceProjection{PriceData: []string{"EX_BEST_OFFERS"}}, "EXECUTABLE", "ROLLED_UP_BY_AVG_PRICE")
		if err != nil {
			return err
		}
		fmt.Printf("Runner Book for match odds for fixture %s\n", catalogue[0].MarketName)
		for _, runner := range runnerBook[0].Runners {
			tw := table.NewWriter()
			name := findSelectionName(catalogue[0].Selections, runner.SelectionID)
			tw.SetTitle(fmt.Sprintf("Prices for selection %s", name))
			tw.AppendHeader(table.Row{"AvailableToBack", "", "AvailableToLay", ""})
			tw.AppendHeader(table.Row{"Odds", "Amount", "Odds", "Amount"})
			for i := range runner.Exchange.AvailableToBack {
				tw.AppendRow([]interface{}{runner.Exchange.AvailableToBack[i].Price, runner.Exchange.AvailableToBack[i].Size, runner.Exchange.AvailableToLay[i].Price, runner.Exchange.AvailableToLay[i].Size})
			}
			fmt.Println(tw.Render())
		}
	default:
		fmt.Printf("The operation %s is not recognised\n", t.Operation)
	}
	return nil
}

func getEventID(client betting.APIInterface) (*types.Detail, error) {
	to := time.Now().Add(7 * 24 * time.Hour).Format(time.RFC3339)
	eventsBefore := types.TimeRange{
		To: to,
	}
	filter := types.MarketFilter{
		CompetitionIds:  []string{"10932509"},
		MarketStartTime: &eventsBefore,
	}
	events, err := client.ListEvents(&filter)
	if err != nil {
		return nil, err
	}
	return events[0].Event, nil
}

func findSelectionName(selections []types.Selection, id int) string {
	for _, selection := range selections {
		if selection.SelectionId == id {
			return selection.Name
		}
	}
	return "Selection Not Found"
}
