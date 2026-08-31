# Go 后端中的设计模式：从“为什么这样设计”到 Agent / RAG 工程实践

> 设计模式真正有价值的地方，不是让代码看起来“高级”，而是帮助我们控制变化、降低耦合、隔离复杂度，并让系统在需求不断增加时仍然可维护。

在学习设计模式时，最常见的问题不是“记不住定义”，而是：

- 为什么这里需要这个模式？
- 不用这个模式会发生什么？
- 两个看起来很像的模式到底有什么区别？
- Go 不是面向对象语言，还需要设计模式吗？
- 真实后端项目里究竟哪些模式值得用？
- Agent、RAG、Tool、Middleware 这种新系统和传统设计模式有什么关系？

这篇文章不按 GoF 23 种模式逐个背定义，而是从**工程问题**出发，重点讲 Go 后端、RAG、Agent 系统中最值得掌握的一组设计思想。

---

# 一、先建立一个核心认知：设计模式是在“管理变化”

很多人第一次学设计模式，会把它理解成：

```text
工厂模式 = 创建对象
策略模式 = 切换算法
适配器模式 = 转换接口
观察者模式 = 发布订阅
```

这些都没错，但还不够。

设计模式真正解决的是：

> **系统中哪些东西会变化，以及我们怎样让这些变化尽量不要扩散。**

例如一个 Agent 平台支持多个模型：

```text
OpenAI
DeepSeek
Qwen
Claude
Ollama
```

真正的问题不是“怎么 new 一个 OpenAIClient”。

而是：

```text
以后再增加一个模型 Provider 时，
我要改多少地方？
```

如果你每增加一个 Provider，都要改：

```text
Handler
Service
Agent Runtime
RAG
配置加载
测试代码
```

说明“Provider 的变化”扩散到了整个系统。

好的设计应该把这种变化收敛到一个边界：

```text
业务代码
   ↓
统一 LLM 接口
   ↓
Provider Factory / Registry
   ↓
OpenAI / DeepSeek / Qwen / Ollama
```

于是新增一个 Provider，主要只增加一个实现，而不是到处修改原有代码。

这就是设计模式最核心的价值：

> **把变化关进笼子里。**

---

# 二、Go 语言中的设计模式和 Java 有什么不同

Go 可以使用设计模式，但实现方式通常比 Java 更轻。

Go 的几个特性会影响设计模式：

```text
interface 隐式实现
函数是一等公民
组合优于继承
struct 很轻量
闭包很好用
泛型可以减少样板代码
```

所以在 Go 中，很多模式并不会长成经典 UML 教材里的样子。

例如 Java 的 Strategy 可能需要：

```text
Strategy interface
ConcreteStrategyA
ConcreteStrategyB
Context
```

Go 很多时候直接：

```go
type RetryStrategy func(ctx context.Context) error
```

就够了。

因此：

> **不要为了“使用设计模式”而制造大量 interface 和 struct。**

设计模式应该服务于复杂度，而不是制造复杂度。

---

# 三、设计模式的三大类别

经典 GoF 模式分成三类：

| 类别 | 核心问题 | 常见模式 |
|---|---|---|
| 创建型 | 对象怎么创建 | Factory、Builder、Singleton |
| 结构型 | 对象怎么组合 | Adapter、Decorator、Proxy |
| 行为型 | 对象怎么协作 | Strategy、Observer、Chain、State、Command、Template Method |

但在真实 Go 工程里，还有几种非常常见的“工程模式”：

```text
Registry
Middleware
Dependency Injection
Pipeline
Functional Options
Repository
```

它们不一定都属于 GoF 23 种模式，但工程价值非常高。

---

# 四、策略模式 Strategy：把“可替换算法”抽出来

## 1. 它解决什么问题

假设一个 RAG 系统支持：

```text
Vector Search
Keyword Search
Hybrid Search
```

最直观的代码可能是：

```go
func Retrieve(ctx context.Context, mode string, query string) ([]Document, error) {
    switch mode {
    case "vector":
        return vectorRetrieve(ctx, query)
    case "keyword":
        return keywordRetrieve(ctx, query)
    case "hybrid":
        return hybridRetrieve(ctx, query)
    default:
        return nil, errors.New("unsupported retrieve mode")
    }
}
```

第一版完全可以。

问题是系统继续增长以后：

```text
vector
keyword
hybrid
graph
elastic
milvus
rerank-only
metadata-filter
```

同一个函数会越来越大。

而且业务逻辑逐渐和“具体检索算法”耦合。

策略模式的核心思想是：

> **把一组可以互换的行为抽象成统一接口，让调用者依赖能力，而不是依赖具体算法。**

---

## 2. Go 示例

```go
type Retriever interface {
    Retrieve(ctx context.Context, query string) ([]Document, error)
}

type VectorRetriever struct {
    store VectorStore
}

func (r *VectorRetriever) Retrieve(
    ctx context.Context,
    query string,
) ([]Document, error) {
    return r.store.Search(ctx, query)
}

type KeywordRetriever struct {
    index KeywordIndex
}

func (r *KeywordRetriever) Retrieve(
    ctx context.Context,
    query string,
) ([]Document, error) {
    return r.index.Search(ctx, query)
}

type HybridRetriever struct {
    vector  Retriever
    keyword Retriever
}

func (r *HybridRetriever) Retrieve(
    ctx context.Context,
    query string,
) ([]Document, error) {
    vectorDocs, err := r.vector.Retrieve(ctx, query)
    if err != nil {
        return nil, err
    }

    keywordDocs, err := r.keyword.Retrieve(ctx, query)
    if err != nil {
        return nil, err
    }

    return RRF(vectorDocs, keywordDocs), nil
}
```

业务层：

```go
type QueryService struct {
    retriever Retriever
}

func (s *QueryService) Query(
    ctx context.Context,
    query string,
) ([]Document, error) {
    return s.retriever.Retrieve(ctx, query)
}
```

---

## 3. 为什么这叫“策略”

因为 QueryService 只知道：

```text
我需要一个 Retriever
```

它并不关心：

```text
具体是 VectorRetriever
还是 KeywordRetriever
还是 HybridRetriever
```

Retriever 本身就是一个策略。

从依赖关系看：

```text
QueryService
    ↓
Retriever interface
    ↓
┌──────────────┬──────────────┬──────────────┐
Vector       Keyword        Hybrid
```

这叫：

> **依赖抽象，而不是依赖具体实现。**

---

