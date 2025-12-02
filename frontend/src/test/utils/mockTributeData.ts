import type { Player } from '../../types';
import type { Card } from '../../types/proto';
import type { TributeStartedPayload, TributeExemptedPayload, TributeCardReturnedPayload } from '../../types/generated/event';
import { TributeType } from '../../types/generated/event';

export const createMockCard = (deckIndex: number, rank: number, suit: number): Card => {
  return {
    deckIndex,
    rank,
    suit,
  };
};

export const createMockPlayers = (): Player[] => [
  { id: 'p1', username: 'Alice', seat: 0, online: true, auto_play: false },
  { id: 'p2', username: 'Bob', seat: 1, online: true, auto_play: false },
  { id: 'p3', username: 'Carol', seat: 2, online: true, auto_play: false },
  { id: 'p4', username: 'Dave', seat: 3, online: true, auto_play: false },
];

export const createMockPlayerHand = (): Card[] => [
  createMockCard(100, 3, 0),
  createMockCard(101, 4, 1),
  createMockCard(102, 5, 2),
  createMockCard(103, 6, 3),
  createMockCard(104, 7, 0),
  createMockCard(105, 8, 1),
  createMockCard(106, 9, 2),
  createMockCard(107, 10, 3),
  createMockCard(108, 11, 0),
  createMockCard(109, 12, 1),
  createMockCard(110, 13, 2),
  createMockCard(111, 14, 3),
];

interface TributeScenarioConfig {
  type: 'DOUBLE_DOWN' | 'SINGLE_LAST' | 'PARTNER_LAST';
}

export const createMockTributeStarted = (config: TributeScenarioConfig): TributeStartedPayload => {
  const { type } = config;

  switch (type) {
    case 'DOUBLE_DOWN':
      return {
        tributeType: TributeType.TRIBUTE_TYPE_DOUBLE_DOWN,
        givers: [0, 1],
        receivers: [2, 3],
      };
    case 'SINGLE_LAST':
      return {
        tributeType: TributeType.TRIBUTE_TYPE_SINGLE_LAST,
        givers: [0],
        receivers: [2],
      };
    case 'PARTNER_LAST':
      return {
        tributeType: TributeType.TRIBUTE_TYPE_PARTNER_LAST,
        givers: [0],
        receivers: [2],
      };
    default:
      return {
        tributeType: TributeType.TRIBUTE_TYPE_NONE,
        givers: [],
        receivers: [],
      };
  }
};

export const createMockTributeExempted = (bigJokerHolders: { [seat: number]: number }): TributeExemptedPayload => ({
  bigJokerHolders,
});

export const createMockCardReturned = (targetPlayer: number, card: Card): TributeCardReturnedPayload => ({
  targetPlayer,
  returnedCard: card,
  isAuto: false,
});

export const MOCK_TRIBUTE_CARDS = {
  doubleDown: [
    createMockCard(1, 14, 0),
    createMockCard(2, 14, 1),
  ],
  singleLast: [
    createMockCard(1, 14, 0),
  ],
};

export const MOCK_RETURN_CARDS = {
  fromReceiver2: createMockCard(50, 3, 0),
  fromReceiver3: createMockCard(51, 4, 1),
};
