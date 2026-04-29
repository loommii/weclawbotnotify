package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	qrterminal "github.com/mdp/qrterminal/v3"
)

type contextTokenEntry struct {
	Token      string
	UserID     string
	ReceivedAt time.Time
}

var (
	creds          *Credentials
	client         *Client
	updatesBuf     string
	currentToken   *contextTokenEntry
	tokenStoreMu   sync.RWMutex
	monitorRunning bool
	monitorCancel  context.CancelFunc
)

func main() {
	fmt.Println("╔══════════════════════════════════════════════╗")
	fmt.Println("║   WeClaw Demo - 微信 Bot 交互测试工具       ║")
	fmt.Println("║   单用户模式 / 测试 token 过期               ║")
	fmt.Println("╚══════════════════════════════════════════════╝")
	fmt.Println()

	loadSavedCreds()

	reader := bufio.NewReader(os.Stdin)
	for {
		printMenu()
		fmt.Print("> ")
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)

		switch line {
		case "1":
			doLogin()
		case "2":
			doStartMonitor()
		case "3":
			doSendMessage(reader)
		case "4":
			doTestTokenExpiry(reader)
		case "5":
			doShowTokenStatus()
		case "6":
			doStopMonitor()
		case "q", "Q":
			doStopMonitor()
			fmt.Println("再见！")
			return
		default:
			fmt.Println("未知命令")
		}
	}
}

func printMenu() {
	fmt.Println()
	fmt.Println("┌──────────────────────────────────────────────┐")
	fmt.Println("│ [1] 扫码登录                                │")
	fmt.Println("│ [2] 开始接收消息（后台长轮询）              │")
	fmt.Println("│ [3] 发送消息给当前用户                      │")
	fmt.Println("│ [4] 测试 context_token 过期                 │")
	fmt.Println("│ [5] 查看 context_token 状态                 │")
	fmt.Println("│ [6] 停止接收消息                            │")
	fmt.Println("│ [q] 退出                                    │")
	fmt.Println("└──────────────────────────────────────────────┘")

	status := "未登录"
	if creds != nil {
		status = fmt.Sprintf("已登录 (用户: %s)", creds.ILinkUserID)
		if monitorRunning {
			status += " | 接收中"
		}
		tokenStoreMu.RLock()
		if currentToken != nil {
			age := time.Since(currentToken.ReceivedAt).Round(time.Second)
			status += fmt.Sprintf(" | token已过%v", age)
		} else {
			status += " | 无token"
		}
		tokenStoreMu.RUnlock()
	}
	fmt.Printf("  状态: %s\n", status)
}

func loadSavedCreds() {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}
	path := home + "/.weclaw-demo-creds.json"
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	var c Credentials
	if json.Unmarshal(data, &c) == nil && c.BotToken != "" {
		creds = &c
		client = NewClient(c.BotToken, c.BaseURL)
		fmt.Printf("✅ 已加载保存的登录信息 (用户: %s)\n", c.ILinkUserID)
	}
}

func saveCreds() {
	if creds == nil {
		return
	}
	home, _ := os.UserHomeDir()
	path := home + "/.weclaw-demo-creds.json"
	data, _ := json.MarshalIndent(creds, "", "  ")
	_ = os.WriteFile(path, data, 0600)
}

func doLogin() {
	ctx := context.Background()

	fmt.Println("\n📡 正在获取二维码...")
	c := NewUnauthenticatedClient()
	var qrResp QRCodeResponse
	if err := c.doGet(ctx, "https://ilinkai.weixin.qq.com/ilink/bot/get_bot_qrcode?bot_type=3", &qrResp); err != nil {
		fmt.Printf("❌ 获取二维码失败: %v\n", err)
		return
	}

	fmt.Println("\n📱 请用微信扫描以下二维码：")
	fmt.Println()
	qrterminal.Generate(qrResp.QRCodeImgContent, qrterminal.L, os.Stdout)
	fmt.Println()
	fmt.Printf("🔗 或用浏览器打开: %s\n", qrResp.QRCodeImgContent)
	fmt.Println("⏳ 等待扫码确认...")

	qrcode := qrResp.QRCode
	for {
		pollCtx, cancel := context.WithTimeout(ctx, 40*time.Second)
		var statusResp QRStatusResponse
		err := c.doGet(pollCtx, "https://ilinkai.weixin.qq.com/ilink/bot/get_qrcode_status?qrcode="+qrcode, &statusResp)
		cancel()

		if err != nil {
			if ctx.Err() != nil {
				fmt.Println("❌ 登录超时")
				return
			}
			continue
		}

		switch statusResp.Status {
		case "wait":
			fmt.Print(".")
		case "scaned":
			fmt.Print("\n👀 已扫码，请在微信上确认...")
		case "confirmed":
			creds = &Credentials{
				BotToken:    statusResp.BotToken,
				ILinkBotID:  statusResp.ILinkBotID,
				BaseURL:     statusResp.BaseURL,
				ILinkUserID: statusResp.ILinkUserID,
			}
			client = NewClient(creds.BotToken, creds.BaseURL)
			saveCreds()
			fmt.Printf("\n✅ 登录成功！\n")
			fmt.Printf("   Bot ID:    %s\n", creds.ILinkBotID)
			fmt.Printf("   User ID:   %s  ← 这就是你，单用户模式下所有消息发给你\n", creds.ILinkUserID)
			fmt.Printf("   Base URL:  %s\n", creds.BaseURL)
			return
		case "expired":
			fmt.Println("\n❌ 二维码已过期，请重新登录")
			return
		}
	}
}