## 4. 策略模式的真正价值

最大的价值不是“消灭 switch”。

而是：

```text
算法变化
   ↓
被限制在策略实现内部
```

上层业务不需要知道内部差异。

这样做的好处：

- 新增策略时旧业务代码通常不用修改；
- 更容易单元测试；
- 可以在运行时根据配置动态选择；
- 不同策略可以独立演进；
- 便于灰度和 A/B Test。

---

## 5. 什么时候不应该用

如果系统永远只有：

```text
A
B
```

两个分支，并且逻辑非常简单：

```go
if mode == "a" {
    return a()
}
return b()
```

完全没有必要为了“模式”引入：

```text
Strategy interface
StrategyFactory
StrategyRegistry
StrategyContext
```

否则反而过度设计。

### 判断标准

当你发现：

> “这里的核心变化点是一组行为，并且未来大概率会增加、替换、测试不同实现。”

再考虑 Strategy。

---

# 五、工厂模式 Factory：把“对象如何创建”从业务中拿出去

## 1. 工厂模式解决什么问题

假设 Agent 平台支持多个模型：

```text
OpenAI
DeepSeek
Qwen
Ollama
```

最直接的实现：

```go
func NewLLM(provider string, cfg Config) (LLM, error) {
    switch provider {
    case "openai":
        return openai.New(cfg.APIKey, cfg.BaseURL), nil
    case "deepseek":
        return deepseek.New(cfg.APIKey, cfg.BaseURL), nil
    case "qwen":
        return qwen.New(cfg.APIKey, cfg.BaseURL), nil
    default:
        return nil, errors.New("unsupported provider")
    }
}
```

这其实已经是一个简单工厂。

---

## 2. 为什么要有工厂

真正的问题不是：

> “new 很麻烦。”

而是不同 Provider 的初始化逻辑可能不同。

例如：

```text
OpenAI
→ API Key
→ Base URL
→ Timeout

Ollama
→ 本地地址
→ Model Name

Azure OpenAI
→ Endpoint
→ Deployment
→ API Version
→ Key

内部模型服务
→ Service Discovery
→ Token
→ TLS
```

如果让业务层知道所有这些初始化细节：

```text
AgentService
   ↓
大量 Provider 判断
   ↓
大量 SDK 初始化代码
```

业务代码就会被“对象构建细节”污染。

工厂的思想：

> **调用者只表达“我要什么”，工厂负责“怎么创建”。**

---

## 3. Go 示例

```go
type LLM interface {
    Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error)
}

func NewLLM(provider string, cfg ProviderConfig) (LLM, error) {
    switch provider {
    case "openai":
        return NewOpenAILLM(cfg)
    case "deepseek":
        return NewDeepSeekLLM(cfg)
    case "ollama":
        return NewOllamaLLM(cfg)
    default:
        return nil, fmt.Errorf("unsupported provider: %s", provider)
    }
}
```

调用者：

```go
llm, err := NewLLM(config.Provider, config.ProviderConfig)
if err != nil {
    return err
}

agent := NewAgent(llm)
```

Agent 并不知道 OpenAI 是怎么初始化的。

---

## 4. 工厂模式和策略模式什么关系

这两个模式经常一起出现：

```text
Factory
负责创建 Strategy

Strategy
负责执行行为
```

例如：

```go
retriever, err := NewRetriever(config.RetrieveMode)
docs, err := retriever.Retrieve(ctx, query)
```

这里：

```text
NewRetriever
→ Factory

Retriever
→ Strategy
```

所以不要把它们混成一个概念。

---

# 六、注册表模式 Registry：比大 switch 更适合插件化系统

Registry 不是 GoF 经典 23 种之一，但 Go 工程里非常常见。

尤其适合：

```text
Tool
Plugin
Parser
LLM Provider
Command
Serializer
Handler
```

---

## 1. 为什么 Factory 继续增长以后会有问题

前面的工厂：

```go
switch provider {
case "openai":
case "deepseek":
case "qwen":
case "ollama":
...
}
```

Provider 越来越多以后：

```text
新增 Provider
↓
必须修改中央 Factory
```

这说明所有插件仍然耦合到一个中心文件。

Registry 的思路是：

> **实现自己注册自己，调用方按名称查找。**

---

## 2. Go 示例

```go
type Tool interface {
    Name() string
    Execute(ctx context.Context, args map[string]any) (any, error)
}

type ToolRegistry struct {
    mu    sync.RWMutex
    tools map[string]Tool
}

func NewToolRegistry() *ToolRegistry {
    return &ToolRegistry{
        tools: make(map[string]Tool),
    }
}

func (r *ToolRegistry) Register(tool Tool) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    name := tool.Name()

    if _, exists := r.tools[name]; exists {
        return fmt.Errorf("tool already exists: %s", name)
    }

    r.tools[name] = tool
    return nil
}

func (r *ToolRegistry) Get(name string) (Tool, bool) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    tool, ok := r.tools[name]
    return tool, ok
}
```

注册：

```go
registry := NewToolRegistry()

_ = registry.Register(&SearchTool{})
_ = registry.Register(&CalculatorTool{})
_ = registry.Register(&KnowledgeBaseTool{})
```

运行时：

```go
tool, ok := registry.Get(call.Name)
if !ok {
    return nil, errors.New("tool not found")
}

return tool.Execute(ctx, call.Arguments)
```

---

## 3. Registry 的设计思想

它把系统从：

```text
中央 switch 决定所有实现
```

变成：

```text
中央只维护注册机制
具体实现独立注册
```

于是扩展系统时：

```text
新增 Tool
→ 实现接口
→ 注册
```

核心 Runtime 不需要知道这个 Tool 是什么。

这对 Agent 平台尤其重要，因为 Tool 天生就是插件化能力。

---

## 4. Registry 为什么经常需要并发控制

如果 Registry：

- 启动时注册一次；
- 运行时只读；

那么初始化完成后理论上可以不再修改。

但如果支持：

```text
动态加载 Tool
动态卸载 MCP Tool
在线启停插件
热更新 Provider
```

Registry 就可能在运行时同时发生：

```text
Register
Get
Remove
List
```

因此需要：

```go
sync.RWMutex
```

或者 Copy-on-Write 等机制。

---

# 七、适配器模式 Adapter：隔离第三方差异

## 1. Adapter 解决的不是“算法选择”

Adapter 解决的是：

> **两个系统接口长得不一样，但我们希望它们能协作。**

