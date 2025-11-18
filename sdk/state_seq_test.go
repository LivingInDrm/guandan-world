package sdk

import (
	"sync/atomic"
	"testing"
)

// TestStateSeqBehavior 验证状态版本号的行为
func TestStateSeqBehavior(t *testing.T) {
	engine := NewGameEngine("test-engine")

	// 初始状态，currentStateSeq应该为0
	initialSeq := atomic.LoadInt64(&engine.currentStateSeq)
	if initialSeq != 0 {
		t.Errorf("初始 currentStateSeq 应该为 0, got %d", initialSeq)
	}

	// 创建玩家
	players := []Player{
		{ID: "p1", Username: "Player 1", Seat: 0},
		{ID: "p2", Username: "Player 2", Seat: 1},
		{ID: "p3", Username: "Player 3", Seat: 2},
		{ID: "p4", Username: "Player 4", Seat: 3},
	}

	// 开始比赛，这会触发事件发射
	err := engine.StartMatch(players)
	if err != nil {
		t.Fatalf("开始比赛失败: %v", err)
	}

	// 开始比赛后，currentStateSeq应该已更新
	afterMatchSeq := atomic.LoadInt64(&engine.currentStateSeq)
	if afterMatchSeq == 0 {
		t.Error("开始比赛后 currentStateSeq 应该已更新")
	}

	// 获取玩家视图
	view1 := engine.GetPlayerView(0)
	if view1 == nil {
		t.Fatal("获取玩家视图失败")
	}

	firstViewSeq := view1.Seq
	if firstViewSeq != afterMatchSeq {
		t.Errorf("视图 seq 应该等于当前状态 seq: view.Seq=%d, currentStateSeq=%d", firstViewSeq, afterMatchSeq)
	}

	// 再次获取相同的视图，seq应该相同（因为状态未变化）
	view2 := engine.GetPlayerView(0)
	if view2 == nil {
		t.Fatal("第二次获取玩家视图失败")
	}

	secondViewSeq := view2.Seq
	if secondViewSeq != firstViewSeq {
		t.Errorf("相同状态下视图 seq 应该相同: first=%d, second=%d", firstViewSeq, secondViewSeq)
	}

	// 开始发牌，这会触发新的事件
	err = engine.StartDeal()
	if err != nil {
		t.Fatalf("开始发牌失败: %v", err)
	}

	// 状态改变后，currentStateSeq应该已更新
	afterDealSeq := atomic.LoadInt64(&engine.currentStateSeq)
	if afterDealSeq <= afterMatchSeq {
		t.Errorf("发牌后 currentStateSeq 应该增加: before=%d, after=%d", afterMatchSeq, afterDealSeq)
	}

	// 获取新的视图，seq应该反映最新状态
	view3 := engine.GetPlayerView(0)
	if view3 == nil {
		t.Fatal("发牌后获取玩家视图失败")
	}

	thirdViewSeq := view3.Seq
	if thirdViewSeq != afterDealSeq {
		t.Errorf("视图 seq 应该反映最新状态: view.Seq=%d, currentStateSeq=%d", thirdViewSeq, afterDealSeq)
	}

	if thirdViewSeq <= secondViewSeq {
		t.Errorf("状态变化后视图 seq 应该增加: before=%d, after=%d", secondViewSeq, thirdViewSeq)
	}
}

// TestMultipleViewCallsSameSeq 验证多次调用GetPlayerView在状态不变时返回相同seq
func TestMultipleViewCallsSameSeq(t *testing.T) {
	engine := NewGameEngine("test-engine")

	players := []Player{
		{ID: "p1", Username: "Player 1", Seat: 0},
		{ID: "p2", Username: "Player 2", Seat: 1},
		{ID: "p3", Username: "Player 3", Seat: 2},
		{ID: "p4", Username: "Player 4", Seat: 3},
	}

	err := engine.StartMatch(players)
	if err != nil {
		t.Fatalf("开始比赛失败: %v", err)
	}

	err = engine.StartDeal()
	if err != nil {
		t.Fatalf("开始发牌失败: %v", err)
	}

	// 获取多次视图，所有seq应该相同
	seqs := make([]int64, 10)
	for i := 0; i < 10; i++ {
		view := engine.GetPlayerView(0)
		if view == nil {
			t.Fatalf("第 %d 次获取视图失败", i+1)
		}
		seqs[i] = view.Seq
	}

	// 验证所有seq相同
	firstSeq := seqs[0]
	for i, seq := range seqs {
		if seq != firstSeq {
			t.Errorf("第 %d 次获取的 seq 不一致: expected=%d, got=%d", i+1, firstSeq, seq)
		}
	}
}
