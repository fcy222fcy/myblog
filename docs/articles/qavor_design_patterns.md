# 从 Qavor 项目理解设计模式：Go 后端、RAG 与 Agent 系统中的工程化设计

> 这篇文章不打算把设计模式写成一份“23 种模式背诵表”，而是尝试回答一个更工程化的问题：**当一个 Go 后端项目逐渐演进成包含 RAG、Agent、Tool、Trace、异步任务和第三方模型服务的复杂系统时，设计模式到底是如何自然出现的？**

对于校招面试来说，真正有价值的也不是说出“我项目用了策略模式、工厂模式、观察者模式”，而是能进一步解释：

- 当时具体遇到了什么变化点；
- 如果不用这个模式，代码会怎样恶化；
- 为什么选择这种抽象，而不是另一种；
- 模式解决了什么耦合问题；
- 它带来了哪些额外复杂度；
- 哪些地方只是“模式思想”，不应该强行套 GoF 名称。

本文以 Qavor 这类 Go + RAG + Agent 平台为背景，分三层理解设计模式：

```text
第一层：代码级设计模式
Strategy / Adapter / Factory / Decorator / Observer / Command ...

第二层：后端工程模式
Registry / Middleware / Pipeline / Repository / Dependency Injection /
State Machine / Retry & Fallback / Facade ...

第三层：Agent / RAG 架构模式
Hybrid Retrieval / ReAct / Workflow / Graph / SubAgent / Tool Use /
Agent Trace ...
```

需要提前说明：真实项目不会严格按照教科书把每一段代码标成“这里是 XX 模式”。成熟工程往往是多个模式思想组合使用。例如 Go Web 中间件既有 Decorator 的包装思想，也有 Chain of Responsibility 的调用链思想；Tool Registry 又经常与 Factory、Dependency Injection 一起出现。

---

# 1. 为什么项目复杂之后需要设计模式

假设最开始只是做一个简单的知识库问答：

```text
用户问题
  ↓
向量检索
  ↓
TopK 文档
  ↓
LLM
  ↓
回答
```

代码甚至可以全部写在一个 Service 中：

```go
func (s *Service) Ask(ctx context.Context, query string) (string, error) {
    embedding, err := s.openAI.Embed(ctx, query)
    if err != nil {
        return "", err
    }

    chunks, err := s.pgVector.Search(ctx, embedding, 10)
    if err != nil {
        return "", err
    }

    prompt := buildPrompt(query, chunks)

    return s.llm.Chat(ctx, prompt)
}
```

对于 Demo，这种代码完全合理。

问题在于，随着需求增长，很快会出现：

```text
Embedding：
OpenAI / Qwen / BGE / Ollama

检索：
Vector / Keyword / Hybrid

Chunk：
Fixed / Markdown / FAQ / Semantic

Rerank：
BGE / LLM / Disabled

Tool：
Knowledge Base / HTTP / MCP / Custom Tool

Agent：
ReAct / Workflow / Graph / SubAgent

可观测性：
Logging / Trace / Metrics

安全：
Auth / Tool Approval / Parameter Validation
```

如果仍然通过大量 `if`、`switch` 把所有变化塞在一起，代码最终会变成：

```text
一个需求变化
   ↓
修改几十个 switch
   ↓
业务层知道所有实现细节
   ↓
模块互相依赖
   ↓
越来越难测试
   ↓
越来越不敢改
```

设计模式本质上就是在寻找：

> **系统中哪些东西稳定，哪些东西容易变化，然后尽量把变化隔离在边界后面。**

这也是本文理解所有模式的一条主线。

---

# 2. Strategy：策略模式 —— 把“可替换算法”从主流程中抽离

## 2.1 它解决的是什么问题

策略模式适用于这样一种场景：

> 业务目标相同，但是完成这个目标的算法存在多种实现，而且这些实现未来可能继续增加。

在 RAG 中，Chunk 就非常典型。

无论使用哪一种分块方法，业务目标都是：

```text
Document
   ↓
Chunk
   ↓
[]Chunk
```

但是算法不同：

```text
FixedSizeChunker
MarkdownChunker
FAQChunker
SemanticChunker
```

最初可能会这样写：

```go
func Chunk(doc *Document, strategy string) ([]Chunk, error) {
    switch strategy {
    case "fixed":
        return fixedChunk(doc)
    case "markdown":
        return markdownChunk(doc)
    case "faq":
        return faqChunk(doc)
    default:
        return nil, errors.New("unsupported chunk strategy")
    }
}
```

代码本身并没有错。

真正的问题是：**变化点被集中写死在业务流程里。**

以后加入 `semantic`：

```go
case "semantic":
    return semanticChunk(doc)
```

再加入代码文档：

```go
case "code":
    return codeChunk(doc)
```

再加入 Parent-Child Chunk：

```go
case "parent_child":
    return parentChildChunk(doc)
```

所有调用方都逐渐开始知道：系统到底有哪些 Chunk 类型。

这就是耦合。

## 2.2 用接口表达“稳定能力”

Go 中通常直接通过 interface 做策略抽象：

```go
type Chunker interface {
    Chunk(ctx context.Context, doc *Document) ([]Chunk, error)
}
```

这个接口只表达一件事：

> 任何 Chunker，只要能够把 Document 转换成 `[]Chunk`，就可以参与索引流程。

注意它没有规定具体怎么切。

固定长度策略：

```go
type FixedSizeChunker struct {
    ChunkSize int
    Overlap   int
}

func (c *FixedSizeChunker) Chunk(
    ctx context.Context,
    doc *Document,
) ([]Chunk, error) {
    if c.ChunkSize <= 0 {
        return nil, errors.New("invalid chunk size")
    }

    // 简化示例：真实实现还需要处理 UTF-8、段落边界、overlap 等问题。
    var chunks []Chunk

    for start := 0; start < len(doc.Content); {
        end := start + c.ChunkSize
        if end > len(doc.Content) {
            end = len(doc.Content)
        }

        chunks = append(chunks, Chunk{
            Content: doc.Content[start:end],
        })

        if end == len(doc.Content) {
            break
        }

        start = end - c.Overlap
    }

    return chunks, nil
}
```

Markdown 策略：

```go
type MarkdownChunker struct {
    MaxTokens int
}

func (c *MarkdownChunker) Chunk(
    ctx context.Context,
    doc *Document,
) ([]Chunk, error) {
    sections := splitByMarkdownHeading(doc.Content)

    var chunks []Chunk

    for _, section := range sections {
        if estimateTokens(section.Content) <= c.MaxTokens {
            chunks = append(chunks, Chunk{
                Content: section.Content,
                Metadata: map[string]any{
                    "heading": section.Heading,
                },
            })
            continue
        }

        // 标题结构切完之后如果仍然超长，再做长度兜底。
        subChunks := splitLongSection(section, c.MaxTokens)
        chunks = append(chunks, subChunks...)
    }

    return chunks, nil
}
```

这里最重要的不是代码，而是调用方从此只依赖：

```go
Chunker
```

而不依赖：

```text
FixedSizeChunker
MarkdownChunker
FAQChunker
SemanticChunker
```

## 2.3 在索引 Pipeline 里怎么使用

```go
type IndexService struct {
    chunker Chunker
}

func NewIndexService(chunker Chunker) *IndexService {
    return &IndexService{
        chunker: chunker,
    }
}

func (s *IndexService) Index(
    ctx context.Context,
    doc *Document,
) error {
    chunks, err := s.chunker.Chunk(ctx, doc)
    if err != nil {
        return fmt.Errorf("chunk document: %w", err)
    }

    // 后续 Embedding / Store 省略。
    _ = chunks

    return nil
}
```

这时候 `IndexService` 的思维是：

```text
我需要一个 Chunker。

至于它是：
FixedSize？
Markdown？
Semantic？

不是我的职责。
```

这就是策略模式带来的**职责边界**。

## 2.4 为什么这比 switch 好

不是因为 `switch` 本身很低级。

真正优势有三个。

### 第一，可替换

测试中可以直接传 Fake：

```go
type FakeChunker struct{}

func (f *FakeChunker) Chunk(
    ctx context.Context,
    doc *Document,
) ([]Chunk, error) {
    return []Chunk{
        {Content: "chunk-1"},
        {Content: "chunk-2"},
    }, nil
}
```

这样测试 IndexService 时，不需要真正运行复杂 Chunk 算法。

### 第二，可扩展

增加 `SemanticChunker` 时：

```go
type SemanticChunker struct {
    model EmbeddingModel
}
```

只需要新增实现，不一定需要修改主流程。

### 第三，隔离变化

IndexService 的职责是“编排索引”，而不是“知道每一种 Chunk 算法”。

## 2.5 策略模式的代价

设计模式从来不是免费的。

如果项目永远只有一种 Chunk：

