import React, { useState, useEffect } from 'react';
import type { Story } from '@ladle/react';
import { MemoryRouter } from 'react-router-dom';
import type { Player } from '../../types';
import type { Card, PlayAction } from '../../types/proto';
import { CompType } from '../../types/proto';
import { useAuthStore } from '../../store/authStore';
import Header from '../layout/Header';
import GameBoard from './GameBoard';
import GameControls from './GameControls';
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

const makeCard = (suit: number, rank: number, deckIndex: number): Card => ({
  suit,
  rank,
  deckIndex,
});

const makePlayer = (id: string, username: string, seat: number): Player => ({
  id,
  username,
  seat,
  online: true,
  auto_play: false,
});

const makePlayAction = (
  playerSeat: number,
  cards: Card[],
  compType: CompType = CompType.COMP_TYPE_UNSPECIFIED,
  isPass = false
): PlayAction => ({
  playerSeat,
  cards,
  compType,
  timestampMs: Date.now(),
  isPass,
});

const mockPlayers: Player[] = [
  makePlayer('p1', '玩家A', 0),
  makePlayer('p2', '玩家B', 1),
  makePlayer('p3', '玩家C', 2),
  makePlayer('p4', '玩家D', 3),
];

const createMockHand = (): Card[] => {
  const cards: Card[] = [];
  let deckIndex = 0;
  cards.push(makeCard(-1, 16, deckIndex++));
  cards.push(makeCard(-1, 15, deckIndex++));
  for (let rank = 14; rank >= 3; rank--) {
    for (let suit = 0; suit < 2; suit++) {
      if (cards.length < 27) {
        cards.push(makeCard(suit, rank, deckIndex++));
      }
    }
  }
  return cards.slice(0, 27);
};

const createFewCards = (): Card[] => [
  makeCard(0, 14, 0),
  makeCard(1, 14, 1),
  makeCard(2, 13, 2),
  makeCard(3, 10, 3),
  makeCard(0, 7, 4),
  makeCard(1, 3, 5),
];

const createManyDuplicatesHand = (): Card[] => {
  const cards: Card[] = [];
  let idx = 0;
  for (let i = 0; i < 8; i++) cards.push(makeCard(i % 4, 9, idx++));
  for (let i = 0; i < 8; i++) cards.push(makeCard(i % 4, 10, idx++));
  for (let rank = 14; rank >= 3 && cards.length < 27; rank--) {
    cards.push(makeCard(0, rank, idx++));
  }
  return cards;
};

const manyPlayedCards: Card[] = [
  ...Array.from({ length: 8 }, (_, i) => makeCard(i % 4, 10, 200 + i)),
  ...Array.from({ length: 8 }, (_, i) => makeCard(i % 4, 11, 208 + i)),
  ...Array.from({ length: 8 }, (_, i) => makeCard(i % 4, 12, 216 + i)),
];

const pairOfKings: Card[] = [makeCard(0, 13, 100), makeCard(1, 13, 101)];
const straight: Card[] = [
  makeCard(0, 5, 110),
  makeCard(1, 6, 111),
  makeCard(2, 7, 112),
  makeCard(3, 8, 113),
  makeCard(0, 9, 114),
];

interface PlayingPhaseWrapperProps {
  hand: Card[];
  plays?: PlayAction[];
  currentTurn: number;
  playerSeat?: number;
  isMyTurn?: boolean;
  canPlay?: boolean;
  finishRank?: number;
  teamLevels?: [number, number];
  dealLevel?: number;
  leader?: number;
}

