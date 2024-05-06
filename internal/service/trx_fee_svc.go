package service

import (
	"context"
	"errors"

	"github.com/jaime1129/fedex/internal/components"
	"github.com/jaime1129/fedex/internal/util"
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

	gasUsed, err := util.HexToInt(trxResp.Result.GasUsed)
	if err != nil {
		return nil, err
	}
	gasPrice, err := util.HexToInt(trxResp.Result.EffectiveGasPrice)
	if err != nil {
		return nil, err
	}

	// calculate gas fee in eth
	gasInETH := util.CalculateFeeInETH(gasUsed, gasPrice)

	// get block timestamp
	blockResp, err := c.ethScanCli.QueryBlock(trxResp.Result.BlockNumber)
	if err != nil {
		return nil, err
	}

	// todo: get eth price in USDT
	trxTime, err := util.HexToInt(blockResp.Result.Timestamp)
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