```text
MarkdownChunker
```

却提前设计：

```text
ChunkerFactory
ChunkerRegistry
ChunkerStrategyResolver
ChunkerContext
```

反而属于过度设计。

因此更准确的原则是：

> **当变化已经出现，或者变化概率足够高时，再对变化点进行抽象。**

## 2.6 在 Qavor / RAG 中还有哪些 Strategy

类似设计还可以用于：

```go
type Embedder interface {
    Embed(ctx context.Context, texts []string) ([][]float32, error)
}
```

```text
OpenAIEmbedder
QwenEmbedder
BGEEmbedder
OllamaEmbedder
```

Retriever：

```go
type Retriever interface {
    Retrieve(
        ctx context.Context,
        query string,
        topK int,
    ) ([]ScoredChunk, error)
}
```

```text
VectorRetriever
KeywordRetriever
HybridRetriever
```

Reranker：

```go
type Reranker interface {
    Rerank(
        ctx context.Context,
        query string,
        chunks []ScoredChunk,
    ) ([]ScoredChunk, error)
}
```

因此 RAG 系统中 Strategy 的核心关键词就是：

> **同一个能力，多种算法实现。**

---

# 3. Adapter：适配器模式 —— 隔离第三方框架和外部协议

## 3.1 为什么 Agent 项目特别需要 Adapter

Qavor 这类系统通常不可能所有能力都自己实现。

可能依赖：

```text
Eino
OpenAI-compatible API
MCP
MinerU
PP-Structure
pgvector
第三方 Embedding API
第三方 Rerank API
```

外部依赖越多，一个问题越明显：

> 如果业务代码直接使用第三方类型，第三方 SDK 的接口就会逐渐渗透到整个项目。

例如平台内部希望统一定义 Tool：

```go
type Tool interface {
    Name() string
    Description() string
    Execute(
        ctx context.Context,
        args map[string]any,
    ) (*ToolResult, error)
}
```

但是 Eino Runtime 要求的接口可能完全不同：

```go
InvokableRun(ctx context.Context, input string) (string, error)
```

两边不兼容。

最简单的办法是：

> 所有 Qavor Tool 直接实现 Eino 接口。

这样确实能跑。

但是依赖关系变成：

```text
QueryKBTool ─────→ Eino
HTTPTool ────────→ Eino
MCPTool ─────────→ Eino
CustomTool ──────→ Eino
```

以后如果 Eino API 改了，所有 Tool 都受影响。

## 3.2 Adapter 的目标

我们希望变成：

```text
业务层

QueryKBTool
HTTPTool
MCPTool
CustomTool
    │
    ↓
Qavor Tool Interface
    │
──── 系统边界 ────
    │
    ↓
EinoToolAdapter
    │
    ↓
Eino Runtime
```

Adapter 就像一个“翻译器”。

内部说自己的语言，外部说第三方语言，由 Adapter 负责转换。

## 3.3 一个完整示例

内部 Tool：

```go
type ToolResult struct {
    Content  string
    Metadata map[string]any
}

type Tool interface {
    Name() string
    Description() string
    Execute(
        ctx context.Context,
        args map[string]any,
    ) (*ToolResult, error)
}
```

Adapter：

```go
type EinoToolAdapter struct {
    tool Tool
}

func NewEinoToolAdapter(tool Tool) *EinoToolAdapter {
    return &EinoToolAdapter{
        tool: tool,
    }
}
```

真正转换发生在：

```go
func (a *EinoToolAdapter) InvokableRun(
    ctx context.Context,
    input string,
) (string, error) {
    var args map[string]any

    if err := json.Unmarshal([]byte(input), &args); err != nil {
        return "", fmt.Errorf("decode tool args: %w", err)
    }

    result, err := a.tool.Execute(ctx, args)
    if err != nil {
        return "", err
    }

    return result.Content, nil
}
```

这几步分别在做什么？

```text
Eino input string
        ↓
JSON Unmarshal
        ↓
Qavor args map
        ↓
Tool.Execute()
        ↓
Qavor ToolResult
        ↓
转换成 Eino 需要的 string
```

Adapter 不应该承载核心业务。

它主要负责：

```text
参数格式转换
返回值格式转换
错误转换
协议转换
第三方上下文转换
```

## 3.4 Parser Service 也可以体现 Adapter 思想

假设文档解析服务可能有：

```text
MinerUParser
PPStructureParser
LocalMarkdownParser
```

可以定义统一接口：

```go
type Parser interface {
    Parse(ctx context.Context, file *File) (*Document, error)
}
```

对于 Python FastAPI Parser：

```go
type HTTPParserAdapter struct {
    baseURL string
    client  *http.Client
}
```

```go
func (p *HTTPParserAdapter) Parse(
    ctx context.Context,
    file *File,
) (*Document, error) {
    reqBody := ParseRequest{
        FileURL: file.URL,
    }

    body, err := json.Marshal(reqBody)
    if err != nil {
        return nil, err
    }

    req, err := http.NewRequestWithContext(
        ctx,
        http.MethodPost,
        p.baseURL+"/parse",
        bytes.NewReader(body),
    )
    if err != nil {
        return nil, err
    }

    resp, err := p.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var result ParseResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }

    return &Document{
        Content: result.Markdown,
    }, nil
}
```

对于上层 IndexService 来说，它只看到：

```go
parser.Parse(ctx, file)
```

它不需要知道底层到底是：

```text
HTTP
Python
FastAPI
MinerU
PP-Structure
```

这就是非常典型的“外部系统边界隔离”。

## 3.5 Strategy 和 Adapter 怎么区分

这两个非常容易混淆。

### Strategy

关注：

> 一件事有多种实现方式。

```text
Retriever
├── VectorRetriever
├── KeywordRetriever
└── HybridRetriever
```

### Adapter

关注：

> 两个原本接口不兼容的系统如何连接。

```text
Qavor Tool Interface
        ↓
      Adapter
        ↓
Eino Tool Interface
```

可以记成：

```text
Strategy = 替换算法
Adapter  = 转换接口
```

---

# 4. Factory：工厂模式 —— 把对象创建逻辑从业务流程中拿出来

## 4.1 为什么“创建对象”也值得抽象

假设平台支持多种 Tool：

```text
HTTP Tool
Knowledge Base Tool
MCP Tool
```

最开始可能直接写：

```go
if config.Type == "http" {
    tool := NewHTTPTool(config)
}

if config.Type == "kb" {
    tool := NewQueryKBTool(config)
}
```

问题不是代码长，而是对象创建通常不仅仅是 `new()`。

创建 MCP Tool 可能需要：

```text
解析配置
建立 Client
鉴权
连接 MCP Server
tools/list
构建 schema
注册资源释放逻辑
```

创建 KB Tool 可能需要：

```text
注入 Retriever
绑定 KnowledgeBase ID
绑定权限信息
加载配置
```

如果这些逻辑散落在 Handler、Service、Agent Runtime 中，会非常难管理。

## 4.2 简单工厂示例

```go
func NewTool(
    cfg ToolConfig,
    deps ToolDependencies,
) (Tool, error) {
    switch cfg.Type {
    case "http":
        return NewHTTPTool(
            cfg,
            deps.HTTPClient,
        ), nil

    case "knowledge_base":
        return NewQueryKBTool(
            cfg,
            deps.Retriever,
        ), nil

    case "mcp":
        return NewMCPTool(
            cfg,
            deps.MCPManager,
        )

    default:
        return nil, fmt.Errorf(
            "unsupported tool type: %s",
            cfg.Type,
        )
    }
}
```

这时候业务层只负责：

```go
tool, err := NewTool(cfg, deps)
```

而不是自己知道每一种 Tool 怎么初始化。

## 4.3 Factory 真正隔离的是“构造复杂度”

可以把依赖理解成：

```text
业务逻辑
  ↓
Factory
  ↓
决定创建哪个实现
  ↓
准备依赖
  ↓
初始化对象
```

这能避免业务代码里充斥：

```text
if provider == ...
if toolType == ...
if model == ...
```

## 4.4 Factory 不等于所有对象都要工厂化

像：

```go
user := &User{}
```

这种简单对象完全不需要：

```go
UserFactory.CreateUser()
```

Go 社区整体也更偏向简洁组合，而不是 Java 风格“万物 Factory”。

因此使用 Factory 的判断标准应该是：

> **对象创建是否已经包含明显的分支、复杂初始化或不同实现选择。**

---

# 5. Registry：注册表模式 —— Tool / Plugin 系统为什么离不开它

## 5.1 Factory 创建完之后，Runtime 怎么找到 Tool？

Agent 收到模型返回：

```json
{
  "name": "query_kb",
  "arguments": {
    "query": "Go GMP 调度"
  }
}
```

Runtime 下一步必须解决：

