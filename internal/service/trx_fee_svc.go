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
	// send query to etherscan api
	url := fmt.Sprintf("https://api.etherscan.io/api?module=proxy&action=eth_getTransactionReceipt&txhash=%s&apikey=%s", req.TrxHash, c.apiKey)
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

	// todo: get eth price in USDT

	return &GetSingleTrxFeeResponse{
		TrxFee: gasInETH.String(),
	}, nil
}

type EthScanTrxResponse struct {
	Result EthScanTrxResult `json:"result"`
	Error  EthScanError     `json:"error"`
}

type EthScanTrxResult struct {
	EffectiveGasPrice string `json:"effectiveGasPrice"`
	GasUsed           string `json:"gasUsed"`
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
