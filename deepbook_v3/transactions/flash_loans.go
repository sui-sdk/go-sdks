package transactions

import (
	"math"

	"github.com/sui-sdks/go-sdks/deepbook_v3/utils"
	stx "github.com/sui-sdks/go-sdks/sui/transactions"
)

type FlashLoanContract struct{ config *utils.DeepBookConfig }

func NewFlashLoanContract(config *utils.DeepBookConfig) *FlashLoanContract {
	return &FlashLoanContract{config: config}
}

func (c *FlashLoanContract) BorrowBaseAsset(tx *stx.Transaction, poolKey string, borrowAmount float64) stx.Argument {
	pool := c.config.GetPool(poolKey)
	base := c.config.GetCoin(pool.BaseCoin)
	quote := c.config.GetCoin(pool.QuoteCoin)
	qty := uint64(math.Round(borrowAmount * base.Scalar))
	return tx.MoveCall(c.config.DeepbookPackageID+"::pool::borrow_flashloan_base", []stx.Argument{tx.Object(pool.Address), pureU64(tx, qty)}, []string{base.Type, quote.Type})
}

func (c *FlashLoanContract) BorrowQuoteAsset(tx *stx.Transaction, poolKey string, borrowAmount float64) stx.Argument {
	pool := c.config.GetPool(poolKey)
	base := c.config.GetCoin(pool.BaseCoin)
	quote := c.config.GetCoin(pool.QuoteCoin)
	qty := uint64(math.Round(borrowAmount * quote.Scalar))
	return tx.MoveCall(c.config.DeepbookPackageID+"::pool::borrow_flashloan_quote", []stx.Argument{tx.Object(pool.Address), pureU64(tx, qty)}, []string{base.Type, quote.Type})
}

func (c *FlashLoanContract) ReturnBaseAsset(tx *stx.Transaction, poolKey string, baseCoin, flashLoan stx.Argument) stx.Argument {
	pool := c.config.GetPool(poolKey)
	base := c.config.GetCoin(pool.BaseCoin)
	quote := c.config.GetCoin(pool.QuoteCoin)
	return tx.MoveCall(c.config.DeepbookPackageID+"::pool::return_flashloan_base", []stx.Argument{
		tx.Object(pool.Address), baseCoin, flashLoan,
	}, []string{base.Type, quote.Type})
}

func (c *FlashLoanContract) ReturnQuoteAsset(tx *stx.Transaction, poolKey string, quoteCoin, flashLoan stx.Argument) stx.Argument {
	pool := c.config.GetPool(poolKey)
	base := c.config.GetCoin(pool.BaseCoin)
	quote := c.config.GetCoin(pool.QuoteCoin)
	return tx.MoveCall(c.config.DeepbookPackageID+"::pool::return_flashloan_quote", []stx.Argument{
		tx.Object(pool.Address), quoteCoin, flashLoan,
	}, []string{base.Type, quote.Type})
}
