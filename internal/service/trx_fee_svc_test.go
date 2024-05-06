package service

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/jaime1129/fedex/internal/components"
	"github.com/jaime1129/fedex/internal/repository"
	"github.com/jaime1129/fedex/internal/util"
	mock_components "github.com/jaime1129/fedex/mock/components"
	mock_repository "github.com/jaime1129/fedex/mock/repository"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestGetSingleTrxFee(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEthScanCli := mock_components.NewMockEthScanCli(ctrl)
	mockBnPriceCli := mock_components.NewMockBnPriceCli(ctrl)
	mockRepo := mock_repository.NewMockRepository(ctrl)

	service := NewTrxService(mockEthScanCli, mockBnPriceCli, mockRepo)

	ctx := context.TODO()
	req := &GetSingleTrxFeeRequest{TrxHash: "hash123"}

	// Setup expectations and return values for the mocks
	mockRepo.EXPECT().GetTrxFee(req.TrxHash).Return(nil, nil) // Simulate no result in DB

	ethScanResp := &components.EthScanTrxResponse{
		Result: components.EthScanTrxResult{
			GasUsed:           "0x5208",
			EffectiveGasPrice: "0x3B9ACA00",
			BlockNumber:       "0x10FB78",
		},
	}
	mockEthScanCli.EXPECT().QueryTrxFee(req.TrxHash).Return(ethScanResp, nil)

	gasUsed, _ := util.HexToInt("0x5208")
	gasPrice, _ := util.HexToInt("0x3B9ACA00")
	gasInETH := util.CalculateFeeInETH(gasUsed, gasPrice)

	blockResp := &components.EthScanBlockResponse{
		Result: components.EthScanBlockResult{Timestamp: "0x5BA46680"},
	}
	mockEthScanCli.EXPECT().QueryBlock("0x10FB78").Return(blockResp, nil)

	trxTime, _ := util.HexToInt("0x5BA46680")
	mockBnPriceCli.EXPECT().QueryETHPrice(trxTime-60, trxTime+60, "1m").Return(decimal.NewFromFloat(2000), nil)

	// Call the function under test
	response, err := service.GetSingleTrxFee(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, gasInETH.Mul(decimal.NewFromFloat(2000)).String(), response.TrxFee)
}

func TestGetTrxFeeList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockRepository(ctrl)
	service := NewTrxService(nil, nil, mockRepo) // nils are safe here since they are not used in this method

	ctx := context.TODO()
	req := &GetTrxFeeListRequest{
		Symbol:    "WETH/USDC",
		StartTime: time.Now().Unix(),
		EndTime:   time.Now().Add(24 * time.Hour).Unix(),
		Page:      1,
		Limit:     10,
	}

	// Mock response from repository
	mockResponse := []repository.UniTrxFee{{TrxFeeUsdt: decimal.NewFromFloat(300.5)}}
	mockRepo.EXPECT().ListTrxFee(req.Symbol, req.StartTime, req.EndTime, req.Page, req.Limit).Return(mockResponse, nil)

	// Call the function under test
	response, err := service.GetTrxFeeList(ctx, req)

	assert.NoError(t, err)
	assert.Len(t, response.Result, 1)
	assert.Equal(t, decimal.NewFromFloat(300.5).String(), response.Result[0].TrxFeeUsdt.String())
}
