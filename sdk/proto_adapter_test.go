package sdk

import (
	"testing"
)

func TestCardAdapter(t *testing.T) {
	// 创建测试卡牌
	card, err := NewCard(14, "Heart", 5)
	if err != nil {
		t.Fatalf("Failed to create card: %v", err)
	}
	card.DeckIndex = 42

	// 测试 Go -> Proto -> Go 转换
	protoCard := ToProtoCard(card)
	if protoCard == nil {
		t.Fatal("ToProtoCard returned nil")
	}

	if protoCard.Number != int32(card.Number) {
		t.Errorf("Number mismatch: got %d, want %d", protoCard.Number, card.Number)
	}
	if protoCard.Color != card.Color {
		t.Errorf("Color mismatch: got %s, want %s", protoCard.Color, card.Color)
	}
	if protoCard.DeckIndex != int32(card.DeckIndex) {
		t.Errorf("DeckIndex mismatch: got %d, want %d", protoCard.DeckIndex, card.DeckIndex)
	}

	// 测试反向转换
	backCard := FromProtoCard(protoCard)
	if backCard == nil {
		t.Fatal("FromProtoCard returned nil")
	}

	if backCard.Number != card.Number {
		t.Errorf("Number mismatch after round-trip: got %d, want %d", backCard.Number, card.Number)
	}
	if backCard.Color != card.Color {
		t.Errorf("Color mismatch after round-trip: got %s, want %s", backCard.Color, card.Color)
	}
	if backCard.DeckIndex != card.DeckIndex {
		t.Errorf("DeckIndex mismatch after round-trip: got %d, want %d", backCard.DeckIndex, card.DeckIndex)
	}
}

func TestPlayerAdapter(t *testing.T) {
	// 创建测试玩家
	player := &Player{
		ID:       "player123",
		Username: "TestUser",
		Seat:     2,
		Online:   true,
		AutoPlay: false,
	}

	// 测试 Go -> Proto -> Go 转换
	protoPlayer := ToProtoPlayer(player)
	if protoPlayer == nil {
		t.Fatal("ToProtoPlayer returned nil")
	}

	if protoPlayer.Id != player.ID {
		t.Errorf("ID mismatch: got %s, want %s", protoPlayer.Id, player.ID)
	}
	if protoPlayer.Username != player.Username {
		t.Errorf("Username mismatch: got %s, want %s", protoPlayer.Username, player.Username)
	}
	if protoPlayer.Seat != int32(player.Seat) {
		t.Errorf("Seat mismatch: got %d, want %d", protoPlayer.Seat, player.Seat)
	}

	// 测试反向转换
	backPlayer := FromProtoPlayer(protoPlayer)
	if backPlayer == nil {
		t.Fatal("FromProtoPlayer returned nil")
	}

	if backPlayer.ID != player.ID {
		t.Errorf("ID mismatch after round-trip: got %s, want %s", backPlayer.ID, player.ID)
	}
	if backPlayer.Seat != player.Seat {
		t.Errorf("Seat mismatch after round-trip: got %d, want %d", backPlayer.Seat, player.Seat)
	}
}

func TestEnumAdapters(t *testing.T) {
	// 测试 VictoryType
	vt := VictoryTypeDoubleDown
	protoVT := ToProtoVictoryType(vt)
	backVT := FromProtoVictoryType(protoVT)
	if backVT != vt {
		t.Errorf("VictoryType round-trip failed: got %v, want %v", backVT, vt)
	}

	// 测试 DealStatus
	ds := DealStatusPlaying
	protoDS := ToProtoDealStatus(ds)
	backDS := FromProtoDealStatus(protoDS)
	if backDS != ds {
		t.Errorf("DealStatus round-trip failed: got %v, want %v", backDS, ds)
	}

	// 测试 MatchStatus
	ms := MatchStatusPlaying
	protoMS := ToProtoMatchStatus(ms)
	backMS := FromProtoMatchStatus(protoMS)
	if backMS != ms {
		t.Errorf("MatchStatus round-trip failed: got %v, want %v", backMS, ms)
	}

	// 测试 CompType
	ct := TypeStraightFlush
	protoCT := ToProtoCompType(ct)
	backCT := FromProtoCompType(protoCT)
	if backCT != ct {
		t.Errorf("CompType round-trip failed: got %v, want %v", backCT, ct)
	}

	// 测试 GameEventType
	get := EventPlayerPlayed
	protoGET := ToProtoGameEventType(get)
	backGET := FromProtoGameEventType(protoGET)
	if backGET != get {
		t.Errorf("GameEventType round-trip failed: got %v, want %v", backGET, get)
	}

	// 测试 TimeoutActionType
	tat := TimeoutActionPlayDecision
	protoTAT := ToProtoTimeoutActionType(tat)
	backTAT := FromProtoTimeoutActionType(protoTAT)
	if backTAT != tat {
		t.Errorf("TimeoutActionType round-trip failed: got %v, want %v", backTAT, tat)
	}

	// 测试所有 TimeoutActionType 值
	timeoutActions := []TimeoutActionType{
		TimeoutActionPlayDecision,
		TimeoutActionTributeSelect,
		TimeoutActionReturnTribute,
	}
	for _, ta := range timeoutActions {
		protoTA := ToProtoTimeoutActionType(ta)
		backTA := FromProtoTimeoutActionType(protoTA)
		if backTA != ta {
			t.Errorf("TimeoutActionType %v round-trip failed: got %v", ta, backTA)
		}
	}
}

