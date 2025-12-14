import React, { useState, useEffect } from 'react';
import type { Story } from '@ladle/react';
import TributeBoard from './TributeBoard';
import type { Player, Card, TributePair } from '../../../types';
import { TributeStatus } from '../../../types/generated/view';
import { TributeType, EventType } from '../../../types/generated/event';
import { wsClient } from '../../../services/websocket';
import { WS_MESSAGE_TYPES } from '../../../types';

function createMockPlayer(username: string, seat: number): Player {
  return {
    id: `player-${seat}`,
    username,
    seat,
    online: true,
    auto_play: false,
  };
}

const mockPlayers: Player[] = [
  createMockPlayer('玩家A', 0),
  createMockPlayer('玩家B', 1),
  createMockPlayer('玩家C', 2),
  createMockPlayer('玩家D', 3),
];

function createCard(rank: number, suit: number, deckIndex: number): Card {
  return { rank, suit, deckIndex };
}

const mockCards = {
  bigJoker: createCard(16, -1, 0),
  smallJoker: createCard(15, -1, 1),
  spadeA: createCard(14, 0, 2),
  heartA: createCard(14, 1, 3),
  clubK: createCard(13, 2, 4),
  diamondK: createCard(13, 3, 5),
  spadeQ: createCard(12, 0, 6),
  heartQ: createCard(12, 1, 7),
  clubJ: createCard(11, 2, 8),
  diamond10: createCard(10, 3, 9),
  spade9: createCard(9, 0, 10),
  heart8: createCard(8, 1, 11),
  club7: createCard(7, 2, 12),
  diamond6: createCard(6, 3, 13),
  spade5: createCard(5, 0, 14),
  heart4: createCard(4, 1, 15),
  club3: createCard(3, 2, 16),
  diamond2: createCard(2, 3, 17),
};

function createTributePair(
  giver: number,
  receiver: number,
  tributeCard?: Card,
  returnCard?: Card
): TributePair {
  return {
    giver,
    receiver,
    tributeCard,
    returnCard,
  };
}

const Wrapper: React.FC<{ children: React.ReactNode }> = ({ children }) => (
  <div className="min-h-screen bg-gray-800 p-4">{children}</div>
);

export const DoubleDownWaiting: Story = () => {
  return (
    <Wrapper>
      <TributeBoard
        tributeData={{
          status: TributeStatus.TRIBUTE_STATUS_WAITING,
          tributeType: TributeType.TRIBUTE_TYPE_DOUBLE_DOWN,
          givers: [0, 2],
          receivers: [1, 3],
          tributePairs: [
            createTributePair(0, 1),
            createTributePair(2, 3),
          ],
          poolCards: [],
          isImmune: false,
        }}
        players={mockPlayers}
        currentPlayerSeat={0}
        teamLevels={[2, 5]}
        currentLevel={2}
        onSelectTribute={(deckIndex) => console.log('Select tribute:', deckIndex)}
      />
    </Wrapper>
  );
};
DoubleDownWaiting.meta = {
  description: '双下 - 待上贡阶段（座位0视角，是进贡方，显示「待上贡」）',
};

export const DoubleDownPartialSubmitted: Story = () => {
  return (
    <Wrapper>
      <TributeBoard
        tributeData={{
          status: TributeStatus.TRIBUTE_STATUS_WAITING,
          tributeType: TributeType.TRIBUTE_TYPE_DOUBLE_DOWN,
          givers: [0, 2],
          receivers: [1, 3],
          tributePairs: [
            createTributePair(0, 1, mockCards.bigJoker),
            createTributePair(2, 3),
          ],
          poolCards: [mockCards.bigJoker],
          isImmune: false,
        }}
        players={mockPlayers}
        currentPlayerSeat={2}
        teamLevels={[2, 5]}
        currentLevel={2}
        onSelectTribute={(deckIndex) => console.log('Select tribute:', deckIndex)}
      />
    </Wrapper>
  );
};
DoubleDownPartialSubmitted.meta = {
  description: '双下 - 部分已上贡（座位2视角，座位0已上贡显示「已上贡」，座位2待上贡）',
};

