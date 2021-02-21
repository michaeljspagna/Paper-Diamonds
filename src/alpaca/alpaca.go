package main

import (
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var cfg *Config

func getAccount() Account {
	reqURL, _ := url.Parse(fmt.Sprintf("%s/%s/account", cfg.URL, cfg.Version))
	req := &http.Request{
		Method: "GET",
		URL:    reqURL,
		Header: map[string][]string{
			"APCA-API-KEY-ID":     {cfg.KeyID},
			"APCA-API-SECRET-KEY": {cfg.Secret},
		},
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	data, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	acc := Account{}
	json.Unmarshal(data, &acc)
	return acc
}

func makeTrade(order OrderBody) {
	reqURL, _ := url.Parse(fmt.Sprintf("%s/%s/orders", cfg.URL, cfg.Version))
	reqBody := ioutil.NopCloser(strings.NewReader(fmt.Sprintf(`
		{
			"symbol":%q,
			"qty":%q,
			"side":%q,
			"type":%q,
			"time_in_force":%q
		}
	`, order.Symbol, order.Qty, order.Side, order.Type, order.TimeInForce)))

	req := &http.Request{
		Method: "POST",
		URL:    reqURL,
		Header: map[string][]string{
			"APCA-API-KEY-ID":     {cfg.KeyID},
			"APCA-API-SECRET-KEY": {cfg.Secret},
		},
		Body: reqBody,
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	data, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	fmt.Printf("Status: %s\n", res.StatusCode)
	fmt.Printf("body: %s\n", data)
}

func getConfig(cfgPath string) error {
	source, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(source, &cfg)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", *cfg)
	return nil
}

func main() {
	err := getConfig("./config.yml")
	if err != nil {
		panic(err)
	}

	order := OrderBody{
		Symbol:      "GOOGL",
		Qty:         "1",
		Side:        "buy",
		Type:        "market",
		TimeInForce: "day",
	}

	makeTrade(order)

}

// Config Struct to store data from congif.yml
type Config struct {
	URL     string `yaml:"AlpacaApiURL"`
	Version string `yaml:"AlpacaApiVersion"`
	KeyID   string `yaml:"AlpacaApiKeyId"`
	Secret  string `yaml:"AlpacaApiSecretKey"`
}

// Account holds account information
type Account struct {
	ID                    string          `json:"id"`
	AccountNumber         string          `json:"account_number"`
	CreatedAt             time.Time       `json:"created_at"`
	UpdatedAt             time.Time       `json:"updated_at"`
	DeletedAt             *time.Time      `json:"deleted_at"`
	Status                string          `json:"status"`
	Currency              string          `json:"currency"`
	Cash                  decimal.Decimal `json:"cash"`
	CashWithdrawable      decimal.Decimal `json:"cash_withdrawable"`
	TradingBlocked        bool            `json:"trading_blocked"`
	TransfersBlocked      bool            `json:"transfers_blocked"`
	AccountBlocked        bool            `json:"account_blocked"`
	ShortingEnabled       bool            `json:"shorting_enabled"`
	BuyingPower           decimal.Decimal `json:"buying_power"`
	PatternDayTrader      bool            `json:"pattern_day_trader"`
	DaytradeCount         int64           `json:"daytrade_count"`
	DaytradingBuyingPower decimal.Decimal `json:"daytrading_buying_power"`
	RegTBuyingPower       decimal.Decimal `json:"regt_buying_power"`
	Equity                decimal.Decimal `json:"equity"`
	LastEquity            decimal.Decimal `json:"last_equity"`
	Multiplier            string          `json:"multiplier"`
	InitialMargin         decimal.Decimal `json:"initial_margin"`
	MaintenanceMargin     decimal.Decimal `json:"maintenance_margin"`
	LastMaintenanceMargin decimal.Decimal `json:"last_maintenance_margin"`
	LongMarketValue       decimal.Decimal `json:"long_market_value"`
	ShortMarketValue      decimal.Decimal `json:"short_market_value"`
	PortfolioValue        decimal.Decimal `json:"portfolio_value"`
}

//OrderBody holds information for orders
type OrderBody struct {
	Symbol      string
	Qty         string
	Side        string
	Type        string
	TimeInForce string
}
