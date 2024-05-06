package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jaime1129/fedex/internal/service"
)

type TrxFeeController interface {
	GetSingleTrxFee(ctx *gin.Context)
}

type trxFeeController struct {
	svc service.TrxFeeService
}

func NewTrxController(svc service.TrxFeeService) TrxFeeController {
	return &trxFeeController{
		svc: svc,
	}
}

func (c *trxFeeController) GetSingleTrxFee(ctx *gin.Context) {
	trxHash := ctx.Param("trx_hash")
	resp, err := c.svc.GetSingleTrxFee(ctx, &service.GetSingleTrxFeeRequest{
		TrxHash: trxHash,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
	return
}
