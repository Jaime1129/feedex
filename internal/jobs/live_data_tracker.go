package jobs

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/jaime1129/fedex/internal/components"
	"github.com/jaime1129/fedex/internal/repository"
	"github.com/jaime1129/fedex/internal/util"
)

type DataTracker interface {
}

type dataTracker struct {
	ethScanCli components.EthScanCli
	bnCli      components.BnPriceCli
	repo       repository.Repository
}

func NewDataTracker(
	ethScanCli components.EthScanCli,
	bnCli components.BnPriceCli,
	repo repository.Repository,
) DataTracker {
	return &dataTracker{
		ethScanCli: ethScanCli,
		bnCli:      bnCli,
		repo:       repo,
	}
}

func (t *dataTracker) Run(ctx context.Context) {
	resp, err := t.ethScanCli.GetLatestBlock()
	if err != nil || resp == 0 {
		log.Fatal("get latest block err: " + err.Error())
		return
	}

	go func() {
		t.TrackLiveData(ctx, resp)
	}()

	go func() {
		t.TrackHistoricalData(ctx, resp)
	}()
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
				Address:    "0x88e6a0c2ddd26feeb64f039a2c41296fcb3f5640",
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
					Symbol:       "WETH/USDC",
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

}
