package types

import (
	"time"
)

const (
	DefaultTimeout = 20 * time.Second
)

type (
	TransportInterface interface {
		Authenticate() (*Authenticate, error)
		SetSessionKey(key string)
		Do(id int, method string, filter *MarketFilter, additionalParams interface{}) ([]byte, error)
	}
	// Login and Authenticate
	Globals struct {
		AppKey string
	}

	Config struct {
		RootCAPath string
		CertPath   string
		KeyPath    string
		User       string
		Password   string
		AppKey     string
	}

	Params struct {
		Filter              *MarketFilter      `json:"filter,omitempty"`
		Instructions        []PlaceInstruction `json:"instructions,omitempty"`
		Granularity         *string            `json:"granularity,omitempty"`
		MaxResults          int                `json:"maxResults,omitempty"`
		MarketId            string             `json:"marketId,omitempty"`
		MarketIds           []string           `json:"marketIds,omitempty"`
		SelectionId         int                `json:"selectionId,omitempty"`
		PriceProjection     *PriceProjection   `json:"priceProjection,omitempty"`
		OrderProjection     string             `json:"orderProjection,omitempty"`
		MatchProjection     string             `json:"matchProjection,omitempty"`
		MarketProjection    []string           `json:"marketProjection,omitempty"`
		Locale              string             `json:"locale,omitempty"`
		CustomerRef         string             `json:"customerRef,omitempty"`
		CustomerStrategyRef string             `json:"customerStrategyRef,omitempty"`
		DateRange           TimeRange          `json:"dateRange"`
	}

	JsonError struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	JsonRPC struct {
		JsonRPC   string `json:"jsonrpc"`
		Method    string `json:"method"`
		RPCParams Params `json:"params"`
		ID        int    `json:"id"`
	}

	JsonRPCResponse struct {
		JsonRPC string      `json:"jsonrpc"`
		Result  interface{} `json:"result"`
		Error   *JsonError  `json:"error,omitempty"`
		ID      int         `json:"id"`
	}

	Authenticate struct {
		SessionToken string `json:"sessionToken"`
		LoginStatus  string `json:"loginStatus"`
	}

	// Betting API Market Information
	MarketBettingType string
	OrderStatus       string

	TimeRange struct {
		From string `json:"from,omitempty"`
		To   string `json:"to,omitempty"`
	}

	MarketFilter struct {
		TextQuery          string              `json:"textQuery,omitempty"`
		EventTypeIds       []string            `json:"eventTypeIds,omitempty"`
		EventIds           []string            `json:"eventIds,omitempty"`
		CompetitionIds     []string            `json:"competitionIds,omitempty"`
		MarketIds          []string            `json:"marketIds,omitempty"`
		Venues             []string            `json:"venues,omitempty"`
		BspOnly            bool                `json:"bspOnly,omitempty"`
		TurnInPlayEnabled  bool                `json:"turnInPlayEnabled,omitempty"`
		InPlayOnly         bool                `json:"inPlayOnly,omitempty"`
		MarketBettingTypes []MarketBettingType `json:"marketBettingTypes,omitempty"`
		MarketTypeCodes    []string            `json:"marketTypeCodes,omitempty"`
		MarketStartTime    *TimeRange          `json:"marketStartTime,omitempty"`
		MarketCountries    []string            `json:"marketCountries"`
		WithOrders         []OrderStatus       `json:"withOrders,omitempty"`
		RaceTypes          []string            `json:"raceTypes,omitempty"`
	}

	MarketFilterParams struct {
		Granularity      string
		MaxResults       int
		MarketId         string
		MarketIds        []string
		SelectionId      int
		MarketProjection []string
		PriceProjection  *PriceProjection
		OrderProjection  string
		MatchProjection  string
		DateRange        *TimeRange
	}

	PlaceInstructionParams struct {
		MarketID            string
		Instructions        []PlaceInstruction
		CustomerRef         string
		CustomerStrategyRef string
	}

	PriceProjection struct {
		PriceData []string `json:"priceData"`
	}

	Detail struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		CountryCode string `json:"countryCode,omitempty"`
		TimeZone    string `json:"timezone,omitempty"`
		OpenDate    string `json:"openDate,omitempty"`
	}

	Selection struct {
		SelectionId int               `json:"selectionId"`
		Name        string            `json:"runnerName"`
		Handicap    float32           `json:"handicap"`
		Ranking     int               `json:"sortPriority"`
		Metadata    map[string]string `json:"metadata"`
	}

	EventTypeWrapper struct {
		EventType   *Detail `json:"eventType"`
		MarketCount int     `json:"marketCount"`
	}

	CompetitionWrapper struct {
		Competition *Detail `json:"competition"`
		MarketCount int     `json:"marketCount"`
		Region      string  `json:"competitionRegion"`
	}

	RangeWrapper struct {
		Range       *TimeRange `json:"timeRange"`
		MarketCount int        `json:"marketCount"`
	}

	EventWrapper struct {
		Event       *Detail `json:"event"`
		MarketCount int     `json:"marketCount"`
	}

	MarketTypeWrapper struct {
		MarketType  string `json:"marketType"`
		MarketCount int    `json:"marketCount"`
	}

	CountryWrapper struct {
		Country     string `json:"countryCode"`
		MarketCount int    `json:"marketCount"`
	}

	VenueWrapper struct {
		Venue       string `json:"venue"`
		MarketCount int    `json:"marketCount"`
	}

	MarketCatalogueWrapper struct {
		MarketId     string      `json:"marketId"`
		MarketName   string      `json:"marketName"`
		TotalMatched float64     `json:"totalMatched"`
		Selections   []Selection `json:"runners"`
	}

	MarketBookWrapper struct {
		MarketId            string   `json:"marketId"`
		IsMarketDataDelayed bool     `json:"isMarketDataDelayed"`
		Status              string   `json:"status"`
		BetDelay            int      `json:"betDelay"`
		BspReconciled       bool     `json:"bspReconciled"`
		Complete            bool     `json:"complete"`
		Inplay              bool     `json:"inplay"`
		NumberOfWinners     int      `json:"numberOfWinners"`
		NumberOfRunners     int      `json:"numberOfRunners"`
		LastMatchTime       string   `json:"lastMatchTime"`
		TotalMatched        float32  `json:"totalMatched"`
		TotalAvailable      float32  `json:"totalAvailable"`
		CrossMatching       bool     `json:"crossMatching"`
		RunnersVoidable     bool     `json:"runnersVoidable"`
		Version             int64    `json:"version"`
		Runners             []Runner `json:"runners"`
	}

	Runner struct {
		SelectionID     int            `json:"selectionId"`
		Handicap        float32        `json:"handicap"`
		Status          string         `json:"status"`
		LastPriceTraded float32        `json:"lastPriceTraded"`
		TotalMatched    float32        `json:"totalMatched"`
		Exchange        ExchangePrices `json:"ex"`
	}

	ExchangePrices struct {
		AvailableToBack []Odds `json:"availableToBack"`
		AvailableToLay  []Odds `json:"availableToLay"`
		TradedVolume    []Odds `json:"tradedVolume"`
	}

	Odds struct {
		Price float32 `json:"price"`
		Size  float32 `json:"size"`
	}

	LimitOrder struct {
		Size            float32
		Price           float32
		PersistanceType string
		TimeInForce     string
		MinFillSize     float32
		BetTargetType   string
		BetTargetSize   float32
	}

	PlaceInstruction struct {
		OrderType        string  `json:"orderType"`
		SelectionId      int     `json:"selectionId"`
		Handicap         float32 `json:"handicap"`
		Side             string  `json:"side"`
		LimitOrder       Price   `json:"limitOrder"`
		CustomerOrderRef string  `json:"customeOrderRef"`
	}

	CurrentOrder struct {
		BetId               string  `json:"betId"`
		MarketId            string  `json:"marketId"`
		SelectionId         string  `json:"selectionId"`
		Handicap            float32 `json:"handicap"`
		PriceSize           Price   `json:"priceSize"`
		BspLiability        float32 `json:"bspLiability"`
		Side                string  `json:"side"`
		Status              string  `json:"status"`
		PersistanceType     string  `json:"persistanceType"`
		OrderType           string  `json:"orderType"`
		PlacedDate          string  `json:"placedDate"`
		MatchedDate         string  `json:"matchedDate"`
		AveragePriceMatched float32 `json:"averagePriceMatched"`
		SizeMatched         float32 `json:"sizeMatched"`
		SizeRemaining       float32 `json:"sizeRemaining"`
		SizeLapsed          float32 `json:"sizeLapsed"`
		SizeCancelled       float32 `json:"sizeCancelled"`
		SizeVoided          float32 `json:"sizeVoided"`
		RegulatorCode       string  `json:"regulatorCode"`
	}

	Price struct {
		Price float32 `json:"price"`
		Size  float32 `json:"size"`
	}

	CurrentOrdersWrapper struct {
		Orders        []CurrentOrder `json:"currentOrders"`
		MoreAvailable bool           `json:"moreAvailable"`
	}

	PlaceInstructionReport struct {
		Status           string           `json:"status"`
		Instruction      PlaceInstruction `json:"instruction"`
		CustomerOrderRef string           `json:"customerOrderRef"`
		OrderType        string           `json:"orderType"`
		Side             string           `json:"side"`
	}
	PlaceExecutionReport struct {
		Status              string                   `json:"status"`
		CustomerRef         string                   `json:"customerRef"`
		MarketID            string                   `json:"marketId"`
		InstructionReports  []PlaceInstructionReport `json:"instructionReports"`
		BetId               string                   `json:"betId"`
		PlacedDate          string                   `json:"placedDate"`
		AveragePriceMatched float32                  `json:"averagePriceMatched"`
		SizeMatched         float32                  `json:"sizeMatched"`
		OrderStatus         string                   `json:"orderStatus"`
	}
)
