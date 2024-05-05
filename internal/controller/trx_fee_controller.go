package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jaime1129/fedex/internal/service"
)

type TrxFeeController interface {
	GetSingleTrxFee(ctx *gin.Context)
}

type trxFeeController struct {
	svc service.TrxFeeService
}

func NewTrxController(ethScanAPIKey string) TrxFeeController {
	return &trxFeeController{
		svc: service.NewTrxService(ethScanAPIKey),
	}
}

type GetSingleTrxFeeRequest struct {
}

type GetSingleTrxFeeResponse struct {
}

func (c *trxFeeController) GetSingleTrxFee(ctx *gin.Context) {
	return
}