// 测试 nil 处理
func TestAdapterNilHandling(t *testing.T) {
	// 测试 nil 输入
	if ToProtoCard(nil) != nil {
		t.Error("ToProtoCard(nil) should return nil")
	}
	
	if FromProtoCard(nil) != nil {
		t.Error("FromProtoCard(nil) should return nil")
	}
	
	if ToProtoPlayer(nil) != nil {
		t.Error("ToProtoPlayer(nil) should return nil")
	}
	
	if FromProtoPlayer(nil) != nil {
		t.Error("FromProtoPlayer(nil) should return nil")
	}

	// 测试 nil 切片
	if ToProtoCards(nil) != nil {
		t.Error("ToProtoCards(nil) should return nil")
	}
	
	if FromProtoCards(nil) != nil {
		t.Error("FromProtoCards(nil) should return nil")
	}
	
	if ToProtoPlayers(nil) != nil {
		t.Error("ToProtoPlayers(nil) should return nil")
	}
	
	if FromProtoPlayers(nil) != nil {
		t.Error("FromProtoPlayers(nil) should return nil")
	}

	// 测试空切片
	emptyCards := ToProtoCards([]*Card{})
	if len(emptyCards) != 0 {
		t.Errorf("ToProtoCards([]) should return empty slice, got length %d", len(emptyCards))
	}
	
	emptyPlayers := ToProtoPlayers([]*Player{})
	if len(emptyPlayers) != 0 {
		t.Errorf("ToProtoPlayers([]) should return empty slice, got length %d", len(emptyPlayers))
	}
}

// 测试批量转换
func TestBatchAdapters(t *testing.T) {
	// 创建测试卡牌列表
	cards := make([]*Card, 3)
	for i := 0; i < 3; i++ {
		card, _ := NewCard(i+2, "Heart", 5)
		card.DeckIndex = i
		cards[i] = card
	}

	// 批量转换
	protoCards := ToProtoCards(cards)
	if len(protoCards) != 3 {
		t.Errorf("ToProtoCards length mismatch: got %d, want 3", len(protoCards))
	}

	// 反向批量转换
	backCards := FromProtoCards(protoCards)
	if len(backCards) != 3 {
		t.Errorf("FromProtoCards length mismatch: got %d, want 3", len(backCards))
	}

	// 验证每张卡片
	for i := 0; i < 3; i++ {
		if backCards[i].Number != cards[i].Number {
			t.Errorf("Card %d number mismatch: got %d, want %d", i, backCards[i].Number, cards[i].Number)
		}
	}
}

// 测试数组转换
func TestArrayAdapters(t *testing.T) {
	// 创建4个玩家
	players := [4]*Player{
		{ID: "p0", Username: "Player0", Seat: 0, Online: true, AutoPlay: false},
		{ID: "p1", Username: "Player1", Seat: 1, Online: true, AutoPlay: false},
		{ID: "p2", Username: "Player2", Seat: 2, Online: false, AutoPlay: true},
		{ID: "p3", Username: "Player3", Seat: 3, Online: true, AutoPlay: false},
	}

	// 数组转换
	protoPlayers := ToProtoPlayersArray(players)
	if len(protoPlayers) != 4 {
		t.Errorf("ToProtoPlayersArray length mismatch: got %d, want 4", len(protoPlayers))
	}

	// 反向转换
	backPlayers := FromProtoPlayersArray(protoPlayers)
	
	// 验证每个玩家
	for i := 0; i < 4; i++ {
		if backPlayers[i].ID != players[i].ID {
			t.Errorf("Player %d ID mismatch: got %s, want %s", i, backPlayers[i].ID, players[i].ID)
		}
		if backPlayers[i].Seat != players[i].Seat {
			t.Errorf("Player %d Seat mismatch: got %d, want %d", i, backPlayers[i].Seat, players[i].Seat)
		}
	}
}

