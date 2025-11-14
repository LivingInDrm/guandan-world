// serializer.go - WebSocket 消息序列化适配器
//
// 职责:
// - WSMessage 类型的双向转换（Backend ↔ Proto）
// - ErrorInfo 类型的双向转换（Backend ↔ Proto）
// - WebSocket 消息的序列化和反序列化
//
// 依赖:
// - SDK proto 适配器（用于 GameEvent 转换）
// - Proto gen/go/messages
//
// 被依赖:
// - backend/websocket/manager.go: WebSocket 消息处理
package websocket

import (
	"encoding/json"
	"fmt"
	"time"

	"guandan-world/sdk"

	pbmsg "guandan-world/proto/gen/go/messages"
	"google.golang.org/protobuf/proto"
)

// ==================== Message Direction Classification ====================

// 客户端请求消息类型（Client → Server）
var clientRequestTypes = map[string]bool{
	MSG_JOIN_ROOM:      true,
	MSG_LEAVE_ROOM:     true,
	MSG_START_GAME:     true,
	MSG_PLAY_CARDS:     true,
	MSG_PASS:           true,
	MSG_TRIBUTE_SELECT: true,
	MSG_TRIBUTE_RETURN: true,
	MSG_PING:           true,
}

// ==================== WSMessage Adapters ====================

// toWSMessageProto 转换 Backend WSMessage 到 Proto WSMessage
// 复杂转换逻辑：
// - Timestamp: time.Time → int64 (毫秒)
// - Data: interface{} → oneof payload
//   - 优先根据 Data 的实际类型判断 (类型安全)
//   - GameEvent → game_event
//   - 错误 map → error
//   - 其他数据根据消息方向（Type）决定用 request 还是 response
func toWSMessageProto(m *WSMessage) (*pbmsg.WSMessage, error) {
	if m == nil {
		return nil, nil
	}

	pm := &pbmsg.WSMessage{
		Type:        m.Type,
		TimestampMs: m.Timestamp.UnixMilli(),
		PlayerId:    m.PlayerID,
	}

	// 根据 Data 的实际类型设置 payload（类型安全优先）
	switch data := m.Data.(type) {
	case *sdk.GameEvent:
		// GameEvent 类型
		pm.Payload = &pbmsg.WSMessage_GameEvent{
			GameEvent: sdk.ToProtoGameEvent(data),
		}

	case map[string]interface{}:
		// map 类型：根据 Type 判断是错误还是其他数据
		if m.Type == MSG_ERROR {
			pm.Payload = &pbmsg.WSMessage_Error{
				Error: toErrorInfoProto(data),
			}
		} else {
			// 其他 map 数据，根据消息方向决定用 request/response
			jsonData, err := json.Marshal(data)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal data: %w", err)
			}

			if clientRequestTypes[m.Type] {
				pm.Payload = &pbmsg.WSMessage_Request{
					Request: &pbmsg.WSRequest{
						Data: jsonData,
					},
				}
			} else {
				pm.Payload = &pbmsg.WSMessage_Response{
					Response: &pbmsg.WSResponse{
						Success: true,
						Data:    jsonData,
					},
				}
			}
		}

	case nil:
		// 无数据的消息（如 ping/pong）
		pm.Payload = nil

	default:
		// 其他类型的数据，序列化为 JSON bytes
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal data: %w", err)
		}

		// 根据消息方向决定用 request/response
		if clientRequestTypes[m.Type] {
			pm.Payload = &pbmsg.WSMessage_Request{
				Request: &pbmsg.WSRequest{
					Data: jsonData,
				},
			}
		} else {
			pm.Payload = &pbmsg.WSMessage_Response{
				Response: &pbmsg.WSResponse{
					Success: true,
					Data:    jsonData,
				},
			}
		}
	}

	return pm, nil
}