假设你的业务层希望统一调用：

```go
type Parser interface {
    Parse(ctx context.Context, file File) (*Document, error)
}
```

但第三方 Python Parser API 是：

```http
POST /parse
```

返回：

```json
{
  "markdown": "...",
  "pages": 20,
  "tables": [...]
}
```

你不希望业务层到处写：

```go
http.NewRequest(...)
json.Marshal(...)
json.Unmarshal(...)
```

于是设计 Adapter。

---

## 2. Go 示例

```go
type Parser interface {
    Parse(ctx context.Context, file File) (*Document, error)
}

type PythonParserAdapter struct {
    client  *http.Client
    baseURL string
}

func (p *PythonParserAdapter) Parse(
    ctx context.Context,
    file File,
) (*Document, error) {
    reqBody := struct {
        URL string `json:"url"`
    }{
        URL: file.URL,
    }

    data, err := json.Marshal(reqBody)
    if err != nil {
        return nil, err
    }

    req, err := http.NewRequestWithContext(
        ctx,
        http.MethodPost,
        p.baseURL+"/parse",
        bytes.NewReader(data),
    )
    if err != nil {
        return nil, err
    }

    resp, err := p.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var apiResp struct {
        Markdown string `json:"markdown"`
        Pages    int    `json:"pages"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
        return nil, err
    }

    return &Document{
        Content: apiResp.Markdown,
        Pages:   apiResp.Pages,
    }, nil
}
```

业务层：

```go
doc, err := parser.Parse(ctx, file)
```

它只认识 Parser。

---

## 3. Adapter 最大的价值：隔离第三方依赖

假设未来从：

```text
Python FastAPI Parser
```

切到：

```text
MinerU Cloud API
PP-Structure
自研 Go Parser
```

如果业务层统一依赖 Parser：

```text
RAG Index Service
      ↓
    Parser
```

那变化主要只发生在 Adapter 层。

否则第三方 SDK 的 DTO、错误码、字段命名会直接污染整个业务。

这叫：

> **Anti-Corruption Layer，防腐层。**

在领域驱动设计中，这也是非常重要的思想。

---

# 八、Strategy vs Adapter：最容易混淆的一组

它们都经常表现成：

```go
interface
+
多个实现
```

但意图完全不同。

## Strategy

关注：

> **同一件事有不同做法。**

例如：

```text
Retrieve
├─ Vector
├─ Keyword
└─ Hybrid
```

这些都是“检索”，只是算法不同。

---

## Adapter

关注：

> **外部系统接口不一样，我把它转换成内部统一接口。**

例如：

```text
内部 Parser
↑
MinerU Adapter
PP-Structure Adapter
第三方 SaaS Adapter
```

这里重点不是算法替换，而是接口兼容和依赖隔离。

---

## 一句话区分

```text
Strategy：
统一“行为选择”

Adapter：
统一“接口差异”
```

---

# 九、建造者模式 Builder：复杂对象不要一次性构造

## 1. 什么场景需要 Builder

假设 Agent 有很多配置：

```text
LLM
System Prompt
Tools
Memory
Retriever
Trace
Max Steps
Timeout
Temperature
Callbacks
```

如果直接：

```go
agent := NewAgent(
    llm,
    prompt,
    tools,
    memory,
    retriever,
    tracer,
    10,
    30*time.Second,
    0.7,
    callbacks,
)
```

很难读。

并且参数顺序容易传错。

---

## 2. Builder 示例

```go
type Agent struct {
    llm       LLM
    prompt    string
    tools     []Tool
    retriever Retriever
    tracer    Tracer
    maxSteps  int
}

type AgentBuilder struct {
    agent Agent
}

func NewAgentBuilder() *AgentBuilder {
    return &AgentBuilder{
        agent: Agent{
            maxSteps: 8,
        },
    }
}

func (b *AgentBuilder) WithLLM(llm LLM) *AgentBuilder {
    b.agent.llm = llm
    return b
}

func (b *AgentBuilder) WithPrompt(prompt string) *AgentBuilder {
    b.agent.prompt = prompt
    return b
}

func (b *AgentBuilder) WithTools(tools ...Tool) *AgentBuilder {
    b.agent.tools = append(b.agent.tools, tools...)
    return b
}

func (b *AgentBuilder) WithRetriever(r Retriever) *AgentBuilder {
    b.agent.retriever = r
    return b
}

func (b *AgentBuilder) Build() (*Agent, error) {
    if b.agent.llm == nil {
        return nil, errors.New("llm is required")
    }

    return &b.agent, nil
}
```

调用：

```go
agent, err := NewAgentBuilder().
    WithLLM(llm).
    WithPrompt(systemPrompt).
    WithTools(searchTool, kbTool).
    WithRetriever(retriever).
    Build()
```

可读性明显更好。

---

# 十、Go 更常用的替代：Functional Options

在 Go 中，Builder 很多时候会被 Functional Options 替代。

---

## 1. 示例

```go
type AgentOption func(*Agent)

func WithPrompt(prompt string) AgentOption {
    return func(a *Agent) {
        a.prompt = prompt
    }
}

func WithTools(tools ...Tool) AgentOption {
    return func(a *Agent) {
        a.tools = append(a.tools, tools...)
    }
}

func WithMaxSteps(n int) AgentOption {
    return func(a *Agent) {
        a.maxSteps = n
    }
}

func NewAgent(llm LLM, opts ...AgentOption) *Agent {
    a := &Agent{
        llm:      llm,
        maxSteps: 8,
    }

    for _, opt := range opts {
        opt(a)
    }

    return a
}
```

调用：

```go
agent := NewAgent(
    llm,
    WithPrompt("You are a helpful assistant"),
    WithTools(searchTool, kbTool),
    WithMaxSteps(10),
)
```

---

## 2. Functional Options 为什么在 Go 很流行

因为它兼顾：

```text
默认值
可选参数
可读性
向后兼容
```

假设最初：

```go
NewClient(baseURL string)
```

后来增加：

```text
timeout
retry
transport
trace
```

如果一直修改函数参数：

```go
NewClient(baseURL, timeout, retry, transport, trace)
```

老代码会全部受影响。

Functional Options 可以：

```go
NewClient(
    baseURL,
    WithTimeout(3*time.Second),
    WithRetry(2),
)
```

旧调用：

```go
NewClient(baseURL)
```

仍然可以继续工作。

---

# 十一、装饰器模式 Decorator：在不修改核心实现的情况下增加能力

## 1. 假设你有 Tool

```go
type Tool interface {
    Execute(ctx context.Context, input string) (string, error)
}
```

SearchTool 本身只负责搜索。

但你还需要：

```text
日志
Trace
权限
超时
指标
重试
```

如果全部塞进 SearchTool：

```go
func (t *SearchTool) Execute(...) {
    // auth
    // trace
    // log
    // metrics
    // timeout
    // retry
    // real search
}
```

核心业务会被横切逻辑淹没。

Decorator 的思想：

> **用一个实现包装另一个实现，在调用前后增加能力。**

---

## 2. Trace Decorator

```go
type TraceTool struct {
    next   Tool
    tracer Tracer
}

