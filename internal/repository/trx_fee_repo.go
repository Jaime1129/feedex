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
	Close()
}

type repository struct {
	db *sql.DB
}

func NewRepository() Repository {
	// Update the DSN as per your user, password, host, and database details
	dsn := "root:@jaime1129@tcp(localhost:3306)/trx_fee"
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
	EthUsdtPrice decimal.Decimal
	TrxFeeUsdt   decimal.Decimal
}

func (r *repository) BatchInsertUniTrxFee(fees []UniTrxFee) error {
	var placeholders []string
	var args []interface{}

	for _, fee := range fees {
		placeholders = append(placeholders, "(?, ?, ?, ?, ?, ?, ?)")
		args = append(args, fee.Symbol, fee.TrxHash, fee.TrxTime, fee.GasUsed, fee.GasPrice, fee.EthUsdtPrice.String(), fee.TrxFeeUsdt.String())
	}

	stmt := fmt.Sprintf("INSERT INTO uni_trx_fee (symbol, trx_hash, trx_time, gas_used, gas_price, eth_usdt_price, trx_fee_usdt) VALUES %s",
		strings.Join(placeholders, ", "))

	_, err := r.db.Exec(stmt, args...)
	if err != nil {
		return err
	}
	return nil
}
