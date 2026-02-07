package utils

import "github.com/sui-sdks/go-sdks/deepbook_v3/types"

type CoinMap map[string]types.Coin
type PoolMap map[string]types.Pool
type MarginPoolMap map[string]types.MarginPool

type DeepbookPackageIDs struct {
	DeepbookPackageID    string
	RegistryID           string
	DeepTreasuryID       string
	MarginPackageID      string
	MarginRegistryID     string
	LiquidationPackageID string
}

var TestnetPackageIDs = DeepbookPackageIDs{
	DeepbookPackageID:    "0x22be4cade64bf2d02412c7e8d0e8beea2f78828b948118d46735315409371a3c",
	RegistryID:           "0x7c256edbda983a2cd6f946655f4bf3f00a41043993781f8674a7046e8c0e11d1",
	DeepTreasuryID:       "0x69fffdae0075f8f71f4fa793549c11079266910e8905169845af1f5d00e09dcb",
	MarginPackageID:      "0xd6a42f4df4db73d68cbeb52be66698d2fe6a9464f45ad113ca52b0c6ebd918b6",
	MarginRegistryID:     "0x48d7640dfae2c6e9ceeada197a7a1643984b5a24c55a0c6c023dac77e0339f75",
	LiquidationPackageID: "0x8d69c3ef3ef580e5bf87b933ce28de19a5d0323588d1a44b9c60b4001741aa24",
}

var MainnetPackageIDs = DeepbookPackageIDs{
	DeepbookPackageID:    "0x337f4f4f6567fcd778d5454f27c16c70e2f274cc6377ea6249ddf491482ef497",
	RegistryID:           "0xaf16199a2dff736e9f07a845f23c5da6df6f756eddb631aed9d24a93efc4549d",
	DeepTreasuryID:       "0x032abf8948dda67a271bcc18e776dbbcfb0d58c8d288a700ff0d5521e57a1ffe",
	MarginPackageID:      "0x97d9473771b01f77b0940c589484184b49f6444627ec121314fae6a6d36fb86b",
	MarginRegistryID:     "0x0e40998b359a9ccbab22a98ed21bd4346abf19158bc7980c8291908086b3a742",
	LiquidationPackageID: "0x73c593882cdb557703e903603f20bd373261fe6ba6e1a40515f4b62f10553e6a",
}

var TestnetCoins = CoinMap{
	"DEEP": {
		Address: "0x36dbef866a1d62bf7328989a10fb2f07d769f4ee587c0de4a0a256e57e0a58a8",
		Type:    "0x36dbef866a1d62bf7328989a10fb2f07d769f4ee587c0de4a0a256e57e0a58a8::deep::DEEP",
		Scalar:  1_000_000,
	},
	"SUI": {
		Address: "0x2",
		Type:    "0x2::sui::SUI",
		Scalar:  1_000_000_000,
	},
	"DBUSDC": {
		Address: "0xf7152c05930480cd740d7311b5b8b45c6f488e3a53a11c3f74a6fac36a52e0d7",
		Type:    "0xf7152c05930480cd740d7311b5b8b45c6f488e3a53a11c3f74a6fac36a52e0d7::DBUSDC::DBUSDC",
		Scalar:  1_000_000,
	},
}

var MainnetCoins = CoinMap{
	"DEEP": {
		Address: "0xdeeb7a4662eec9f2f3def03fb937a663dddaa2e215b8078a284d026b7946c270",
		Type:    "0xdeeb7a4662eec9f2f3def03fb937a663dddaa2e215b8078a284d026b7946c270::deep::DEEP",
		Scalar:  1_000_000,
	},
	"SUI": {
		Address: "0x2",
		Type:    "0x2::sui::SUI",
		Scalar:  1_000_000_000,
	},
	"USDC": {
		Address: "0xdba34672e30cb065b1f93e3ab55318768fd6fef66c15942c9f7cb846e2f900e7",
		Type:    "0xdba34672e30cb065b1f93e3ab55318768fd6fef66c15942c9f7cb846e2f900e7::usdc::USDC",
		Scalar:  1_000_000,
	},
}

var TestnetPools = PoolMap{
	"DEEP_SUI": {
		Address:   "0x48c95963e9eac37a316b7ae04a0deb761bcdcc2b67912374d6036e7f0e9bae9f",
		BaseCoin:  "DEEP",
		QuoteCoin: "SUI",
	},
	"SUI_DBUSDC": {
		Address:   "0x1c19362ca52b8ffd7a33cee805a67d40f31e6ba303753fd3a4cfdfacea7163a5",
		BaseCoin:  "SUI",
		QuoteCoin: "DBUSDC",
	},
}

var MainnetPools = PoolMap{
	"DEEP_SUI": {
		Address:   "0xb663828d6217467c8a1838a03793da896cbe745b150ebd57d82f814ca579fc22",
		BaseCoin:  "DEEP",
		QuoteCoin: "SUI",
	},
	"SUI_USDC": {
		Address:   "0xe05dafb5133bcffb8d59f4e12465dc0e9faeaa05e3e342a08fe135800e3e4407",
		BaseCoin:  "SUI",
		QuoteCoin: "USDC",
	},
}

var TestnetMarginPools = MarginPoolMap{}
var MainnetMarginPools = MarginPoolMap{}

var TestnetPythConfigs = map[string]string{}
var MainnetPythConfigs = map[string]string{}