func (t *TraceTool) Execute(
    ctx context.Context,
    input string,
) (string, error) {
    ctx, span := t.tracer.Start(ctx, "tool.execute")
    defer span.End()

    result, err := t.next.Execute(ctx, input)

    if err != nil {
        span.RecordError(err)
    }

    return result, err
}
```

Timeout Decorator：

```go
type TimeoutTool struct {
    next    Tool
    timeout time.Duration
}

func (t *TimeoutTool) Execute(
    ctx context.Context,
    input string,
) (string, error) {
    ctx, cancel := context.WithTimeout(ctx, t.timeout)
    defer cancel()

    return t.next.Execute(ctx, input)
}
```

组合：

```go
tool :=
    &TraceTool{
        tracer: tracer,
        next: &TimeoutTool{
            timeout: 3 * time.Second,
            next:    searchTool,
        },
    }
```

调用链：

```text
Trace
 ↓
Timeout
 ↓
SearchTool
```

---

# 十二、Middleware：本质上经常就是函数式 Decorator

Go Web 开发里最常见的 Middleware：

```go
type Handler func(ctx context.Context) error

type Middleware func(next Handler) Handler
```

例如：

```go
func LoggingMiddleware(next Handler) Handler {
    return func(ctx context.Context) error {
        start := time.Now()

        err := next(ctx)

        log.Printf("cost=%s err=%v", time.Since(start), err)
        return err
    }
}
```

组合：

```go
handler := coreHandler

handler = LoggingMiddleware(handler)
handler = TraceMiddleware(handler)
handler = AuthMiddleware(handler)
```

最终调用：

```text
Auth
 ↓
Trace
 ↓
Logging
 ↓
Core
```

这和 Decorator 的思想高度一致。

---

# 十三、为什么 Middleware 常常要倒序包装

假设配置：

```go
middlewares := []Middleware{
    Auth,
    Trace,
    Logging,
}
```

我们希望执行：

```text
Auth → Trace → Logging → Core
```

如果直接正序：

```go
for _, m := range middlewares {
    handler = m(handler)
}
```

最终会变成：

```text
Logging → Trace → Auth → Core
```

因为后包装的在最外层。

所以一般要：

```go
for i := len(middlewares) - 1; i >= 0; i-- {
    handler = middlewares[i](handler)
}
```

这样才能保持配置顺序与执行顺序一致。

这是面试里一个很好用的工程细节。

---

# 十四、Decorator vs Middleware vs Chain of Responsibility

这三个很容易被混在一起。

## Decorator

关注：

> **给一个对象叠加能力。**

每一层一般都会调用 next。

```text
Trace
 ↓
Metrics
 ↓
Timeout
 ↓
Core
```

---

## Middleware

在 Web / RPC / Agent Runtime 中，通常就是：

> **一种工程化的 Decorator / Interceptor 机制。**

它不是 GoF 独立模式，但思想高度接近 Decorator。

---

## Chain of Responsibility

关注：

> **请求沿着多个处理器传递，每个处理器决定处理、继续、停止。**

区别在于：

Decorator 通常：

```text
每一层都在增强同一个调用
```

责任链可以：

```text
某一层直接终止
某一层处理后不再向后传
某一层选择下一个节点
```

---

# 十五、责任链模式 Chain of Responsibility

## 1. 一个典型场景：请求校验

例如创建知识库：

```text
参数格式校验
 ↓
权限校验
 ↓
额度校验
 ↓
内容安全校验
 ↓
真正创建
```

每一步都可能阻断请求。

---

## 2. Go 示例

```go
type Handler interface {
    Handle(ctx context.Context, req *Request) error
    SetNext(next Handler)
}

type BaseHandler struct {
    next Handler
}

func (h *BaseHandler) SetNext(next Handler) {
    h.next = next
}

func (h *BaseHandler) nextHandle(
    ctx context.Context,
    req *Request,
) error {
    if h.next == nil {
        return nil
    }

    return h.next.Handle(ctx, req)
}
```

权限处理器：

```go
type AuthHandler struct {
    BaseHandler
}

func (h *AuthHandler) Handle(
    ctx context.Context,
    req *Request,
) error {
    if !req.User.HasPermission {
        return errors.New("permission denied")
    }

    return h.nextHandle(ctx, req)
}
```

额度处理器：

```go
type QuotaHandler struct {
    BaseHandler
}

func (h *QuotaHandler) Handle(
    ctx context.Context,
    req *Request,
) error {
    if req.User.Quota <= 0 {
        return errors.New("quota exceeded")
    }

    return h.nextHandle(ctx, req)
}
```

---

# 十六、观察者模式 Observer：一个事件，多方响应

## 1. 场景

假设知识库文档索引完成后，需要：

```text
更新任务状态
发送 SSE
记录审计日志
统计指标
通知用户
```

最差的写法是：

```go
func FinishIndex() {
    updateTask()
    sendSSE()
    writeAudit()
    reportMetrics()
    notifyUser()
}
```

这样索引业务和所有下游逻辑直接耦合。

Observer 的思想：

> **发布事件，订阅者自己决定如何响应。**

---

## 2. Go 示例

```go
type Event struct {
    Type string
    Data any
}

type Listener interface {
    OnEvent(ctx context.Context, event Event)
}

type EventBus struct {
    listeners map[string][]Listener
}

func NewEventBus() *EventBus {
    return &EventBus{
        listeners: make(map[string][]Listener),
    }
}

func (b *EventBus) Subscribe(eventType string, listener Listener) {
    b.listeners[eventType] = append(
        b.listeners[eventType],
        listener,
    )
}