export const DoubleDownSelecting: Story = () => {
  return (
    <Wrapper>
      <TributeBoard
        tributeData={{
          status: TributeStatus.TRIBUTE_STATUS_SELECTING,
          tributeType: TributeType.TRIBUTE_TYPE_DOUBLE_DOWN,
          givers: [0, 2],
          receivers: [1, 3],
          tributePairs: [
            createTributePair(0, 1, mockCards.bigJoker),
            createTributePair(2, 3, mockCards.smallJoker),
          ],
          poolCards: [mockCards.bigJoker, mockCards.smallJoker],
          isImmune: false,
        }}
        players={mockPlayers}
        currentPlayerSeat={1}
        teamLevels={[2, 5]}
        currentLevel={2}
        onSelectTribute={(deckIndex) => console.log('Select tribute:', deckIndex)}
      />
    </Wrapper>
  );
};
DoubleDownSelecting.meta = {
  description: '双下 - 选牌中（座位1视角，是收贡方，显示「选牌中」带脉冲动画，可点击贡牌池选牌）',
};

export const DoubleDownReturning: Story = () => {
  return (
    <Wrapper>
      <TributeBoard
        tributeData={{
          status: TributeStatus.TRIBUTE_STATUS_RETURNING,
          tributeType: TributeType.TRIBUTE_TYPE_DOUBLE_DOWN,
          givers: [0, 2],
          receivers: [1, 3],
          tributePairs: [
            createTributePair(0, 1, mockCards.bigJoker),
            createTributePair(2, 3, mockCards.smallJoker),
          ],
          poolCards: [],
          isImmune: false,
        }}
        players={mockPlayers}
        currentPlayerSeat={1}
        teamLevels={[2, 5]}
        currentLevel={2}
        onSelectTribute={(deckIndex) => console.log('Select tribute:', deckIndex)}
      />
    </Wrapper>
  );
};
DoubleDownReturning.meta = {
  description: '双下 - 还贡中（座位1视角，收贡方需选择一张牌还贡，底部显示手牌和"确认还贡"按钮）',
};

export const SingleLastWaiting: Story = () => {
  return (
    <Wrapper>
      <TributeBoard
        tributeData={{
          status: TributeStatus.TRIBUTE_STATUS_WAITING,
          tributeType: TributeType.TRIBUTE_TYPE_SINGLE_LAST,
          givers: [2],
          receivers: [0],
          tributePairs: [createTributePair(2, 0)],
          poolCards: [],
          isImmune: false,
        }}
        players={mockPlayers}
        currentPlayerSeat={0}
        teamLevels={[5, 2]}
        currentLevel={5}
        onSelectTribute={(deckIndex) => console.log('Select tribute:', deckIndex)}
      />
    </Wrapper>
  );
};
SingleLastWaiting.meta = {
  description: '单下 - 待上贡（座位0视角是收贡方，座位2是进贡方显示「待上贡」）',
};

export const PartnerLastWaiting: Story = () => {
  return (
    <Wrapper>
      <TributeBoard
        tributeData={{
          status: TributeStatus.TRIBUTE_STATUS_WAITING,
          tributeType: TributeType.TRIBUTE_TYPE_PARTNER_LAST,
          givers: [1],
          receivers: [3],
          tributePairs: [createTributePair(1, 3)],
          poolCards: [],
          isImmune: false,
        }}
        players={mockPlayers}
        currentPlayerSeat={0}
        teamLevels={[5, 3]}
        currentLevel={5}
        onSelectTribute={(deckIndex) => console.log('Select tribute:', deckIndex)}
      />
    </Wrapper>
  );
};
PartnerLastWaiting.meta = {
  description: '末游 - 待上贡（座位0视角是旁观者，座位1进贡给座位3）',
};

