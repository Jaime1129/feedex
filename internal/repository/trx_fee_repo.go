package repository

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/shopspring/decimal"
)

type Repository interface {
	BatchInsertUniTrxFee(fees []UniTrxFee) error
	GetMaxBlockNum(symbol string) (uint64, error)
	BatchRecordHistoricalTrx(fees []UniTrxFee, symbol string, maxBlock uint64) error
	GetTrxFee(txHash string) (*UniTrxFee, error)
	ListTrxFee(symbol string, startTime int64, endTime int64, page int, limit int) ([]UniTrxFee, error)
	Close()
}

type repository struct {
	db *sql.DB
}

func NewRepository(dsn string) Repository {
	// Update the DSN as per your user, password, host, and database details
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}
	return &repository{db: db}
}

func (r *repository) Close() {
	r.db.Close()
}

type UniTrxFee struct {
	Symbol       string
	TrxHash      string
	TrxTime      uint64
	GasUsed      uint64
	GasPrice     uint64
	BlockNumber  uint64
	EthUsdtPrice decimal.Decimal
	TrxFeeUsdt   decimal.Decimal
}

func (r *repository) BatchInsertUniTrxFee(fees []UniTrxFee) error {
	var placeholders []string
	var args []interface{}

	for _, fee := range fees {
		placeholders = append(placeholders, "(?, ?, ?, ?, ?, ?, ?, ?)")
		args = append(args, fee.Symbol, fee.TrxHash, fee.TrxTime, fee.GasUsed, fee.GasPrice, fee.EthUsdtPrice.String(), fee.TrxFeeUsdt.String(), fee.BlockNumber)
	}

	// use ignore to avoid dup key conflict error
	stmt := fmt.Sprintf("INSERT IGNORE INTO uni_trx_fee (symbol, trx_hash, trx_time, gas_used, gas_price, eth_usdt_price, trx_fee_usdt, block_num) VALUES %s",
		strings.Join(placeholders, ", "))

	_, err := r.db.Exec(stmt, args...)
	if err != nil {
		return err
	}
	return nil
}

// batch insert historical trxs and record the maximum block number
func (r *repository) BatchRecordHistoricalTrx(fees []UniTrxFee, symbol string, maxBlock uint64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	err = r.BatchInsertUniTrxFee(fees)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = r.db.Exec("INSERT INTO block_num_record (symbol, max_block) VALUES (?,?) ON DUPLICATE KEY UPDATE max_block=?", symbol, maxBlock, maxBlock)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) GetMaxBlockNum(symbol string) (uint64, error) {
	var blockNum uint64
	err := r.db.QueryRow("SELECT max_block FROM block_num_record WHERE symbol = ?", symbol).Scan(&blockNum)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	return blockNum, nil
}

func (r *repository) GetTrxFee(txHash string) (*UniTrxFee, error) {
	query := "SELECT symbol, trx_hash, trx_time, gas_used, gas_price, eth_usdt_price, trx_fee_usdt, block_num FROM uni_trx_fee where tx_hash=?"
	rows, err := r.db.Query(query, txHash)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var fee UniTrxFee
	err = rows.Scan(&fee.Symbol, &fee.TrxHash, &fee.TrxTime, &fee.GasUsed, &fee.GasPrice, &fee.EthUsdtPrice, &fee.TrxFeeUsdt, &fee.BlockNumber)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return &fee, nil
}

func (r *repository) ListTrxFee(symbol string, startTime int64, endTime int64, page int, limit int) ([]UniTrxFee, error) {
	if limit == 0 || limit > 50 {
		limit = 20
	}
	query := "SELECT symbol, trx_hash, trx_time, gas_used, gas_price, eth_usdt_price, trx_fee_usdt, block_num FROM uni_trx_fee where " +
		"symbol=? and trx_time >= ? and trx_time <= ? limit ? offset ?"
	rows, err := r.db.Query(query, symbol, startTime, endTime, limit, page*limit)
	if err == sql.ErrNoRows {
		return []UniTrxFee{}, nil
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fees []UniTrxFee
	for rows.Next() {
		var fee UniTrxFee
		err := rows.Scan(&fee.Symbol, &fee.TrxHash, &fee.TrxTime, &fee.GasUsed, &fee.GasPrice, &fee.EthUsdtPrice, &fee.TrxFeeUsdt, &fee.BlockNumber)
		if err != nil {
			return nil, err
		}
		fees = append(fees, fee)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return fees, nil
}