func (b *EventBus) Publish(ctx context.Context, event Event) {
    listeners := b.listeners[event.Type]

    for _, listener := range listeners {
        listener.OnEvent(ctx, event)
    }
}
```

索引完成：

```go
bus.Publish(ctx, Event{
    Type: "document.indexed",
    Data: docID,
})
```

订阅者：

```text
SSEListener
AuditListener
MetricsListener
```

---

## 3. Observer 和消息队列有什么区别

Observer 是设计模式。

Kafka / RabbitMQ / Redis Streams 是具体基础设施。

你完全可以实现：

```text
Observer
↓
进程内 EventBus
```

也可以：

```text
Observer
↓
Kafka
```

如果只是单进程内部解耦：

```text
EventBus
```

可能足够。

如果需要：

```text
跨服务
持久化
重试
消费进度
削峰
```

就更适合 MQ。

---

# 十七、状态模式 State：让状态迁移成为显式模型

## 1. 为什么一堆 if/switch 会越来越危险

假设文档任务有：

```text
Pending
Parsing
Parsed
Indexing
Completed
Failed
```

最初：

```go
if task.Status == "pending" {
    ...
}
```

后来：

```go
if status == "pending" {
...
} else if status == "parsing" {
...
} else if status == "parsed" {
...
}
```

再后来各种业务都在判断状态：

```text
API
Worker
Retry
Cancel
Resume
Admin
```

就容易出现非法迁移：

```text
Completed → Parsing
Failed → Completed
Pending → Completed
```

---

## 2. 状态机的核心思想

不要只存“当前状态”。

还要显式定义：

> **什么状态可以迁移到什么状态。**

例如：

```go
var allowedTransitions = map[Status]map[Status]bool{
    StatusPending: {
        StatusParsing: true,
    },
    StatusParsing: {
        StatusParsed: true,
        StatusFailed: true,
    },
    StatusParsed: {
        StatusIndexing: true,
    },
    StatusIndexing: {
        StatusCompleted: true,
        StatusFailed:    true,
    },
}
```

迁移：

```go
func CanTransition(from, to Status) bool {
    return allowedTransitions[from][to]
}
```

---

## 3. 为什么数据库还要做 CAS

仅仅在 Go 内存里检查：

```go
if task.Status == oldStatus {
    task.Status = newStatus
}
```

并不能解决并发问题。

两个 Worker 可能同时读到：

```text
status = parsed
```

然后同时开始 indexing。

更稳妥的方式：

```sql
UPDATE document_tasks
SET status = 'indexing'
WHERE id = ?
  AND status = 'parsed';
```

然后检查：

```text
RowsAffected == 1
```

这其实是数据库层的 CAS：

> Compare-And-Swap。

状态机负责：

```text
业务合法性
```

数据库条件更新负责：

```text
并发一致性
```

两者是不同层次的问题。

---

# 十八、命令模式 Command：把“要执行的操作”变成对象

Agent Tool Call 天然很接近 Command。

LLM 返回：

```json
{
  "name": "search_web",
  "arguments": {
    "query": "Go GMP"
  }
}
```

这本质上就是：

```text
Command Name
+
Command Arguments
```

Runtime 再根据命令执行。

---

## Go 示例

```go
type Command interface {
    Execute(ctx context.Context) (any, error)
}

type SearchCommand struct {
    query  string
    search SearchService
}

func (c *SearchCommand) Execute(ctx context.Context) (any, error) {
    return c.search.Search(ctx, c.query)
}
```

Command 的价值是：

```text
请求发送方
与
真正执行方
解耦
```

它还非常适合：

```text
任务队列
撤销
重放
审计
异步执行
```

---

# 十九、代理模式 Proxy：在“访问目标对象”前加一道代理

## 1. Proxy 和 Decorator 很像

两者结构都可能是：

```text
Wrapper
 ↓
Real Object
```

但意图不一样。

Proxy 强调：

> **控制访问。**

Decorator 强调：

> **增加能力。**

---

## 2. 例子：LLM Proxy

```go
type LLM interface {
    Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error)
}

type RateLimitLLMProxy struct {
    target  LLM
    limiter *rate.Limiter
}

func (p *RateLimitLLMProxy) Chat(
    ctx context.Context,
    req ChatRequest,
) (*ChatResponse, error) {
    if err := p.limiter.Wait(ctx); err != nil {
        return nil, err
    }

    return p.target.Chat(ctx, req)
}
```

这里代理层控制：

```text
是否允许访问真实 LLM
```

常见 Proxy：

```text
权限代理
缓存代理
限流代理
远程代理
懒加载代理
```

---

# 二十、Proxy vs Adapter

非常容易混。

## Proxy

接口通常保持一致：

```text
Client
 ↓
LLM Interface
 ↓
Proxy
 ↓
Real LLM
```

调用者甚至不一定知道中间有代理。

目的：

```text
控制访问
```

---

## Adapter

接口发生转换：

```text
内部 Parser Interface
 ↓
Adapter
 ↓
第三方 HTTP API
```

目的：

```text
接口兼容
```

一句话：

```text
Proxy：
挡在前面

Adapter：
夹在中间翻译
```

---

# 二十一、模板方法 Template Method：固定流程，开放部分步骤

经典 Template Method 常依赖继承。

Go 没有传统继承，所以一般使用：

```text
组合
+
接口
+
高阶函数
```

---

## 1. RAG Index Pipeline

索引流程可能固定：

```text
Parse
 ↓
Chunk
 ↓
Embed
 ↓
Store
```

但每一步实现可以替换。

这本质上就是模板方法 / Pipeline 思想。

---

## 2. Go 版实现

```go
type Parser interface {
    Parse(ctx context.Context, file File) (*Document, error)
}

type Chunker interface {
    Chunk(ctx context.Context, doc *Document) ([]Chunk, error)
}

type Embedder interface {
    Embed(ctx context.Context, chunks []Chunk) ([]Vector, error)
}

type VectorStore interface {
    Save(ctx context.Context, vectors []Vector) error
}

type IndexPipeline struct {
    parser   Parser
    chunker  Chunker
    embedder Embedder
    store    VectorStore
}

func (p *IndexPipeline) Run(
    ctx context.Context,
    file File,
) error {
    doc, err := p.parser.Parse(ctx, file)
    if err != nil {
        return err
    }

    chunks, err := p.chunker.Chunk(ctx, doc)
    if err != nil {
        return err
    }

    vectors, err := p.embedder.Embed(ctx, chunks)
    if err != nil {
        return err
    }

    return p.store.Save(ctx, vectors)
}
```

这里：

```text
流程固定
实现可替换
```

是非常典型的 Template / Pipeline 思想。

---

# 二十二、Pipeline 模式：把复杂流程拆成独立阶段

Pipeline 在后端和 AI 工程里极其常见。

例如 RAG：

```text
Parser
 ↓
