package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jaime1129/fedex/internal/service"
)

type TrxFeeController interface {
	GetSingleTrxFee(ctx *gin.Context)
	GetTrxFeeList(ctx *gin.Context)
}

type trxFeeController struct {
	svc service.TrxFeeService
}

func NewTrxController(svc service.TrxFeeService) TrxFeeController {
	return &trxFeeController{
		svc: svc,
	}
}

// GetSingleTrxFee godoc
//	@Summary		Get trx fee of single trx
//	@Description	get trx fee by trx hash
//	@Accept			json
//	@Produce		json
//	@Param			trx_hash	path	string	true	"trx hash"
//	@Success		200			string	trx_fee
//	@Failure		500			string	msg
//	@Router			/trxfee/{trx_hash} [get]
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
}

// GetTrxFeeList godoc
//	@Summary		Get a list of trx fee
//	@Description	get trx fee by given time period
//	@Accept			json
//	@Produce		json
//	@Param			symbol		query		string	true	"symbol"
//	@Param			start_time	query		int		true	"start timestamp"
//	@Param			end_time	query		int		true	"end timestamp"
//	@Param			page		query		int		true	"page starting from 0"
//	@Param			limit		query		int		true	"20 by default"
//	@Success		200			{object}	service.GetTrxFeeListResponse
//	@Failure		500			string		msg
//	@Router			/trxfee/list [get]
func (c *trxFeeController) GetTrxFeeList(ctx *gin.Context) {
	symbol := ctx.DefaultQuery("symbol", "WETH/USDC")
	startTime, _ := strconv.ParseInt(ctx.DefaultQuery("start_time", "0"), 10, 64)
	endTime, _ := strconv.ParseInt(ctx.DefaultQuery("end_time", "0"), 10, 64)
	page, _ := strconv.ParseInt(ctx.DefaultQuery("page", "0"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.DefaultQuery("limit", "20"), 10, 64)
	if endTime == 0 {
		endTime = time.Now().Unix()
	}
	resp, err := c.svc.GetTrxFeeList(ctx, &service.GetTrxFeeListRequest{
		Symbol:    symbol,
		StartTime: startTime,
		EndTime:   endTime,
		Page:      int(page),
		Limit:     int(limit),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
