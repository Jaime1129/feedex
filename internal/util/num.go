package util

import (
	"log"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

func HexToInt(hexStr string) (int64, error) {
	hexStr = strings.TrimPrefix(hexStr, "0x")
	// base 16 for hexadecimal, 64 bits
	decimalValue, err := strconv.ParseInt(hexStr, 16, 64)
	if err != nil {
		log.Println("Error converting hex to decimal:", err)
		return 0, err
	}
	return decimalValue, nil
}

func CalculateFeeInETH(gasUsed int64, gasPrice int64) decimal.Decimal {
	// convert gasPrice in Wei to Eth, then multiply with gasUsed
	return decimal.NewFromInt(gasPrice).Div(decimal.NewFromInt(1e18)).Mul(decimal.NewFromInt(gasUsed))
}