Normalizer
 ↓
Chunker
 ↓
Embedder
 ↓
Indexer
```

每个阶段只做一件事。

优点：

- 责任清晰；
- 阶段可以替换；
- 阶段可以复用；
- 更容易观测；
- 更容易断点续跑；
- 更方便针对单阶段做性能优化。

---

## 1. 为什么 Pipeline 适合 Trace

如果整个流程是一坨代码：

```go
func IndexDocument() {
    // 500 lines
}
```

Trace 很难知道到底慢在哪里。

如果拆成：

```text
parse span
chunk span
embedding span
store span
```

就可以得到：

```text
IndexDocument 3200ms
├─ Parse       1800ms
├─ Chunk        120ms
├─ Embed       1100ms
└─ Store        180ms
```

这也是“设计模式”和“可观测性”发生联系的地方。

---

# 二十三、单例 Singleton：Go 里不要滥用

Singleton 的目标：

> **保证某种对象只有一个实例，并提供统一访问方式。**

例如：

```text
配置
连接池
Logger
Registry
```

但在 Go 里最容易被滥用。

---

## 1. sync.Once

```go
var (
    client *Client
    once   sync.Once
)

func GetClient() *Client {
    once.Do(func() {
        client = NewClient()
    })

    return client
}
```

`sync.Once` 保证初始化函数只执行一次。

---

## 2. 为什么不建议到处全局 Singleton

因为全局对象会带来：

```text
隐藏依赖
难测试
难替换
生命周期模糊
全局状态污染
```

例如：

```go
func CreateUser() {
    db := GetDB()
}
```

函数签名看不出来它依赖数据库。

而依赖注入：

```go
type UserService struct {
    db *gorm.DB
}
```

依赖关系更明确。

所以：

> **Singleton 可以控制实例数量，但不意味着一定要通过全局变量访问。**

---

# 二十四、依赖注入 Dependency Injection：比 Singleton 更重要的工程思想

依赖注入不是 GoF 模式，但非常重要。

核心思想：

> **对象不自己创建依赖，而是从外部接收依赖。**

---

## 不好的方式

```go
func NewUserService() *UserService {
    db := global.GetDB()
    redis := global.GetRedis()

    return &UserService{
        db:    db,
        redis: redis,
    }
}
```

Service 自己偷偷拿依赖。

---

## 更好的方式

```go
type UserService struct {
    db    UserRepository
    cache Cache
}

func NewUserService(
    db UserRepository,
    cache Cache,
) *UserService {
    return &UserService{
        db:    db,
        cache: cache,
    }
}
```

调用者负责组装：

```go
repo := NewMySQLUserRepository(db)
cache := NewRedisCache(redis)

service := NewUserService(repo, cache)
```

优势：

```text
依赖显式
更容易测试
实现可替换
生命周期更清楚
```

---

# 二十五、Repository：隔离业务与数据访问

Repository 也不属于 GoF，但 Go Web 项目很常见。

核心思想：

> **业务层依赖“数据访问能力”，而不是依赖具体 SQL / ORM。**

---

## 示例

```go
type UserRepository interface {
    FindByID(ctx context.Context, id int64) (*User, error)
    Save(ctx context.Context, user *User) error
}
```

MySQL：

```go
type MySQLUserRepository struct {
    db *gorm.DB
}

func (r *MySQLUserRepository) FindByID(
    ctx context.Context,
    id int64,
) (*User, error) {
    var user User

    err := r.db.WithContext(ctx).
        First(&user, id).
        Error

    if err != nil {
        return nil, err
    }

    return &user, nil
}
```

Service：

```go
type UserService struct {
    repo UserRepository
}
```

这样 Service 不需要关心：

```text
GORM
SQL
MySQL
PostgreSQL
```

但也要注意：

> Repository 并不是越抽象越好。

如果只是非常简单 CRUD，过度抽象会增加无意义样板代码。

---

# 二十六、Null Object：用一个“什么都不做”的实现消灭 nil 判断

这是一个很实用但经常被忽视的模式。

假设 Trace 可选。

最初：

```go
if tracer != nil {
    tracer.Start(...)
}
```

所有地方都要判 nil。

可以做：

```go
type Tracer interface {
    Start(ctx context.Context, name string) (context.Context, Span)
}

type NoopTracer struct{}

func (NoopTracer) Start(
    ctx context.Context,
    name string,
) (context.Context, Span) {
    return ctx, NoopSpan{}
}
```

然后无论有没有启用 Trace：

```go
ctx, span := tracer.Start(ctx, "retrieve")
defer span.End()
```

不用：

```go
if tracer != nil
```

这可以让核心代码保持统一。

---

# 二十七、Circuit Breaker：不是 GoF，但真实后端非常重要

调用第三方服务时：

```text
LLM
Embedding
Parser
Rerank
支付
短信
```

如果下游已经故障，你还继续疯狂请求：

```text
请求
↓
超时
↓
请求堆积
↓
Goroutine 增多
↓
连接池耗尽
↓
上游也被拖死
```

熔断器思想：

```text
Closed
↓ 失败率达到阈值
Open
↓ 一段时间
Half-Open
↓ 探测成功
Closed
```

它本质上也是一种状态机。

---

# 二十八、重试 Retry：不能“无脑 retry”

Retry 通常和：

```text
Strategy
Decorator
Middleware
```

组合使用。

例如：

```go
type RetryLLM struct {
    next       LLM
    maxRetries int
}

