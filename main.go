package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// 版本信息（由 goreleaser 注入）
var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

// getVersionInfo 获取版本信息
func getVersionInfo() string {
	return fmt.Sprintf("gobuilder-mcp version %s (commit: %s, built: %s)", version, commit, date)
}

// 编译目标配置
type BuildTarget struct {
	GOOS   string `json:"goos" jsonschema:"目标操作系统"`
	GOARCH string `json:"goarch" jsonschema:"目标架构"`
	Name   string `json:"name" jsonschema:"平台名称"`
}

var supportedTargets = []BuildTarget{
	{GOOS: "windows", GOARCH: "amd64", Name: "Windows x64"},
	{GOOS: "darwin", GOARCH: "amd64", Name: "macOS x64"},
	{GOOS: "darwin", GOARCH: "arm64", Name: "macOS ARM64"},
	{GOOS: "linux", GOARCH: "amd64", Name: "Linux x64"},
}

// ============================================
// 工具参数定义
// ============================================

// BuildParam 构建参数
type BuildParam struct {
	SourceFile string        `json:"source_file" jsonschema:"Go源文件路径，支持绝对路径和相对路径，默认为 ./main.go"`
	AppName    string        `json:"app_name" jsonschema:"生成的可执行文件名称，默认为 app"`
	OutputDir  string        `json:"output_dir" jsonschema:"输出目录路径，支持绝对路径和相对路径，默认为 bin"`
	Targets    []BuildTarget `json:"targets" jsonschema:"编译目标平台列表，为空则默认编译 Linux x64 和 Windows x64"`
	Jobs       int           `json:"jobs" jsonschema:"并行构建任务数，默认为 4，建议不超过 CPU 核心数"`
}

// BuildResult 构建结果
type BuildResult struct {
	Success       bool     `json:"success" jsonschema:"整体构建是否成功，true表示所有目标都构建成功"`
	BuiltTargets  []string `json:"built_targets" jsonschema:"成功构建的目标列表，包含平台名称和文件路径"`
	FailedTargets []string `json:"failed_targets" jsonschema:"构建失败的目标列表，包含平台名称和错误信息"`
	OutputDir     string   `json:"output_dir" jsonschema:"实际的输出目录路径"`
}

// QuickBuildParam 快速构建参数
type QuickBuildParam struct {
	SourceFile string `json:"source_file,omitempty" jsonschema:"Go源文件路径，支持绝对路径和相对路径，默认为 ./main.go"`
	AppName    string `json:"app_name,omitempty" jsonschema:"生成的可执行文件名称，默认为 app"`
	OutputDir  string `json:"output_dir,omitempty" jsonschema:"输出目录路径，支持绝对路径和相对路径，默认为 bin"`
}

// ListTargetsParam 列出目标参数
type ListTargetsParam struct{}

// TargetsOutput 目标列表输出
type TargetsOutput struct {
	Targets string `json:"targets" jsonschema:"支持的编译目标列表"`
}

// ============================================
// 工具实现
// ============================================

// CrossPlatformBuild 跨平台构建
func CrossPlatformBuild(ctx context.Context, req *mcp.CallToolRequest, param BuildParam) (
	*mcp.CallToolResult,
	BuildResult,
	error,
) {
	// 设置默认值
	if param.SourceFile == "" {
		param.SourceFile = "./main.go"
	}
	if param.AppName == "" {
		param.AppName = "app"
	}
	if param.OutputDir == "" {
		param.OutputDir = "bin"
	}
	if param.Jobs <= 0 {
		param.Jobs = 4 // 默认4个并行任务
	}

	// 检查源文件是否存在
	if !fileExists(param.SourceFile) {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("❌ 源文件不存在: %s", param.SourceFile)},
			},
		}, BuildResult{Success: false, FailedTargets: []string{}}, nil
	}

	// 检查Go是否安装
	if !isGoInstalled() {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "❌ Go 未安装或不在 PATH 中"},
			},
		}, BuildResult{Success: false, FailedTargets: []string{}}, nil
	}

	// 创建输出目录
	if err := os.MkdirAll(param.OutputDir, 0755); err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("❌ 创建输出目录失败: %v", err)},
			},
		}, BuildResult{Success: false, FailedTargets: []string{}}, nil
	}

	// 确定编译目标
	var targets []BuildTarget
	if len(param.Targets) == 0 {
		// 默认构建 Linux 和 Windows
		targets = []BuildTarget{
			{GOOS: "linux", GOARCH: "amd64", Name: "Linux x64"},
			{GOOS: "windows", GOARCH: "amd64", Name: "Windows x64"},
		}
	} else {
		targets = param.Targets
	}

	// 并行构建
	semaphore := make(chan struct{}, param.Jobs)
	results := make(chan buildTaskResult, len(targets))
	var wg sync.WaitGroup

	for _, t := range targets {
		wg.Add(1)
		go func(t BuildTarget) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			suffix := ""
			if t.GOOS == "windows" {
				suffix = ".exe"
			}
			outputName := fmt.Sprintf("%s_%s_%s%s", param.AppName, t.GOOS, t.GOARCH, suffix)
			outputPath := filepath.Join(param.OutputDir, outputName)

			err := compileTarget(t, param.SourceFile, outputPath)
			results <- buildTaskResult{
				target:  t,
				path:    outputPath,
				success: err == nil,
				error:   err,
			}
		}(t)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	// 收集结果
	builtTargets := make([]string, 0)
	failedTargets := make([]string, 0)
	for result := range results {
		if result.success {
			builtTargets = append(builtTargets, fmt.Sprintf("%s: %s", result.target.Name, result.path))
		} else {
			failedTargets = append(failedTargets, fmt.Sprintf("%s: %v", result.target.Name, result.error))
		}
	}

	// 格式化输出
	var output strings.Builder
	fmt.Fprintf(&output, "🚀 跨平台构建完成\n\n")
	fmt.Fprintf(&output, "📁 源文件: %s\n", param.SourceFile)
	fmt.Fprintf(&output, "📦 应用名称: %s\n", param.AppName)
	fmt.Fprintf(&output, "📂 输出目录: %s\n", param.OutputDir)
	fmt.Fprintf(&output, "⚡ 并行任务数: %d\n\n", param.Jobs)

	if len(builtTargets) > 0 {
		fmt.Fprintf(&output, "✅ 成功构建:\n")
		for _, target := range builtTargets {
			fmt.Fprintf(&output, "  • %s\n", target)
		}
	}

	if len(failedTargets) > 0 {
		fmt.Fprintf(&output, "\n❌ 构建失败:\n")
		for _, target := range failedTargets {
			fmt.Fprintf(&output, "  • %s\n", target)
		}
	}

	success := len(failedTargets) == 0
	if success {
		fmt.Fprintf(&output, "\n🎉 所有目标构建成功！共 %d 个\n", len(builtTargets))
	} else {
		fmt.Fprintf(&output, "\n⚠️ 部分构建失败。成功: %d，失败: %d\n", len(builtTargets), len(failedTargets))
	}

	return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: output.String()},
			},
		}, BuildResult{
			Success:       success,
			BuiltTargets:  builtTargets,
			FailedTargets: failedTargets,
			OutputDir:     param.OutputDir,
		}, nil
}

