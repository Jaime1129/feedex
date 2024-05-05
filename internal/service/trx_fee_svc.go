package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/jaime1129/fedex/internal/components"
	"github.com/shopspring/decimal"
)

type TrxFeeService interface {
	GetSingleTrxFee(ctx context.Context, req *GetSingleTrxFeeRequest) (*GetSingleTrxFeeResponse, error)
}

type trxFeeService struct {
	apiKey     string
	ethScanCli components.EthScanCli
	bnPriceCli components.BnPriceCli
}

func NewTrxService(apiKey string) TrxFeeService {
	return &trxFeeService{
		apiKey:     apiKey,
		ethScanCli: components.NewEthScanCli(apiKey),
		bnPriceCli: components.NewBnPriceCLi(),
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

	trxResp, err := c.ethScanCli.QueryTrxFee(req.TrxHash)
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
	blockResp, err := c.ethScanCli.QueryBlock(trxResp.Result.BlockNumber)
	if err != nil {
		return nil, err
	}

	// todo: get eth price in USDT
	trxTime, err := hexToInt(blockResp.Result.Timestamp)
	if err != nil {
		return nil, err
	}

	// fetch the average price of [trxTime-60, trxTime+60]
	price, err := c.bnPriceCli.QueryETHPrice(trxTime-60, trxTime+60)
	if err != nil {
		return nil, err
	}

	return &GetSingleTrxFeeResponse{
		TrxFee: gasInETH.Mul(price).String(),
	}, nil
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
