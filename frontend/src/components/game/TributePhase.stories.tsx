import React, { useState, useEffect } from 'react';
import type { Story } from '@ladle/react';
import { MemoryRouter } from 'react-router-dom';
import type { Player, Card, TributePair } from '../../types';
import { TributeStatus } from '../../types/generated/view';
import { TributeType } from '../../types/generated/event';
import { useAuthStore } from '../../store/authStore';
import UserMenuFab from '../layout/UserMenuFab';
import TributeBoard from './tribute/TributeBoard';
import TributeControls from './tribute/TributeControls';
import PlayerHand from './PlayerHand';

const MockAuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { login, logout } = useAuthStore();
  useEffect(() => {
    login(
      { id: '1', username: 'testuser', nickname: '测试用户', online: true },
      { access_token: 'mock-token', refresh_token: 'mock-refresh', expires_at: new Date(Date.now() + 3600000).toISOString(), user_id: '1' }
    );
    return () => logout();
  }, [login, logout]);
  return <>{children}</>;
};

const withProviders = (Story: React.ComponentType) => (
  <MemoryRouter>
    <MockAuthProvider>
      <Story />
    </MockAuthProvider>
  </MemoryRouter>
);

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
  return { giver, receiver, tributeCard, returnCard };
}

function createMockHand(): Card[] {
  const cards: Card[] = [];
  let deckIndex = 100;
  cards.push(createCard(16, -1, deckIndex++));
  cards.push(createCard(15, -1, deckIndex++));
  for (let rank = 14; rank >= 3; rank--) {
    for (let suit = 0; suit < 2; suit++) {
      if (cards.length < 27) {
        cards.push(createCard(rank, suit, deckIndex++));
      }
    }
  }
  return cards.slice(0, 27);
}

interface TributeData {
  status: TributeStatus;
  tributeType: TributeType;
  givers: number[];
  receivers: number[];
  tributePairs: TributePair[];
  poolCards: Card[];
  isImmune: boolean;
}

interface TributePhaseWrapperProps {
  tributeData: TributeData;
  hand: Card[];
  playerSeat?: number;
  canReturnTribute?: boolean;
  teamLevels?: [number, number];
  dealLevel?: number;
}

const TributePhaseWrapper: React.FC<TributePhaseWrapperProps> = ({
  tributeData,
  hand,
  playerSeat = 0,
  canReturnTribute = false,
  teamLevels = [2, 5],
  dealLevel = 2,
}) => {
  const [selectedCards, setSelectedCards] = useState<Card[]>([]);

  const handleSelectTribute = (deckIndex: number) => {
    console.log('Select tribute:', deckIndex);
  };

  const handleReturnTribute = () => {
    if (selectedCards.length === 1) {
      console.log('Return tribute:', selectedCards[0].deckIndex);
      setSelectedCards([]);
    }
  };

  const handleHint = (cards: Card[]) => {
    setSelectedCards(cards);
  };

  return (
    <div className="fixed inset-0 z-40 overflow-hidden bg-gradient-to-br from-[hsl(40,8%,96%)] via-[hsl(38,6%,94%)] to-[hsl(35,8%,91%)]">
      <UserMenuFab />
      <div className="absolute inset-0 p-2">
        <TributeBoard
          tributeData={tributeData}
          players={mockPlayers}
          currentPlayerSeat={playerSeat}
          teamLevels={teamLevels}
          currentLevel={dealLevel}
          onSelectTribute={handleSelectTribute}
          className="h-full"
        />
      </div>
      {canReturnTribute && (
        <div
          className="absolute left-1/2 z-20"
          style={{
            transform: 'translateX(-50%)',
            top: 'calc(50% + var(--play-area-offset-y, 0) + var(--table-center-height) / 2)',
          }}
        >
          <TributeControls
            selectedCards={selectedCards}
            canReturnTribute={true}
            turnDeadlineAtMs={Date.now() + 30000}
            onReturnTribute={handleReturnTribute}
            onHint={handleHint}
            handCards={hand}
            dealLevel={dealLevel}
          />
        </div>
      )}
      {hand.length > 0 && (
        <div className="absolute bottom-0 left-0 right-0 z-10">
          <PlayerHand
            cards={hand}
            selectedCards={selectedCards}
            onCardSelect={setSelectedCards}
            currentLevel={dealLevel}
          />
        </div>
      )}
    </div>
  );
};