// MCPQuickBuild MCP 快速构建
func MCPQuickBuild(ctx context.Context, req *mcp.CallToolRequest, param QuickBuildParam) (
	*mcp.CallToolResult,
	BuildResult,
	error,
) {
	// 转换为完整构建参数
	buildParam := BuildParam{
		SourceFile: param.SourceFile,
		AppName:    param.AppName,
		OutputDir:  param.OutputDir,
		Jobs:       2, // MCP 模式使用2个并行任务
		Targets: []BuildTarget{
			{GOOS: "linux", GOARCH: "amd64", Name: "Linux x64"},
			{GOOS: "windows", GOARCH: "amd64", Name: "Windows x64"},
		},
	}

	return CrossPlatformBuild(ctx, req, buildParam)
}

// ListBuildTargets 列出支持的构建目标
func ListBuildTargets(ctx context.Context, req *mcp.CallToolRequest, param ListTargetsParam) (
	*mcp.CallToolResult,
	TargetsOutput,
	error,
) {
	var output strings.Builder
	fmt.Fprintf(&output, "🎯 支持的编译目标:\n\n")
	fmt.Fprintf(&output, "| 平台 | GOOS | GOARCH | 说明 |\n")
	fmt.Fprintf(&output, "|------|------|--------|------|\n")

	for _, t := range supportedTargets {
		fmt.Fprintf(&output, "| %s | %s | %s | %s |\n", t.Name, t.GOOS, t.GOARCH, t.Name)
	}

	fmt.Fprintf(&output, "\n💡 提示:\n")
	fmt.Fprintf(&output, "• 使用 cross_platform_build 工具进行自定义构建\n")
	fmt.Fprintf(&output, "• 使用 mcp_quick_build 工具快速构建 Linux 和 Windows 版本\n")

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: output.String()},
		},
	}, TargetsOutput{Targets: output.String()}, nil
}

// ============================================
// 辅助函数和类型
// ============================================

type buildTaskResult struct {
	target  BuildTarget
	path    string
	success bool
	error   error
}

func compileTarget(target BuildTarget, sourceFile, outputPath string) error {
	cmd := exec.CommandContext(context.Background(), "go", "build", "-o", outputPath, sourceFile)

	// 设置环境变量
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("GOOS=%s", target.GOOS),
		fmt.Sprintf("GOARCH=%s", target.GOARCH),
		"CGO_ENABLED=0",
	)

	// 执行命令并捕获输出
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v\nOutput: %s", err, string(output))
	}

	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func isGoInstalled() bool {
	_, err := exec.LookPath("go")
	return err == nil
}

func main() {
	// 创建 MCP 服务器实例
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "gobuilder-mcp",
		Version: version,
	}, nil)

	// 添加跨平台构建工具
	mcp.AddTool(server, &mcp.Tool{
		Name:        "cross_platform_build",
		Description: "跨平台构建Go应用程序，支持Windows、macOS、Linux多目标平台并行构建，自动创建输出目录",
	}, CrossPlatformBuild)

	// 添加 MCP 快速构建工具
	mcp.AddTool(server, &mcp.Tool{
		Name:        "mcp_quick_build",
		Description: "快速构建Go应用程序，专门为MCP服务优化，同时构建Linux x64和Windows x64版本",
	}, MCPQuickBuild)

	// 添加列出构建目标工具
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_build_targets",
		Description: "列出所有支持的编译目标平台，包含GOOS、GOARCH和平台说明信息",
	}, ListBuildTargets)

	// 启动服务器，通过 stdio 传输
	log.Printf("Starting %s...", getVersionInfo())
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