export const ViewFromSeat0: Story = () => {
  return (
    <Wrapper>
      <TributeBoard
        tributeData={{
          status: TributeStatus.TRIBUTE_STATUS_WAITING,
          tributeType: TributeType.TRIBUTE_TYPE_DOUBLE_DOWN,
          givers: [0, 2],
          receivers: [1, 3],
          tributePairs: [
            createTributePair(0, 1),
            createTributePair(2, 3),
          ],
          poolCards: [],
          isImmune: false,
        }}
        players={mockPlayers}
        currentPlayerSeat={0}
        teamLevels={[2, 5]}
        currentLevel={2}
        onSelectTribute={(deckIndex) => console.log('Select tribute:', deckIndex)}
      />
    </Wrapper>
  );
};
ViewFromSeat0.meta = {
  description: '座位0视角 - 玩家B在左，玩家C在上，玩家D在右',
};

export const ViewFromSeat1: Story = () => {
  return (
    <Wrapper>
      <TributeBoard
        tributeData={{
          status: TributeStatus.TRIBUTE_STATUS_WAITING,
          tributeType: TributeType.TRIBUTE_TYPE_DOUBLE_DOWN,
          givers: [0, 2],
          receivers: [1, 3],
          tributePairs: [
            createTributePair(0, 1),
            createTributePair(2, 3),
          ],
          poolCards: [],
          isImmune: false,
        }}
        players={mockPlayers}
        currentPlayerSeat={1}
        teamLevels={[2, 5]}
        currentLevel={2}
        onSelectTribute={(deckIndex) => console.log('Select tribute:', deckIndex)}
      />
    </Wrapper>
  );
};
ViewFromSeat1.meta = {
  description: '座位1视角 - 玩家C在左，玩家D在上，玩家A在右',
};

export const ViewFromSeat2: Story = () => {
  return (
    <Wrapper>
      <TributeBoard
        tributeData={{
          status: TributeStatus.TRIBUTE_STATUS_WAITING,
          tributeType: TributeType.TRIBUTE_TYPE_DOUBLE_DOWN,
          givers: [0, 2],
          receivers: [1, 3],
          tributePairs: [
            createTributePair(0, 1),
            createTributePair(2, 3),
          ],
          poolCards: [],
          isImmune: false,
        }}
        players={mockPlayers}
        currentPlayerSeat={2}
        teamLevels={[2, 5]}
        currentLevel={2}
        onSelectTribute={(deckIndex) => console.log('Select tribute:', deckIndex)}
      />
    </Wrapper>
  );
};
ViewFromSeat2.meta = {
  description: '座位2视角 - 玩家D在左，玩家A在上，玩家B在右',
};

export const ViewFromSeat3: Story = () => {
  return (
    <Wrapper>
      <TributeBoard
        tributeData={{
          status: TributeStatus.TRIBUTE_STATUS_WAITING,
          tributeType: TributeType.TRIBUTE_TYPE_DOUBLE_DOWN,
          givers: [0, 2],
          receivers: [1, 3],
          tributePairs: [
            createTributePair(0, 1),
            createTributePair(2, 3),
          ],
          poolCards: [],
          isImmune: false,
        }}
        players={mockPlayers}
        currentPlayerSeat={3}
        teamLevels={[2, 5]}
        currentLevel={2}
        onSelectTribute={(deckIndex) => console.log('Select tribute:', deckIndex)}
      />
    </Wrapper>
  );
};
ViewFromSeat3.meta = {
  description: '座位3视角 - 玩家A在左，玩家B在上，玩家C在右',
};

