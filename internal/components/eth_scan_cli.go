package components

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/ratelimit"
)

type EthScanCli interface {
	QueryTrxFee(trxHash string) (*EthScanTrxResponse, error)
	QueryBlock(blockNumber string) (*EthScanBlockResponse, error)
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
	// send query to etherscan api
	url := fmt.Sprintf("https://api.etherscan.io/api?module=proxy&action=eth_getTransactionReceipt&txhash=%s&apikey=%s", trxHash, c.apiKey)
	fmt.Println("ethscan api url: " + url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println("ethscan api resp body: " + string(body))

	trxResp := &EthScanTrxResponse{}
	err = json.Unmarshal(body, trxResp)
	if err != nil {
		return nil, err
	}

	// check if api call returns error
	if trxResp.Error.Code != 0 {
		fmt.Println("ethscan api call returns error: " + trxResp.Error.Message)
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
	// send query to etherscan api
	url := fmt.Sprintf("https://api.etherscan.io/api?module=proxy&action=eth_getBlockByNumber&tag=%s&boolean=true&apikey=%s", blockNumber, c.apiKey)
	fmt.Println("ethscan api url: " + url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println("ethscan api resp body: " + string(body))

	trxResp := &EthScanBlockResponse{}
	err = json.Unmarshal(body, trxResp)
	if err != nil {
		return nil, err
	}

	// check if api call returns error
	if trxResp.Error.Code != 0 {
		fmt.Println("ethscan api call returns error: " + trxResp.Error.Message)
		return nil, errors.New("ethscan api call returns error")
	}

	return trxResp, nil
}