> `query_kb` 对应哪个 Tool 实例？

最粗暴的方法仍然是 switch：

```go
switch call.Name {
case "query_kb":
    return queryKBTool.Execute(ctx, call.Arguments)
case "http_request":
    return httpTool.Execute(ctx, call.Arguments)
case "github":
    return githubTool.Execute(ctx, call.Arguments)
}
```

Tool 一多就会变成中心化巨型分支。

Registry 的思路是：

> 启动阶段把 Tool 注册起来，运行阶段根据名字查找。

## 5.2 Registry 实现

```go
type ToolRegistry struct {
    mu    sync.RWMutex
    tools map[string]Tool
}

func NewToolRegistry() *ToolRegistry {
    return &ToolRegistry{
        tools: make(map[string]Tool),
    }
}
```

注册：

```go
func (r *ToolRegistry) Register(tool Tool) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    name := tool.Name()

    if _, exists := r.tools[name]; exists {
        return fmt.Errorf(
            "tool %q already registered",
            name,
        )
    }

    r.tools[name] = tool
    return nil
}
```

读取：

```go
func (r *ToolRegistry) Get(name string) (Tool, bool) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    tool, ok := r.tools[name]
    return tool, ok
}
```

Runtime：

```go
func (e *Executor) Execute(
    ctx context.Context,
    call ToolCall,
) (*ToolResult, error) {
    tool, ok := e.registry.Get(call.Name)
    if !ok {
        return nil, fmt.Errorf(
            "tool %q not found",
            call.Name,
        )
    }

    return tool.Execute(ctx, call.Arguments)
}
```

现在 Tool 数量增加时：

```text
Agent Executor
      ↓
Registry.Get(name)
      ↓
Tool
```

Executor 不需要改。

## 5.3 为什么 map 不等于“设计模式很简单所以没价值”

Registry 代码确实经常只有一个 map。

真正价值不是数据结构，而是架构含义：

> Runtime 不再直接依赖所有具体 Tool。

原来：

```text
Executor
├── QueryKBTool
├── HTTPTool
├── MCPTool
├── GithubTool
└── ...
```

现在：

```text
Executor
   ↓
Registry
   ↓
Tool Interface
```

依赖从“一对多具体依赖”变成“依赖抽象”。

## 5.4 为什么需要 RWMutex

如果 Registry 只在启动阶段初始化，之后完全只读，可以不用复杂并发控制。

但是 Agent 平台可能支持：

```text
用户动态创建 Tool
MCP Server 动态加载
Tool Enable / Disable
运行时刷新 Tool
```

这时候可能同时发生：

```text
Goroutine A：Get("query_kb")
Goroutine B：Register("github")
Goroutine C：Remove("old_tool")
```

Go 原生 map 并发读写可能直接 panic，因此需要同步机制。

这里体现的是一个很重要的工程原则：

> **模式解决架构问题，但真正落地时仍然要考虑语言和运行时的具体约束。**

## 5.5 Factory 与 Registry 的区别

```text
Factory
解决：对象怎么创建？

Registry
解决：对象创建后怎么管理、怎么找到？
```

典型组合：

```text
Config
  ↓
Factory
  ↓
Tool Instance
  ↓
Registry.Register()
  ↓
Agent Runtime
  ↓
Registry.Get(name)
```

---

# 6. Middleware / Decorator / Chain of Responsibility：如何处理横切逻辑

## 6.1 什么叫“横切逻辑”

假设 `query_kb` Tool 最开始只需要：

```go
func (t *QueryKBTool) Execute(
    ctx context.Context,
    args map[string]any,
) (*ToolResult, error) {
    return t.retriever.Search(ctx, args)
}
```

业务很纯粹：查询知识库。

后来陆续增加：

```text
日志
Trace
耗时统计
权限校验
参数校验
超时控制
Panic Recovery
Tool 审批
敏感参数脱敏
Metrics
限流
```

如果全部写进 Execute：

```go
func (t *QueryKBTool) Execute(...) {
    checkPermission()
    validateArgs()
    startTrace()
    logRequest()
    startTimer()
    createTimeoutContext()
    defer recoverPanic()

    result := search()

    recordMetrics()
    finishTrace()
    logResult()

    return result
}
```

那么 Tool 的业务逻辑会被基础设施代码淹没。

这些逻辑有一个共同特点：

> 它们不是某一个 Tool 独有，而是很多 Tool 都需要。

这就是横切关注点。

## 6.2 定义 Handler 和 Middleware

```go
type ToolRequest struct {
    Name      string
    Arguments map[string]any
    Risk      RiskLevel
}

type ToolHandler func(
    ctx context.Context,
    req *ToolRequest,
) (*ToolResult, error)
```

Middleware 的本质是：

```go
type ToolMiddleware func(next ToolHandler) ToolHandler
```

第一次看到这个类型可能比较绕。

拆开理解：

```text
输入一个 Handler
      ↓
包装它
      ↓
返回一个新的 Handler
```

例如：

```text
原始：
Tool.Execute

包装 Timeout 后：
Timeout(Tool.Execute)

继续包装 Logging：
Logging(Timeout(Tool.Execute))
```

这就是 Decorator 的典型思想：

> 不修改原对象，通过外层包装给它增加能力。

## 6.3 Timeout Middleware

```go
func TimeoutMiddleware(
    timeout time.Duration,
) ToolMiddleware {
    return func(next ToolHandler) ToolHandler {
        return func(
            ctx context.Context,
            req *ToolRequest,
        ) (*ToolResult, error) {
            ctx, cancel := context.WithTimeout(
                ctx,
                timeout,
            )
            defer cancel()

            return next(ctx, req)
        }
    }
}
```

逐层解释。

第一层：

```go
func TimeoutMiddleware(timeout time.Duration) ToolMiddleware
```

表示先配置一个超时时间。

例如：

```go
TimeoutMiddleware(5 * time.Second)
```

得到一个 Middleware。

第二层：

```go
return func(next ToolHandler) ToolHandler
```

表示它需要接收“下一个处理器”。

第三层：

```go
ctx, cancel := context.WithTimeout(ctx, timeout)
```

生成带 Deadline 的子 Context。

然后：

```go
return next(ctx, req)
```

把新的 Context 继续向下传。

因此链路变成：

```text
Agent Request
    ↓
Timeout Middleware
    ↓
带 Deadline 的 ctx
    ↓
Tool
    ↓
Retriever
    ↓
Database / Model API
```

只要下游正确使用同一个 Context，取消信号就能够一路传播。

这也是 Go `context.Context` 在工程里最重要的价值之一。

## 6.4 Logging Middleware

```go
func LoggingMiddleware(
    logger *zap.Logger,
) ToolMiddleware {
    return func(next ToolHandler) ToolHandler {
        return func(
            ctx context.Context,
            req *ToolRequest,
        ) (*ToolResult, error) {
            start := time.Now()

            logger.Info(
                "tool started",
                zap.String("tool", req.Name),
            )

            result, err := next(ctx, req)

            logger.Info(
                "tool finished",
                zap.String("tool", req.Name),
                zap.Duration("duration", time.Since(start)),
                zap.Error(err),
            )

            return result, err
        }
    }
}
```

核心点在：

```go
result, err := next(ctx, req)
```

它把 Middleware 分成两个阶段：

```text
next() 之前
→ 请求前逻辑

next()
→ 真正执行下游

next() 之后
→ 请求后逻辑
```

因此很多能力都可以通过它实现：

```text
请求前：
权限校验
开始 Trace
记录开始时间
参数检查

请求后：
记录耗时
结束 Span
Metrics
错误转换
```

## 6.5 Panic Recovery

对于 Agent Tool，来源可能非常复杂：

```text
平台内置 Tool
第三方 Tool
MCP Tool
用户配置 Tool
```

一个 Tool panic 不应该轻易把整个 HTTP 请求甚至服务进程打崩。

```go
func RecoveryMiddleware() ToolMiddleware {
    return func(next ToolHandler) ToolHandler {
        return func(
            ctx context.Context,
            req *ToolRequest,
        ) (
            result *ToolResult,
            err error,
        ) {
            defer func() {
                if r := recover(); r != nil {
                    err = fmt.Errorf(
                        "tool panic: %v",
                        r,
                    )
                }
            }()

            return next(ctx, req)
        }
    }
}
```

注意 `recover()` 只能在 `defer` 中生效。

这段代码本质是把：

```text
panic
```

转换为：

```text
error
```

让上层 Agent Runtime 决定：

```text
终止 Agent？
把 Tool Error 交回 LLM？
尝试其他 Tool？
执行降级？
```

## 6.6 Middleware Chain

```go
func Chain(
    handler ToolHandler,
    middlewares ...ToolMiddleware,
) ToolHandler {
    for i := len(middlewares) - 1; i >= 0; i-- {
        handler = middlewares[i](handler)
    }

    return handler
}
```

