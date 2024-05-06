package components

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/jaime1129/fedex/internal/util"
	"go.uber.org/ratelimit"
)

type EthScanCli interface {
	QueryTrxFee(trxHash string) (*EthScanTrxResponse, error)
	QueryBlock(blockNumber string) (*EthScanBlockResponse, error)
	GetLatestBlock() (int64, error)
	QueryHistoricalTrxs(req *QueryHistoricalTrxsReq) (*QueryHistoricalTrxsResp, error)
}

type ethScanCli struct {
	rl     ratelimit.Limiter
	apiKey string
}

func NewEthScanCli(apiKey string) EthScanCli {
	return &ethScanCli{
		rl:     ratelimit.New(10),
		apiKey: apiKey,
	}
}

type EthScanTrxResponse struct {
	Result EthScanTrxResult `json:"result"`
	Error  EthScanError     `json:"error"`
}

type EthScanTrxResult struct {
	EffectiveGasPrice string `json:"effectiveGasPrice"`
	GasUsed           string `json:"gasUsed"`
	BlockNumber       string `json:"blockNumber"`
}

type EthScanError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (c *ethScanCli) QueryTrxFee(trxHash string) (*EthScanTrxResponse, error) {
	c.rl.Take()
	// send query to etherscan api
	url := fmt.Sprintf("https://api.etherscan.io/api?module=proxy&action=eth_getTransactionReceipt&txhash=%s&apikey=%s", trxHash, c.apiKey)
	log.Println("ethscan api url: " + url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Println("ethscan api resp body: " + string(body))

	trxResp := &EthScanTrxResponse{}
	err = json.Unmarshal(body, trxResp)
	if err != nil {
		return nil, err
	}

	// check if api call returns error
	if trxResp.Error.Code != 0 {
		// return empty result if no transactions found
		if trxResp.Error.Message == "No transactions found" {
			return trxResp, nil
		}
		log.Println("ethscan api call returns error: " + trxResp.Error.Message)
		return nil, errors.New("ethscan api call returns error")
	}

	return trxResp, nil
}

type EthScanBlockResponse struct {
	Result EthScanBlockResult `json:"result"`
	Error  EthScanError       `json:"error"`
}

type EthScanBlockResult struct {
	Timestamp string `json:"timestamp"`
}

func (c *ethScanCli) QueryBlock(blockNumber string) (*EthScanBlockResponse, error) {
	c.rl.Take()
	// send query to etherscan api
	url := fmt.Sprintf("https://api.etherscan.io/api?module=proxy&action=eth_getBlockByNumber&tag=%s&boolean=true&apikey=%s", blockNumber, c.apiKey)
	log.Println("ethscan api url: " + url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Println("ethscan api resp body: " + string(body))

	trxResp := &EthScanBlockResponse{}
	err = json.Unmarshal(body, trxResp)
	if err != nil {
		return nil, err
	}

	// check if api call returns error
	if trxResp.Error.Code != 0 {
		log.Println("ethscan api call returns error: " + trxResp.Error.Message)
		return nil, errors.New("ethscan api call returns error")
	}

	return trxResp, nil
}

type QueryHistoricalTrxsReq struct {
	Address    string
	StartBlock int64
	EndBlock   *int64 // optional
	Page       int64
	Offset     int64
	Sort       string
}

type QueryHistoricalTrxsResp struct {
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Result  []Transaction `json:"result"`
}

const (
	SortAsc  = "asc"
	SortDesc = "desc"
	StatusOK = "1"
)

// Transaction details within the 'result' array
type Transaction struct {
	Hash              string `json:"hash"`
	BlockNumber       string `json:"blockNumber"`
	TimeStamp         string `json:"timeStamp"`
	Gas               string `json:"gas"`
	GasPrice          string `json:"gasPrice"`
	ContractAddress   string `json:"contractAddress"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	GasUsed           string `json:"gasUsed"`
}

func (c *ethScanCli) QueryHistoricalTrxs(req *QueryHistoricalTrxsReq) (*QueryHistoricalTrxsResp, error) {
	c.rl.Take()
	if req == nil {
		return nil, errors.New("nil req")
	}

	url := fmt.Sprintf(
		"https://api.etherscan.io/api?module=account&action=txlist&address=%s&startblock=%d&page=%d&offset=%d&sort=%s&apikey=%s",
		req.Address,
		req.StartBlock,
		req.Page,
		req.Offset,
		req.Sort,
		c.apiKey,
	)
	if req.EndBlock != nil {
		url = url + fmt.Sprintf("&endblock=%d", *req.EndBlock)
	}

	log.Println("ethscan api url: " + url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Println("ethscan api resp body: " + string(body))

	trxResp := &QueryHistoricalTrxsResp{}
	err = json.Unmarshal(body, trxResp)
	if err != nil {
		return nil, err
	}

	// check if api call returns error
	if trxResp.Status != StatusOK {
		log.Printf("ethscan api call not ok: %s, %s\n", trxResp.Status, trxResp.Message)
		return nil, errors.New("ethscan api call returns error")
	}

	return trxResp, nil
}

type GetLatestBlockResp struct {
	BlockNumber string       `json:"result"`
	Error       EthScanError `json:"error"`
}

func (c *ethScanCli) GetLatestBlock() (int64, error) {
	c.rl.Take()
	// send query to etherscan api
	url := fmt.Sprintf("https://api.etherscan.io/api?module=proxy&action=eth_blockNumber&apikey=%s", c.apiKey)
	log.Println("ethscan api url: " + url)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	log.Println("ethscan api resp body: " + string(body))

	trxResp := &GetLatestBlockResp{}
	err = json.Unmarshal(body, trxResp)
	if err != nil {
		return 0, err
	}

	// check if api call returns error
	if trxResp.Error.Code != 0 {
		log.Println("ethscan api call returns error: " + trxResp.Error.Message)
		return 0, errors.New("ethscan api call returns error")
	}

	blockNum, err := util.HexToInt(trxResp.BlockNumber)
	if err != nil {
		return 0, err
	}
	return blockNum, nil
}