export const ReceiverSelectingHighlight: Story = () => {
  return (
    <Wrapper>
      <TributeBoard
        tributeData={{
          status: TributeStatus.TRIBUTE_STATUS_SELECTING,
          tributeType: TributeType.TRIBUTE_TYPE_DOUBLE_DOWN,
          givers: [0, 2],
          receivers: [1, 3],
          tributePairs: [
            createTributePair(0, 1, mockCards.bigJoker),
            createTributePair(2, 3, mockCards.smallJoker),
          ],
          poolCards: [mockCards.bigJoker, mockCards.smallJoker],
          isImmune: false,
        }}
        players={mockPlayers}
        currentPlayerSeat={0}
        teamLevels={[2, 5]}
        currentLevel={2}
        onSelectTribute={(deckIndex) => console.log('Select tribute:', deckIndex)}
      />
    </Wrapper>
  );
};
ReceiverSelectingHighlight.meta = {
  description: '选牌高亮 - 座位0视角看座位1选牌中（座位1有高亮边框+「选牌中」脉冲动画）',
};

export const FlyingCardSubmit: Story = () => {
  const [poolCards, setPoolCards] = useState<Card[]>([]);
  const [tributePairs, setTributePairs] = useState<TributePair[]>([
    createTributePair(0, 1),
    createTributePair(2, 3),
  ]);

  const handleSubmitCard = (fromSeat: number, card: Card, toSlot: number) => {
    // 1. 先更新 poolCards（模拟 TRIBUTE_VIEW 消息到达）
    setPoolCards(prev => {
      const newPool = [...prev];
      newPool[toSlot] = card;
      return newPool;
    });
    setTributePairs(prev =>
      prev.map(pair =>
        pair.giver === fromSeat ? { ...pair, tributeCard: card } : pair
      )
    );

    // 2. 然后触发 GAME_EVENT（飞牌动画）
    wsClient.__mockEmit__(WS_MESSAGE_TYPES.GAME_EVENT, {
      type: EventType.EVENT_TYPE_TRIBUTE_CARD_SUBMITTED,
      actorSeat: fromSeat,
      tributeCardSubmitted: { submittedCard: card },
    });
  };

  return (
    <Wrapper>
      <div className="space-y-4">
        <div className="flex gap-2 justify-center">
          <button
            onClick={() => handleSubmitCard(0, mockCards.bigJoker, 0)}
            className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
          >
            座位0上贡大王
          </button>
          <button
            onClick={() => handleSubmitCard(2, mockCards.smallJoker, 1)}
            className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
          >
            座位2上贡小王
          </button>
          <button
            onClick={() => {
              setPoolCards([]);
              setTributePairs([
                createTributePair(0, 1),
                createTributePair(2, 3),
              ]);
            }}
            className="px-4 py-2 bg-gray-500 text-white rounded hover:bg-gray-600"
          >
            重置
          </button>
        </div>

        <TributeBoard
          tributeData={{
            status: TributeStatus.TRIBUTE_STATUS_WAITING,
            tributeType: TributeType.TRIBUTE_TYPE_DOUBLE_DOWN,
            givers: [0, 2],
            receivers: [1, 3],
            tributePairs,
            poolCards,
            isImmune: false,
          }}
          players={mockPlayers}
          currentPlayerSeat={0}
          teamLevels={[2, 5]}
          currentLevel={2}
          onSelectTribute={(deckIndex) => console.log('Select:', deckIndex)}
        />
      </div>
    </Wrapper>
  );
};
FlyingCardSubmit.meta = {
  description: '飞牌动画 - 上贡（点击按钮触发卡牌从玩家飞向贡牌池）',
};

