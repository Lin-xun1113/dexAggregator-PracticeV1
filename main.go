package main

import (
	"fmt"
	"time"

	"dex-aggregator/dex"
	"dex-aggregator/model"
)

func main() {
	// 初始化随机数种子 (Go 1.20 之前需要手动初始化)
	// rand.Seed(time.Now().UnixNano())

	fmt.Println("🚀 启动简易多 DEX 并发询价器...\n")

	fromToken := "ETH"
	toToken := "USDC"
	amount := 1.5 // 想用 1.5 个 ETH 换 USDC

	// 1. 初始化我们要查询的 DEX 列表
	// 这里放的虽然是 MockDEX，但它们都实现了 DEX 接口。
	// 当你写好真正的 uniswap.go 后，只需要在这里把 NewMockDEX 换成 NewUniswap(api_url) 即可！
	dexes := []dex.DEX{
		dex.NewBinanceDEX(),
		dex.NewMockDEX("🦄 Uniswap V3", 100, 800),    // 模拟 100~800ms 延迟
		dex.NewMockDEX("🍣 Sushiswap", 300, 1500),    // 模拟 300~1500ms 延迟 (稍微慢点)
		dex.NewMockDEX("📈 Curve Finance", 200, 500), // 模拟 200~500ms 延迟
	}

	// // （挑战A版本已废弃）2. 准备接收结果的 Slice (切片)
	// var results []*model.QuoteResult
	// // （挑战A版本已废弃）我们需要一个锁来保护 results 切片，因为多个 goroutine 可能会同时往里面 append
	// var mu sync.Mutex
	// // （挑战A版本已废弃）3. 使用 WaitGroup 等待所有的询价结束
	// var wg sync.WaitGroup

	// 2. 准备接收结果的 Channel (通道)
	// 使用带缓冲的 Channel，大小等于 DEX 数量，这样即使所有 goroutine 同时写入也不会阻塞
	ch := make(chan *model.QuoteResult, len(dexes))

	startTime := time.Now() // 记录开始时间

	fmt.Printf("开始询价: 用 %.1f %s 兑换 %s\n", amount, fromToken, toToken)
	fmt.Println("--------------------------------------------------")

	// 遍历每个 DEX，为每个 DEX 启动一个专门的 goroutine
	for _, d := range dexes {
		// wg.Add(1) // （挑战A版本已废弃）启动前计数器 +1

		// 启动 Goroutine
		// ⚠️ 注意：这里把 `d` 作为参数传进去，为了避免闭包捕获循环变量的问题 (虽然 Go 1.22 修复了，但显式传递更好)
		go func(exchange dex.DEX) {
			// （挑战A版本已废弃）defer wg.Done() // 函数退出时计数器 -1

			// 调用接口的询价方法 (此时它们会各自去 sleep 模拟网络耗时)
			res, err := exchange.GetQuote(fromToken, toToken, amount)
			if err != nil {
				fmt.Printf("❌ [%s] 询价失败: %v\n", exchange.Name(), err)
				return // 提早退出这个 goroutine
			}

			// （挑战A版本已废弃）成功拿到结果，我们需要加锁写入公共的 results 切片中
			// mu.Lock()
			// results = append(results, res)
			// mu.Unlock()
			ch <- res

			fmt.Printf("✅ [%s]\t 返回报价: 耗时 %dms\n", res.DEXName, res.Latency)

		}(d)
	}

	// （挑战A版本已废弃）阻塞这里，直到所有加入了 WG 的 goroutine 报告 Done
	// wg.Wait()
	timeout := time.After(1000 * time.Millisecond)

	var results []*model.QuoteResult

Loop:
	for i := 0; i < len(dexes); i++ {
		select {
		case res := <-ch:
			results = append(results, res)
		case <-timeout:
			fmt.Println("已达500ms等待时间上限，超时返回，舍弃未返回的gorountine结果")
			break Loop
		}
	}

	totalTime := time.Since(startTime)
	fmt.Println("--------------------------------------------------")
	fmt.Printf("所有人完成询价！总计耗时: %v\n", totalTime)

	// 4. 选出最优报价 (换到最多 USDC 的那一家)
	if len(results) == 0 {
		fmt.Println("所有 DEX 都询价失败了 😢")
		return
	}

	var bestQuote *model.QuoteResult
	for _, res := range results {
		// 第一次循环，先假定第一个是最好的
		if bestQuote == nil {
			bestQuote = res
			continue
		}
		// 如果找到换出代币比当前最好记录更多的，就更新 bestQuote
		if res.AmountOut > bestQuote.AmountOut {
			bestQuote = res
		}
	}

	// 打印最终结果
	fmt.Printf("\n🏆 最优推荐路径 🏆\n")
	fmt.Printf("请去 【%s】 交易！\n", bestQuote.DEXName)
	fmt.Printf("预期可换得: %.2f %s\n", bestQuote.AmountOut, toToken)
	fmt.Printf("均价: %.2f %s/%s\n", bestQuote.AmountOut/amount, toToken, fromToken)
}
