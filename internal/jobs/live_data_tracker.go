package jobs

import (
	"context"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/jaime1129/fedex/internal/components"
	"github.com/jaime1129/fedex/internal/repository"
	"github.com/jaime1129/fedex/internal/util"
)

const WETHUSDC = "WETH/USDC"
const WETHUSDCPOOLADDRESS = "0x88e6a0c2ddd26feeb64f039a2c41296fcb3f5640"

type DataTracker interface {
	Run()
	Stop()
}

type dataTracker struct {
	ctx        context.Context
	cancel     func()
	ethScanCli components.EthScanCli
	bnCli      components.BnPriceCli
	repo       repository.Repository
}

func NewDataTracker(
	ctx context.Context,
	ethScanCli components.EthScanCli,
	bnCli components.BnPriceCli,
	repo repository.Repository,
) DataTracker {
	ctx, cancel := context.WithCancel(ctx)
	return &dataTracker{
		ctx:        ctx,
		cancel:     cancel,
		ethScanCli: ethScanCli,
		bnCli:      bnCli,
		repo:       repo,
	}
}

func (t *dataTracker) Run() {
	resp, err := t.ethScanCli.GetLatestBlock()
	if err != nil || resp == 0 {
		log.Fatal("get latest block err: " + err.Error())
		return
	}

	go func() {
		t.TrackLiveData(t.ctx, resp)
	}()

	go func() {
		t.TrackHistoricalData(t.ctx, resp)
	}()
}

func (t *dataTracker) Stop() {
	t.cancel()
}

func (t *dataTracker) TrackLiveData(ctx context.Context, latestBlockNumber int64) {
	ticker := time.NewTicker(time.Second)
	initialPage := int64(1)
	offset := int64(20)
	for {
		select {
		case <-ticker.C:
			log.Println("refresh live data...")
			currentUnix := time.Now().Unix()
			price, err := t.bnCli.QueryETHPrice(currentUnix-60, currentUnix)
			if err != nil {
				log.Println("query price err: " + err.Error())
				continue
			}
			resp, err := t.ethScanCli.QueryHistoricalTrxs(&components.QueryHistoricalTrxsReq{
				Address:    WETHUSDCPOOLADDRESS,
				StartBlock: latestBlockNumber,
				// EndBlock   : nil
				Page:   initialPage,
				Offset: offset,
				Sort:   components.SortAsc,
			})
			if err != nil {
				log.Println("query history trx err: " + err.Error())
				continue
			}
			if len(resp.Result) == 0 {
				log.Println("no transaction to record")
				continue
			}

			// save transaction to db
			res := make([]repository.UniTrxFee, len(resp.Result))
			for i, r := range resp.Result {
				timeStamp, _ := strconv.Atoi(r.TimeStamp)
				gasUsed, _ := strconv.Atoi(r.GasUsed)
				gasPrice, _ := strconv.Atoi(r.GasPrice)
				res[i] = repository.UniTrxFee{
					Symbol:       WETHUSDC,
					TrxHash:      r.Hash,
					TrxTime:      uint64(timeStamp),
					GasUsed:      uint64(gasUsed),
					GasPrice:     uint64(gasPrice),
					EthUsdtPrice: price,
					TrxFeeUsdt:   util.CalculateFeeInETH(int64(gasUsed), int64(gasPrice)).Mul(price),
				}
			}

			err = t.repo.BatchInsertUniTrxFee(res)
			if err != nil {
				log.Println("batch insertion err: " + err.Error())
				continue
			}

			initialPage++
		case <-ctx.Done():
			log.Println("live data tracker stopped")
		}

	}
}

func (t *dataTracker) TrackHistoricalData(ctx context.Context, latestBlockNumber int64) {
	ticker := time.NewTicker(time.Second)
	initialPage := int64(1)
	offset := int64(20)

	maxBlock, err := t.repo.GetMaxBlockNum(WETHUSDC)
	if err != nil {
		log.Println("fail to get maxBlock: " + err.Error())
		return
	}

	for {
		select {
		case <-ticker.C:
			log.Println("tracking historical data...")
			resp, err := t.ethScanCli.QueryHistoricalTrxs(&components.QueryHistoricalTrxsReq{
				Address:    WETHUSDCPOOLADDRESS,
				StartBlock: int64(maxBlock),
				EndBlock:   &latestBlockNumber,
				Page:       initialPage,
				Offset:     offset,
				Sort:       components.SortAsc,
			})
			if err != nil {
				log.Println("query history trx err: " + err.Error())
				continue
			}
			if len(resp.Result) == 0 {
				log.Println("no transaction to record")
				continue
			}

			maxBlockNum := uint64(0)
			minTime := int64(math.MaxInt64)
			maxTime := int64(0)
			// save transaction to db
			res := make([]repository.UniTrxFee, len(resp.Result))
			for i, r := range resp.Result {
				timeStamp, _ := strconv.Atoi(r.TimeStamp)
				gasUsed, _ := strconv.Atoi(r.GasUsed)
				gasPrice, _ := strconv.Atoi(r.GasPrice)
				blockNum, _ := strconv.Atoi(r.BlockNumber)
				maxBlockNum = uint64(math.Max(float64(maxBlockNum), float64(blockNum)))
				minTime = int64(math.Min(float64(minTime), float64(timeStamp)))
				maxTime = int64(math.Max(float64(maxTime), float64(timeStamp)))
				res[i] = repository.UniTrxFee{
					Symbol:   WETHUSDC,
					TrxHash:  r.Hash,
					TrxTime:  uint64(timeStamp),
					GasUsed:  uint64(gasUsed),
					GasPrice: uint64(gasPrice),
				}
			}

			price, err := t.bnCli.QueryETHPrice(minTime, maxTime)
			if err != nil {
				log.Println("query price err: " + err.Error())
				continue
			}

			for i := range res {
				res[i].EthUsdtPrice = price
				res[i].TrxFeeUsdt = util.CalculateFeeInETH(int64(res[i].GasUsed), int64(res[i].GasPrice)).Mul(price)
			}

			err = t.repo.BatchRecordHistoricalTrx(res, WETHUSDC, maxBlockNum)
			if err != nil {
				log.Println("batch insertion err: " + err.Error())
				continue
			}

			initialPage++
		case <-ctx.Done():
			log.Println("live data tracker stopped")
		}

	}
}
