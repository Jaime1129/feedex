package service

import (
	"context"
	"errors"

	"github.com/jaime1129/fedex/internal/components"
	"github.com/jaime1129/fedex/internal/repository"
	"github.com/jaime1129/fedex/internal/util"
)

type TrxFeeService interface {
	GetSingleTrxFee(ctx context.Context, req *GetSingleTrxFeeRequest) (*GetSingleTrxFeeResponse, error)
	GetTrxFeeList(ctx context.Context, req *GetTrxFeeListRequest) (*GetTrxFeeListResponse, error)
}

type trxFeeService struct {
	ethScanCli components.EthScanCli
	bnPriceCli components.BnPriceCli
	repo       repository.Repository
}

func NewTrxService(
	ethScanCli components.EthScanCli,
	bnPriceCli components.BnPriceCli,
	repo repository.Repository,
) TrxFeeService {
	return &trxFeeService{
		ethScanCli: ethScanCli,
		bnPriceCli: bnPriceCli,
		repo:       repo,
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

	// prefering directly querying from db
	res, err := c.repo.GetTrxFee(req.TrxHash)
	if err != nil {
		return nil, err
	}
	if res != nil {
		return &GetSingleTrxFeeResponse{
			TrxFee: res.TrxFeeUsdt.String(),
		}, nil
	}

	// alternatively querying from etherscan api
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

	// get eth price in USDT
	trxTime, err := util.HexToInt(blockResp.Result.Timestamp)
	if err != nil {
		return nil, err
	}

	// fetch the average price of [trxTime-60, trxTime+60]
	price, err := c.bnPriceCli.QueryETHPrice(trxTime-60, trxTime+60, components.INTERVAL_1MIN)
	if err != nil {
		return nil, err
	}

	return &GetSingleTrxFeeResponse{
		TrxFee: gasInETH.Mul(price).String(),
	}, nil
}

type GetTrxFeeListRequest struct {
	Symbol    string
	StartTime int64
	EndTime   int64
	Page      int
	Limit     int
}

type GetTrxFeeListResponse struct {
	Result []repository.UniTrxFee
}

func (c *trxFeeService) GetTrxFeeList(ctx context.Context, req *GetTrxFeeListRequest) (*GetTrxFeeListResponse, error) {
	if req == nil {
		return nil, errors.New("nil req")
	}
	res, err := c.repo.ListTrxFee(req.Symbol, req.StartTime, req.EndTime, req.Page, req.Limit)
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return &GetTrxFeeListResponse{}, nil
	}

	return &GetTrxFeeListResponse{Result: res}, nil
}
