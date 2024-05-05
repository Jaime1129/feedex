package controller

import (
	"github.com/gin-gonic/gin"
)

type TrxFeeController interface {
	GetSingleTrxFee(ctx *gin.Context)
}

type trxFeeController struct {
}

func NewTrxController() TrxFeeController {
	return &trxFeeController{}
}

type GetSingleTrxFeeRequest struct {
}

type GetSingleTrxFeeResponse struct {
}

func (c *trxFeeController) GetSingleTrxFee(ctx *gin.Context) {
	return
}
