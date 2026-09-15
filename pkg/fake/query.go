// Copyright 2022 Guy Barden
// store.go - loads or creates a new store of fixture match odds and price information

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package fake

import (
	"fmt"

	"github.com/guysports/go-betfair-api/pkg/types"
)

type (
	FakeQuery struct {
		InjectListEventsError          bool
		InjectListMarketCatalogueError bool
		InjectListMarketBookError      bool
		AppendPrices                   bool
	}
)

func (f *FakeQuery) ListEvents(filter *types.MarketFilter) ([]types.EventWrapper, error) {
	if f.InjectListEventsError {
		return nil, fmt.Errorf("error listing events")
	}

	// "event": {
	// 	"id": "31317592",
	// 	"name": "Augsburg v Mainz",
	// 	"countryCode": "DE",
	// 	"timezone": "GMT",
	// 	"openDate": "2022-04-06T16:30:00.000Z"
	// },
	// "marketCount": 15

	events := []types.EventWrapper{
		{
			Event: &types.Detail{
				ID:          "fixture1",
				Name:        "Mainz v Dortmund",
				CountryCode: "DE",
				TimeZone:    "GMT",
				OpenDate:    "2022-04-06T16:30:00.000Z",
			},
			MarketCount: 1,
		},
	}
	return events, nil
}
func (f *FakeQuery) ListMarketCatalogue(filter *types.MarketFilter, numEvts int, marketType []string) ([]types.MarketCatalogueWrapper, error) {
	if f.InjectListMarketCatalogueError {
		return nil, fmt.Errorf("error listing marketcatalogue")
	}

	// {
	// 	"marketId": "1.195693926",
	// 	"marketName": "Match Odds",
	// 	"totalMatched": 53933.43,
	// 	"runners": [{
	// 		"selectionId": 64374,
	// 		"runnerName": "Mainz",
	// 		"handicap": 0.0,
	// 		"sortPriority": 1,
	// 		"metadata": {
	// 			"runnerId": "64374"
	// 		}
	// 	}, {
	// 		"selectionId": 44785,
	// 		"runnerName": "Dortmund",
	// 		"handicap": 0.0,
	// 		"sortPriority": 2,
	// 		"metadata": {
	// 			"runnerId": "44785"
	// 		}
	// 	}, {
	// 		"selectionId": 58805,
	// 		"runnerName": "The Draw",
	// 		"handicap": 0.0,
	// 		"sortPriority": 3,
	// 		"metadata": {
	// 			"runnerId": "58805"
	// 		}
	// 	}]
	// }
	catalog := []types.MarketCatalogueWrapper{
		{
			MarketId:     "1.195693926",
			MarketName:   "Match Odds",
			TotalMatched: 1000.00,
			Selections: []types.Selection{
				{
					SelectionId: 64374,
					Name:        "Mainz",
					Handicap:    0.0,
					Ranking:     1,
					Metadata: map[string]string{
						"runnerId": "64374",
					},
				},
				{
					SelectionId: 44785,
					Name:        "Dortmund",
					Handicap:    0.0,
					Ranking:     2,
					Metadata: map[string]string{
						"runnerId": "64374",
					},
				},
				{
					SelectionId: 58805,
					Name:        "The Draw",
					Handicap:    0.0,
					Ranking:     3,
					Metadata: map[string]string{
						"runnerId": "58805",
					},
				},
			},
		},
	}
	return catalog, nil
}
func (f *FakeQuery) ListMarketBook(marketIds []string, priceProjection *types.PriceProjection, orderProjection string, matchProjection string) ([]types.MarketBookWrapper, error) {
	if f.InjectListMarketBookError {
		return nil, fmt.Errorf("error listing marketbook")
	}
	firstHomeBackOdds := []types.Odds{{Price: 3.7, Size: 777.45}, {Price: 3.65, Size: 1164.38}, {Price: 3.6, Size: 1019.41}}
	firstHomeLayOdds := []types.Odds{{Price: 3.75, Size: 718.45}, {Price: 3.8, Size: 1145.15}, {Price: 3.85, Size: 1505.97}}
	secondHomeBackOdds := []types.Odds{{Price: 3.75, Size: 777.45}, {Price: 3.7, Size: 1164.38}, {Price: 3.65, Size: 1019.41}}
	secondHomeLayOdds := []types.Odds{{Price: 3.8, Size: 718.45}, {Price: 3.85, Size: 1145.15}, {Price: 3.95, Size: 1505.97}}
	firstAwayBackOdds := []types.Odds{{Price: 2.86, Size: 537.2}, {Price: 2.84, Size: 167.17}, {Price: 2.82, Size: 833.53}}
	firstAwayLayOdds := []types.Odds{{Price: 2.9, Size: 194.9}, {Price: 2.92, Size: 1168.09}, {Price: 2.94, Size: 526.12}}
	secondAwayBackOdds := []types.Odds{{Price: 2.84, Size: 537.2}, {Price: 2.82, Size: 167.17}, {Price: 2.8, Size: 833.53}}
	secondAwayLayOdds := []types.Odds{{Price: 2.92, Size: 194.9}, {Price: 2.94, Size: 1168.09}, {Price: 2.96, Size: 526.12}}
	drawBackOdds := []types.Odds{{Price: 2.56, Size: 70.4}, {Price: 2.54, Size: 1882.05}, {Price: 2.52, Size: 2229.28}}
	drawLayOdds := []types.Odds{{Price: 2.6, Size: 171.1}, {Price: 2.62, Size: 904.11}, {Price: 2.64, Size: 73.12}}

	// {
	// 	"marketId": "1.195693926",
	// 	"isMarketDataDelayed": true,
	// 	"status": "OPEN",
	// 	"betDelay": 5,
	// 	"bspReconciled": false,
	// 	"complete": true,
	// 	"inplay": true,
	// 	"numberOfWinners": 1,
	// 	"numberOfRunners": 3,
	// 	"numberOfActiveRunners": 3,
	// 	"lastMatchTime": "2022-03-16T18:34:06.924Z",
	// 	"totalMatched": 946450.74,
	// 	"totalAvailable": 173585.43,
	// 	"crossMatching": true,
	// 	"runnersVoidable": false,
	// 	"version": 4417661231,
	// 	"runners": [{
	// 		"selectionId": 64374,
	// 		"handicap": 0.0,
	// 		"status": "ACTIVE",
	// 		"lastPriceTraded": 3.75,
	// 		"totalMatched": 0.0,
	// 		"ex": {
	// 			"availableToBack": [{
	// 				"price": 3.7,
	// 				"size": 777.45
	// 			}, {
	// 				"price": 3.65,
	// 				"size": 1164.38
	// 			}, {
	// 				"price": 3.6,
	// 				"size": 1019.41
	// 			}],
	// 			"availableToLay": [{
	// 				"price": 3.75,
	// 				"size": 718.89
	// 			}, {
	// 				"price": 3.8,
	// 				"size": 1145.15
	// 			}, {
	// 				"price": 3.85,
	// 				"size": 1505.97
	// 			}],
	// 			"tradedVolume": []
	// 		}
	// 	}, {
	// 		"selectionId": 44785,
	// 		"handicap": 0.0,
	// 		"status": "ACTIVE",
	// 		"lastPriceTraded": 2.9,
	// 		"totalMatched": 0.0,
	// 		"ex": {
	// 			"availableToBack": [{
	// 				"price": 2.88,
	// 				"size": 537.2
	// 			}, {
	// 				"price": 2.86,
	// 				"size": 167.17
	// 			}, {
	// 				"price": 2.84,
	// 				"size": 833.53
	// 			}],
	// 			"availableToLay": [{
	// 				"price": 2.9,
	// 				"size": 194.9
	// 			}, {
	// 				"price": 2.92,
	// 				"size": 1168.09
	// 			}, {
	// 				"price": 2.94,
	// 				"size": 526.12
	// 			}],
	// 			"tradedVolume": []
	// 		}
	// 	}, {
	// 		"selectionId": 58805,
	// 		"handicap": 0.0,
	// 		"status": "ACTIVE",
	// 		"lastPriceTraded": 2.6,
	// 		"totalMatched": 0.0,
	// 		"ex": {
	// 			"availableToBack": [{
	// 				"price": 2.58,
	// 				"size": 70.4
	// 			}, {
	// 				"price": 2.56,
	// 				"size": 1882.05
	// 			}, {
	// 				"price": 2.54,
	// 				"size": 2229.28
	// 			}],
	// 			"availableToLay": [{
	// 				"price": 2.6,
	// 				"size": 171.1
	// 			}, {
	// 				"price": 2.62,
	// 				"size": 904.11
	// 			}, {
	// 				"price": 2.64,
	// 				"size": 73.12
	// 			}],
	// 			"tradedVolume": []
	// 		}
	// 	}]
	// }]
	hbOdds := firstHomeBackOdds
	hlOdds := firstHomeLayOdds
	abOdds := firstAwayBackOdds
	alOdds := firstAwayLayOdds
	if f.AppendPrices {
		hbOdds = secondHomeBackOdds
		hlOdds = secondHomeLayOdds
		abOdds = secondAwayBackOdds
		alOdds = secondAwayLayOdds
	}
	marketbook := []types.MarketBookWrapper{
		{
			MarketId:            "1.195693926",
			IsMarketDataDelayed: true,
			Status:              "OPEN",
			BetDelay:            5,
			BspReconciled:       false,
			Complete:            true,
			Inplay:              false,
			NumberOfWinners:     1,
			NumberOfRunners:     3,
			LastMatchTime:       "2022-03-16T18:34:06.924Z",
			TotalMatched:        946450.74,
			TotalAvailable:      173585.43,
			CrossMatching:       true,
			RunnersVoidable:     false,
			Version:             4417661231,
			Runners: []types.Runner{
				{
					SelectionID: 64374,
					Handicap:    0.0,
					Status:      "ACTIVE",
					Exchange: types.ExchangePrices{
						AvailableToBack: hbOdds,
						AvailableToLay:  hlOdds,
					},
				},
				{
					SelectionID: 44785,
					Handicap:    0.0,
					Status:      "ACTIVE",
					Exchange: types.ExchangePrices{
						AvailableToBack: abOdds,
						AvailableToLay:  alOdds,
					},
				},
				{
					SelectionID: 58805,
					Handicap:    0.0,
					Status:      "ACTIVE",
					Exchange: types.ExchangePrices{
						AvailableToBack: drawBackOdds,
						AvailableToLay:  drawLayOdds,
					},
				},
			},
		},
	}
	return marketbook, nil
}
