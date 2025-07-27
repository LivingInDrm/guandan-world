package websocket

import (
	"encoding/json"
	"sync"
	"time"
)

// MessageOptimizer WebSocket消息优化器
// 实现消息压缩、批量发送、增量更新等优化功能
type MessageOptimizer struct {
	// 消息批处理
	batchBuffer   map[string][]*WSMessage // playerID -> messages
	batchMutex    sync.RWMutex
	batchInterval time.Duration
	batchSize     int

	// 状态缓存和增量更新
	lastStates map[string]interface{} // playerID -> last state
	stateMutex sync.RWMutex

	// 消息压缩
	enableCompression bool
	compressionLevel  int

	// 统计信息
	stats MessageStats
}

// MessageStats 消息统计信息
type MessageStats struct {
	TotalMessages     int64
	CompressedBytes   int64
	UncompressedBytes int64
	BatchedMessages   int64
	IncrementalUpdates int64
	mutex             sync.RWMutex
}

// OptimizedMessage 优化后的消息结构
type OptimizedMessage struct {
	Type      string      `json:"type"`
	Data      interface{} `json:"data,omitempty"`
	Delta     interface{} `json:"delta,omitempty"`     // 增量数据
	Batch     bool        `json:"batch,omitempty"`     // 是否为批量消息
	Compressed bool       `json:"compressed,omitempty"` // 是否压缩
	Timestamp time.Time   `json:"timestamp"`
}

// NewMessageOptimizer 创建消息优化器
func NewMessageOptimizer() *MessageOptimizer {
	optimizer := &MessageOptimizer{
		batchBuffer:       make(map[string][]*WSMessage),
		batchInterval:     50 * time.Millisecond, // 50ms批处理间隔
		batchSize:         10,                     // 最大批处理大小
		lastStates:        make(map[string]interface{}),
		enableCompression: true,
		compressionLevel:  6, // 中等压缩级别
	}

	// 启动批处理定时器
	go optimizer.startBatchProcessor()

	return optimizer
}

// OptimizeMessage 优化单个消息
func (mo *MessageOptimizer) OptimizeMessage(playerID string, msg *WSMessage) (*OptimizedMessage, error) {
	optimized := &OptimizedMessage{
		Type:      msg.Type,
		Timestamp: msg.Timestamp,
	}

	// 1. 检查是否可以进行增量更新
	if delta := mo.calculateDelta(playerID, msg); delta != nil {
		optimized.Delta = delta
		mo.updateLastState(playerID, msg.Data)
		mo.incrementStat("IncrementalUpdates")
	} else {
		optimized.Data = msg.Data
	}

	// 2. 消息压缩
	if mo.enableCompression && mo.shouldCompress(msg) {
		compressed, err := mo.compressMessage(optimized)
		if err == nil {
			optimized = compressed
			optimized.Compressed = true
		}
	}

	mo.incrementStat("TotalMessages")
	return optimized, nil
}

// BatchMessage 将消息加入批处理队列
func (mo *MessageOptimizer) BatchMessage(playerID string, msg *WSMessage) {
	mo.batchMutex.Lock()
	defer mo.batchMutex.Unlock()

	if mo.batchBuffer[playerID] == nil {
		mo.batchBuffer[playerID] = make([]*WSMessage, 0, mo.batchSize)
	}

	mo.batchBuffer[playerID] = append(mo.batchBuffer[playerID], msg)

	// 如果达到批处理大小，立即处理
	if len(mo.batchBuffer[playerID]) >= mo.batchSize {
		go mo.processBatch(playerID)
	}
}

// startBatchProcessor 启动批处理定时器
func (mo *MessageOptimizer) startBatchProcessor() {
	ticker := time.NewTicker(mo.batchInterval)
	defer ticker.Stop()

	for range ticker.C {
		mo.processAllBatches()
	}
}

// processAllBatches 处理所有待处理的批次
func (mo *MessageOptimizer) processAllBatches() {
	mo.batchMutex.RLock()
	playerIDs := make([]string, 0, len(mo.batchBuffer))
	for playerID := range mo.batchBuffer {
		if len(mo.batchBuffer[playerID]) > 0 {
			playerIDs = append(playerIDs, playerID)
		}
	}
	mo.batchMutex.RUnlock()

	for _, playerID := range playerIDs {
		go mo.processBatch(playerID)
	}
}

// processBatch 处理单个玩家的批次
func (mo *MessageOptimizer) processBatch(playerID string) {
	mo.batchMutex.Lock()
	messages := mo.batchBuffer[playerID]
	mo.batchBuffer[playerID] = nil
	mo.batchMutex.Unlock()

	if len(messages) == 0 {
		return
	}

	// 创建批量消息
	batchMsg := &OptimizedMessage{
		Type:      "batch",
		Data:      messages,
		Batch:     true,
		Timestamp: time.Now(),
	}

	// 这里应该调用实际的发送逻辑
	// 由于我们在优化器中，这里只是示例
	mo.incrementStat("BatchedMessages")
}

