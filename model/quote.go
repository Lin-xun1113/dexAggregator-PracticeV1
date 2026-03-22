// Package model 定义了聚合器使用的数据结构
package model

// QuoteResult 代表了一个 DEX 返回的询价结果
type QuoteResult struct {
	DEXName   string  // DEX 的名称，例如 "Uniswap V3", "Sushiswap"
	AmountOut float64 // 换到的目标代币数量
	Latency   int64   // 查询耗费的时间 (毫秒)
	Err       error   // 如果查询失败，记录具体的错误信息
}
