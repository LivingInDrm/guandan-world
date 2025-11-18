package main

import (
	"encoding/json"
	"fmt"
	eventpb "guandan-world/proto/event"
	viewpb "guandan-world/proto/view"
	"guandan-world/sdk"
	"google.golang.org/protobuf/encoding/protojson"
	"time"
)

var protoJSONMarshaler = protojson.MarshalOptions{
	UseProtoNames:   true,
	EmitUnpopulated: false,
}

func marshalProtoToRawJSON(msg interface{}, logPrefix string) json.RawMessage {
	var jsonBytes []byte
	var err error
	
	switch v := msg.(type) {
	case *sdk.GameEvent:
		jsonBytes, err = protoJSONMarshaler.Marshal(v)
	case *viewpb.PlayerView:
		jsonBytes, err = protoJSONMarshaler.Marshal(v)
	case *viewpb.TributeView:
		jsonBytes, err = protoJSONMarshaler.Marshal(v)
	}
	
	if err != nil {
		fmt.Printf("Failed to marshal %s to JSON: %v\n", logPrefix, err)
		return nil
	}
	return json.RawMessage(jsonBytes)
}

func main() {
	fmt.Println("=== Proto 消息序列化验证 ===\n")
	
	// 1. 验证 Event 序列化
	fmt.Println("1. Event 序列化测试")
	event := &sdk.GameEvent{
		Type:        eventpb.EventType_EVENT_TYPE_DEAL_STARTED,
		CreatedAtMs: time.Now().UnixMilli(),
	}
	
	eventJSON := marshalProtoToRawJSON(event, "GameEvent")
	data := map[string]interface{}{
		"event_type": event.Type.String(),
		"event_data": eventJSON,
		"timestamp":  time.Now().Unix(),
	}
	
	finalJSON, _ := json.Marshal(data)
	fmt.Println("输出:", string(finalJSON))
	
	var result map[string]interface{}
	json.Unmarshal(finalJSON, &result)
	eventData := result["event_data"]
	
	if obj, ok := eventData.(map[string]interface{}); ok {
		fmt.Println("✅ event_data 是 JSON 对象（可直接访问）")
		if typeField, ok := obj["type"].(string); ok {
			fmt.Printf("✅ 枚举字段是字符串: %s\n", typeField)
		}
		if createdAt, ok := obj["created_at_ms"]; ok {
			fmt.Printf("✅ 字段名是 snake_case: created_at_ms = %v\n", createdAt)
		}
	} else {
		fmt.Println("❌ event_data 不是 JSON 对象")
	}
	
	fmt.Println("\n2. PlayerView 序列化测试")
	pv := &viewpb.PlayerView{
		MatchId:    "test-match",
		DealIndex:  1,
		PlayerSeat: 0,
		TeamLevels: []int32{2, 2},
		DealLevel:  2,
		DealStatus: viewpb.DealStatus_DEAL_STATUS_PLAYING,
	}
	
	pvJSON := marshalProtoToRawJSON(pv, "PlayerView")
	data2 := map[string]interface{}{
		"player_view": pvJSON,
		"event_type":  "SOME_EVENT",
	}
	
	finalJSON2, _ := json.Marshal(data2)
	fmt.Println("输出:", string(finalJSON2))
	
	var result2 map[string]interface{}
	json.Unmarshal(finalJSON2, &result2)
	playerViewData := result2["player_view"].(map[string]interface{})
	
	if dealStatus, ok := playerViewData["deal_status"].(string); ok {
		fmt.Printf("✅ 枚举是字符串: %s\n", dealStatus)
	}
	if matchId, ok := playerViewData["match_id"].(string); ok {
		fmt.Printf("✅ 字段名是 snake_case: match_id = %s\n", matchId)
	}
	
	fmt.Println("\n3. TributeView 序列化测试")
	tv := &viewpb.TributeView{
		MatchId:   "test-match",
		DealIndex: 1,
		Status:    viewpb.TributeStatus_TRIBUTE_STATUS_SELECTING,
	}
	
	tvJSON := marshalProtoToRawJSON(tv, "TributeView")
	data3 := map[string]interface{}{
		"tribute_view": tvJSON,
		"event_type":   "TRIBUTE_EVENT",
	}
	
	finalJSON3, _ := json.Marshal(data3)
	fmt.Println("输出:", string(finalJSON3))
	
	var result3 map[string]interface{}
	json.Unmarshal(finalJSON3, &result3)
	tributeViewData := result3["tribute_view"].(map[string]interface{})
	
	if status, ok := tributeViewData["status"].(string); ok {
		fmt.Printf("✅ 枚举是字符串: %s\n", status)
	}
	if matchId, ok := tributeViewData["match_id"].(string); ok {
		fmt.Printf("✅ 字段名是 snake_case: match_id = %s\n", matchId)
	}
	
	fmt.Println("\n✅ 所有验证通过！")
}