func (r *RetryLLM) Chat(
    ctx context.Context,
    req ChatRequest,
) (*ChatResponse, error) {
    var lastErr error

    for i := 0; i <= r.maxRetries; i++ {
        resp, err := r.next.Chat(ctx, req)
        if err == nil {
            return resp, nil
        }

        if !isRetryable(err) {
            return nil, err
        }

        lastErr = err

        select {
        case <-ctx.Done():
            return nil, ctx.Err()
        case <-time.After(backoff(i)):
        }
    }

    return nil, lastErr
}
```

真正重要的是：

> Retry 必须考虑幂等性。

例如：

```text
GET 查询
```

一般更容易重试。

但：

```text
创建订单
扣款
发送消息
```

如果没有幂等机制，无脑重试可能造成重复副作用。

---

# 二十九、Fallback：降级模式

在 RAG 里经常遇到：

```text
Rerank Service
↓
超时 / 失败
```

此时可以：

```text
跳过 Rerank
↓
直接返回 RRF 融合结果
```

这是典型降级。

但要注意：

> **降级不等于吞错。**

正确做法：

```text
用户请求仍然成功
+
系统记录 rerank_failed
+
Trace 标记降级
+
Metrics 统计失败率
```

否则你虽然“看起来稳定”，实际上系统长期处于退化状态却没人知道。

---

# 三十、Facade：给复杂子系统提供一个简单入口

假设 Agent Runtime 内部涉及：

```text
Prompt
Memory
RAG
Tool
LLM
Trace
Streaming
```

上层 API 不应该知道所有内部细节。

可以提供：

```go
type AgentFacade struct {
    runtime *Runtime
}

func (f *AgentFacade) Chat(
    ctx context.Context,
    req ChatRequest,
) (<-chan Event, error) {
    return f.runtime.Run(ctx, req)
}
```

API 层只负责：

```text
解析请求
鉴权
调用 AgentFacade
输出 SSE
```

内部复杂性都藏在 Facade 后面。

---

# 三十一、Facade vs Adapter

两者都“包一层”，但意图不同。

```text
Facade：
把复杂系统变简单

Adapter：
把不兼容接口变兼容
```

例如：

```text
AgentFacade
→ 把复杂 Agent Runtime 暴露成 Chat()

MinerUAdapter
→ 把第三方解析接口适配成 Parser
```

---

# 三十二、组合模式 Composite：Agent / Workflow 中很自然

Composite 适合：

> **把单个对象和组合对象统一对待。**

例如工作流节点：

```go
type Node interface {
    Run(ctx context.Context, input Input) (Output, error)
}
```

普通节点：

```text
LLMNode
ToolNode
RetrieverNode
```

组合节点：

```text
SequenceNode
ParallelNode
SubWorkflowNode
```

它们都实现 Node。

于是 Runtime 不需要关心：

```text
这是一个单节点
还是一个子图
```

都可以：

```go
output, err := node.Run(ctx, input)
```

这在 Agent Graph、Workflow Engine 中非常常见。

---

# 三十三、策略 + 工厂 + Registry 为什么经常一起出现

在一个插件化系统中，最常见的组合是：

```text
Strategy
定义统一能力

Factory
负责复杂构建

Registry
负责运行时发现
```

例如 LLM Provider：

```text
LLM interface
    ↓
OpenAI / Qwen / DeepSeek
        ↑
Provider Factory
        ↑
Provider Registry
```

可能的实现：

```go
type ProviderFactory func(cfg Config) (LLM, error)

type LLMRegistry struct {
    factories map[string]ProviderFactory
}
```

注册：

```go
registry.Register("openai", NewOpenAI)
registry.Register("qwen", NewQwen)
```

创建：

```go
factory, ok := registry.Get(provider)
if !ok {
    return nil, errors.New("unknown provider")
}

llm, err := factory(cfg)
```

这种设计非常适合：

```text
插件系统
Agent Tool
Parser
LLM Provider
Embedding Provider
Storage Driver
```

---

# 三十四、一个完整 Agent Tool Runtime 可以组合多少模式

下面看一个比较真实的结构：

```text
Agent Runtime
   ↓
Tool Registry
   ↓
根据 tool name 找 Tool
   ↓
Tool Adapter
   ↓
Middleware
   ├─ Auth
   ├─ Trace
   ├─ Timeout
   └─ Metrics
   ↓
Tool.Execute
```

对应模式：

```text
Tool interface
→ Strategy

Tool Registry
→ Registry

MCP / HTTP / Local Tool
→ Adapter

Auth / Trace / Timeout
→ Decorator / Middleware

Tool Call
→ Command

Tool Event
→ Observer

Tool 状态
→ State
```

这说明：

> 设计模式不是互斥的。

一个真实系统往往是多个模式组合使用。

---

# 三十五、RAG Pipeline 中的设计模式

再看 RAG：

```text
文件
 ↓
Parser
 ↓
Markdown
 ↓
Chunker
 ↓
Embedding
 ↓
Vector Store
```

这里可以识别：

```text
Parser interface
→ Strategy / Adapter

不同 Parser Provider
→ Factory / Registry

整个索引流程
→ Pipeline / Template Method

Embedding API Wrapper
→ Adapter

Trace / Metrics
→ Decorator

解析状态
→ State Machine

异步事件
→ Observer / MQ
```

检索阶段：

```text
Vector Search ──┐
                ├→ RRF → Rerank → LLM
Keyword Search ─┘
```

对应：

```text
Retriever
→ Strategy

RetrieverFactory
→ Factory

RRF
→ Fusion Strategy

Reranker
→ Strategy / Adapter

Rerank fallback
→ Resilience Pattern
```

---

# 三十六、设计模式不是越多越好

一个很常见的错误是：

> “我学了设计模式，所以我要把项目全部改成模式。”

这是错误方向。

真正的判断应该是：

```text
有没有真实变化点？
有没有重复结构？
有没有难以隔离的第三方依赖？
有没有跨模块横切逻辑？
有没有状态爆炸？
有没有新增功能导致大量修改？
```

如果没有，不要强行套模式。

---

# 三十七、如何判断是否过度设计

可以问自己四个问题。

## 1. 当前是否真的有第二个实现？

如果只有：

```text
OneRetriever
```

却提前设计：

```text
Retriever
RetrieverFactory
RetrieverRegistry
RetrieverBuilder
RetrieverManager
RetrieverProvider
```

很可能过度设计。

---

## 2. 变化是不是已经出现，而不是“也许十年后会出现”

合理抽象通常来自：

```text
已有两个实现
已有重复逻辑
已有变化压力
```

而不是：

```text
我猜未来可能需要 20 种实现
```

---

## 3. 抽象有没有降低认知成本

好的抽象：

```text
看到 Retriever
就知道“这是检索能力”
```

坏的抽象：

```text
RetrieverProviderFactoryManager
```

看完更不知道代码在干什么。

---

## 4. 新增功能是否真的变简单

如果新增一个 Tool 仍然要修改：

```text
5 个 switch
3 个配置文件
4 个 Service
```

那所谓 Registry 可能并没有真正解决问题。

---

# 三十八、面试中不要这样回答设计模式

非常低质量的回答：

> 工厂模式就是创建对象，策略模式就是封装算法，适配器模式就是接口转换，我项目里用了很多设计模式。

面试官会马上追：

> 那你为什么要用？

如果说不出来，就会显得像背八股。

---

# 三十九、面试官真正想听什么

高质量回答至少应该包含：

```text
1. 当时的问题是什么
2. 什么地方会变化
3. 如果不抽象会出现什么问题
4. 为什么选这个模式
5. 怎么落地
6. 有什么 Trade-off
```

例如策略模式：

> 我在 RAG 检索层没有直接把 Vector、Keyword、Hybrid 写成一个大 switch，而是抽成统一 Retriever 接口。原因不是为了套设计模式，而是这些检索方式本身就是可替换策略，上层 Query Service 只关心 Retrieve 能力。这样以后增加新的检索策略时，对业务层影响会更小。不过如果策略数量很少，我不会为了模式强行做多层抽象。

这就比：

> “策略模式封装一组算法。”

强很多。

---

# 四十、几个高频模式的面试速记

## Factory

```text
核心问题：
对象创建过程复杂或具体类型经常变化

