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
	txHash string
}

type GetSingleTrxFeeResponse struct {
	trxFee string
}

func (c *trxFeeService) GetSingleTrxFee(ctx context.Context, req *GetSingleTrxFeeRequest) (*GetSingleTrxFeeResponse, error) {
	if req == nil {
		return nil, errors.New("nil req")
	}
	// send query to etherscan api
	url := fmt.Sprintf("https://api.etherscan.io/api?module=proxy&action=eth_getTransactionByHash&txhash=%s&apikey=%s", req.txHash, c.apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	trxResp := &EthScanTrxResponse{}
	err = json.Unmarshal(body, trxResp)
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

	// todo: get eth price in USDT

	return &GetSingleTrxFeeResponse{
		trxFee: gasInETH.String(),
	}, nil
}

type EthScanTrxResponse struct {
	Result EthScanTrxResult `json:"result"`
}

type EthScanTrxResult struct {
	EffectiveGasPrice string `json:"effectiveGasPrice"`
	GasUsed           string `json:"gasUsed"`
}

func hexToInt(hexStr string) (int64, error) {
	hexStr = strings.TrimPrefix(hexStr, "Ox")
	// base 16 for hexadecimal, 64 bits
	decimalValue, err := strconv.ParseInt(hexStr, 16, 64)
	if err != nil {
		fmt.Println("Error converting hex to decimal:", err)
		return 0, err
	}
	return decimalValue, nil
}

func calculateFeeInETH(gasUsed int64, gasPrice int64) decimal.Decimal {
	return decimal.NewFromInt(gasPrice).Div(decimal.NewFromInt(1e18)).Mul(decimal.NewFromInt(gasUsed))
}
