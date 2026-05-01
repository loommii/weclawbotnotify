package ilinkhttp

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

const defaultBaseURL = "https://ilinkai.weixin.qq.com"

type Client struct {
	baseURL    string
	botToken   string
	httpClient *http.Client
	wechatUIN  string
}

func NewClient(botToken, baseURL string) *Client {
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	return &Client{
		baseURL:    baseURL,
		botToken:   botToken,
		httpClient: &http.Client{Timeout: 45 * time.Second},
		wechatUIN:  generateWechatUIN(),
	}
}

func NewUnauthenticatedClient(baseURL string) *Client {
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 45 * time.Second},
		wechatUIN:  generateWechatUIN(),
	}
}

func (c *Client) GetQRCode(ctx context.Context, botType int32) (*QRCodeResponse, error) {
	url := fmt.Sprintf("%s/ilink/bot/get_bot_qrcode?bot_type=%d", c.baseURL, botType)
	var resp QRCodeResponse
	if err := c.doGet(ctx, url, &resp); err != nil {
		logx.Errorf("[ilink] GetQRCode failed: %v", err)
		return nil, err
	}
	logx.Infof("[ilink] GetQRCode success")
	return &resp, nil
}

func (c *Client) PollQRCodeStatus(ctx context.Context, qrcode string) (*QRStatusResponse, error) {
	url := fmt.Sprintf("%s/ilink/bot/get_qrcode_status?qrcode=%s", c.baseURL, qrcode)
	var resp QRStatusResponse
	if err := c.doGet(ctx, url, &resp); err != nil {
		logx.Errorf("[ilink] PollQRCodeStatus failed: %v", err)
		return nil, err
	}
	logx.Infof("[ilink] PollQRCodeStatus result: status=%s", resp.Status)
	return &resp, nil
}

func (c *Client) SendMessage(ctx context.Context, toUserID, text, contextToken string) (*SendMessageResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	itemList := []map[string]interface{}{
		{"type": 1, "text_item": map[string]string{"text": text}},
	}
	msg := map[string]interface{}{
		"from_user_id":  "",
		"to_user_id":    toUserID,
		"client_id":     generateClientID(),
		"message_type":  2,
		"message_state": 2,
		"item_list":     itemList,
	}
	if contextToken != "" {
		msg["context_token"] = contextToken
	}
	reqBody := map[string]interface{}{
		"msg":       msg,
		"base_info": map[string]string{"channel_version": "2.1.7"},
	}

	var resp SendMessageResponse
	if err := c.doPost(ctx, "/ilink/bot/sendmessage", reqBody, &resp); err != nil {
		logx.Errorf("[ilink] SendMessage failed: toUserID=%s, err=%v", toUserID, err)
		return nil, err
	}
	logx.Infof("[ilink] SendMessage success: toUserID=%s, ret=%d", toUserID, resp.Ret)
	return &resp, nil
}

func (c *Client) GetUpdates(ctx context.Context, buf string) (*GetUpdatesResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 40*time.Second)
	defer cancel()

	reqBody := map[string]interface{}{
		"get_updates_buf": buf,
		"base_info":       map[string]string{"channel_version": "2.1.7"},
	}

	var resp GetUpdatesResponse
	if err := c.doPost(ctx, "/ilink/bot/getupdates", reqBody, &resp); err != nil {
		logx.Errorf("[ilink] GetUpdates failed: %v", err)
		return nil, err
	}
	logx.Infof("[ilink] GetUpdates success: ret=%d, msgCount=%d", resp.Ret, len(resp.Msgs))
	return &resp, nil
}

func (c *Client) doPost(ctx context.Context, path string, body interface{}, result interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	c.setHeaders(req, data)

	logx.Debugf("[ilink] POST %s", c.baseURL+path)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		logx.Errorf("[ilink] POST %s returned HTTP %d: %s", c.baseURL+path, resp.StatusCode, string(respBody))
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	if err := json.Unmarshal(respBody, result); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}
	return nil
}

func (c *Client) doGet(ctx context.Context, url string, result interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	logx.Debugf("[ilink] GET %s", url)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		logx.Errorf("[ilink] GET %s returned HTTP %d: %s", url, resp.StatusCode, string(respBody))
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	if err := json.Unmarshal(respBody, result); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}
	return nil
}

func (c *Client) setHeaders(req *http.Request, body []byte) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("AuthorizationType", "ilink_bot_token")
	if c.botToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.botToken)
	}
	req.Header.Set("X-WECHAT-UIN", c.wechatUIN)
	req.Header.Set("iLink-App-Id", "bot")
	req.Header.Set("iLink-App-ClientVersion", "131083")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(body)))
}

func generateWechatUIN() string {
	var n uint32
	_ = binary.Read(rand.Reader, binary.LittleEndian, &n)
	s := fmt.Sprintf("%d", n)
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func generateClientID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("rpc-%x", b)
}

type GetUpdatesResponse struct {
	Ret                  int             `json:"ret"`
	ErrCode              int             `json:"errcode,omitempty"`
	ErrMsg               string          `json:"errmsg,omitempty"`
	Msgs                 []WeixinMessage `json:"msgs"`
	GetUpdatesBuf        string          `json:"get_updates_buf"`
	LongPollingTimeoutMs int             `json:"longpolling_timeout_ms,omitempty"`
}

type WeixinMessage struct {
	Seq          int           `json:"seq,omitempty"`
	MessageID    int64         `json:"message_id,omitempty"`
	FromUserID   string        `json:"from_user_id"`
	ToUserID     string        `json:"to_user_id"`
	MessageType  int           `json:"message_type"`
	MessageState int           `json:"message_state"`
	ItemList     []MessageItem `json:"item_list"`
	ContextToken string        `json:"context_token"`
}

type MessageItem struct {
	Type     int       `json:"type"`
	TextItem *TextItem `json:"text_item,omitempty"`
}

type TextItem struct {
	Text string `json:"text"`
}

type SendMessageResponse struct {
	Ret    int    `json:"ret"`
	ErrMsg string `json:"errmsg,omitempty"`
}

type QRCodeResponse struct {
	QRCode           string `json:"qrcode"`
	QRCodeImgContent string `json:"qrcode_img_content"`
}

type QRStatusResponse struct {
	Status      string `json:"status"`
	BotToken    string `json:"bot_token"`
	ILinkBotID  string `json:"ilink_bot_id"`
	BaseURL     string `json:"baseurl"`
	ILinkUserID string `json:"ilink_user_id"`
}