为什么倒序？

假设：

```go
Chain(
    toolHandler,
    Recovery,
    Trace,
    Logging,
    Timeout,
)
```

希望执行顺序是：

```text
Recovery
  ↓
Trace
  ↓
Logging
  ↓
Timeout
  ↓
Tool
```

构造时必须从最内层开始：

```text
handler = Timeout(Tool)
handler = Logging(handler)
handler = Trace(handler)
handler = Recovery(handler)
```

最终得到：

```text
Recovery(
    Trace(
        Logging(
            Timeout(
                Tool
            )
        )
    )
)
```

这也是为什么循环需要倒序。

## 6.7 Decorator 和责任链到底是什么关系

现实中的 Middleware 往往同时体现两种思想。

Decorator 强调：

> 包装已有对象，在不修改原对象的前提下增加能力。

Chain of Responsibility 强调：

> 一个请求沿着多个处理节点依次传递，每个节点可以处理、修改或决定是否继续。

Go Middleware：

```text
Request
  ↓
Recovery
  ↓
Trace
  ↓
Auth
  ↓
Timeout
  ↓
Handler
```

因此既像 Decorator，也像责任链。

面试时不要死抠：

> “Gin Middleware 到底只能算哪一种 GoF 模式？”

更加工程化的回答是：

> Middleware 通常使用函数包装形成调用链，在结构上体现 Decorator 的包装思想，同时请求又逐层经过多个处理节点，因此也具有 Chain of Responsibility 的特点。

---

# 7. Pipeline：Pipeline 模式 —— RAG 索引为什么天然适合流水线

## 7.1 RAG 索引本身就是阶段式处理

文档进入知识库后：

```text
Upload
  ↓
Parse
  ↓
Normalize
  ↓
Chunk
  ↓
Embedding
  ↓
Store
```

每一步的共同特征：

```text
接收上一步输出
      ↓
完成自己的职责
      ↓
产生下一步输入
```

这就是 Pipeline。

## 7.2 最简单的 Pipeline 编排

```go
type IndexPipeline struct {
    parser   Parser
    chunker  Chunker
    embedder Embedder
    store    VectorStore
}
```

```go
func (p *IndexPipeline) Run(
    ctx context.Context,
    file *File,
) error {
    doc, err := p.parser.Parse(ctx, file)
    if err != nil {
        return fmt.Errorf("parse document: %w", err)
    }

    chunks, err := p.chunker.Chunk(ctx, doc)
    if err != nil {
        return fmt.Errorf("chunk document: %w", err)
    }

    texts := make([]string, 0, len(chunks))
    for _, chunk := range chunks {
        texts = append(texts, chunk.Content)
    }

    vectors, err := p.embedder.Embed(ctx, texts)
    if err != nil {
        return fmt.Errorf("embed chunks: %w", err)
    }

    for i := range chunks {
        chunks[i].Vector = vectors[i]
    }

    if err := p.store.Save(ctx, chunks); err != nil {
        return fmt.Errorf("save chunks: %w", err)
    }

    return nil
}
```

这段代码虽然简单，但已经体现了很明确的架构：

```text
Pipeline 负责流程

Parser 负责解析
Chunker 负责分块
Embedder 负责向量化
VectorStore 负责持久化
```

每一层职责都清晰。

## 7.3 Pipeline 和 Strategy 是怎么组合的

这两个模式经常一起出现。

Pipeline 决定：

```text
流程顺序是什么？
```

Strategy 决定：

```text
某一个阶段具体怎么实现？
```

例如：

```text
Pipeline

Parse
 ↓
Chunk  ← Strategy：Markdown / Fixed / Semantic
 ↓
Embed  ← Strategy：OpenAI / Qwen / BGE
 ↓
Store  ← Strategy：pgvector / Milvus / Qdrant
```

因此设计模式并不是互斥的。

成熟系统往往是：

> **大流程用 Pipeline，阶段内部的变化点用 Strategy。**

## 7.4 为什么统一 Markdown 是 Pipeline 边界设计

如果 PDF、DOCX、HTML、Markdown 各走各的：

```text
PDF → PDFChunker
DOCX → DOCXChunker
HTML → HTMLChunker
MD → MarkdownChunker
```

那么文件格式差异会一直传递到下游。

如果 Parser 统一输出：

```text
Document {
    Content: Markdown
    Metadata: ...
}
```

就变成：

```text
PDF ─────┐
DOCX ────┤
HTML ────┼→ Parser → Unified Document → Chunk → Embed → Store
Markdown ┘
```

这实际上是在 Pipeline 中人为建立一个稳定的“中间表示（Intermediate Representation）”。

它的工程价值是：

> **上游负责解决格式差异，下游只处理统一结构。**

这与编译器中的 IR 思想非常相似。

---

# 8. Observer / Pub-Sub：异步任务和 SSE 为什么天然存在“观察者”思想

## 8.1 一个很典型的问题

文档解析、OCR、Embedding 都可能耗时较长。

如果 HTTP 请求直接同步等待：

```text
Client
  ↓
POST /upload
  ↓
Parse 20s
  ↓
Chunk
  ↓
Embedding
  ↓
Store
  ↓
Response
```

请求会一直挂着。

更合理的设计通常是：

```text
Client
  ↓
创建任务
  ↓
返回 task_id
  ↓
Worker 后台处理
```

但用户还希望知道：

```text
正在解析
解析 50%
正在分块
正在 Embedding
完成
```

因此系统需要一种“状态发生变化后通知订阅者”的机制。

这就非常接近 Observer / Publish-Subscribe 思想。

## 8.2 定义事件

```go
type TaskEvent struct {
    TaskID   string
    Stage    string
    Progress int
    Message  string
}
```

例如：

```go
TaskEvent{
    TaskID:   "task-123",
    Stage:    "parsing",
    Progress: 30,
    Message:  "正在解析 PDF",
}
```

## 8.3 EventBus

```go
type EventBus struct {
    mu          sync.RWMutex
    subscribers map[string][]chan TaskEvent
}
```

订阅：

```go
func (b *EventBus) Subscribe(taskID string) <-chan TaskEvent {
    ch := make(chan TaskEvent, 16)

    b.mu.Lock()
    defer b.mu.Unlock()

    b.subscribers[taskID] = append(
        b.subscribers[taskID],
        ch,
    )

    return ch
}
```

发布：

```go
func (b *EventBus) Publish(event TaskEvent) {
    b.mu.RLock()
    defer b.mu.RUnlock()

    for _, ch := range b.subscribers[event.TaskID] {
        select {
        case ch <- event:
        default:
            // 简化处理：消费者太慢时避免阻塞生产者。
        }
    }
}
```

Worker 只负责发布：

```go
bus.Publish(TaskEvent{
    TaskID:   task.ID,
    Stage:    "embedding",
    Progress: 70,
})
```

它不需要知道：

```text
前端是不是 SSE
有没有日志消费者
有没有 Metrics 消费者
有没有 WebSocket 消费者
```

## 8.4 SSE 只是事件的一种出口

SSE Handler：

```go
func StreamTaskEvents(
    w http.ResponseWriter,
    r *http.Request,
    bus *EventBus,
    taskID string,
) {
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")

    flusher, ok := w.(http.Flusher)
    if !ok {
        http.Error(w, "stream unsupported", http.StatusInternalServerError)
        return
    }

    events := bus.Subscribe(taskID)

    for {
        select {
        case <-r.Context().Done():
            return

        case event := <-events:
            data, _ := json.Marshal(event)

            fmt.Fprintf(w, "data: %s\n\n", data)
            flusher.Flush()
        }
    }
}
```

链路：

```text
Worker
  ↓
Publish TaskEvent
  ↓
EventBus
  ↓
Subscriber
  ↓
SSE Handler
  ↓
Browser
```

这时候 Worker 与 SSE 被解耦了。

## 8.5 Observer 和 MQ 有什么关系

Observer / Pub-Sub 是一种思想。

实现可以是：

```text
进程内 channel
Redis Pub/Sub
Redis Streams
Kafka
NATS
RabbitMQ
```

如果是单实例原型：

```text
Go channel EventBus
```

可能够用。

如果是多实例部署：

```text
API Instance A
Worker Instance B
SSE Instance C
```

进程内 channel 就不够了，因为它们不共享内存。

这时候需要外部消息系统。

所以面试时一定要体现：

> 模式是一种设计思想，但实现方案要根据部署模型和可靠性要求选择。

---

# 9. State Machine：状态机 —— 异步任务不能只靠几个 bool

## 9.1 为什么任务状态容易失控

一个文档索引任务可能经历：

