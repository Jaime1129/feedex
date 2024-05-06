package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jaime1129/fedex/docs"
	"github.com/jaime1129/fedex/internal/components"
	"github.com/jaime1129/fedex/internal/controller"
	"github.com/jaime1129/fedex/internal/jobs"
	"github.com/jaime1129/fedex/internal/repository"
	"github.com/jaime1129/fedex/internal/service"
	"github.com/swaggo/files"       // swagger embed files
	"github.com/swaggo/gin-swagger" // gin-swagger middleware
)

func main() {
	ctx := context.Background()

	ethscanAPIKey := "28I2TNRKKSUMM4QNMP6B9AICEQTJIE3978"
	dsn := "root:@jaime1129@tcp(localhost:3306)/trx_fee"
	ethScanCli := components.NewEthScanCli(ethscanAPIKey)
	bnPriceCli := components.NewBnPriceCLi()

	repo := repository.NewRepository(dsn)
	svc := service.NewTrxService(ethScanCli, bnPriceCli, repo)

	t := jobs.NewDataTracker(
		ctx,
		ethScanCli,
		bnPriceCli,
		repo,
	)
	t.Run()

	c := controller.NewTrxController(svc)
	router := setupRouter(c)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		// Start the server
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	{
		t.Stop()
		repo.Close()
	}

	time.Sleep(5 * time.Second)
	log.Println("Server exited")
}

func setupRouter(c controller.TrxFeeController) *gin.Engine {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group("/api/v1")
	{
		trxFee := v1.Group("/trxfee")
		trxFee.GET(":trx_hash", c.GetSingleTrxFee)
		trxFee.GET("/list", c.GetTrxFeeList)
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}