export const FlyingCardSelect: Story = () => {
  const [poolCards, setPoolCards] = useState<Card[]>([
    mockCards.bigJoker,
    mockCards.smallJoker,
  ]);

  const handleSelectCard = (actorSeat: number, card: Card) => {
    // 1. 先更新 poolCards（模拟 TRIBUTE_VIEW 消息到达）
    setPoolCards(prev => prev.filter(c => c.deckIndex !== card.deckIndex));

    // 2. 然后触发 GAME_EVENT（飞牌动画）
    wsClient.__mockEmit__(WS_MESSAGE_TYPES.GAME_EVENT, {
      type: EventType.EVENT_TYPE_TRIBUTE_CARD_SELECTED,
      actorSeat,
      tributeCardSelected: { selectedCard: card },
    });
  };

  return (
    <Wrapper>
      <div className="space-y-4">
        <div className="flex gap-2 justify-center">
          <button
            onClick={() => handleSelectCard(1, mockCards.bigJoker)}
            className="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600"
          >
            座位1选大王
          </button>
          <button
            onClick={() => handleSelectCard(3, mockCards.smallJoker)}
            className="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600"
          >
            座位3选小王
          </button>
          <button
            onClick={() => setPoolCards([mockCards.bigJoker, mockCards.smallJoker])}
            className="px-4 py-2 bg-gray-500 text-white rounded hover:bg-gray-600"
          >
            重置
          </button>
        </div>

        <TributeBoard
          tributeData={{
            status: TributeStatus.TRIBUTE_STATUS_SELECTING,
            tributeType: TributeType.TRIBUTE_TYPE_DOUBLE_DOWN,
            givers: [0, 2],
            receivers: [1, 3],
            tributePairs: [
              createTributePair(0, 1, mockCards.bigJoker),
              createTributePair(2, 3, mockCards.smallJoker),
            ],
            poolCards,
            isImmune: false,
          }}
          players={mockPlayers}
          currentPlayerSeat={0}
          teamLevels={[2, 5]}
          currentLevel={2}
          onSelectTribute={(deckIndex) => console.log('Select:', deckIndex)}
        />
      </div>
    </Wrapper>
  );
};
FlyingCardSelect.meta = {
  description: '飞牌动画 - 选牌（点击按钮触发卡牌从贡牌池飞向收贡方）',
};

export const FlyingCardReturn: Story = () => {
  const handleReturnCard = (fromSeat: number, toSeat: number) => {
    wsClient.__mockEmit__(WS_MESSAGE_TYPES.GAME_EVENT, {
      type: EventType.EVENT_TYPE_TRIBUTE_CARD_RETURNED,
      actorSeat: fromSeat,
      tributeCardReturned: {
        returnedCard: mockCards.diamond2,
        targetPlayer: toSeat,
      },
    });
  };

  return (
    <Wrapper>
      <div className="space-y-4">
        <div className="flex gap-2 justify-center">
          <button
            onClick={() => handleReturnCard(1, 0)}
            className="px-4 py-2 bg-purple-500 text-white rounded hover:bg-purple-600"
          >
            座位1还贡给座位0
          </button>
          <button
            onClick={() => handleReturnCard(3, 2)}
            className="px-4 py-2 bg-purple-500 text-white rounded hover:bg-purple-600"
          >
            座位3还贡给座位2
          </button>
        </div>

        <TributeBoard
          tributeData={{
            status: TributeStatus.TRIBUTE_STATUS_RETURNING,
            tributeType: TributeType.TRIBUTE_TYPE_DOUBLE_DOWN,
            givers: [0, 2],
            receivers: [1, 3],
            tributePairs: [
              createTributePair(0, 1, mockCards.bigJoker),
              createTributePair(2, 3, mockCards.smallJoker),
            ],
            poolCards: [],
            isImmune: false,
          }}
          players={mockPlayers}
          currentPlayerSeat={0}
          teamLevels={[2, 5]}
          currentLevel={2}
          onSelectTribute={(deckIndex) => console.log('Select:', deckIndex)}
        />
      </div>
    </Wrapper>
  );
};
FlyingCardReturn.meta = {
  description: '飞牌动画 - 还贡（点击按钮触发卡牌从收贡方飞向进贡方）',
};