```text
PENDING
  ↓
PARSING
  ↓
CHUNKING
  ↓
EMBEDDING
  ↓
INDEXING
  ↓
SUCCEEDED
```

中途也可能：

```text
FAILED
CANCELED
ROLLING_BACK
```

如果只定义：

```go
IsParsed bool
IsIndexed bool
IsFailed bool
```

马上会产生非法组合：

```text
IsParsed = false
IsIndexed = true
```

到底是什么意思？

或者：

```text
IsFailed = true
IsIndexed = true
```

是失败了还是成功了？

状态机的价值就是：

> 把“系统现在处于什么状态”和“允许从哪里跳到哪里”显式建模。

## 9.2 定义状态

```go
type TaskStatus string

const (
    StatusPending     TaskStatus = "pending"
    StatusParsing     TaskStatus = "parsing"
    StatusChunking    TaskStatus = "chunking"
    StatusEmbedding   TaskStatus = "embedding"
    StatusIndexing    TaskStatus = "indexing"
    StatusSucceeded   TaskStatus = "succeeded"
    StatusFailed      TaskStatus = "failed"
    StatusCanceled    TaskStatus = "canceled"
)
```

## 9.3 显式定义合法迁移

```go
var transitions = map[TaskStatus]map[TaskStatus]bool{
    StatusPending: {
        StatusParsing:  true,
        StatusCanceled: true,
    },

    StatusParsing: {
        StatusChunking: true,
        StatusFailed:   true,
        StatusCanceled: true,
    },

    StatusChunking: {
        StatusEmbedding: true,
        StatusFailed:    true,
        StatusCanceled:  true,
    },

    StatusEmbedding: {
        StatusIndexing: true,
        StatusFailed:   true,
        StatusCanceled: true,
    },

    StatusIndexing: {
        StatusSucceeded: true,
        StatusFailed:    true,
    },
}
```

更新时验证：

```go
func CanTransition(from, to TaskStatus) bool {
    next, ok := transitions[from]
    if !ok {
        return false
    }

    return next[to]
}
```

```go
func (s *TaskService) Transition(
    ctx context.Context,
    taskID string,
    to TaskStatus,
) error {
    task, err := s.repo.Get(ctx, taskID)
    if err != nil {
        return err
    }

    if !CanTransition(task.Status, to) {
        return fmt.Errorf(
            "invalid task transition: %s -> %s",
            task.Status,
            to,
        )
    }

    return s.repo.UpdateStatus(ctx, taskID, to)
}
```

## 9.4 为什么状态机对重试和断点续跑很重要

Worker 崩溃后重启：

```text
Task Status = EMBEDDING
```

系统就可以判断：

> 解析和 Chunk 已经完成，是否可以从 Embedding 阶段恢复？

如果只是一个：

```go
Progress = 65
```

恢复语义就很弱。

所以状态不仅用于 UI 展示，也可以直接参与：

```text
幂等
失败恢复
断点续跑
补偿
任务重试
```

## 9.5 数据库并发更新还要考虑 CAS

假设两个 Worker 意外消费同一个任务：

```text
Worker A：PENDING → PARSING
Worker B：PENDING → PARSING
```

仅在内存验证状态不够。

更稳妥的 SQL 类似：

```sql
UPDATE tasks
SET status = 'parsing'
WHERE id = $1
  AND status = 'pending';
```

检查 affected rows：

```text
1 → 抢占成功
0 → 状态已被其他 Worker 修改
```

这相当于数据库层面的 Compare-And-Set 思想。

这也是一个非常好的面试深入点：

> 状态机负责定义合法状态变化，但并发安全仍然需要数据库原子更新、事务或分布式锁等机制保证。

---

# 10. Command：命令模式 —— Tool Call 本质上就是结构化命令

## 10.1 LLM 为什么不能直接执行 Tool

模型返回 Function Calling：

```json
{
  "name": "query_kb",
  "arguments": {
    "query": "Redis 为什么快"
  }
}
```

这不是执行结果。

它只是一个：

> “我希望执行 query_kb，并且参数是这些。”

换句话说，模型生成的是一个**命令描述**。

真正执行发生在 Agent Runtime。

## 10.2 抽象 Command

```go
type ToolCallCommand struct {
    Name      string
    Arguments map[string]any
    CallID    string
}
```

Executor：

```go
type CommandExecutor struct {
    registry *ToolRegistry
}

func (e *CommandExecutor) Execute(
    ctx context.Context,
    cmd ToolCallCommand,
) (*ToolResult, error) {
    tool, ok := e.registry.Get(cmd.Name)
    if !ok {
        return nil, ErrToolNotFound
    }

    return tool.Execute(ctx, cmd.Arguments)
}
```

现在 Tool Call 就可以：

```text
生成
↓
记录
↓
排队
↓
审批
↓
执行
↓
重试
↓
审计
```

因为“请求执行什么”已经被封装成了独立数据对象。

## 10.3 Command 在 Agent 系统中的实际价值

比如高风险工具需要人工审批。

如果模型一生成 Tool Call 就立即执行：

```text
LLM
 ↓
delete_file
 ↓
直接删
```

很危险。

如果先转换成 Command：

```text
LLM
 ↓
ToolCallCommand
 ↓
Risk Check
 ↓
Approval
 ↓
Executor
```

系统就有了一个清晰的控制边界。

## 10.4 还可以做审计日志

```go
type CommandLog struct {
    CallID     string
    ToolName   string
    Arguments  json.RawMessage
    Status     string
    StartedAt  time.Time
    FinishedAt *time.Time
}
```

以后 Trace 页面就可以展示：

```text
Agent Step #3
Tool: query_kb
Args: {query: "Redis 为什么快"}
Duration: 82ms
Status: success
```

因此 Command 模式与 Agent Trace 也非常容易组合。

---

# 11. Facade：外观模式 —— 不要让 Controller 知道整个 RAG 内部结构

## 11.1 一个典型的错误依赖

如果 API Handler 直接完成：

```text
查询向量
查询关键词
执行 RRF
调用 Rerank
拼 Prompt
调用 LLM
```

Controller 会变得非常重。

```go
func AskHandler(c *gin.Context) {
    // 解析参数
    // vector search
    // keyword search
    // RRF
    // rerank
    // prompt
    // LLM
    // SSE
}
```

这样 HTTP 层知道太多内部实现细节。

Facade 的思路是给复杂子系统提供一个更简单的入口。

例如：

```go
type RAGService struct {
    vectorRetriever  Retriever
    keywordRetriever Retriever
    reranker         Reranker
    llm              ChatModel
}
```

对外只暴露：

```go
func (s *RAGService) Ask(
    ctx context.Context,
    query string,
) (*Answer, error)
```

内部自己完成：

```text
Vector Search
Keyword Search
RRF
Rerank
Prompt
LLM
```

Controller 只负责：

```go
answer, err := ragService.Ask(c.Request.Context(), req.Query)
```

## 11.2 Facade 的核心价值

它不是为了把所有代码塞进一个大 Service。

而是：

> 对外暴露稳定、简单的使用入口，把内部多个组件的协作关系隐藏起来。

例如：

```text
Controller
    ↓
RAG Facade
    ↓
┌───────────┬──────────────┬─────────┐
Retriever   Fusion         Reranker  LLM
```

Controller 不需要知道这些组件的组合方式。

## 11.3 Facade 和 Service Layer 很相似

现代后端项目中，Facade 经常不会真的命名成：

```text
RAGFacade
```

而是：

```text
RAGService
AgentService
KnowledgeBaseService
```

所以不要为了套模式强行改名。

面试时可以说：

> Service 层在架构上承担了类似 Facade 的职责，对 API 层提供较稳定的业务入口，把底层 Retriever、Reranker、Model 等组件的组合细节隐藏在内部。

比说“我的 Service 就是标准 GoF Facade”更稳妥。

---

# 12. Repository：把数据访问和业务逻辑分离

Repository 并不是 GoF 23 种之一，但后端面试非常重要。

## 12.1 为什么 Service 不应该到处写 GORM

例如：

```go
func (s *TaskService) GetTask(
    ctx context.Context,
    id uint,
) (*Task, error) {
    var task Task

    err := s.db.WithContext(ctx).
        Where("id = ?", id).
        First(&task).
        Error

    return &task, err
}
```

小项目完全能用。

但随着复杂度增加：

```text
事务
锁
复杂查询
分页
软删除
缓存
数据库切换
测试 Mock
```

业务层会越来越依赖 ORM 细节。

## 12.2 Repository 接口

```go
type TaskRepository interface {
    GetByID(
        ctx context.Context,
        id uint,
    ) (*Task, error)

    UpdateStatus(
        ctx context.Context,
        id uint,
        status TaskStatus,
    ) error
}
```

GORM 实现：

