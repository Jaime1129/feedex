package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

const ETHUSDT = "ETHUSDT"

type TrxFeeService interface {
	GetSingleTrxFee(ctx context.Context, req *GetSingleTrxFeeRequest) (*GetSingleTrxFeeResponse, error)
}

type trxFeeService struct {
	apiKey string
}

func NewTrxService(apiKey string) TrxFeeService {
	return &trxFeeService{
		apiKey: apiKey,
	}
}

type GetSingleTrxFeeRequest struct {
	TrxHash string
}

type GetSingleTrxFeeResponse struct {
	TrxFee string `json:"trx_fee"`
}

func (c *trxFeeService) GetSingleTrxFee(ctx context.Context, req *GetSingleTrxFeeRequest) (*GetSingleTrxFeeResponse, error) {
	if req == nil {
		return nil, errors.New("nil req")
	}

	trxResp, err := c.queryTrxFee(req.TrxHash)
	if err != nil {
		return nil, err
	}

	gasUsed, err := hexToInt(trxResp.Result.GasUsed)
	if err != nil {
		return nil, err
	}
	gasPrice, err := hexToInt(trxResp.Result.EffectiveGasPrice)
	if err != nil {
		return nil, err
	}

	// calculate gas fee in eth
	gasInETH := calculateFeeInETH(gasUsed, gasPrice)

	// get block timestamp
	blockResp, err := c.queryBlock(trxResp.Result.BlockNumber)
	if err != nil {
		return nil, err
	}

	// todo: get eth price in USDT
	trxTime, err := hexToInt(blockResp.Result.Timestamp)
	if err != nil {
		return nil, err
	}

	// fetch the average price of [trxTime-60, trxTime+60]
	price, err := c.queryETHPrice(trxTime-60, trxTime+60)
	if err != nil {
		return nil, err
	}

	return &GetSingleTrxFeeResponse{
		TrxFee: gasInETH.Mul(price).String(),
	}, nil
}

type EthScanTrxResponse struct {
	Result EthScanTrxResult `json:"result"`
	Error  EthScanError     `json:"error"`
}

type EthScanTrxResult struct {
	EffectiveGasPrice string `json:"effectiveGasPrice"`
	GasUsed           string `json:"gasUsed"`
	BlockNumber       string `json:"blockNumber"`
}

type EthScanError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func hexToInt(hexStr string) (int64, error) {
	hexStr = strings.TrimPrefix(hexStr, "0x")
	// base 16 for hexadecimal, 64 bits
	decimalValue, err := strconv.ParseInt(hexStr, 16, 64)
	if err != nil {
		fmt.Println("Error converting hex to decimal:", err)
		return 0, err
	}
	return decimalValue, nil
}

func calculateFeeInETH(gasUsed int64, gasPrice int64) decimal.Decimal {
	// convert gasPrice in Wei to Eth, then multiply with gasUsed
	return decimal.NewFromInt(gasPrice).Div(decimal.NewFromInt(1e18)).Mul(decimal.NewFromInt(gasUsed))
}

func (c *trxFeeService) queryTrxFee(trxHash string) (*EthScanTrxResponse, error) {
	// send query to etherscan api
	url := fmt.Sprintf("https://api.etherscan.io/api?module=proxy&action=eth_getTransactionReceipt&txhash=%s&apikey=%s", trxHash, c.apiKey)
	fmt.Println("ethscan api url: " + url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println("ethscan api resp body: " + string(body))

	trxResp := &EthScanTrxResponse{}
	err = json.Unmarshal(body, trxResp)
	if err != nil {
		return nil, err
	}

	// check if api call returns error
	if trxResp.Error.Code != 0 {
		fmt.Println("ethscan api call returns error: " + trxResp.Error.Message)
		return nil, errors.New("ethscan api call returns error")
	}

	return trxResp, nil
}

type EthScanBlockResponse struct {
	Result EthScanBlockResult `json:"result"`
	Error  EthScanError       `json:"error"`
}

type EthScanBlockResult struct {
	Timestamp string `json:"timestamp"`
}

func (c *trxFeeService) queryBlock(blockNumber string) (*EthScanBlockResponse, error) {
	// send query to etherscan api
	url := fmt.Sprintf("https://api.etherscan.io/api?module=proxy&action=eth_getBlockByNumber&tag=%s&boolean=true&apikey=%s", blockNumber, c.apiKey)
	fmt.Println("ethscan api url: " + url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println("ethscan api resp body: " + string(body))

	trxResp := &EthScanBlockResponse{}
	err = json.Unmarshal(body, trxResp)
	if err != nil {
		return nil, err
	}

	// check if api call returns error
	if trxResp.Error.Code != 0 {
		fmt.Println("ethscan api call returns error: " + trxResp.Error.Message)
		return nil, errors.New("ethscan api call returns error")
	}

	return trxResp, nil
}

// https://api.binance.com/api/v3/klines?symbol=ETHUSDT&interval=1m&startTime=1714906101&endTime=1714906191
func (c *trxFeeService) queryETHPrice(start int64, end int64) (price decimal.Decimal, err error) {
	url := fmt.Sprintf("https://api.binance.com/api/v3/klines?symbol=%s&interval=%s&startTime=%d&endTime=%d", ETHUSDT, "1m", start*1e3, end*1e3)
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	fmt.Println("binance api resp body: " + string(body))

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