export const FlyingCardAutoDemo: Story = () => {
  const [poolCards, setPoolCards] = useState<Card[]>([]);
  const [phase, setPhase] = useState<'idle' | 'submit' | 'select' | 'return' | 'done'>('idle');
  const [tributePairs, setTributePairs] = useState<TributePair[]>([
    createTributePair(0, 1),
    createTributePair(2, 3),
  ]);

  const startDemo = () => {
    setPoolCards([]);
    setTributePairs([createTributePair(0, 1), createTributePair(2, 3)]);
    setPhase('submit');
  };

  useEffect(() => {
    if (phase === 'submit') {
      const timer = setTimeout(() => {
        // 1. 先更新 poolCards（模拟 TRIBUTE_VIEW 消息到达）
        setPoolCards([mockCards.bigJoker]);
        setTributePairs([
          createTributePair(0, 1, mockCards.bigJoker),
          createTributePair(2, 3),
        ]);

        // 2. 然后触发 GAME_EVENT（飞牌动画）
        wsClient.__mockEmit__(WS_MESSAGE_TYPES.GAME_EVENT, {
          type: EventType.EVENT_TYPE_TRIBUTE_CARD_SUBMITTED,
          actorSeat: 0,
          tributeCardSubmitted: { submittedCard: mockCards.bigJoker },
        });

        // 3. 动画结束后进入下一阶段
        setTimeout(() => {
          setPhase('select');
        }, 600);
      }, 500);

      return () => clearTimeout(timer);
    } else if (phase === 'select') {
      const timer = setTimeout(() => {
        // 1. 先更新 poolCards（模拟 TRIBUTE_VIEW 消息到达）
        setPoolCards([]);

        // 2. 然后触发 GAME_EVENT（飞牌动画）
        wsClient.__mockEmit__(WS_MESSAGE_TYPES.GAME_EVENT, {
          type: EventType.EVENT_TYPE_TRIBUTE_CARD_SELECTED,
          actorSeat: 1,
          tributeCardSelected: { selectedCard: mockCards.bigJoker },
        });

        // 3. 动画结束后进入下一阶段
        setTimeout(() => {
          setPhase('return');
        }, 600);
      }, 1500);

      return () => clearTimeout(timer);
    } else if (phase === 'return') {
      const timer = setTimeout(() => {
        wsClient.__mockEmit__(WS_MESSAGE_TYPES.GAME_EVENT, {
          type: EventType.EVENT_TYPE_TRIBUTE_CARD_RETURNED,
          actorSeat: 1,
          tributeCardReturned: {
            returnedCard: mockCards.diamond2,
            targetPlayer: 0,
          },
        });

        setTimeout(() => {
          setPhase('done');
        }, 600);
      }, 1500);

      return () => clearTimeout(timer);
    }
  }, [phase]);

  const getPhaseText = () => {
    switch (phase) {
      case 'idle': return '点击开始';
      case 'submit': return '上贡中...';
      case 'select': return '选牌中...';
      case 'return': return '还贡中...';
      case 'done': return '完成';
    }
  };

  return (
    <Wrapper>
      <div className="space-y-4">
        <div className="flex gap-2 justify-center items-center">
          <div className="px-4 py-2 bg-gray-700 text-white rounded min-w-[100px] text-center">
            {getPhaseText()}
          </div>
          <button
            onClick={startDemo}
            disabled={phase !== 'idle' && phase !== 'done'}
            className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {phase === 'done' ? '重新播放' : '开始演示'}
          </button>
        </div>

        <TributeBoard
          tributeData={{
            status:
              phase === 'submit' || phase === 'idle'
                ? TributeStatus.TRIBUTE_STATUS_WAITING
                : phase === 'select'
                ? TributeStatus.TRIBUTE_STATUS_SELECTING
                : TributeStatus.TRIBUTE_STATUS_RETURNING,
            tributeType: TributeType.TRIBUTE_TYPE_DOUBLE_DOWN,
            givers: [0, 2],
            receivers: [1, 3],
            tributePairs,
            poolCards,
            isImmune: false,
          }}
          players={mockPlayers}
          currentPlayerSeat={0}
          teamLevels={[2, 5]}
          currentLevel={2}
          onSelectTribute={(deckIndex) => console.log('Select:', deckIndex)}
        />
      </div>
    </Wrapper>
  );
};
FlyingCardAutoDemo.meta = {
  description: '飞牌动画 - 自动演示（自动播放完整流程：上贡→选牌→还贡）',
};
