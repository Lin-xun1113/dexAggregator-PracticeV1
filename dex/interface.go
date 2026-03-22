// Package dex 定义了与各种去中心化交易所交互的接口和具体实现
package dex

import (
	"dex-aggregator/model"
)

// DEX 接口定义了所有去中心化交易所必须实现的方法
type DEX interface {
	// Name 返回该 DEX 的名字
	Name() string
	
	// GetQuote 向该 DEX 发起询价请求。
	// 参数:
	//   fromToken: 源代币的名称或地址 (例如 "ETH")
	//   toToken:   目标代币的名称或地址 (例如 "USDC")
	//   amount:    源代币的数量 (例如 1.5)
	// 返回值:
	//   *model.QuoteResult: 包含询价结果的结构体指针
	//   error: 如果查询过程中发生网络或解析错误，返回 error
	GetQuote(fromToken, toToken string, amount float64) (*model.QuoteResult, error)
}