```go
type GormTaskRepository struct {
    db *gorm.DB
}

func (r *GormTaskRepository) GetByID(
    ctx context.Context,
    id uint,
) (*Task, error) {
    var task Task

    if err := r.db.WithContext(ctx).
        First(&task, id).
        Error; err != nil {
        return nil, err
    }

    return &task, nil
}
```

Service：

```go
type TaskService struct {
    repo TaskRepository
}
```

业务层看到的是：

```text
GetByID
UpdateStatus
```

而不是：

```text
GORM
SQL
表结构细节
```

## 12.3 Repository 不是“为了换数据库”才存在

很多八股回答会说：

> “Repository 让我们以后从 MySQL 换 PostgreSQL 很方便。”

这不是错，但真实项目里数据库并不会天天换。

Repository 更重要的意义是：

```text
业务语义
和
数据持久化细节
分离
```

例如业务层关心：

```go
repo.GetRunnableTasks(ctx)
```

而不是关心：

```sql
SELECT ... WHERE status IN (...) AND retry_count < ...
```

这才是最实际的收益。

---

# 13. Dependency Injection：依赖注入 —— 谁负责组装这些模式？

前面定义了：

```text
Parser
Chunker
Embedder
Retriever
Reranker
Repository
ToolRegistry
```

问题来了：

> 这些对象到底在哪里创建并连接起来？

如果 Service 自己创建依赖：

```go
type RAGService struct {
    retriever *PgVectorRetriever
}

func NewRAGService() *RAGService {
    db := connectDB()
    retriever := NewPgVectorRetriever(db)

    return &RAGService{
        retriever: retriever,
    }
}
```

RAGService 既负责业务，又负责：

```text
数据库连接
具体实现选择
对象生命周期
```

这叫依赖隐藏。

更好的方式：

```go
func NewRAGService(
    retriever Retriever,
    reranker Reranker,
    llm ChatModel,
) *RAGService {
    return &RAGService{
        retriever: retriever,
        reranker:  reranker,
        llm:       llm,
    }
}
```

RAGService 明确声明：

> 我需要这些东西，但是我不负责创建它们。

## 13.1 在 main / wire 层组装

```go
func buildApplication(cfg Config) (*Application, error) {
    db, err := initDB(cfg.Database)
    if err != nil {
        return nil, err
    }

    vectorRepo := NewPgVectorRepository(db)

    retriever := NewHybridRetriever(
        NewVectorRetriever(vectorRepo),
        NewKeywordRetriever(db),
    )

    reranker := NewBGEReranker(cfg.Rerank)
    llm := NewOpenAICompatibleModel(cfg.Model)

    ragService := NewRAGService(
        retriever,
        reranker,
        llm,
    )

    return &Application{
        RAGService: ragService,
    }, nil
}
```

这块叫 Composition Root：

> 对象依赖关系在系统边界统一组装。

## 13.2 为什么 DI 对测试很重要

生产：

```go
NewRAGService(
    realRetriever,
    realReranker,
    realLLM,
)
```

测试：

```go
NewRAGService(
    fakeRetriever,
    fakeReranker,
    fakeLLM,
)
```

业务 Service 不需要改。

所以 Dependency Injection 与 Strategy、Repository、Adapter 等模式的关系非常紧密。

没有 DI，虽然定义了接口，最终还是可能在 Service 内部写死具体实现。

---

# 14. Proxy：代理模式 —— 给模型客户端统一增加缓存、限流、计费或观测

## 14.1 为什么 Model Client 很适合做代理

假设定义：

```go
type ChatModel interface {
    Generate(
        ctx context.Context,
        req ChatRequest,
    ) (*ChatResponse, error)
}
```

真实实现：

```go
type OpenAIModel struct {
    client *http.Client
}
```

现在想增加：

```text
调用统计
Token 计费
限流
缓存
Trace
Fallback
```

可以通过代理层包一层，而不修改原 Model Client。

## 14.2 Metrics Proxy

```go
type MetricsModelProxy struct {
    next    ChatModel
    metrics Metrics
}

func (p *MetricsModelProxy) Generate(
    ctx context.Context,
    req ChatRequest,
) (*ChatResponse, error) {
    start := time.Now()

    resp, err := p.next.Generate(ctx, req)

    p.metrics.ObserveModelCall(
        req.Model,
        time.Since(start),
        err,
    )

    return resp, err
}
```

调用方仍然认为自己拿到的是：

```go
ChatModel
```

实际上是：

```text
RAGService
   ↓
MetricsModelProxy
   ↓
OpenAIModel
```

## 14.3 Proxy 和 Decorator 的区别

这两个结构很像：都包一层。

Decorator 更强调：

> 动态增加功能。

Proxy 更强调：

> 控制对真实对象的访问。

例如：

```text
RateLimitProxy
CacheProxy
PermissionProxy
RemoteProxy
```

现实代码并不一定要严格区分命名。

只要你能解释：

> 为什么包这一层、控制了什么访问行为、隔离了什么复杂度。

就已经足够。

---

# 15. Retry / Timeout / Fallback：工程中的弹性设计模式

这类不属于传统 GoF，但在 Agent / RAG 工程中非常重要。

## 15.1 Rerank 为什么需要 Fallback

完整检索链：

```text
Vector Retrieval
      +
Keyword Retrieval
      ↓
RRF Fusion
      ↓
Rerank API
      ↓
TopK
```

如果 Rerank 服务挂了，是否整个问答都失败？

通常不一定值得。

因为 RRF 已经产生了一份“可用但可能没那么优”的结果。

因此：

```go
ranked, err := s.reranker.Rerank(
    ctx,
    query,
    candidates,
)

if err != nil {
    logger.Warn(
        "rerank failed, fallback to fusion result",
        zap.Error(err),
    )

    ranked = candidates
}
```

这就是 Fallback：

```text
最佳路径失败
      ↓
退化到次优但可用路径
```

## 15.2 Fallback 为什么不是“偷偷吃掉 error”

错误仍然应该：

```text
记录日志
记录 Trace
Metrics +1
必要时告警
```

只是产品层面不一定让用户请求完全失败。

因此正确语义是：

> **错误被观测，但用户路径降级继续执行。**

而不是：

```go
if err != nil {
    // ignore
}
```

## 15.3 Retry 什么时候能用

例如模型 API 短暂返回：

```text
502
503
connection reset
```

可以有限重试。

但不是所有错误都应该重试。

```text
401 Unauthorized
→ 重试没意义

400 Invalid Request
→ 重试没意义

Context Canceled
→ 不应该继续重试

503 Service Unavailable
→ 可以考虑指数退避重试
```

一个简单判断：

```go
func shouldRetry(err error) bool {
    if errors.Is(err, context.Canceled) ||
        errors.Is(err, context.DeadlineExceeded) {
        return false
    }

    var httpErr *HTTPError
    if errors.As(err, &httpErr) {
        return httpErr.StatusCode == 502 ||
            httpErr.StatusCode == 503 ||
            httpErr.StatusCode == 504
    }

    return true
}
```

## 15.4 为什么 Retry 必须有 Timeout

没有总超时：

```text
一次调用 10 秒
失败
重试 10 秒
失败
重试 10 秒
```

用户请求可能挂几十秒。

因此应当有：

```text
Request Deadline
      ↓
每次调用 Timeout
      ↓
有限 Retry
```

它们不是独立功能，而是组合设计。

---

# 16. Hybrid Retrieval：组合模式思想在 RAG 中的体现

混合检索并不一定对应严格 GoF Composite，但非常适合用“组合小组件形成更强能力”的思想理解。

## 16.1 两路检索

向量：

```go
type VectorRetriever struct {
    repo VectorRepository
}
```

关键词：

```go
type KeywordRetriever struct {
    repo KeywordRepository
}
```

混合：

```go
type HybridRetriever struct {
    vector  Retriever
    keyword Retriever
}
```

但它自己仍然实现：

```go
Retriever
```

```go
func (r *HybridRetriever) Retrieve(
    ctx context.Context,
    query string,
    topK int,
) ([]ScoredChunk, error) {
    vectorResult, err := r.vector.Retrieve(
        ctx,
        query,
        topK,
    )
    if err != nil {
        return nil, err
    }

    keywordResult, err := r.keyword.Retrieve(
        ctx,
        query,
        topK,
    )
    if err != nil {
        return nil, err
    }

    return rrfFuse(
        vectorResult,
        keywordResult,
        topK,
    ), nil
}
```

调用方仍然只看到：

```go
Retriever
```

这意味着未来可以继续组合：

```text
HybridRetriever
├── VectorRetriever
└── KeywordRetriever
```

甚至：

```text
MultiRetriever
├── Dense Vector
├── Sparse Vector
├── BM25
└── Graph Retrieval
```

这种设计有一个很重要的特点：

> 组合对象和叶子对象对调用方暴露相同接口。