func doStartMonitor() {
	if creds == nil {
		fmt.Println("❌ 请先登录 (命令 1)")
		return
	}
	if monitorRunning {
		fmt.Println("⚠️  接收已在运行中")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	monitorCancel = cancel
	monitorRunning = true

	go func() {
		defer func() {
			monitorRunning = false
			fmt.Println("\n🛑 消息接收已停止")
		}()

		buf := updatesBuf
		failures := 0

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			resp, err := client.GetUpdates(ctx, buf)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				failures++
				fmt.Printf("\n⚠️  GetUpdates 错误 (%d): %v\n", failures, err)
				time.Sleep(time.Duration(min(failures, 5)) * 3 * time.Second)
				continue
			}

			failures = 0

			if resp.ErrCode == -14 {
				fmt.Println("\n⚠️  Session 过期 (errcode=-14)，重置 sync buf")
				buf = ""
				time.Sleep(5 * time.Second)
				continue
			}

			if resp.Ret != 0 && resp.ErrCode != 0 {
				fmt.Printf("\n⚠️  服务器错误: ret=%d errcode=%d errmsg=%s\n", resp.Ret, resp.ErrCode, resp.ErrMsg)
				continue
			}

			if resp.GetUpdatesBuf != "" {
				buf = resp.GetUpdatesBuf
				updatesBuf = buf
			}

			for _, msg := range resp.Msgs {
				handleInboundMessage(msg)
			}
		}
	}()

	fmt.Println("✅ 开始接收消息（后台运行中）")
}

func doStopMonitor() {
	if !monitorRunning {
		fmt.Println("⚠️  没有在运行中的接收")
		return
	}
	monitorCancel()
	monitorRunning = false
	fmt.Println("✅ 已停止接收消息")
}

func handleInboundMessage(msg WeixinMessage) {
	if msg.MessageType != 1 {
		return
	}

	text := ""
	for _, item := range msg.ItemList {
		if item.Type == 1 && item.TextItem != nil {
			text = item.TextItem.Text
		}
	}

	fmt.Printf("\n📩 收到消息 from=%s text=%q context_token=%s\n",
		msg.FromUserID, text, shortToken(msg.ContextToken))

	if msg.ContextToken != "" {
		tokenStoreMu.Lock()
		currentToken = &contextTokenEntry{
			Token:      msg.ContextToken,
			UserID:     msg.FromUserID,
			ReceivedAt: time.Now(),
		}
		tokenStoreMu.Unlock()
		fmt.Printf("   💾 已更新 context_token\n")
	}
}

func doSendMessage(reader *bufio.Reader) {
	if creds == nil {
		fmt.Println("❌ 请先登录 (命令 1)")
		return
	}

	toUserID := creds.ILinkUserID

	fmt.Printf("发送给: %s (扫码用户)\n", toUserID)
	fmt.Print("消息内容: ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	if text == "" {
		fmt.Println("❌ 消息不能为空")
		return
	}

	tokenStoreMu.RLock()
	var contextToken string
	var tokenAge time.Duration
	if currentToken != nil {
		tokenAge = time.Since(currentToken.ReceivedAt)
		contextToken = currentToken.Token
	}
	tokenStoreMu.RUnlock()

	if contextToken != "" {
		fmt.Printf("📋 使用 context_token (已过 %v)\n", tokenAge.Round(time.Second))
	} else {
		fmt.Println("⚠️  没有 context_token，将不带 token 发送（可能失败）")
		fmt.Println("   💡 提示：先用微信给 Bot 发一条消息即可获取 context_token")
	}

	fmt.Printf("📤 正在发送...\n")
	resp, err := client.SendMessage(context.Background(), toUserID, text, contextToken)
	if err != nil {
		fmt.Printf("❌ 发送失败: %v\n", err)
		return
	}

	if resp.Ret != 0 {
		fmt.Printf("❌ 发送失败: ret=%d errmsg=%s\n", resp.Ret, resp.ErrMsg)
		if resp.Ret == -14 || strings.Contains(resp.ErrMsg, "expired") || strings.Contains(resp.ErrMsg, "context") {
			fmt.Println("💡 这可能是 context_token 过期导致的")
			fmt.Print("   是否尝试不带 token 重新发送？(y/n): ")
			ans, _ := reader.ReadString('\n')
			if strings.TrimSpace(strings.ToLower(ans)) == "y" {
				fmt.Println("📤 不带 token 重新发送...")
				resp2, err2 := client.SendMessage(context.Background(), toUserID, text, "")
				if err2 != nil {
					fmt.Printf("❌ 不带 token 发送也失败: %v\n", err2)
				} else if resp2.Ret != 0 {
					fmt.Printf("❌ 不带 token 发送也失败: ret=%d errmsg=%s\n", resp2.Ret, resp2.ErrMsg)
				} else {
					fmt.Println("✅ 不带 token 发送成功！说明微信允许无 context_token 推送")
				}
			}
		}
		return
	}

	fmt.Println("✅ 发送成功！")
}

