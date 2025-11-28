import React, { useState, useCallback } from 'react';
import { useTributeStore } from '../../store/tributeStore';
import { useTributeData } from '../../hooks/useTributeData';
import { useRoomStore } from '../../store/roomStore';
import TributeBoard from '../game/tribute/TributeBoard';
import type { Card, Room } from '../../types';
import { RoomStatus } from '../../types';
import {
  createMockPlayers,
  createMockPlayerHand,
  createMockTributeStarted,
  createMockTributeExempted,
  createMockCardReturned,
  MOCK_TRIBUTE_CARDS,
  createMockCard,
} from '../../test/utils/mockTributeData';

const MOCK_ROOM: Room = {
  id: 'test-room-001',
  status: RoomStatus.PLAYING,
  players: createMockPlayers(),
  owner: 'p1',
  created_at: new Date().toISOString(),
};

const TributeTestPage: React.FC = () => {
  const [currentPlayerSeat, setCurrentPlayerSeat] = useState(2);
  const [playerHand, setPlayerHand] = useState<Card[]>(createMockPlayerHand());
  const [selectedCards, setSelectedCards] = useState<Card[]>([]);
  const [showStatePanel, setShowStatePanel] = useState(false);

  const tributeState = useTributeStore();
  const tributeData = useTributeData();
  const setCurrentRoom = useRoomStore((s) => s.setCurrentRoom);

  React.useEffect(() => {
    setCurrentRoom(MOCK_ROOM);
    return () => {
      tributeState.reset();
      setCurrentRoom(null);
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const handleDoubleDown = useCallback(() => {
    tributeState.reset();
    const payload = createMockTributeStarted({ type: 'DOUBLE_DOWN' });
    tributeState.handleTributeStarted(payload);
  }, [tributeState]);

  const handleSingleLast = useCallback(() => {
    tributeState.reset();
    const payload = createMockTributeStarted({ type: 'SINGLE_LAST' });
    tributeState.handleTributeStarted(payload);
  }, [tributeState]);

  const handlePartnerLast = useCallback(() => {
    tributeState.reset();
    const payload = createMockTributeStarted({ type: 'PARTNER_LAST' });
    tributeState.handleTributeStarted(payload);
  }, [tributeState]);

  const handleExemption = useCallback(() => {
    tributeState.reset();
    const payload = createMockTributeStarted({ type: 'DOUBLE_DOWN' });
    tributeState.handleTributeStarted(payload);
    setTimeout(() => {
      const exemptedPayload = createMockTributeExempted({ 0: 2 });
      tributeState.handleTributeExempted(exemptedPayload);
    }, 500);
  }, [tributeState]);

  const handleReset = useCallback(() => {
    tributeState.reset();
    setSelectedCards([]);
  }, [tributeState]);

  const handleSubmitCard = useCallback((seat: number, cardIndex: number) => {
    const card = MOCK_TRIBUTE_CARDS.doubleDown[cardIndex] || createMockCard(cardIndex, 14, cardIndex);
    tributeState.handleCardSubmitted(seat, card);
  }, [tributeState]);

  const handleSelectTribute = useCallback((deckIndex: number) => {
    const poolCard = tributeState.poolCards.find((c) => c !== null && c.deckIndex === deckIndex);
    if (poolCard && tributeState.currentSelectingSeat !== null) {
      tributeState.handleCardSelected(tributeState.currentSelectingSeat, poolCard);
    }
  }, [tributeState]);

  const handleReturnTribute = useCallback((deckIndex: number) => {
    const card = playerHand.find((c) => c.deckIndex === deckIndex);
    if (!card) return;

    const receivers = tributeState.tributeStarted?.receivers || [];
    const givers = tributeState.tributeStarted?.givers || [];
    const receiverIndex = receivers.indexOf(currentPlayerSeat);
    const targetGiver = givers[receiverIndex];

    if (targetGiver !== undefined) {
      const payload = createMockCardReturned(targetGiver, card);
      tributeState.handleCardReturned(currentPlayerSeat, payload);
      setPlayerHand((prev) => prev.filter((c) => c.deckIndex !== deckIndex));
      setSelectedCards([]);
    }
  }, [tributeState, playerHand, currentPlayerSeat]);

  const handleComplete = useCallback(() => {
    tributeState.handleCompleted();
  }, [tributeState]);

  const handleAutoSubmitAll = useCallback(() => {
    const givers = tributeState.tributeStarted?.givers || [];
    givers.forEach((giver, idx) => {
      setTimeout(() => {
        handleSubmitCard(giver, idx);
      }, idx * 600);
    });
  }, [tributeState, handleSubmitCard]);

  return (
    <div className="min-h-screen bg-gray-100 p-6">
      <div className="max-w-7xl mx-auto space-y-6">
        <div className="bg-white rounded-lg shadow p-6">
          <h1 className="text-2xl font-bold mb-4">Tribute 测试页面</h1>
          
          <div className="flex flex-wrap gap-4 mb-6">
            <div className="flex items-center gap-2">
              <span className="text-sm text-gray-600">当前视角:</span>
              <select
                value={currentPlayerSeat}
                onChange={(e) => setCurrentPlayerSeat(Number(e.target.value))}
                className="border rounded px-2 py-1"
              >
                {[0, 1, 2, 3].map((seat) => (
                  <option key={seat} value={seat}>
                    座位 {seat} ({MOCK_ROOM.players[seat]?.username})
                  </option>
                ))}
              </select>
            </div>
          </div>

          <div className="space-y-4">
            <div>
              <h3 className="text-sm font-medium text-gray-700 mb-2">场景预设</h3>
              <div className="flex flex-wrap gap-2">
                <button
                  onClick={handleDoubleDown}
                  className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
                >
                  双下
                </button>
                <button
                  onClick={handleSingleLast}
                  className="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600"
                >
                  单下
                </button>
                <button
                  onClick={handlePartnerLast}
                  className="px-4 py-2 bg-purple-500 text-white rounded hover:bg-purple-600"
                >
                  末游
                </button>
                <button
                  onClick={handleExemption}
                  className="px-4 py-2 bg-yellow-500 text-white rounded hover:bg-yellow-600"
                >
                  抗贡
                </button>
                <button
                  onClick={handleReset}
                  className="px-4 py-2 bg-gray-500 text-white rounded hover:bg-gray-600"
                >
                  重置
                </button>
              </div>
            </div>

            <div>
              <h3 className="text-sm font-medium text-gray-700 mb-2">手动步骤控制</h3>
              <div className="flex flex-wrap gap-2">
                <button
                  onClick={handleAutoSubmitAll}
                  disabled={tributeState.step !== 'started' && tributeState.step !== 'submitting'}
                  className="px-3 py-1 bg-orange-500 text-white rounded text-sm hover:bg-orange-600 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  自动提交贡牌
                </button>
                <button
                  onClick={() => handleSubmitCard(0, 0)}
                  disabled={tributeState.step !== 'started' && tributeState.step !== 'submitting'}
                  className="px-3 py-1 bg-red-400 text-white rounded text-sm hover:bg-red-500 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  座位0提交
                </button>
                <button
                  onClick={() => handleSubmitCard(1, 1)}
                  disabled={tributeState.step !== 'started' && tributeState.step !== 'submitting'}
                  className="px-3 py-1 bg-red-400 text-white rounded text-sm hover:bg-red-500 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  座位1提交
                </button>
                <button
                  onClick={handleComplete}
                  disabled={tributeState.step !== 'returning'}
                  className="px-3 py-1 bg-green-600 text-white rounded text-sm hover:bg-green-700 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  完成
                </button>
              </div>
            </div>

            <div className="flex items-center gap-4 text-sm">
              <span className="text-gray-600">当前步骤:</span>
              <span className="px-2 py-1 bg-blue-100 text-blue-800 rounded font-mono">
                {tributeState.step}
              </span>
              <span className="text-gray-600">当前选择者:</span>
              <span className="px-2 py-1 bg-yellow-100 text-yellow-800 rounded font-mono">
                {tributeState.currentSelectingSeat ?? 'null'}
              </span>
              <span className="text-gray-600">飞行动画:</span>
              <span className="px-2 py-1 bg-purple-100 text-purple-800 rounded font-mono">
                {tributeState.flyingCards.length}
              </span>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-lg shadow">
          {tributeData ? (
            <TributeBoard
              tributeData={tributeData}
              players={MOCK_ROOM.players}
              currentPlayerSeat={currentPlayerSeat}
              playerHand={playerHand}
              selectedCards={selectedCards}
              onCardSelect={setSelectedCards}
              onSelectTribute={handleSelectTribute}
              onReturnTribute={handleReturnTribute}
            />
          ) : (
            <div className="p-12 text-center text-gray-500">
              <div className="text-lg mb-2">请选择一个测试场景</div>
              <div className="text-sm">点击上方按钮开始测试</div>
            </div>
          )}
        </div>

        <div className="bg-white rounded-lg shadow">
          <button
            onClick={() => setShowStatePanel(!showStatePanel)}
            className="w-full p-4 text-left font-medium text-gray-700 flex justify-between items-center"
          >
            <span>状态查看</span>
            <span className="text-gray-400">{showStatePanel ? '▲' : '▼'}</span>
          </button>
          {showStatePanel && (
            <div className="p-4 border-t">
              <pre className="text-xs bg-gray-50 p-4 rounded overflow-auto max-h-96">
                {JSON.stringify(
                  {
                    step: tributeState.step,
                    tributeStarted: tributeState.tributeStarted,
                    tributeExempted: tributeState.tributeExempted,
                    submittedCards: tributeState.submittedCards,
                    poolCards: tributeState.poolCards,
                    selectedCards: tributeState.selectedCards,
                    returnedCards: tributeState.returnedCards,
                    currentSelectingSeat: tributeState.currentSelectingSeat,
                    flyingCards: tributeState.flyingCards,
                    messages: tributeState.messages,
                  },
                  null,
                  2
                )}
              </pre>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default TributeTestPage;