const PlayingPhaseWrapper: React.FC<PlayingPhaseWrapperProps> = ({
  hand,
  plays = [],
  currentTurn,
  playerSeat = 0,
  isMyTurn = false,
  canPlay = false,
  finishRank,
  teamLevels = [2, 2],
  dealLevel = 2,
  leader = 0,
}) => {
  const [selectedCards, setSelectedCards] = useState<Card[]>([]);

  const handlePlayCards = (cards: Card[]) => {
    console.log('Play cards:', cards);
    setSelectedCards([]);
  };

  const handlePass = () => {
    console.log('Pass');
    setSelectedCards([]);
  };

  const handleHint = (cards: Card[]) => {
    setSelectedCards(cards);
  };

  return (
    <div className="fixed inset-0 z-40 overflow-hidden bg-gradient-to-br from-[hsl(40,8%,96%)] via-[hsl(38,6%,94%)] to-[hsl(35,8%,91%)]">
      <Header collapsible />
      <div className="absolute inset-0 p-2">
        <GameBoard
          teamLevels={teamLevels}
          currentLevel={dealLevel}
          plays={plays}
          currentTurn={currentTurn}
          players={mockPlayers}
          currentPlayerSeat={playerSeat}
          turnDeadline={null}
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
        <GameControls
          selectedCards={selectedCards}
          canPlay={canPlay}
          isMyTurn={isMyTurn}
          turnDeadlineAtMs={isMyTurn ? Date.now() + 30000 : 0}
          onPlayCards={handlePlayCards}
          onPass={handlePass}
          onHint={handleHint}
          handCards={hand}
          plays={plays}
          leader={leader}
          playerSeat={playerSeat}
          dealLevel={dealLevel}
          disabled={false}
        />
      </div>
      <div className="absolute bottom-0 left-0 right-0 z-10">
        <PlayerHand
          cards={hand}
          selectedCards={selectedCards}
          onCardSelect={setSelectedCards}
          currentLevel={dealLevel}
          finishRank={finishRank}
        />
      </div>
    </div>
  );
};

export const MyTurn: Story = () => (
  <PlayingPhaseWrapper
    hand={createMockHand()}
    currentTurn={0}
    playerSeat={0}
    isMyTurn={true}
    canPlay={true}
    leader={0}
  />
);
MyTurn.meta = { description: '轮到我出牌' };
MyTurn.decorators = [withProviders];

export const WaitingOthers: Story = () => (
  <PlayingPhaseWrapper
    hand={createMockHand()}
    currentTurn={1}
    playerSeat={0}
    isMyTurn={false}
    canPlay={false}
  />
);
WaitingOthers.meta = { description: '等待其他玩家出牌' };
WaitingOthers.decorators = [withProviders];

export const WithSelectedCards: Story = () => {
  const hand = createMockHand();
  const Wrapper = () => {
    const [selected, setSelected] = useState<Card[]>([hand[2], hand[3]]);
    return (
      <div className="fixed inset-0 z-40 overflow-hidden bg-gradient-to-br from-[hsl(40,8%,96%)] via-[hsl(38,6%,94%)] to-[hsl(35,8%,91%)]">
        <Header collapsible />
        <div className="absolute inset-0 p-2">
          <GameBoard
            teamLevels={[5, 3]}
            currentLevel={5}
            plays={[]}
            currentTurn={0}
            players={mockPlayers}
            currentPlayerSeat={0}
            turnDeadline={null}
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
          <GameControls
            selectedCards={selected}
            canPlay={true}
            isMyTurn={true}
            turnDeadlineAtMs={Date.now() + 30000}
            onPlayCards={() => setSelected([])}
            onPass={() => setSelected([])}
            onHint={setSelected}
            handCards={hand}
            plays={[]}
            leader={0}
            playerSeat={0}
            dealLevel={5}
            disabled={false}
          />
        </div>
        <div className="absolute bottom-0 left-0 right-0 z-10">
          <PlayerHand
            cards={hand}
            selectedCards={selected}
            onCardSelect={setSelected}
            currentLevel={5}
          />
        </div>
      </div>
    );
  };
  return <Wrapper />;
};
WithSelectedCards.meta = { description: '已选中手牌' };
WithSelectedCards.decorators = [withProviders];