export const DoubleDownWaiting: Story = () => (
  <TributePhaseWrapper
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
    hand={createMockHand()}
    playerSeat={0}
  />
);
DoubleDownWaiting.meta = { description: '双下 - 待上贡（进贡方视角，座位0需要选择贡牌）' };
DoubleDownWaiting.decorators = [withProviders];

export const DoubleDownSelecting: Story = () => (
  <TributePhaseWrapper
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
    hand={createMockHand()}
    playerSeat={1}
  />
);
DoubleDownSelecting.meta = { description: '双下 - 选牌中（收贡方视角，座位1可点击贡牌池选牌）' };
DoubleDownSelecting.decorators = [withProviders];

export const DoubleDownReturning: Story = () => (
  <TributePhaseWrapper
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
    hand={createMockHand()}
    playerSeat={1}
    canReturnTribute={true}
  />
);
DoubleDownReturning.meta = { description: '双下 - 还贡中（收贡方需选择一张牌还贡，显示TributeControls）' };
DoubleDownReturning.decorators = [withProviders];

export const SingleLastReturning: Story = () => (
  <TributePhaseWrapper
    tributeData={{
      status: TributeStatus.TRIBUTE_STATUS_RETURNING,
      tributeType: TributeType.TRIBUTE_TYPE_SINGLE_LAST,
      givers: [2],
      receivers: [0],
      tributePairs: [createTributePair(2, 0, mockCards.bigJoker)],
      poolCards: [],
      isImmune: false,
    }}
    hand={createMockHand()}
    playerSeat={0}
    canReturnTribute={true}
    teamLevels={[5, 2]}
    dealLevel={5}
  />
);
SingleLastReturning.meta = { description: '单下 - 还贡阶段（座位0是收贡方，需要还贡）' };
SingleLastReturning.decorators = [withProviders];

export const PartnerLastObserver: Story = () => (
  <TributePhaseWrapper
    tributeData={{
      status: TributeStatus.TRIBUTE_STATUS_WAITING,
      tributeType: TributeType.TRIBUTE_TYPE_PARTNER_LAST,
      givers: [1],
      receivers: [3],
      tributePairs: [createTributePair(1, 3)],
      poolCards: [],
      isImmune: false,
    }}
    hand={createMockHand()}
    playerSeat={0}
    teamLevels={[5, 3]}
    dealLevel={5}
  />
);
PartnerLastObserver.meta = { description: '末游 - 旁观者视角（座位0不参与进贡/收贡）' };
PartnerLastObserver.decorators = [withProviders];

export const WithCardSelection: Story = () => {
  const hand = createMockHand();
  const Wrapper = () => {
    const [selected, setSelected] = useState<Card[]>([hand[5]]);
    return (
      <div className="fixed inset-0 z-40 overflow-hidden bg-gradient-to-br from-[hsl(40,8%,96%)] via-[hsl(38,6%,94%)] to-[hsl(35,8%,91%)]">
        <UserMenuFab />
        <div className="absolute inset-0 p-2">
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
            onSelectTribute={(idx) => console.log('Select:', idx)}
            className="h-full"
          />
        </div>
        <div
          className="absolute left-1/2 z-20"
          style={{
            transform: 'translateX(-50%)',
            top: 'calc(50% + var(--play-area-offset-y, 0) + var(--table-center-height) / 2)',
          }}
        >
          <TributeControls
            selectedCards={selected}
            canReturnTribute={true}
            turnDeadlineAtMs={Date.now() + 30000}
            onReturnTribute={() => {
              console.log('Return:', selected[0]?.deckIndex);
              setSelected([]);
            }}
            onHint={setSelected}
            handCards={hand}
            dealLevel={2}
          />
        </div>
        <div className="absolute bottom-0 left-0 right-0 z-10">
          <PlayerHand
            cards={hand}
            selectedCards={selected}
            onCardSelect={setSelected}
            currentLevel={2}
          />
        </div>
      </div>
    );
  };
  return <Wrapper />;
};
WithCardSelection.meta = { description: '还贡时已选中手牌（测试卡牌选择交互）' };
WithCardSelection.decorators = [withProviders];
