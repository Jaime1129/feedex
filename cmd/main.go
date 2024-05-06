package main

import (
	"context"
	"fmt"
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
	// Open a file for logging
	logFile, err := os.OpenFile("application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	// Set the output of the standard logger to the file
	log.SetOutput(logFile)

	ethscanAPIKey := os.Getenv("ETHSCAN_API_KEY")
	dsn := fmt.Sprintf("%s:%s@tcp(localhost:%s)/%s", os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_ROOT_PASSWORD"), os.Getenv("MYSQL_PORT"), os.Getenv("MYSQL_DATABASE"))
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
