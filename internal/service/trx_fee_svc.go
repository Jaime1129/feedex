package service

import "context"

type TrxFeeService interface {
	GetSingleTrxFee(ctx context.Context, req *GetSingleTrxFeeRequest) (*GetSingleTrxFeeResponse, error)
}

type trxFeeService struct {
}

func NewTrxService() TrxFeeService {
	return &trxFeeService{}
}

type GetSingleTrxFeeRequest struct {
}

type GetSingleTrxFeeResponse struct {
}

func (c *trxFeeService) GetSingleTrxFee(ctx context.Context, req *GetSingleTrxFeeRequest) (*GetSingleTrxFeeResponse, error) {
	return nil, nil
}
