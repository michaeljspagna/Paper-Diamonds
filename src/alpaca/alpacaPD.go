package main

import (
	"fmt"
	"github.com/alpacahq/alpaca-trade-api-go/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/common"
	"github.com/shopspring/decimal"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

var (
	cfg    *Config
	client *alpaca.Client
	acct   *alpaca.Account
)

//loadConfig loads cofnig data into global struct
func loadConfig(cfgPath string) error {
	source, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		panic(err)
		return err
	}
	err = yaml.Unmarshal(source, &cfg)
	if err != nil {
		panic(err)
		return err
	}
	return nil
}

func loadClient() {
	client = alpaca.NewClient(common.Credentials())
}

func loadAccount() {
	alpacaAccount, err := client.GetAccount()
	acct = alpacaAccount
	if err != nil {
		panic(err)
	}
}

func checkCost(ticker string) (decimal.Decimal, error) {
	quote, err := client.GetLastQuote(ticker)
	if err != nil {
		return decimal.NewFromInt(-1), err
	}
	return decimal.NewFromFloat32(quote.Last.AskPrice), nil
}

func tickerAndQty(tickers []string, cashToSpend decimal.Decimal, maxPrice decimal.Decimal) (string, decimal.Decimal) {
	var bestTicker string = ""
	var bestQty decimal.Decimal = decimal.NewFromInt(-1.0)
	for _, ticker := range tickers {
		quote, err := checkCost(ticker)
		if quote.GreaterThanOrEqual(maxPrice.Div(decimal.NewFromInt(2))) || err != nil {
			continue
		}
		bestTicker = ticker
		bestQty = cashToSpend.Div(quote)
		break
	}
	return bestTicker, bestQty
}

func buy(ticker string, quantity int64) *alpaca.Order {
	var request alpaca.PlaceOrderRequest
	var trailPercent = decimal.NewFromFloat(0.9)
	tickerPointer := &ticker
	request.AccountID = acct.ID
	request.AssetKey = tickerPointer
	request.Qty = decimal.NewFromInt(quantity)
	request.Side = alpaca.Buy
	request.Type = alpaca.TrailingStop
	request.TimeInForce = alpaca.Day
	request.TrailPercent = &trailPercent
	order, err := client.PlaceOrder(request)
	if err != nil {
		panic(err)
	}
	loadAccount()
	return order
}

func cancelAll() {
	err := client.CancelAllOrders()
	if err != nil {
		panic err
	}
}

func closeAllPositions() {
	err := client.CloseAllPositions()
	if err != nil {
		panic err
	}
}

func run(percentToSpend float32, minimumSharecount int64) {
	percent := decimal.NewFromFloat32(percentToSpend)
	cashToSpend := acct.Cash.Mul(percent)
	maxPricePerShare := cashToSpend.Div(decimal.NewFromInt(minimumSharecount))
	fmt.Println(maxPricePerShare)
	loadAccount()
}

func init() {
	cfgErr := loadConfig("./config.yml")
	if cfgErr != nil {
		panic(cfgErr)
	}
	os.Setenv(common.EnvApiKeyID, cfg.KeyID)
	os.Setenv(common.EnvApiSecretKey, cfg.Secret)
	alpaca.SetBaseUrl(cfg.URL)
	loadClient()
	loadAccount()

}

func main() {
	ticker, qty := tickerAndQty([]string{"AAPL", "GOOG", "PLUG"}, decimal.NewFromInt(500), decimal.NewFromInt(200))
	fmt.Println(ticker)
	fmt.Println(qty)
}

// Config Struct to store data from congif.yml
type Config struct {
	URL     string `yaml:"AlpacaApiURL"`
	Version string `yaml:"AlpacaApiVersion"`
	KeyID   string `yaml:"AlpacaApiKeyId"`
	Secret  string `yaml:"AlpacaApiSecretKey"`
}