func doTestTokenExpiry(reader *bufio.Reader) {
	if creds == nil {
		fmt.Println("❌ 请先登录 (命令 1)")
		return
	}

	tokenStoreMu.RLock()
	token := currentToken
	tokenStoreMu.RUnlock()

	if token == nil {
		fmt.Println("❌ 没有 context_token，请先用微信给 Bot 发一条消息")
		return
	}

	age := time.Since(token.ReceivedAt).Round(time.Second)

	fmt.Println("\n📊 当前 context_token：")
	fmt.Println("─────────────────────────────────────────────────")
	fmt.Printf("  用户:     %s\n", token.UserID)
	fmt.Printf("  Token:    %s\n", shortToken(token.Token))
	fmt.Printf("  获取时间: %s\n", token.ReceivedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("  已过:     %v\n", age)
	fmt.Printf("  状态:     %s\n", tokenAgeStatus(age))
	fmt.Println("─────────────────────────────────────────────────")

	fmt.Println()
	fmt.Println("测试选项：")
	fmt.Println("  [a] 立即测试发送（带 token）")
	fmt.Println("  [b] 立即测试发送（不带 token）")
	fmt.Println("  [c] 对比测试（先带 token，再不带 token）")
	fmt.Println("  [d] 等待 N 分钟后自动测试")
	fmt.Print("> ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(strings.ToLower(choice))

	switch choice {
	case "a":
		testSendWithToken(true)
	case "b":
		testSendWithToken(false)
	case "c":
		testCompare()
	case "d":
		testWaitAndSend(reader)
	default:
		fmt.Println("未知选项")
	}
}

func testSendWithToken(withToken bool) {
	tokenStoreMu.RLock()
	token := currentToken
	tokenStoreMu.RUnlock()

	if token == nil {
		fmt.Println("❌ 没有可用的 token")
		return
	}

	contextToken := ""
	label := "不带 token"
	if withToken {
		contextToken = token.Token
		label = fmt.Sprintf("带 token (已过 %v)", time.Since(token.ReceivedAt).Round(time.Second))
	}

	text := fmt.Sprintf("[WeClaw Demo 测试] %s - %s", label, time.Now().Format("15:04:05"))
	fmt.Printf("📤 %s 发送消息给 %s...\n", label, creds.ILinkUserID)

	resp, err := client.SendMessage(context.Background(), creds.ILinkUserID, text, contextToken)
	if err != nil {
		fmt.Printf("❌ 发送失败: %v\n", err)
		return
	}
	if resp.Ret != 0 {
		fmt.Printf("❌ 发送失败: ret=%d errmsg=%s\n", resp.Ret, resp.ErrMsg)
		return
	}
	fmt.Printf("✅ 发送成功！(%s)\n", label)
}

func testCompare() {
	tokenStoreMu.RLock()
	token := currentToken
	tokenStoreMu.RUnlock()

	if token == nil {
		fmt.Println("❌ 没有可用的 token")
		return
	}

	age := time.Since(token.ReceivedAt).Round(time.Second)

	fmt.Println("═══════════════════════════════════════════════")
	fmt.Println("  对比测试：带 token vs 不带 token")
	fmt.Println("═══════════════════════════════════════════════")

	text1 := fmt.Sprintf("[对比测试-带token] token已过%v - %s", age, time.Now().Format("15:04:05"))
	fmt.Printf("\n📤 [1/2] 带 token 发送 (token 已过 %v)...\n", age)
	resp1, err1 := client.SendMessage(context.Background(), creds.ILinkUserID, text1, token.Token)
	printResult("带 token", resp1, err1)

	time.Sleep(2 * time.Second)

	text2 := fmt.Sprintf("[对比测试-不带token] - %s", time.Now().Format("15:04:05"))
	fmt.Printf("\n📤 [2/2] 不带 token 发送...\n")
	resp2, err2 := client.SendMessage(context.Background(), creds.ILinkUserID, text2, "")
	printResult("不带 token", resp2, err2)

	fmt.Println("\n📊 结论：")
	fmt.Println("  如果两者都成功 → 微信不强制要求 context_token")
	fmt.Println("  如果只有带 token 成功 → context_token 是必须的")
	fmt.Println("  如果两者都失败 → 可能是其他问题（如 session 过期）")
}

func testWaitAndSend(reader *bufio.Reader) {
	tokenStoreMu.RLock()
	token := currentToken
	tokenStoreMu.RUnlock()

	if token == nil {
		fmt.Println("❌ 没有可用的 token")
		return
	}

	fmt.Print("等待多少分钟后测试？(如 5, 10, 30, 60): ")
	ans, _ := reader.ReadString('\n')
	var minutes int
	fmt.Sscanf(strings.TrimSpace(ans), "%d", &minutes)
	if minutes <= 0 {
		fmt.Println("❌ 无效时间")
		return
	}

	waitDuration := time.Duration(minutes) * time.Minute
	targetAge := time.Since(token.ReceivedAt) + waitDuration
	savedToken := token.Token

	fmt.Printf("⏰ 将在 %v 后自动测试 (届时 token 已过约 %v)\n", waitDuration, targetAge.Round(time.Minute))
	fmt.Println("   程序将在后台等待，你可以继续使用其他功能...")

	go func() {
		time.Sleep(waitDuration)

		tokenStoreMu.RLock()
		currentAge := time.Since(token.ReceivedAt).Round(time.Second)
		tokenStoreMu.RUnlock()

		fmt.Printf("\n🔔 定时测试触发！token 已过 %v\n", currentAge)

		text := fmt.Sprintf("[定时测试] token已过%v - %s", currentAge, time.Now().Format("15:04:05"))

		fmt.Printf("📤 带 token 发送...\n")
		resp1, err1 := client.SendMessage(context.Background(), creds.ILinkUserID, text, savedToken)
		printResult("带 token", resp1, err1)

		time.Sleep(2 * time.Second)

		fmt.Printf("📤 不带 token 发送...\n")
		text2 := fmt.Sprintf("[定时测试-无token] - %s", time.Now().Format("15:04:05"))
		resp2, err2 := client.SendMessage(context.Background(), creds.ILinkUserID, text2, "")
		printResult("不带 token", resp2, err2)

		fmt.Println("🔔 定时测试完成！")
	}()
}

func doShowTokenStatus() {
	tokenStoreMu.RLock()
	defer tokenStoreMu.RUnlock()

	if currentToken == nil {
		fmt.Println("📋 没有 context_token")
		fmt.Println("   💡 提示：用微信给 Bot 发一条消息即可获取")
		return
	}

	age := time.Since(currentToken.ReceivedAt)
	fmt.Println("\n📋 context_token 状态：")
	fmt.Println("═════════════════════════════════════════════════════")
	fmt.Printf("  用户 ID:     %s\n", currentToken.UserID)
	fmt.Printf("  Token:       %s\n", shortToken(currentToken.Token))
	fmt.Printf("  获取时间:    %s\n", currentToken.ReceivedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("  已过时长:    %v\n", age.Round(time.Second))
	fmt.Printf("  状态:        %s\n", tokenAgeStatus(age))
	fmt.Println("═════════════════════════════════════════════════════")
}

func tokenAgeStatus(age time.Duration) string {
	switch {
	case age < 5*time.Minute:
		return "🟢 很新"
	case age < 30*time.Minute:
		return "🟡 较新"
	case age < 1*time.Hour:
		return "🟠 可能过期"
	case age < 24*time.Hour:
		return "🔴 很可能过期"
	default:
		return "⛔ 大概率已过期"
	}
}

func printResult(label string, resp *SendMessageResponse, err error) {
	if err != nil {
		fmt.Printf("  ❌ %s 发送失败: %v\n", label, err)
		return
	}
	if resp.Ret != 0 {
		fmt.Printf("  ❌ %s 发送失败: ret=%d errmsg=%s\n", label, resp.Ret, resp.ErrMsg)
		return
	}
	fmt.Printf("  ✅ %s 发送成功\n", label)
}

func shortToken(token string) string {
	if token == "" {
		return "(空)"
	}
	if len(token) <= 16 {
		return token
	}
	return token[:8] + "..." + token[len(token)-8:]
}
