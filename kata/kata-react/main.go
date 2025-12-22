// main.go
package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

// System prompt：要求模型按 ReAct 格式输出
const systemPrompt = `你是一个 AI 编程老师（中文），要以“ReAct”模式工作（Reason + Act）。
当你需要调用外部工具时，请严格按以下格式输出（不要输出其他没有说明的内容）：
Thought: <你当前的思路简短描述>
Action: <tool_name>
Action Input: <JSON 格式或纯文本的工具输入>

如果你已经准备好最终答案（不需要调用工具），请直接输出：
FinalAnswer: <你的答案>

可用工具：
- search    : 在本地知识库中按关键字检索（返回若干条简短片段）
- run_code  : 在受限环境下运行 Go 代码片段（默认关闭，需环境变量 ALLOW_CODE_EXEC=1 开启）

注意：只在确实需要“外部数据/运行”时才调用工具；否则直接给出教学建议或学习步骤。
`

// Simple in-memory docs作为 demo（实际应替换为向量检索）
var knowledgeBase = []string{
	"入门路线：先学 Python 基础 -> 数据结构与算法 -> 线性代数/概率基础 -> PyTorch/TensorFlow -> 实践项目",
	"Golang 与 AI：Go 常用于后端与工具链，训练/研究通常用 Python；你可以用 Go 做模型服务和工程化。",
	"常见算法：反向传播、梯度下降、交叉熵损失、卷积神经网络、Transformer 架构。",
	"实践建议：做 3 个项目：1) 图像分类 2) 文本分类 3) 小型生成模型微调。",
	"工具建议：PyTorch、Hugging Face Transformers、Weights & Biases（实验跟踪）",
}

type Tool func(ctx context.Context, input string) (string, error)

func main() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("请先设置 OPENAI_API_KEY 环境变量")
	}
	client := openai.NewClient(option.WithAPIKey(apiKey))
	ctx := context.Background()

	tools := map[string]Tool{
		"search":   searchDocs,
		"run_code": runCode,
	}

	fmt.Println("=== AI 编程老师（ReAct 示例） ===")
	fmt.Println("请输入你的问题（例如：如何从零开始学 AI 编程？）:")
	reader := bufio.NewReader(os.Stdin)
	userQ, _ := reader.ReadString('\n')
	userQ = strings.TrimSpace(userQ)
	if userQ == "" {
		log.Fatal("问题为空，退出")
	}

	if err := reactAgentLoop(ctx, client, tools, userQ); err != nil {
		log.Fatalf("agent 错误: %v", err)
	}
}

// reactAgentLoop：核心 ReAct 循环实现
func reactAgentLoop(ctx context.Context, client *openai.Client, tools map[string]Tool, userQuestion string) error {
	// messages 用于保存对话历史（系统 + 交互）
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(systemPrompt),
		openai.UserMessage(userQuestion),
	}

	// 最多循环 N 次以避免无限循环
	const maxSteps = 6

	for step := 0; step < maxSteps; step++ {
		// 调用 OpenAI Chat Completion
		resp, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
			Messages: messages,
			Model:    openai.ChatModelGPT4o, // 可替换为适合的模型
		})
		if err != nil {
			return fmt.Errorf("chat completion 调用失败: %w", err)
		}
		if len(resp.Choices) == 0 {
			return fmt.Errorf("无返回 choices")
		}

		content := resp.Choices[0].Message.Content
		fmt.Printf("\n=== LLM 输出（step %d）===\n%s\n", step+1, content)

		// 解析是否为 Action（ReAct 格式）
		action, actionInput, isAction := parseAction(content)
		if !isAction {
			// 不是 action，检查是否包含 FinalAnswer
			final := parseFinalAnswer(content)
			if final != "" {
				fmt.Printf("\n--- Agent 最终答案 ---\n%s\n", final)
				return nil
			}
			// 既没有 action 也没有 final，直接作为答案结束
			fmt.Printf("\n--- Agent 回答（非结构化）---\n%s\n", content)
			return nil
		}

		// 有 action，查找工具
		tool, ok := tools[action]
		if !ok {
			obs := fmt.Sprintf("错误：未找到工具 '%s'。", action)
			// 将模型原始输出和 observation 加入 history，继续循环
			messages = append(messages, openai.AssistantMessage(content))
			messages = append(messages, openai.AssistantMessage("Observation: "+obs))
			fmt.Println("工具未注册，继续让模型决定下一步。")
			continue
		}

		// 执行工具（注意：run_code 默认受限，需要额外允许）
		obs, err := tool(ctx, actionInput)
		if err != nil {
			obs = fmt.Sprintf("工具执行失败: %v", err)
		}

		// 把模型的 action 输出和工具的 observation 加入对话历史，让模型看见结果
		messages = append(messages, openai.AssistantMessage(content))
		// observation 也以 assistant 消息加入（合适地告诉模型工具结果）
		messages = append(messages, openai.AssistantMessage("Observation: "+obs))
		// 下一轮模型会基于 observation 继续推理或给出 FinalAnswer
	}

	fmt.Println("达到最大交互步数，停止。")
	return nil
}