// TestCardCompNormalizedCardsRoundTrip 测试 CardComp 的 NormalizedCards 字段在序列化和反序列化后保持一致
func TestCardCompNormalizedCardsRoundTrip(t *testing.T) {
	// 创建一个普通对子（无万能牌）
	card1, _ := NewCard(10, "Heart", 2)
	card2, _ := NewCard(10, "Spade", 2)
	pair := NewPair([]*Card{card1, card2})

	if !pair.IsValid() {
		t.Fatalf("Pair should be valid")
	}

	// 序列化到 Proto
	protoPair := ToProtoCardComp(pair)
	if protoPair == nil {
		t.Fatalf("ToProtoCardComp returned nil")
	}

	// 反序列化回 SDK
	restoredPair := FromProtoCardComp(protoPair)
	if restoredPair == nil {
		t.Fatalf("FromProtoCardComp returned nil")
	}

	// 验证类型
	restoredPairTyped, ok := restoredPair.(*Pair)
	if !ok {
		t.Fatalf("Expected *Pair, got %T", restoredPair)
	}

	// 验证 Cards 一致
	if len(restoredPairTyped.Cards) != len(pair.Cards) {
		t.Errorf("Cards length mismatch: got %d, want %d", len(restoredPairTyped.Cards), len(pair.Cards))
	}

	// 验证 NormalizedCards 一致（对于无万能牌的对子，NormalizedCards 可能为 nil）
	if pair.NormalizedCards == nil && restoredPairTyped.NormalizedCards != nil && len(restoredPairTyped.NormalizedCards) > 0 {
		t.Errorf("NormalizedCards mismatch: original is nil, restored has %d cards", len(restoredPairTyped.NormalizedCards))
	}
	if pair.NormalizedCards != nil && restoredPairTyped.NormalizedCards == nil {
		t.Errorf("NormalizedCards mismatch: original has %d cards, restored is nil", len(pair.NormalizedCards))
	}
	if pair.NormalizedCards != nil && restoredPairTyped.NormalizedCards != nil {
		if len(restoredPairTyped.NormalizedCards) != len(pair.NormalizedCards) {
			t.Errorf("NormalizedCards length mismatch: got %d, want %d",
				len(restoredPairTyped.NormalizedCards), len(pair.NormalizedCards))
		}
	}

	// 验证 Valid 一致
	if restoredPairTyped.Valid != pair.Valid {
		t.Errorf("Valid mismatch: got %v, want %v", restoredPairTyped.Valid, pair.Valid)
	}

	// 验证 Type 一致
	if restoredPairTyped.Type != pair.Type {
		t.Errorf("Type mismatch: got %v, want %v", restoredPairTyped.Type, pair.Type)
	}

	t.Logf("✅ Round-trip test passed for Pair with NormalizedCards field")
}

// TestCardCompAllTypesRoundTrip 测试所有牌型的序列化和反序列化
func TestCardCompAllTypesRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		comp CardComp
	}{
		{
			name: "Fold",
			comp: &Fold{BaseComp: BaseComp{Cards: []*Card{}, Valid: true, Type: TypeFold}},
		},
		{
			name: "Single",
			comp: NewSingle([]*Card{mustNewCard(10, "Heart", 2)}),
		},
		{
			name: "Pair",
			comp: NewPair([]*Card{mustNewCard(10, "Heart", 2), mustNewCard(10, "Spade", 2)}),
		},
		{
			name: "Triple",
			comp: NewTriple([]*Card{mustNewCard(10, "Heart", 2), mustNewCard(10, "Spade", 2), mustNewCard(10, "Club", 2)}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 序列化
			proto := ToProtoCardComp(tt.comp)
			if proto == nil {
				t.Fatalf("ToProtoCardComp returned nil for %s", tt.name)
			}

			// 反序列化
			restored := FromProtoCardComp(proto)
			if restored == nil {
				t.Fatalf("FromProtoCardComp returned nil for %s", tt.name)
			}

			// 验证类型
			if restored.GetType() != tt.comp.GetType() {
				t.Errorf("Type mismatch: got %v, want %v", restored.GetType(), tt.comp.GetType())
			}

			// 验证有效性
			if restored.IsValid() != tt.comp.IsValid() {
				t.Errorf("Valid mismatch: got %v, want %v", restored.IsValid(), tt.comp.IsValid())
			}

			// 验证牌数
			if len(restored.GetCards()) != len(tt.comp.GetCards()) {
				t.Errorf("Cards count mismatch: got %d, want %d",
					len(restored.GetCards()), len(tt.comp.GetCards()))
			}

			t.Logf("✅ %s round-trip successful", tt.name)
		})
	}
}

