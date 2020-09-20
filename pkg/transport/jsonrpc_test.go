package transport

import (
	"reflect"
	"testing"

	"github.com/guysports/go-betfair-api/pkg/types"
)

func strPtr(str string) *string {
	return &str
}

func Test_createParams(t *testing.T) {
	type args struct {
		filter       *types.MarketFilter
		marketParams *types.MarketFilterParams
	}
	tests := []struct {
		name string
		args args
		want types.Params
	}{
		{
			name: "filter only",
			args: args{
				filter: &types.MarketFilter{
					EventIds: []string{"1"},
				},
			},
			want: types.Params{
				Filter: &types.MarketFilter{
					EventIds: []string{"1"},
				},
				Locale: "en",
			},
		},
		{
			name: "market parameters supplied",
			args: args{
				marketParams: &types.MarketFilterParams{
					Granularity:      "DAY",
					MaxResults:       1,
					MarketIds:        []string{"123", "456", "678"},
					MarketProjection: []string{"EVENT"},
					PriceProjection: &types.PriceProjection{
						PriceData: []string{"EX_BEST_OFFERS"},
					},
					OrderProjection: "EXECUTABLE",
					MatchProjection: "ROLLED_UP_BY_AVG_PRICE",
				},
			},
			want: types.Params{
				Granularity:      strPtr("DAY"),
				MaxResults:       1,
				MarketIds:        []string{"123", "456", "678"},
				MarketProjection: []string{"EVENT"},
				PriceProjection: &types.PriceProjection{
					PriceData: []string{"EX_BEST_OFFERS"},
				},
				OrderProjection: "EXECUTABLE",
				MatchProjection: "ROLLED_UP_BY_AVG_PRICE",
				Locale:          "en",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createParams(tt.args.filter, tt.args.marketParams); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createParams() = %v, want %v", got, tt.want)
			}
		})
	}
}
