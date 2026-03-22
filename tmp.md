# 入门指南：Go 工程化思维与多 DEX 并发询价器实战

## 0. 心态转变：从“码农”到“工程师”

当你不再问“这行代码怎么敲”，而是问“这个工程怎么开局”时，你就已经脱离了“学基础语法的初学者”，开始像一个真正的“后端工程师”思考了。

对于一个从零开始的 Go 工程项目，我们通常遵循以下工程化流程。

---

## 1. 需求拆解：确定输入、处理与输出

在写代码之前，先明确项目的核心链路：

**输入 (Input)：用户提供的参数**

- 源代币 (From Token，如 ETH)
- 目标代币 (To Token，如 USDC)
- 金额 (Amount)

**处理 (Process)：核心业务逻辑**

- **数据源获取**：同时向 Uniswap V3、Sushiswap 等 DEX 发起查询。
- **性能优化**：由于网络请求耗时，必须使用 **Goroutine** 并发执行。
- **数据标准化**：定义统一的结构来存放不同 DEX 的报价。

**输出 (Output)：呈现给用户的结果**

- **超时控制**：超出预设时间的数据直接丢弃。
- **结果汇总**：对比报价，输出最优路径。

---

## 2. 数据结构设计 (Model / Entity)

后端开发应遵循“结构体优先”原则。首先定义统一的报价返回格式：

```go
// model/quote.go
package model

type QuoteResult struct {
    DEXName     string  // 例如 "Uniswap"
    AmountOut   float64 // 换到的金额，例如 1850.5
    Latency     int64   // 耗时，比如 200ms
    Err         error   // 用于记录查询失败的错误
}
```

---

## 3. 面向接口编程 (Interface)

为了方便后续接入更多 DEX（如 1inch, Curve），应将“询价”动作抽象为接口：

```go
// dex/interface.go
package dex

import "dex-aggregator/model"

type DEX interface {
    Name() string
    GetQuote(fromToken, toToken string, amount float64) (*model.QuoteResult, error)
}
```

> [!TIP]
> 只要实现了该接口，任何新的交易所都能轻松无缝地接入核心逻辑。

---

## 4. 目录结构规范 (Project Layout)

遵循 Go 标准目录规范（Project Layout）进行模块拆分：

```text
aggregator-cli/            # 项目根目录
├── go.mod                 # 依赖管理
├── main.go                # 入口：组装模块，发号施令
├── dex/                   # 业务逻辑包
│   ├── interface.go       # 接口定义
│   ├── uniswap.go         # Uniswap 实现
│   ├── sushiswap.go       # Sushiswap 实现
│   └── mock_dex.go        # 测试用的模拟数据 (打桩)
└── model/                 # 数据结构定义
    └── quote.go           
```

---

## 5. 敏捷开发：分阶段执行计划

1. **搭建框架**：建立目录，定义 `model` 和 `interface`。
2. **Mock 验证**：编写 `mock_dex.go` 模拟随机响应和延迟，调通主循环逻辑。
3. **并发实装**：在 `main.go` 中使用 `sync.WaitGroup` 启动并发查询，实现超时控制。
4. **接入实测**：最后填充 `uniswap.go` 等真实 API 调用逻辑。

---

## 6. 即刻开始：项目初始化

> [!IMPORTANT]
> 请在指定目录下完成以下操作。

1. **环境准备**：确认进入 `d:\Code\dex-aggregator` 目录。
2. **模块初始化**：

   ```bash
   go mod init dex-aggregator
   ```

3. **完成反馈**：搞定后请告知我，我们将开始建立上述目录结构！
