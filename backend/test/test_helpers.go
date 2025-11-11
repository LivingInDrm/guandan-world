package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"guandan-world/backend/handlers"
	"guandan-world/sdk"
)

// HTTPClient 封装HTTP客户端功能
type HTTPClient struct {
	serverURL  string
	authToken  string
	httpClient *http.Client
}

// NewHTTPClient 创建HTTP客户端
func NewHTTPClient(serverURL string, httpClient *http.Client) *HTTPClient {
	return &HTTPClient{
		serverURL:  serverURL,
		httpClient: httpClient,
	}
}

// SetAuthToken 设置认证token
func (c *HTTPClient) SetAuthToken(token string) {
	c.authToken = token
}

// CallAPI 调用HTTP API
func (c *HTTPClient) CallAPI(method, path string, body interface{}, result interface{}) error {
	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request: %w", err)
		}
	}

	url := fmt.Sprintf("http://%s%s", c.serverURL, path)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	if c.authToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.authToken)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var errResp handlers.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return fmt.Errorf("HTTP %d: failed to parse error response", resp.StatusCode)
		}
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, errResp.Error)
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}
	}

	return nil
}

// ParseCard 解析卡牌数据（从WebSocket消息）
func ParseCard(cardMap map[string]interface{}, currentRank int) *sdk.Card {
	cardID, ok := cardMap["id"].(string)
	if !ok {
		return nil
	}

	card, err := sdk.ParseCardFromID(cardID, currentRank)
	if err != nil {
		return nil
	}

	return card
}

// ParseCards 解析卡牌数组
func ParseCards(cardsData []interface{}, currentRank int) []*sdk.Card {
	cards := make([]*sdk.Card, 0, len(cardsData))
	for _, cardData := range cardsData {
		cardMap, ok := cardData.(map[string]interface{})
		if !ok {
			continue
		}
		card := ParseCard(cardMap, currentRank)
		if card != nil {
			cards = append(cards, card)
		}
	}
	return cards
}

// SelectMaxCard 从卡牌列表中选择最大的牌
func SelectMaxCard(cards []*sdk.Card) *sdk.Card {
	if len(cards) == 0 {
		return nil
	}

	maxCard := cards[0]
	for _, card := range cards[1:] {
		if card.GreaterThan(maxCard) {
			maxCard = card
		}
	}
	return maxCard
}

// SelectMinCard 从卡牌列表中选择最小的牌
func SelectMinCard(cards []*sdk.Card) *sdk.Card {
	if len(cards) == 0 {
		return nil
	}

	minCard := cards[0]
	for _, card := range cards[1:] {
		if card.LessThan(minCard) {
			minCard = card
		}
	}
	return minCard
}