// parseAction：从 LLM 文本中提取 Action 与 Action Input（简单正则）
func parseAction(s string) (action string, input string, ok bool) {
	reAct := regexp.MustCompile(`(?m)^Action:\s*([a-zA-Z0-9_-]+)`)
	reInput := regexp.MustCompile(`(?ms)^Action Input:\s*(.+)$`)
	actM := reAct.FindStringSubmatch(s)
	if len(actM) < 2 {
		return "", "", false
	}
	action = strings.TrimSpace(actM[1])
	inM := reInput.FindStringSubmatch(s)
	if len(inM) >= 2 {
		input = strings.TrimSpace(inM[1])
	} else {
		input = ""
	}
	return action, input, true
}

// parseFinalAnswer：查找 FinalAnswer:
func parseFinalAnswer(s string) string {
	re := regexp.MustCompile(`(?ms)FinalAnswer:\s*(.+)$`)
	m := re.FindStringSubmatch(s)
	if len(m) >= 2 {
		return strings.TrimSpace(m[1])
	}
	return ""
}

// searchDocs：一个非常简单的全文 substring 检索示例（replace by vector DB）
func searchDocs(ctx context.Context, input string) (string, error) {
	q := strings.ToLower(input)
	var results []string
	for _, d := range knowledgeBase {
		if strings.Contains(strings.ToLower(d), q) {
			results = append(results, d)
		}
	}
	if len(results) == 0 {
		// 近似匹配：返回 top-3（演示用）
		for i := 0; i < len(knowledgeBase) && i < 3; i++ {
			results = append(results, knowledgeBase[i])
		}
	}
	out := "检索到：" + strings.Join(results, " || ")
	return out, nil
}

// runCode：在受限环境下运行 go 代码（默认禁用，需设置环境变量 ALLOW_CODE_EXEC=1）
func runCode(ctx context.Context, input string) (string, error) {
	if os.Getenv("ALLOW_CODE_EXEC") != "1" {
		return "EXEC_DISABLED: 本地代码执行已禁用（设置 ALLOW_CODE_EXEC=1 允许）", nil
	}

	// 将用户输入包装在 main 模板中（注意：这只是示例，真实执行需更严格的沙箱）
	template := `package main
import (
	"fmt"
)
func main() {
%s
}
`
	code := fmt.Sprintf(template, input)

	tmpDir, err := os.MkdirTemp("", "react-code-*")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmpDir)

	fpath := filepath.Join(tmpDir, "snippet.go")
	if err := os.WriteFile(fpath, []byte(code), 0644); err != nil {
		return "", err
	}

	// 使用 go run 执行，限制超时
	cmd := exec.CommandContext(ctx, "go", "run", fpath)
	// 设置超时（总共 6s）
	ctx2, cancel := context.WithTimeout(ctx, 6*time.Second)
	defer cancel()
	cmd = exec.CommandContext(ctx2, "go", "run", fpath)

	out, err := cmd.CombinedOutput()
	if ctx2.Err() == context.DeadlineExceeded {
		return "EXEC_TIMEOUT: 代码运行超时", nil
	}
	if err != nil {
		return fmt.Sprintf("EXEC_ERROR: %v\n%s", err, string(out)), nil
	}
	return fmt.Sprintf("EXEC_OK:\n%s", string(out)), nil
}
