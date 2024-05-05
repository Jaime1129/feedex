package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jaime1129/fedex/internal/controller"
)

func main() {
	r := gin.Default()
	ethscanAPIKey := "28I2TNRKKSUMM4QNMP6B9AICEQTJIE3978"
	c := controller.NewTrxController(ethscanAPIKey)

	v1 := r.Group("/api/v1")
	{
		trxFee := v1.Group("/trxfee")
		trxFee.GET(":trx_hash", c.GetSingleTrxFee)
	}

	r.Run(":8080")
}
