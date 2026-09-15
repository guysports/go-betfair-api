// Copyright 2022 Guy Barden
// MarketQuery provides the league information with optional time and odds information
// to retrieve from Betfair

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package access

import (
	"encoding/json"
	"io/ioutil"

	"github.com/guysports/go-betfair-api/pkg/types"
)

type (
	// MarketQuery parameters for looking at fixtures to obtain pricing information on
	MarketQuery struct {
		LeagueIds         []string `json:"leagueids"`
		MinDaysToFixtures int      `json:"mindays"`
		MaxDaysToFixtures int      `json:"maxdays"`
		MinOdds           float32  `json:"minodds"`
		MaxOdds           float32  `json:"maxodds"`
	}

	QueryInterface interface {
		ListEvents(filter *types.MarketFilter) ([]types.EventWrapper, error)
		ListMarketCatalogue(filter *types.MarketFilter, numEvts int, marketType []string) ([]types.MarketCatalogueWrapper, error)
		ListMarketBook(marketIds []string, priceProjection *types.PriceProjection, orderProjection string, matchProjection string) ([]types.MarketBookWrapper, error)
	}
)

func NewQuery(queryPath string) (*MarketQuery, error) {
	query := MarketQuery{}
	queryData, err := ioutil.ReadFile(queryPath)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(queryData, &query); err != nil {
		return nil, err
	}

	return &query, nil
}
