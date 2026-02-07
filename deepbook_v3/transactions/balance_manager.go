package transactions

import (
	"math"

	"github.com/sui-sdks/go-sdks/deepbook_v3/utils"
	stx "github.com/sui-sdks/go-sdks/sui/transactions"
)

type BalanceManagerContract struct {
	config *utils.DeepBookConfig
}

func NewBalanceManagerContract(config *utils.DeepBookConfig) *BalanceManagerContract {
	return &BalanceManagerContract{config: config}
}

func (c *BalanceManagerContract) target(fn string) string {
	return c.config.DeepbookPackageID + "::balance_manager::" + fn
}

func (c *BalanceManagerContract) CreateAndShareBalanceManager(tx *stx.Transaction) stx.Argument {
	manager := tx.MoveCall(c.target("new"), nil, nil)
	tx.MoveCall("0x2::transfer::public_share_object", []stx.Argument{manager}, []string{c.config.DeepbookPackageID + "::balance_manager::BalanceManager"})
	return manager
}

func (c *BalanceManagerContract) CreateBalanceManagerWithOwner(tx *stx.Transaction, ownerAddress string) stx.Argument {
	return tx.MoveCall(c.target("new_with_custom_owner"), []stx.Argument{tx.PureBytes([]byte(ownerAddress))}, nil)
}

func (c *BalanceManagerContract) ShareBalanceManager(tx *stx.Transaction, manager stx.Argument) stx.Argument {
	return tx.MoveCall("0x2::transfer::public_share_object", []stx.Argument{manager}, []string{c.config.DeepbookPackageID + "::balance_manager::BalanceManager"})
}

func (c *BalanceManagerContract) DepositIntoManager(tx *stx.Transaction, managerKey, coinKey string, amountToDeposit float64) stx.Argument {
	managerID := c.config.GetBalanceManager(managerKey).Address
	coin := c.config.GetCoin(coinKey)
	amount := uint64(math.Round(amountToDeposit * coin.Scalar))
	return tx.MoveCall(c.target("deposit"), []stx.Argument{
		tx.Object(managerID),
		pureU64(tx, amount),
	}, []string{coin.Type})
}

func (c *BalanceManagerContract) WithdrawFromManager(tx *stx.Transaction, managerKey, coinKey string, amountToWithdraw float64) stx.Argument {
	managerID := c.config.GetBalanceManager(managerKey).Address
	coin := c.config.GetCoin(coinKey)
	amount := uint64(math.Round(amountToWithdraw * coin.Scalar))
	return tx.MoveCall(c.target("withdraw"), []stx.Argument{
		tx.Object(managerID), pureU64(tx, amount),
	}, []string{coin.Type})
}

func (c *BalanceManagerContract) WithdrawAllFromManager(tx *stx.Transaction, managerKey, coinKey string) stx.Argument {
	managerID := c.config.GetBalanceManager(managerKey).Address
	coin := c.config.GetCoin(coinKey)
	return tx.MoveCall(c.target("withdraw_all"), []stx.Argument{tx.Object(managerID)}, []string{coin.Type})
}

func (c *BalanceManagerContract) CheckManagerBalance(tx *stx.Transaction, managerKey, coinKey string) stx.Argument {
	managerID := c.config.GetBalanceManager(managerKey).Address
	coin := c.config.GetCoin(coinKey)
	return tx.MoveCall(c.target("balance"), []stx.Argument{tx.Object(managerID)}, []string{coin.Type})
}

func (c *BalanceManagerContract) GenerateProof(tx *stx.Transaction, managerKey string) stx.Argument {
	manager := c.config.GetBalanceManager(managerKey)
	if manager.TradeCap != "" {
		return c.GenerateProofAsTrader(tx, manager.Address, manager.TradeCap)
	}
	return c.GenerateProofAsOwner(tx, manager.Address)
}

func (c *BalanceManagerContract) GenerateProofAsOwner(tx *stx.Transaction, managerID string) stx.Argument {
	return tx.MoveCall(c.target("generate_proof_as_owner"), []stx.Argument{tx.Object(managerID)}, nil)
}

func (c *BalanceManagerContract) GenerateProofAsTrader(tx *stx.Transaction, managerID, tradeCapID string) stx.Argument {
	return tx.MoveCall(c.target("generate_proof_as_trader"), []stx.Argument{tx.Object(managerID), tx.Object(tradeCapID)}, nil)
}

func (c *BalanceManagerContract) MintTradeCap(tx *stx.Transaction, managerKey string) stx.Argument {
	managerID := c.config.GetBalanceManager(managerKey).Address
	return tx.MoveCall(c.target("mint_trade_cap"), []stx.Argument{tx.Object(managerID)}, nil)
}

func (c *BalanceManagerContract) MintDepositCap(tx *stx.Transaction, managerKey string) stx.Argument {
	managerID := c.config.GetBalanceManager(managerKey).Address
	return tx.MoveCall(c.target("mint_deposit_cap"), []stx.Argument{tx.Object(managerID)}, nil)
}

func (c *BalanceManagerContract) MintWithdrawalCap(tx *stx.Transaction, managerKey string) stx.Argument {
	managerID := c.config.GetBalanceManager(managerKey).Address
	return tx.MoveCall(c.target("mint_withdraw_cap"), []stx.Argument{tx.Object(managerID)}, nil)
}

func (c *BalanceManagerContract) RegisterBalanceManager(tx *stx.Transaction, managerKey string) stx.Argument {
	managerID := c.config.GetBalanceManager(managerKey).Address
	return tx.MoveCall(c.target("register_balance_manager"), []stx.Argument{
		tx.Object(c.config.RegistryID), tx.Object(managerID),
	}, nil)
}

func (c *BalanceManagerContract) Owner(tx *stx.Transaction, managerKey string) stx.Argument {
	managerID := c.config.GetBalanceManager(managerKey).Address
	return tx.MoveCall(c.target("owner"), []stx.Argument{tx.Object(managerID)}, nil)
}

func (c *BalanceManagerContract) ID(tx *stx.Transaction, managerKey string) stx.Argument {
	managerID := c.config.GetBalanceManager(managerKey).Address
	return tx.MoveCall(c.target("id"), []stx.Argument{tx.Object(managerID)}, nil)
}

func (c *BalanceManagerContract) BalanceManagerReferralOwner(tx *stx.Transaction, referralID string) stx.Argument {
	return tx.MoveCall(c.target("balance_manager_referral_owner"), []stx.Argument{tx.Object(referralID)}, nil)
}

func (c *BalanceManagerContract) BalanceManagerReferralPoolID(tx *stx.Transaction, referralID string) stx.Argument {
	return tx.MoveCall(c.target("balance_manager_referral_pool_id"), []stx.Argument{tx.Object(referralID)}, nil)
}

func (c *BalanceManagerContract) GetBalanceManagerReferralID(tx *stx.Transaction, managerKey, poolKey string) stx.Argument {
	managerID := c.config.GetBalanceManager(managerKey).Address
	poolID := c.config.GetPool(poolKey).Address
	return tx.MoveCall(c.target("get_balance_manager_referral_id"), []stx.Argument{tx.Object(managerID), tx.Object(poolID)}, nil)
}

func (c *BalanceManagerContract) RevokeTradeCap(tx *stx.Transaction, managerKey, tradeCapID string) stx.Argument {
	managerID := c.config.GetBalanceManager(managerKey).Address
	return tx.MoveCall(c.target("revoke_trade_cap"), []stx.Argument{tx.Object(managerID), tx.Object(tradeCapID)}, nil)
}