export const WithPlayHistory: Story = () => {
  const plays: PlayAction[] = [
    makePlayAction(1, pairOfKings, CompType.COMP_TYPE_PAIR),
    makePlayAction(2, [], CompType.COMP_TYPE_UNSPECIFIED, true),
    makePlayAction(3, [], CompType.COMP_TYPE_UNSPECIFIED, true),
  ];
  return (
    <PlayingPhaseWrapper
      hand={createMockHand()}
      plays={plays}
      currentTurn={0}
      playerSeat={0}
      isMyTurn={true}
      canPlay={true}
      teamLevels={[7, 5]}
      dealLevel={7}
      leader={1}
    />
  );
};
WithPlayHistory.meta = { description: '桌面有已出牌记录' };
WithPlayHistory.decorators = [withProviders];

export const AfterPass: Story = () => {
  const plays: PlayAction[] = [
    makePlayAction(3, straight, CompType.COMP_TYPE_STRAIGHT),
    makePlayAction(0, [], CompType.COMP_TYPE_UNSPECIFIED, true),
    makePlayAction(1, [], CompType.COMP_TYPE_UNSPECIFIED, true),
    makePlayAction(2, [], CompType.COMP_TYPE_UNSPECIFIED, true),
  ];
  return (
    <PlayingPhaseWrapper
      hand={createMockHand()}
      plays={plays}
      currentTurn={3}
      playerSeat={3}
      isMyTurn={true}
      canPlay={true}
      teamLevels={[10, 8]}
      dealLevel={10}
      leader={3}
    />
  );
};
AfterPass.meta = { description: '其他玩家都不出，新一轮' };
AfterPass.decorators = [withProviders];

export const PlayerFinishedFirst: Story = () => (
  <PlayingPhaseWrapper
    hand={[]}
    currentTurn={1}
    playerSeat={0}
    isMyTurn={false}
    canPlay={false}
    finishRank={1}
  />
);
PlayerFinishedFirst.meta = { description: '头游 - 第一名完成' };
PlayerFinishedFirst.decorators = [withProviders];

export const PlayerFinishedLast: Story = () => (
  <PlayingPhaseWrapper
    hand={[]}
    currentTurn={1}
    playerSeat={0}
    isMyTurn={false}
    canPlay={false}
    finishRank={4}
  />
);
PlayerFinishedLast.meta = { description: '末游 - 第四名完成' };
PlayerFinishedLast.decorators = [withProviders];

export const FewCardsLeft: Story = () => (
  <PlayingPhaseWrapper
    hand={createFewCards()}
    currentTurn={0}
    playerSeat={0}
    isMyTurn={true}
    canPlay={true}
    teamLevels={[12, 11]}
    dealLevel={12}
  />
);
FewCardsLeft.meta = { description: '手牌只剩几张' };
FewCardsLeft.decorators = [withProviders];

export const ManyDuplicateCards: Story = () => (
  <PlayingPhaseWrapper
    hand={createManyDuplicatesHand()}
    currentTurn={0}
    playerSeat={0}
    isMyTurn={true}
    canPlay={true}
    teamLevels={[9, 7]}
    dealLevel={9}
  />
);
ManyDuplicateCards.meta = { description: '极端场景 - 8张9、8张10' };
ManyDuplicateCards.decorators = [withProviders];

export const ManyPlayedCards: Story = () => {
  const plays: PlayAction[] = [
    makePlayAction(1, manyPlayedCards, CompType.COMP_TYPE_UNSPECIFIED),
  ];
  return (
    <PlayingPhaseWrapper
      hand={createFewCards()}
      plays={plays}
      currentTurn={2}
      playerSeat={0}
      isMyTurn={false}
      canPlay={false}
      teamLevels={[10, 10]}
      dealLevel={10}
      leader={1}
    />
  );
};
ManyPlayedCards.meta = { description: '极端场景 - 出牌区24张牌' };
ManyPlayedCards.decorators = [withProviders];