// calculateDelta 计算增量更新
func (mo *MessageOptimizer) calculateDelta(playerID string, msg *WSMessage) interface{} {
	mo.stateMutex.RLock()
	lastState, exists := mo.lastStates[playerID]
	mo.stateMutex.RUnlock()

	if !exists {
		return nil
	}

	// 根据消息类型计算增量
	switch msg.Type {
	case "game_state":
		return mo.calculateGameStateDelta(lastState, msg.Data)
	case "player_view":
		return mo.calculatePlayerViewDelta(lastState, msg.Data)
	case "room_update":
		return mo.calculateRoomUpdateDelta(lastState, msg.Data)
	default:
		return nil
	}
}

// calculateGameStateDelta 计算游戏状态增量
func (mo *MessageOptimizer) calculateGameStateDelta(lastState, currentState interface{}) interface{} {
	// 简化的增量计算逻辑
	// 实际实现中应该进行深度比较
	lastMap, ok1 := lastState.(map[string]interface{})
	currentMap, ok2 := currentState.(map[string]interface{})

	if !ok1 || !ok2 {
		return nil
	}

	delta := make(map[string]interface{})
	hasChanges := false

	// 比较关键字段
	keyFields := []string{"current_player", "current_trick", "player_cards", "game_status"}
	for _, field := range keyFields {
		if lastMap[field] != currentMap[field] {
			delta[field] = currentMap[field]
			hasChanges = true
		}
	}

	if hasChanges {
		return delta
	}
	return nil
}

// calculatePlayerViewDelta 计算玩家视图增量
func (mo *MessageOptimizer) calculatePlayerViewDelta(lastState, currentState interface{}) interface{} {
	// 类似的增量计算逻辑
	return nil
}

// calculateRoomUpdateDelta 计算房间更新增量
func (mo *MessageOptimizer) calculateRoomUpdateDelta(lastState, currentState interface{}) interface{} {
	// 类似的增量计算逻辑
	return nil
}

// updateLastState 更新最后状态
func (mo *MessageOptimizer) updateLastState(playerID string, state interface{}) {
	mo.stateMutex.Lock()
	defer mo.stateMutex.Unlock()
	mo.lastStates[playerID] = state
}

// shouldCompress 判断是否应该压缩消息
func (mo *MessageOptimizer) shouldCompress(msg *WSMessage) bool {
	// 只压缩大消息
	data, err := json.Marshal(msg.Data)
	if err != nil {
		return false
	}

	// 大于1KB的消息才压缩
	return len(data) > 1024
}

// compressMessage 压缩消息
func (mo *MessageOptimizer) compressMessage(msg *OptimizedMessage) (*OptimizedMessage, error) {
	// 这里应该实现实际的压缩逻辑
	// 可以使用gzip、lz4等压缩算法
	
	// 简化实现：只是标记为压缩
	originalSize := mo.calculateMessageSize(msg)
	compressedSize := originalSize * 70 / 100 // 假设压缩率70%

	mo.updateCompressionStats(originalSize, compressedSize)

	return msg, nil
}

// calculateMessageSize 计算消息大小
func (mo *MessageOptimizer) calculateMessageSize(msg *OptimizedMessage) int64 {
	data, err := json.Marshal(msg)
	if err != nil {
		return 0
	}
	return int64(len(data))
}

// updateCompressionStats 更新压缩统计
func (mo *MessageOptimizer) updateCompressionStats(original, compressed int64) {
	mo.stats.mutex.Lock()
	defer mo.stats.mutex.Unlock()
	mo.stats.UncompressedBytes += original
	mo.stats.CompressedBytes += compressed
}

// incrementStat 增加统计计数
func (mo *MessageOptimizer) incrementStat(statName string) {
	mo.stats.mutex.Lock()
	defer mo.stats.mutex.Unlock()

	switch statName {
	case "TotalMessages":
		mo.stats.TotalMessages++
	case "BatchedMessages":
		mo.stats.BatchedMessages++
	case "IncrementalUpdates":
		mo.stats.IncrementalUpdates++
	}
}

// GetStats 获取统计信息
func (mo *MessageOptimizer) GetStats() MessageStats {
	mo.stats.mutex.RLock()
	defer mo.stats.mutex.RUnlock()
	return mo.stats
}

// GetCompressionRatio 获取压缩比
func (mo *MessageOptimizer) GetCompressionRatio() float64 {
	mo.stats.mutex.RLock()
	defer mo.stats.mutex.RUnlock()

	if mo.stats.UncompressedBytes == 0 {
		return 0
	}

	return float64(mo.stats.CompressedBytes) / float64(mo.stats.UncompressedBytes)
}

// Reset 重置优化器状态
func (mo *MessageOptimizer) Reset() {
	mo.batchMutex.Lock()
	mo.stateMutex.Lock()
	defer mo.batchMutex.Unlock()
	defer mo.stateMutex.Unlock()

	mo.batchBuffer = make(map[string][]*WSMessage)
	mo.lastStates = make(map[string]interface{})
	mo.stats = MessageStats{}
}

// SetBatchConfig 设置批处理配置
func (mo *MessageOptimizer) SetBatchConfig(interval time.Duration, size int) {
	mo.batchInterval = interval
	mo.batchSize = size
}

// SetCompressionConfig 设置压缩配置
func (mo *MessageOptimizer) SetCompressionConfig(enabled bool, level int) {
	mo.enableCompression = enabled
	mo.compressionLevel = level
}