这就是 Composite 思想最核心的部分。

---

# 17. Agent Runtime：ReAct 本身也是一种架构模式

传统设计模式解决的是代码组织问题。

Agent 里还有一种更高层次的“行为模式”：

> 模型应该以什么控制流程完成任务？

ReAct 就是一种非常典型的 Agent Pattern。

## 17.1 ReAct 核心循环

可以抽象成：

```text
用户问题
   ↓
LLM Reason
   ↓
是否需要 Tool？
   ├── 否 → Final Answer
   │
   └── 是
        ↓
     Tool Call
        ↓
   Tool Execution
        ↓
    Observation
        ↓
    重新交给 LLM
        ↓
       循环
```

伪代码：

```go
for step := 0; step < maxSteps; step++ {
    response, err := llm.Chat(ctx, messages, tools)
    if err != nil {
        return nil, err
    }

    if len(response.ToolCalls) == 0 {
        return &Answer{
            Content: response.Content,
        }, nil
    }

    messages = append(messages, response.Message)

    for _, call := range response.ToolCalls {
        result, err := executor.Execute(ctx, call)

        toolMessage := buildToolMessage(
            call.ID,
            result,
            err,
        )

        messages = append(
            messages,
            toolMessage,
        )
    }
}
```

## 17.2 为什么要 maxSteps

没有限制：

```text
LLM → Tool → LLM → Tool → ...
```

可能无限循环。

因此 Runtime 必须有硬约束：

```go
const maxSteps = 10
```

这体现了一个 Agent 工程的重要思想：

> Prompt 是软约束，Runtime 才是硬约束。

不能只在 Prompt 写：

```text
请不要调用超过十次工具。
```

真正限制必须在程序中。

---

# 18. Workflow / Graph：当 Agent 流程比 ReAct 更确定时

并不是所有 AI 任务都应该让 LLM 自由决定下一步。

例如知识库索引：

```text
Parse
 ↓
Chunk
 ↓
Embedding
 ↓
Store
```

流程非常确定。

这种场景用 Workflow 往往比 ReAct 合适。

## 18.1 ReAct 的优势

```text
动态
灵活
模型可以自行决定 Tool
适合开放问题
```

## 18.2 Workflow 的优势

```text
确定性强
容易测试
成本可控
延迟稳定
错误边界清晰
```

因此系统中可能同时存在：

```text
确定流程
→ Workflow / Pipeline

动态决策
→ ReAct Agent
```

这是非常重要的架构 Trade-off。

不要因为“Agent 很高级”就所有逻辑都让模型决定。

---

# 19. SubAgent：把复杂职责拆给不同 Agent

当一个 Agent 同时需要：

```text
知识库搜索
网页搜索
代码分析
报告生成
数据库查询
```

System Prompt 会越来越长，Tool 数量越来越多。

可能产生：

```text
Tool Selection 难度增加
Prompt 复杂
上下文膨胀
错误工具调用增加
```

于是可以把不同职责交给不同 SubAgent。

```text
                Main Agent
          ┌────────┼────────┐
          ↓        ↓        ↓
       RAG Agent  Web Agent  Code Agent
          ↓        ↓        ↓
        KB Tool   Web Tool  GitHub Tool
```

Main Agent 不直接管理所有底层 Tool，而是调用更高层能力。

## 19.1 这里体现了什么设计思想

首先是职责拆分：

```text
Single Responsibility
```

其次也可以看作组合结构：

```text
Agent
├── SubAgent
│   ├── Tool
│   └── Tool
└── SubAgent
```

但要注意：

> SubAgent 不是 Tool 越多越应该拆。

因为 SubAgent 会引入：

```text
额外模型调用
更高 Token 成本
更长延迟
更复杂 Trace
跨 Agent 上下文同步
```

因此真正的判断标准应该是：

> 职责边界是否已经明显不同，并且拆分后能显著降低单 Agent 的认知复杂度。

---

# 20. Trace：为什么可观测性也可以通过 Middleware / Observer 思想实现

Agent 与普通 HTTP 最大区别之一是调用链很长。

一个请求可能：

```text
HTTP Request
   ↓
Agent Run
   ↓
LLM Call #1
   ↓
Tool Call
   ↓
Query KB
   ↓
Vector Search
   ↓
Keyword Search
   ↓
Rerank
   ↓
LLM Call #2
   ↓
SSE Output
```

如果最终回答有问题，仅记录：

```text
HTTP 200
```

几乎没有意义。

所以需要 Trace：

```text
Trace
└── Agent Span
    ├── LLM Span #1
    ├── Tool Span
    │   └── RAG Span
    │       ├── Vector Span
    │       ├── Keyword Span
    │       └── Rerank Span
    └── LLM Span #2
```

## 20.1 Context 传播

Go 中非常适合把 Span Context 放进 `context.Context`：

```go
ctx, span := tracer.Start(ctx, "agent.run")
defer span.End()

result, err := agent.Run(ctx, req)
```

下游继续：

```go
func (t *QueryKBTool) Execute(
    ctx context.Context,
    args map[string]any,
) (*ToolResult, error) {
    ctx, span := tracer.Start(ctx, "tool.query_kb")
    defer span.End()

    return t.retriever.Search(ctx, args)
}
```

只要传递的是派生后的 `ctx`，父子关系就可以继续建立。

## 20.2 一个非常常见的错误

```go
go func() {
    worker.Run(context.Background())
}()
```

如果这是从当前请求派生出的异步任务，直接使用 `context.Background()` 会把 Trace Context 丢掉。

于是：

```text
HTTP Trace
   ↓
创建异步任务

Worker Trace
```

两条链路断开。

真实异步场景通常需要把：

```text
trace_id
span_id / parent context
```

显式写入任务消息，再由 Worker 恢复。

这已经属于分布式 Trace 的工程问题，而不是简单的 GoF 模式。

这也是校招项目里非常有技术含量的深入点。

---

# 21. 一个完整的 Qavor 风格架构：这些模式是怎么组合起来的

把前面所有东西组合起来，就会发现设计模式从来不是孤立存在的。

## 21.1 文档索引链

```text
                Factory
                   ↓
              创建 Parser
                   ↓
                   │
PDF / DOCX / HTML  │
       ↓           │
     Adapter ──────┘
       ↓
Unified Document
       ↓
Pipeline
       ↓
Chunk Stage
       ↓
Strategy
├── MarkdownChunker
├── FixedChunker
└── SemanticChunker
       ↓
Embedding Stage
       ↓
Strategy
├── Qwen
├── OpenAI
└── BGE
       ↓
Repository
       ↓
pgvector
```

## 21.2 Agent Tool 链

```text
LLM Function Calling
        ↓
ToolCallCommand
        ↓
Command Executor
        ↓
Tool Registry
        ↓
找到 Tool
        ↓
Middleware Chain
├── Recovery
├── Trace
├── Logging
├── Auth
├── Approval
└── Timeout
        ↓
Tool Adapter
        ↓
真实 Tool
├── Query KB
├── HTTP
└── MCP
```

## 21.3 RAG 查询链

```text
RAG Facade / Service
        ↓
HybridRetriever
        ↓
 ┌──────┴──────┐
 ↓             ↓
Vector       Keyword
 ↓             ↓
 └──────┬──────┘
        ↓
      RRF
        ↓
Reranker Strategy
        ↓
失败？
├── 否 → Rerank Result
└── 是 → Fallback to RRF Result
        ↓
Prompt
        ↓
Model Proxy
├── Trace
├── Metrics
└── Rate Limit
        ↓
真实 Model Client
```

这时候可以发现：

> 一个成熟架构不是“用了某一个设计模式”，而是不同层次的问题由不同模式分别解决。

---

# 22. 最容易犯的错误：为了设计模式而设计模式

设计模式本身也会产生复杂度。

例如一个函数原本只有：

```go
func ParseMarkdown(content string) []Chunk
```

却设计成：

```text
AbstractChunkFactory
ChunkBuilder
ChunkStrategy
ChunkStrategyFactory
ChunkRegistry
ChunkManager
```

代码量可能翻几倍。

这不是架构能力，而可能是过度设计。

## 22.1 什么时候暂时不要抽象

如果：

```text
只有一个实现
未来变化概率低
逻辑只有十几行
没有复用需求
没有测试隔离需求
```

直接写可能最好。

## 22.2 什么时候应该考虑抽象

如果开始出现：

```text
大量 switch / if
相同逻辑复制多份
第三方 SDK 类型到处传播
一个 Service 知道过多实现细节
为了测试不得不启动所有外部依赖
新增功能需要修改很多旧代码
```

通常说明已经出现明确的变化边界。

这时候设计模式才开始真正产生收益。

---

# 23. 从 SOLID 再理解一次这些模式

