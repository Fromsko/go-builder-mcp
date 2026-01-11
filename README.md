# GoBuilder MCP Server

一个用于跨平台构建 Go 应用程序的 MCP (Model Context Protocol) 服务器。

## 功能特性

- 🚀 跨平台构建：支持 Windows、macOS、Linux
- ⚡ 并行编译：提高构建效率
- 📦 灵活配置：自定义应用名称、输出目录等
- 🔧 MCP 集成：通过 MCP 协议提供工具接口

## 安装和使用

### 1. 构建服务器

```bash
cd gobuilder-mcp
go mod tidy
go build -o gobuilder-mcp main.go
```

### 2. 在 MCP 客户端中配置

在你的 MCP 客户端配置文件中添加：

```json
{
  "mcpServers": {
    "gobuilder": {
      "command": "/path/to/gobuilder-mcp/gobuilder-mcp"
    }
  }
}
```

## 可用工具

### 1. `cross_platform_build`

跨平台构建 Go 应用程序，支持多个目标平台并行构建。

**参数：**
- `source_file` (string): 源文件路径，默认为 `./main.go`
- `app_name` (string): 应用名称，默认为 `app`
- `output_dir` (string): 输出目录，默认为 `bin`
- `targets` (array): 编译目标列表，为空则默认编译 Linux 和 Windows
- `jobs` (int): 并行任务数，默认为 CPU 核心数

**示例：**
```json
{
  "source_file": "./main.go",
  "app_name": "myapp",
  "output_dir": "dist",
  "targets": [
    {"goos": "linux", "goarch": "amd64", "name": "Linux x64"},
    {"goos": "windows", "goarch": "amd64", "name": "Windows x64"},
    {"goos": "darwin", "goarch": "amd64", "name": "macOS x64"}
  ],
  "jobs": 4
}
```

### 2. `mcp_quick_build`

MCP 快速构建，专门为 MCP 服务构建 Linux 和 Windows 版本。

**参数：**
- `source_file` (string, 可选): 源文件路径，默认为 `./main.go`
- `app_name` (string, 可选): 应用名称，默认为 `app`
- `output_dir` (string, 可选): 输出目录，默认为 `bin`

**示例：**
```json
{
  "app_name": "mcp-service",
  "output_dir": "bin"
}
```

### 3. `list_build_targets`

列出所有支持的编译目标平台。

**参数：** 无

## 支持的目标平台

| 平台 | GOOS | GOARCH | 说明 |
|------|------|--------|------|
| Windows x64 | windows | amd64 | 64位 Windows |
| macOS x64 | darwin | amd64 | Intel Mac |
| macOS ARM64 | darwin | arm64 | Apple Silicon Mac |
| Linux x64 | linux | amd64 | 64位 Linux |

## 使用示例

### 在 Claude Desktop 中使用

1. 配置 MCP 服务器后，你可以在对话中直接使用工具：

```
请帮我构建我的 Go 应用，需要支持 Windows 和 Linux 平台
```

2. Claude 会自动调用 `mcp_quick_build` 工具：

```
🚀 跨平台构建完成

📁 源文件: ./main.go
📦 应用名称: app
📂 输出目录: bin
⚡ 并行任务数: 2

✅ 成功构建:
  • Linux x64: bin/app_linux_amd64
  • Windows x64: bin/app_windows_amd64.exe

🎉 所有目标构建成功！共 2 个
```

### 自定义构建

```
我需要构建一个名为 "my-tool" 的应用，支持所有平台，输出到 dist 目录
```

Claude 会调用 `cross_platform_build` 工具并传入相应参数。

## 开发说明

### 项目结构

```
gobuilder-mcp/
├── main.go          # MCP 服务器主文件
├── go.mod           # Go 模块文件
├── go.sum           # 依赖锁定文件
├── README.md        # 说明文档
└── gobuilder-mcp    # 编译后的可执行文件
```

### 依赖项

- `github.com/modelcontextprotocol/go-sdk`: MCP Go SDK

### 扩展功能

可以通过修改 `supportedTargets` 数组来添加更多支持的平台：

```go
var supportedTargets = []BuildTarget{
    {GOOS: "windows", GOARCH: "amd64", Name: "Windows x64"},
    {GOOS: "darwin", GOARCH: "amd64", Name: "macOS x64"},
    {GOOS: "darwin", GOARCH: "arm64", Name: "macOS ARM64"},
    {GOOS: "linux", GOARCH: "amd64", Name: "Linux x64"},
    // 添加更多平台...
    {GOOS: "freebsd", GOARCH: "amd64", Name: "FreeBSD x64"},
}
```

## 许可证

MIT License
