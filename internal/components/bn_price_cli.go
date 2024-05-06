package components

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/shopspring/decimal"
	"go.uber.org/ratelimit"
)

const ETHUSDT = "ETHUSDT"
const INTERVAL_1MIN = "1m"
const INTERVAL_12HOUR = "12h"

type BnPriceCli interface {
	QueryETHPrice(start int64, end int64, interval string) (price decimal.Decimal, err error)
}

type bnPriceCli struct {
	rl ratelimit.Limiter
}

func NewBnPriceCLi() BnPriceCli {
	return &bnPriceCli{
		rl: ratelimit.New(10),
	}
}

func (c *bnPriceCli) QueryETHPrice(start int64, end int64, interval string) (price decimal.Decimal, err error) {
	c.rl.Take()
	url := fmt.Sprintf("https://api.binance.com/api/v3/klines?symbol=%s&interval=%s&startTime=%d&endTime=%d", ETHUSDT, interval, start*1e3, end*1e3)
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	log.Println("binance api resp body: " + string(body))

	var rawCandlesticks [][]interface{}
	err = json.Unmarshal(body, &rawCandlesticks)
	if err != nil {
		return
	}

	if len(rawCandlesticks) == 0 {
		err = errors.New("price not found")
		return
	}

	for _, c := range rawCandlesticks {
		closePrice, _ := decimal.NewFromString(c[4].(string))
		openPrice, _ := decimal.NewFromString(c[1].(string))
		price = price.Add(closePrice).Add(openPrice)
	}
	price = price.Div(decimal.NewFromInt(int64(len(rawCandlesticks)) * 2))
	return
}