// 辅助函数：创建卡片，出错时 panic
func mustNewCard(number int, color string, level int) *Card {
	card, err := NewCard(number, color, level)
	if err != nil {
		panic(err)
	}
	return card
}

// TestConnectionEventsAdapter 测试连接事件的 Proto 适配器
func TestConnectionEventsAdapter(t *testing.T) {
	tests := []struct {
		name      string
		eventType GameEventType
		data      map[string]interface{}
		validate  func(t *testing.T, data map[string]interface{})
	}{
		{
			name:      "PlayerTimeout_PlayDecision",
			eventType: EventPlayerTimeout,
			data: map[string]interface{}{
				"action": "play_decision",
			},
			validate: func(t *testing.T, data map[string]interface{}) {
				if action, ok := data["action"].(string); !ok || action != "play_decision" {
					t.Errorf("action mismatch: got %v, want play_decision", data["action"])
				}
			},
		},
		{
			name:      "PlayerTimeout_TributeSelect",
			eventType: EventPlayerTimeout,
			data: map[string]interface{}{
				"action": "tribute_select",
			},
			validate: func(t *testing.T, data map[string]interface{}) {
				if action, ok := data["action"].(string); !ok || action != "tribute_select" {
					t.Errorf("action mismatch: got %v, want tribute_select", data["action"])
				}
			},
		},
		{
			name:      "PlayerTimeout_ReturnTribute",
			eventType: EventPlayerTimeout,
			data: map[string]interface{}{
				"action": "return_tribute",
			},
			validate: func(t *testing.T, data map[string]interface{}) {
				if action, ok := data["action"].(string); !ok || action != "return_tribute" {
					t.Errorf("action mismatch: got %v, want return_tribute", data["action"])
				}
			},
		},
		{
			name:      "PlayerDisconnect",
			eventType: EventPlayerDisconnect,
			data: map[string]interface{}{
				"player_seat": 2,
				"auto_play":   true,
			},
			validate: func(t *testing.T, data map[string]interface{}) {
				if autoPlay, ok := data["auto_play"].(bool); !ok || autoPlay != true {
					t.Errorf("auto_play mismatch: got %v, want true", data["auto_play"])
				}
			},
		},
		{
			name:      "PlayerReconnect",
			eventType: EventPlayerReconnect,
			data: map[string]interface{}{
				"player_seat": 2,
				"auto_play":   false,
			},
			validate: func(t *testing.T, data map[string]interface{}) {
				if autoPlay, ok := data["auto_play"].(bool); !ok || autoPlay != false {
					t.Errorf("auto_play mismatch: got %v, want false", data["auto_play"])
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建原始事件
			event := &GameEvent{
				Type:       tt.eventType,
				PlayerSeat: 2,
				Data:       tt.data,
			}

			// Go -> Proto
			protoEvent := ToProtoGameEvent(event)
			if protoEvent == nil {
				t.Fatal("ToProtoGameEvent returned nil")
			}

			if protoEvent.PlayerSeat != 2 {
				t.Errorf("PlayerSeat mismatch: got %d, want 2", protoEvent.PlayerSeat)
			}

			// Proto -> Go
			restored := FromProtoGameEvent(protoEvent)
			if restored == nil {
				t.Fatal("FromProtoGameEvent returned nil")
			}

			if restored.Type != tt.eventType {
				t.Errorf("Type mismatch: got %v, want %v", restored.Type, tt.eventType)
			}

			if restored.PlayerSeat != 2 {
				t.Errorf("PlayerSeat mismatch: got %d, want 2", restored.PlayerSeat)
			}

			// 验证 Data 字段
			data, ok := restored.Data.(map[string]interface{})
			if !ok {
				t.Fatal("Data is not a map")
			}

			// 使用自定义验证函数
			tt.validate(t, data)

			t.Logf("✅ %s round-trip successful", tt.name)
		})
	}
}
