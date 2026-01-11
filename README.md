# gobuilder-mcp

一个用于跨平台构建 Go 应用程序的 MCP (Model Context Protocol) 服务器。

## 功能特性

- 🚀 **跨平台构建**: 支持 Windows、macOS、Linux 多平台并行构建
- ⚡ **快速构建**: 专门为 MCP 服务优化的快速构建模式
- 📋 **目标列表**: 列出所有支持的编译目标平台
- 🐳 **Docker 支持**: 提供 Docker 镜像用于容器化部署
- 🔄 **CI/CD**: 完整的 GitHub Actions 工作流

## 安装

### 从源码构建

```bash
git clone https://github.com/fromsko/gobuilder-mcp.git
cd gobuilder-mcp
go build -o gobuilder-mcp main.go
```

### 使用 Docker

```bash
# 生产环境
docker pull fromsko/gobuilder-mcp:latest
docker run --rm fromsko/gobuilder-mcp:latest

# 开发环境
docker build -f Dockerfile.dev -t gobuilder-mcp:dev .
docker run --rm gobuilder-mcp:dev
```

## 使用方法

### 作为 MCP 服务器

启动 MCP 服务器：

```bash
./gobuilder-mcp
```

然后在支持 MCP 的客户端中配置服务器路径。

### 可用工具

#### 1. cross_platform_build
跨平台构建 Go 应用程序，支持 Windows、macOS、Linux 多目标平台并行构建。

**参数:**
- `source_file`: Go源文件路径，支持绝对路径和相对路径，默认为 ./main.go
- `app_name`: 生成的可执行文件名称，默认为 app
- `output_dir`: 输出目录路径，支持绝对路径和相对路径，默认为 bin
- `targets`: 编译目标平台列表，为空则默认编译 Linux x64 和 Windows x64
- `jobs`: 并行构建任务数，默认为 4，建议不超过 CPU 核心数

#### 2. mcp_quick_build
快速构建 Go 应用程序，专门为 MCP 服务优化，同时构建 Linux x64 和 Windows x64 版本。

**参数:**
- `source_file`: Go源文件路径，支持绝对路径和相对路径，默认为 ./main.go
- `app_name`: 生成的可执行文件名称，默认为 app
- `output_dir`: 输出目录路径，支持绝对路径和相对路径，默认为 bin

#### 3. list_build_targets
列出所有支持的编译目标平台，包含 GOOS、GOARCH 和平台说明信息。

## 支持的平台

| 平台 | GOOS | GOARCH | 说明 |
|------|------|--------|------|
| Windows x64 | windows | amd64 | Windows x64 |
| macOS x64 | darwin | amd64 | macOS x64 |
| macOS ARM64 | darwin | arm64 | macOS ARM64 |
| Linux x64 | linux | amd64 | Linux x64 |

## 开发

### 本地开发

```bash
# 克隆仓库
git clone https://github.com/fromsko/gobuilder-mcp.git
cd gobuilder-mcp

# 安装依赖
go mod tidy

# 运行测试
go test ./...

# 构建
go build -o gobuilder-mcp main.go

# 运行
./gobuilder-mcp
```

### 使用 Docker 开发

```bash
# 构建开发镜像
docker build -f Dockerfile.dev -t gobuilder-mcp:dev .

# 运行开发容器
docker run --rm -v $(pwd):/app gobuilder-mcp:dev
```

## 发布

### 自动发布

项目使用 GitHub Actions 进行自动发布：

1. 创建版本标签：`git tag v1.0.0`
2. 推送标签：`git push origin v1.0.0`
3. GitHub Actions 将自动：
   - 运行测试和代码检查
   - 构建多平台二进制文件
   - 创建 GitHub Release
   - 构建 and push Docker 镜像

### 手动发布

```bash
# 构建多平台版本
GOOS=linux GOARCH=amd64 go build -o gobuilder-mcp-linux-amd64 main.go
GOOS=windows GOARCH=amd64 go build -o gobuilder-mcp-windows-amd64.exe main.go
GOOS=darwin GOARCH=amd64 go build -o gobuilder-mcp-darwin-amd64 main.go
GOOS=darwin GOARCH=arm64 go build -o gobuilder-mcp-darwin-arm64 main.go
```

## 配置

### MCP 客户端配置

在你的 MCP 客户端配置文件中添加：

```json
{
  "mcpServers": {
    "gobuilder-mcp": {
      "command": "/path/to/gobuilder-mcp"
    }
  }
}
```

### 环境变量

- `GO_VERSION`: 指定 Go 版本（默认：1.21）
- `GOMAXPROCS`: 设置最大 CPU 核心数

## 贡献

欢迎贡献代码！请遵循以下步骤：

1. Fork 项目
2. 创建功能分支：`git checkout -b feature/amazing-feature`
3. 提交更改：`git commit -m 'Add amazing feature'`
4. 推送分支：`git push origin feature/amazing-feature`
5. 创建 Pull Request

## 许可证

本项目采用 MIT 许可证。详见 [LICENSE](LICENSE) 文件。

## 更新日志

### v1.0.0
- 初始版本发布
- 支持跨平台构建
- 集成 MCP 协议
- 添加 Docker 支持
- 完整的 CI/CD 工作流

## 支持

如果你遇到问题或有建议，请：

1. 查看 [Issues](https://github.com/fromsko/gobuilder-mcp/issues)
2. 创建新的 Issue
3. 参与 [Discussions](https://github.com/fromsko/gobuilder-mcp/discussions)

## 作者

- [fromsko](https://github.com/fromsko)

## 致谢

感谢 [Model Context Protocol](https://modelcontextprotocol.io/) 项目提供的 MCP 框架。

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