一句话：
调用者只表达“我要什么”，工厂负责“怎么创建”
```

---

## Strategy

```text
核心问题：
同一个目标有多种可替换实现

一句话：
把变化的算法收敛在统一接口后面
```

---

## Adapter

```text
核心问题：
外部系统接口和内部约定不一致

一句话：
隔离第三方差异，避免外部 DTO/SDK 污染业务层
```

---

## Decorator

```text
核心问题：
给对象动态叠加横切能力

一句话：
不修改核心实现，通过包装增加 Trace、Log、Timeout 等能力
```

---

## Chain

```text
核心问题：
请求需要经过多个可中断处理步骤

一句话：
每个节点只处理自己职责，并决定是否继续向后
```

---

## Observer

```text
核心问题：
一个事件发生后有多个下游动作

一句话：
事件生产者不直接依赖所有消费者
```

---

## State

```text
核心问题：
对象行为高度依赖当前状态

一句话：
把合法状态迁移显式化，避免大量 scattered if/switch
```

---

## Builder / Functional Options

```text
核心问题：
复杂对象有很多可选配置

一句话：
提升构造可读性，并降低后续参数扩展成本
```

---

# 四十一、面试官可能继续追问什么

如果你说项目用了 Strategy，可能追问：

```text
为什么不用 switch？
什么时候 switch 更好？
策略实例什么时候创建？
是否支持运行时切换？
策略有状态吗？
并发安全吗？
```

如果说 Factory：

```text
Factory 和 Builder 区别？
Factory 和 Registry 区别？
新增 Provider 是否还需要修改 Factory？
```

如果说 Adapter：

```text
Adapter 和 Proxy 区别？
Adapter 和 Facade 区别？
为什么不直接使用第三方 SDK？
```

如果说 Middleware：

```text
为什么倒序包装？
怎么保证 Context 传递？
一个中间件 panic 怎么办？
顺序错误有什么问题？
```

如果说 State：

```text
状态迁移怎么保证并发安全？
数据库怎么防止重复消费？
CAS 失败怎么办？
失败是否允许重试？
```

这些才是设计模式在真实面试中的深入方向。

---

# 四十二、作为应届生应该掌握到什么深度

## 必须掌握

```text
Factory
Strategy
Adapter
Decorator
Chain
Observer
Builder / Functional Options
State
```

并且一定要能：

```text
举一个真实工程例子
解释为什么需要
说出不用时的问题
说出一个 Trade-off
```

---

## 加分项

```text
Registry
Composite
Command
Facade
Null Object
Dependency Injection
Repository
Pipeline
```

这些非常适合 Go 后端 / Agent 工程。

---

## 知道思想即可

很多极冷门模式没有必要为了校招无限深挖。

例如：

```text
Flyweight
Memento
Visitor
Prototype
Interpreter
Mediator
```

如果你的项目没有真实使用场景，知道定义和典型用途即可。

不要为了“我会 23 种设计模式”去硬背。

---

# 四十三、一个完整的面试回答示范

假设面试官问：

> 你项目里用过什么设计模式？

可以这样回答：

我项目里没有为了套模式专门设计，但随着 Agent 和 RAG 模块变复杂以后，确实自然用了几种设计思想。

一个比较典型的是策略模式。RAG 检索层有向量检索、关键词检索和混合检索，它们目标相同，但实现不同，所以我会抽象成统一 Retriever 接口，上层 Query Service 只依赖 Retrieve 能力，而不关心底层到底是 pgvector、关键词还是 RRF 融合。这样后续如果增加新的检索策略，修改面会更小。

另外 Tool 和模型 Provider 这种插件化能力，我更倾向于用 Registry 配合 Factory。Registry 负责按名字发现实现，Factory 负责创建实例，避免所有类型都集中在一个很大的 switch 里。

对于第三方 Parser、Embedding、Rerank 这类服务，我会用 Adapter 做隔离，因为我不希望第三方 SDK 的请求结构、错误码和字段直接渗透到业务层。以后替换 Provider 时，主要改适配层。

Trace、Timeout、Metrics 这种横切逻辑则更适合 Middleware 或 Decorator，这样核心业务只做自己的事情。

我觉得设计模式对我最大的价值不是“用了多少种”，而是帮助我识别系统里的变化点，再决定把变化隔离在哪一层。如果系统非常简单，我反而不会为了模式做额外抽象。

---

# 四十四、最后总结：真正应该记住的是“设计思想”

设计模式不应该背成：

```text
Factory
Strategy
Adapter
Builder
Observer
...
```

更应该建立下面这些思维：

```text
行为会变化
→ Strategy

创建过程会变化
→ Factory

插件种类不断增加
→ Registry

第三方接口不稳定
→ Adapter

横切逻辑很多
→ Decorator / Middleware

请求需要逐层处理
→ Chain

一个事件有多个响应者
→ Observer

状态迁移复杂
→ State Machine

构造参数太复杂
→ Builder / Functional Options

复杂子系统对外暴露过多
→ Facade

依赖需要可替换、可测试
→ Dependency Injection
```

最终你会发现：

> **设计模式本质上不是“代码模板”，而是一套控制复杂度的经验。**

好的设计不会让你觉得：

> “这里用了一个 Strategy Pattern。”

而会让你觉得：

> “这个变化被很好地隔离了。”

这才是设计模式真正值得学习的地方。