设计模式和 SOLID 原则不是两套完全无关的东西。

很多模式其实是在实现这些原则。

## 单一职责原则 SRP

```text
Parser 只解析
Chunker 只分块
Retriever 只检索
Reranker 只排序
Repository 只处理持久化
Middleware 处理横切逻辑
```

## 开闭原则 OCP

例如新增一个 Retriever：

```go
type GraphRetriever struct{}
```

理想情况下应该主要是增加新代码，而不是修改大量旧代码。

## 依赖倒置原则 DIP

高层业务：

```text
RAGService
```

不直接依赖：

```text
PgVectorRetriever
OpenAIClient
```

而依赖：

```text
Retriever
ChatModel
```

这就是 Dependency Inversion。

## 接口隔离原则 ISP

不要设计：

```go
type AIComponent interface {
    Parse()
    Chunk()
    Embed()
    Retrieve()
    Rerank()
    Generate()
}
```

因为大部分实现只需要其中一两个方法。

应该拆成小接口：

```text
Parser
Chunker
Embedder
Retriever
Reranker
ChatModel
```

Go 特别适合“小接口”。

---

# 24. 面试官真正想听的不是“你用了什么模式”

如果面试官问：

> 你的项目中用过哪些设计模式？

不建议回答：

> 我用了工厂模式、策略模式、适配器模式、观察者模式、责任链模式、命令模式……

这听起来非常像刚背完设计模式。

更好的回答方式是从真实问题出发。

例如：

> 我项目里没有为了套模式专门设计类层次，更多是遇到变化点以后自然做抽象。比如 RAG 的 Chunk、Retriever、Reranker 都存在多种实现，所以我通过接口把它们抽成可替换策略；Agent Tool 这边因为 Tool 是动态扩展的，我用了 Registry 根据 Tool Name 做运行时查找；日志、Trace、超时和权限属于横切逻辑，所以更适合放在 Middleware Chain，而不是散在每一个 Tool 里面。另外，外部的 Eino、MCP 或 Python Parser，我会尽量通过 Adapter 隔离第三方协议，避免业务层直接依赖外部 SDK。

这段话相比单纯罗列模式，多了三个非常重要的信息：

```text
为什么出现这个模式
模式放在哪一层
它解决了什么耦合
```

这才是面试官真正想判断的工程能力。

---

# 25. 面试官可能继续追问什么

## 追问 1：策略模式和工厂模式有什么区别？

可以回答：

> 策略模式解决的是“同一种能力有多种算法实现，运行时选择哪一种”；工厂模式解决的是“具体对象怎么创建”。比如 Chunker 有 Markdown、Fixed、Semantic，这是 Strategy；根据配置创建具体 Chunker 实例的那段逻辑则更接近 Factory。两者经常一起使用。

## 追问 2：为什么不用 switch？

不要说：

> switch 不优雅。

更好的回答：

> 如果实现数量很少而且稳定，我会直接用 switch，因为更简单。只有当分支持续增加、调用方开始依赖实现细节，或者我需要通过接口进行测试替换时，我才会抽 Strategy 或 Factory。设计模式本身也有复杂度，不应该为了模式消灭所有 switch。

## 追问 3：Middleware 属于责任链还是装饰器？

可以回答：

> 从实现结构上看，Go Middleware 通常通过 `func(next Handler) Handler` 不断包装下一个 Handler，这体现了 Decorator 的包装思想；从请求执行过程看，又是一个请求依次经过 Trace、Logging、Auth、Timeout 等节点，因此也具有责任链的特点。实际工程里我不会强行只给它贴一个标签，更关注它解决横切关注点的问题。

## 追问 4：为什么 Adapter 很重要？

> 因为 Agent 项目第三方依赖很多，比如模型 SDK、Eino、MCP、Parser 服务。如果业务代码直接依赖这些外部类型，SDK 变化会传播到很多模块。我更希望系统内部先定义自己的稳定接口，然后在边界通过 Adapter 做协议和类型转换，把第三方变化控制在较小范围内。

## 追问 5：Registry 并发安全吗？

> 如果 Registry 只在启动阶段注册、运行阶段只读，那么并发问题比较简单；如果支持 MCP Tool 动态加载、启停或者运行时刷新，就会存在 map 并发读写，需要 RWMutex、copy-on-write 或其他同步方案。同时还需要处理重复注册、删除和正在执行 Tool 之间的生命周期问题。

## 追问 6：Pipeline 某一步失败怎么办？

> 要看这个阶段是否可重试、是否有副作用。例如 Parse 失败可以标记任务失败；Embedding API 临时失败可以有限重试；Store 阶段如果已经写入部分 Chunk，就要考虑事务、幂等写或者失败补偿。Pipeline 只描述阶段式结构，真正工程落地还要配合状态机、幂等、重试和补偿机制。

---

# 26. 作为应届生应该掌握到什么深度

如果目标是 Go 后端 / Agent 应用校招，不需要把 23 种 GoF 模式全部背到 UML 级别。

建议优先真正掌握：

```text
必须掌握
────────────
Strategy
Factory
Adapter
Observer / Pub-Sub 思想
Decorator / Middleware
Chain of Responsibility

后端工程必须理解
────────────────
Repository
Dependency Injection
Registry
Pipeline
State Machine
Retry / Timeout / Fallback

Agent 项目重点理解
────────────────
Tool Registry
Tool Adapter
Middleware
Command / Function Calling
ReAct
Workflow / Graph
SubAgent
Trace Context
```

加分项：

```text
Composite
Proxy
Circuit Breaker
Saga / Compensation
Outbox Pattern
CQRS
Event Sourcing
```

但如果项目里没有真实使用，不建议为了面试主动声称自己“实践过”。

知道原理和适用场景即可。

---

# 27. 最后总结：设计模式的本质是管理变化

如果整篇文章只记住一句话，可以记：

> **设计模式的本质不是记住类图，而是识别系统里的变化点，然后通过稳定的抽象把变化限制在局部。**

重新看 Qavor 这样的系统：

```text
Chunk 方法会变化
→ Strategy

第三方接口会变化
→ Adapter

Tool 类型会变化
→ Factory

Tool 数量会动态变化
→ Registry

日志 / Trace / Auth 是横切逻辑
→ Middleware + Decorator + Chain

RAG 有明确阶段
→ Pipeline

异步任务状态复杂
→ State Machine

任务进度需要通知前端
→ Observer / Pub-Sub

Tool Call 需要排队、审批、审计
→ Command

复杂子系统需要给上层简单入口
→ Facade

数据访问不应该污染业务层
→ Repository

实现需要可替换、可测试
→ Dependency Injection

Rerank / 模型服务可能失败
→ Retry + Timeout + Fallback

开放式 Tool 决策
→ ReAct

确定性任务流
→ Workflow / Graph
```

它们最终解决的其实都是同一件事：

```text
今天系统能运行
        ↓
明天需求变化
        ↓
代码仍然能继续演进
```

这才是设计模式真正的工程价值。

---

# 附：面试中可以直接说的一版

如果面试官问：

> “你项目里用过哪些设计模式？”

可以这样回答：

> 我项目里其实没有为了套 GoF 模式刻意设计很多类，更多是遇到明确的变化点以后做抽象。比如 RAG 的 Chunk、Embedding、Retriever 和 Reranker 都可能有多种实现，所以我会通过 Go interface 做 Strategy，让索引或者查询 Pipeline 依赖抽象而不是具体实现。Tool 这一层因为运行时要根据模型返回的 Tool Name 动态找到工具，所以比较适合 Registry，Tool 的创建逻辑再交给 Factory。像 Trace、Logging、Timeout、Auth 这种很多 Tool 都需要的横切逻辑，我不会写进每个 Tool，而是通过 Middleware Chain 统一处理。对于 Eino、MCP、Python Parser 这种外部框架和服务，我更倾向在系统边界增加 Adapter，避免第三方接口直接污染核心业务层。异步文档索引这边则更像 Pipeline 加 State Machine，因为任务会经过 Parse、Chunk、Embedding、Index 等阶段，还要处理失败、重试和状态恢复。对我来说这些模式最重要的不是名称，而是把变化点隔离起来，让后续增加实现或者替换组件时不需要大面积改业务代码。

如果面试官继续追问：

> “那你觉得设计模式越多越好吗？”

可以继续说：

> 不是。我觉得设计模式本身也会增加抽象层和理解成本。如果当前只有一个实现，而且未来也没有明显变化，我会优先保持代码简单。一般是在看到重复逻辑、大量 switch、第三方接口向业务层扩散，或者测试很难替换依赖时，我才认为这个变化点值得抽象。对于应届生项目，我更希望能解释清楚为什么抽象，而不是单纯说用了多少种模式。

这类回答通常比背诵“23 种设计模式定义”更能体现真实工程思考。