// fromWSMessageProto 转换 Proto WSMessage 到 Backend WSMessage
// 复杂转换逻辑：
// - TimestampMs: int64 → time.Time
// - Payload: oneof → interface{}
// - 验证 type 和 payload 的一致性
func fromWSMessageProto(pm *pbmsg.WSMessage) (*WSMessage, error) {
	if pm == nil {
		return nil, nil
	}

	m := &WSMessage{
		Type:      pm.Type,
		Timestamp: time.UnixMilli(pm.TimestampMs),
		PlayerID:  pm.PlayerId,
	}

	// 根据 payload 类型设置 Data
	switch payload := pm.Payload.(type) {
	case *pbmsg.WSMessage_GameEvent:
		// 验证 type 和 payload 的一致性
		if pm.Type != MSG_GAME_EVENT {
			return nil, fmt.Errorf("type mismatch: type=%s but payload is GameEvent", pm.Type)
		}
		// GameEvent 转换
		event := sdk.FromProtoGameEvent(payload.GameEvent)
		m.Data = event

	case *pbmsg.WSMessage_Error:
		// 验证 type 和 payload 的一致性
		if pm.Type != MSG_ERROR {
			return nil, fmt.Errorf("type mismatch: type=%s but payload is Error", pm.Type)
		}
		// 错误信息转换
		m.Data = fromErrorInfoProto(payload.Error)

	case *pbmsg.WSMessage_Request:
		// 验证 type 应该是客户端请求类型
		if !clientRequestTypes[pm.Type] && pm.Type != "" {
			// 允许空 type（兼容性）或未知请求类型（向前兼容）
			// 但记录警告日志
			// 注：这里不返回错误，因为可能是新增的请求类型
		}
		// Request 数据，反序列化 JSON
		var data interface{}
		if len(payload.Request.Data) > 0 {
			if err := json.Unmarshal(payload.Request.Data, &data); err != nil {
				return nil, fmt.Errorf("failed to unmarshal request data: %w", err)
			}
		}
		m.Data = data

	case *pbmsg.WSMessage_Response:
		// 验证 type 不应该是客户端请求类型
		if clientRequestTypes[pm.Type] {
			return nil, fmt.Errorf("type mismatch: type=%s is a request type but payload is Response", pm.Type)
		}
		// Response 数据，反序列化 JSON
		var data interface{}
		if len(payload.Response.Data) > 0 {
			if err := json.Unmarshal(payload.Response.Data, &data); err != nil {
				return nil, fmt.Errorf("failed to unmarshal response data: %w", err)
			}
		}
		m.Data = data

	case nil:
		// 无 payload（如 ping/pong）
		m.Data = nil

	default:
		return nil, fmt.Errorf("unknown payload type: %T", pm.Payload)
	}

	return m, nil
}

// ==================== ErrorInfo Adapters ====================

// toErrorInfoProto 转换 Backend 错误数据到 Proto ErrorInfo
// 输入: map[string]interface{} (包含 "message", "original_type" 等字段)
// 输出: *pbmsg.ErrorInfo
func toErrorInfoProto(errData map[string]interface{}) *pbmsg.ErrorInfo {
	errorInfo := &pbmsg.ErrorInfo{
		Details: make(map[string]string),
	}

	// 提取 message
	if msg, ok := errData["message"].(string); ok {
		errorInfo.Message = msg
	}

	// 提取 code（如果有）
	if code, ok := errData["code"].(string); ok {
		errorInfo.Code = code
	}

	// 将其他字段放入 details
	for key, value := range errData {
		if key != "message" && key != "code" {
			if strValue, ok := value.(string); ok {
				errorInfo.Details[key] = strValue
			} else {
				// 非字符串值，转换为 JSON 字符串
				if jsonValue, err := json.Marshal(value); err == nil {
					errorInfo.Details[key] = string(jsonValue)
				}
			}
		}
	}

	return errorInfo
}

// fromErrorInfoProto 转换 Proto ErrorInfo 到 Backend 错误数据
// 输入: *pbmsg.ErrorInfo
// 输出: map[string]interface{}
func fromErrorInfoProto(pe *pbmsg.ErrorInfo) map[string]interface{} {
	if pe == nil {
		return nil
	}

	errData := make(map[string]interface{})

	if pe.Code != "" {
		errData["code"] = pe.Code
	}

	if pe.Message != "" {
		errData["message"] = pe.Message
	}

	// 将 details 合并到 errData
	for key, value := range pe.Details {
		errData[key] = value
	}

	return errData
}

// ==================== Serialization Utilities ====================

// SerializeWSMessage 序列化 WSMessage 到 Protobuf bytes
func SerializeWSMessage(m *WSMessage) ([]byte, error) {
	pm, err := toWSMessageProto(m)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to proto: %w", err)
	}

	return proto.Marshal(pm)
}

// DeserializeWSMessage 反序列化 Protobuf bytes 到 WSMessage
func DeserializeWSMessage(data []byte) (*WSMessage, error) {
	var pm pbmsg.WSMessage

	if err := proto.Unmarshal(data, &pm); err != nil {
		return nil, fmt.Errorf("failed to unmarshal proto: %w", err)
	}

	return fromWSMessageProto(&pm)
}
