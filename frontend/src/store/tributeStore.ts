import { create } from 'zustand';
import type { Card } from '../types';
import type {
  TributeStartedPayload,
  TributeExemptedPayload,
  TributeCardReturnedPayload,
} from '../types/generated/event';
import { TributeType } from '../types/generated/event';

export type TributeStep =
  | 'idle'
  | 'started'
  | 'exempted'
  | 'submitting'
  | 'selecting'
  | 'returning'
  | 'completed';

interface ReturnedCardInfo {
  fromSeat: number;
  toSeat: number;
  card: Card;
}

export interface FlyingCard {
  id: string;
  card: Card;
  fromSeat: number | 'pool';
  toSeat: number | 'pool';
  fromPoolSlot?: number;
  toPoolSlot?: number;
}

interface TributeState {
  step: TributeStep;
  tributeStarted: TributeStartedPayload | null;
  tributeExempted: TributeExemptedPayload | null;
  submittedCards: { [giverSeat: number]: Card };
  poolCards: (Card | null)[];
  selectedCards: { [receiverSeat: number]: Card };
  returnedCards: ReturnedCardInfo[];
  messages: string[];
  currentSelectingSeat: number | null;
  flyingCards: FlyingCard[];
}

interface TributeActions {
  handleTributeStarted: (payload: TributeStartedPayload) => void;
  handleTributeExempted: (payload: TributeExemptedPayload) => void;
  handleCardSubmitted: (actorSeat: number, card: Card) => void;
  handleCardSelected: (actorSeat: number, card: Card) => void;
  handleCardReturned: (actorSeat: number, payload: TributeCardReturnedPayload) => void;
  handleCompleted: () => void;
  removeFlyingCard: (id: string) => void;
  reset: () => void;
}

type TributeStore = TributeState & TributeActions;

const initialState: TributeState = {
  step: 'idle',
  tributeStarted: null,
  tributeExempted: null,
  submittedCards: {},
  poolCards: [],
  selectedCards: {},
  returnedCards: [],
  messages: [],
  currentSelectingSeat: null,
  flyingCards: [],
};

const getTributeTypeName = (type: TributeType): string => {
  switch (type) {
    case TributeType.TRIBUTE_TYPE_DOUBLE_DOWN:
      return '双下';
    case TributeType.TRIBUTE_TYPE_SINGLE_LAST:
      return '单下';
    case TributeType.TRIBUTE_TYPE_PARTNER_LAST:
      return '末游';
    case TributeType.TRIBUTE_TYPE_NONE:
      return '无需进贡';
    default:
      return '未知';
  }
};

const getNextSelectingSeat = (
  receivers: number[],
  selectedCards: { [receiverSeat: number]: Card }
): number | null => {
  for (const receiver of receivers) {
    if (!(receiver in selectedCards)) {
      return receiver;
    }
  }
  return null;
};

export const useTributeStore = create<TributeStore>((set, get) => ({
  ...initialState,

  handleTributeStarted: (payload: TributeStartedPayload) => {
    const tributeTypeName = getTributeTypeName(payload.tributeType);
    const messages = [
      `进贡类型: ${tributeTypeName}`,
      `进贡方: 座位 ${payload.givers.join(', ')}`,
      `收贡方: 座位 ${payload.receivers.join(', ')}`,
    ];

    const maxSlots = payload.tributeType === TributeType.TRIBUTE_TYPE_DOUBLE_DOWN ? 2 : 1;

    set({
      step: 'started',
      tributeStarted: payload,
      tributeExempted: null,
      submittedCards: {},
      poolCards: Array(maxSlots).fill(null),
      selectedCards: {},
      returnedCards: [],
      messages,
      currentSelectingSeat: null,
      flyingCards: [],
    });
  },

  handleTributeExempted: (payload: TributeExemptedPayload) => {
    const state = get();
    const holders = Object.entries(payload.bigJokerHolders || {});
    const holderMessages = holders.map(
      ([seat, count]) => `座位 ${seat} 持有 ${count} 张大王`
    );
    const messages = [
      ...state.messages,
      ...holderMessages,
      holders.length > 0 ? '抗贡成功' : '抗贡失败；开始进贡',
    ];

    set({
      step: holders.length > 0 ? 'exempted' : 'submitting',
      tributeExempted: payload,
      messages,
    });
  },

  handleCardSubmitted: (actorSeat: number, card: Card) => {
    const state = get();
    const newSubmittedCards = { ...state.submittedCards, [actorSeat]: card };
    const messages = [...state.messages, `座位 ${actorSeat} 提交了贡牌`];

    const tributeStarted = state.tributeStarted;
    const allSubmitted =
      tributeStarted && tributeStarted.givers.every((g) => g in newSubmittedCards);

    const occupiedByFlying = new Set(
      state.flyingCards
        .filter(fc => fc.toSeat === 'pool' && fc.toPoolSlot !== undefined)
        .map(fc => fc.toPoolSlot)
    );
    const toPoolSlot = state.poolCards.findIndex(
      (c, idx) => c === null && !occupiedByFlying.has(idx)
    );
    const targetSlot = toPoolSlot >= 0 ? toPoolSlot : state.poolCards.length;

    const flyingCard: FlyingCard = {
      id: `submit-${actorSeat}-${Date.now()}`,
      card,
      fromSeat: actorSeat,
      toSeat: 'pool',
      toPoolSlot: targetSlot,
    };

    const nextSelectingSeat = allSubmitted && tributeStarted
      ? getNextSelectingSeat(tributeStarted.receivers, state.selectedCards)
      : null;

    set({
      step: allSubmitted ? 'selecting' : 'submitting',
      submittedCards: newSubmittedCards,
      messages,
      currentSelectingSeat: nextSelectingSeat,
      flyingCards: [...state.flyingCards, flyingCard],
    });
  },

  handleCardSelected: (actorSeat: number, card: Card) => {
    const state = get();
    const newSelectedCards = { ...state.selectedCards, [actorSeat]: card };
    const fromPoolSlot = state.poolCards.findIndex(
      (c) => c !== null && c.deckIndex === card.deckIndex
    );
    const messages = [...state.messages, `座位 ${actorSeat} 选择了贡牌`];

    const tributeStarted = state.tributeStarted;
    const allSelected =
      tributeStarted &&
      tributeStarted.receivers.every((r) => r in newSelectedCards);

    const flyingCard: FlyingCard = {
      id: `select-${actorSeat}-${Date.now()}`,
      card,
      fromSeat: 'pool',
      toSeat: actorSeat,
      fromPoolSlot: fromPoolSlot >= 0 ? fromPoolSlot : undefined,
    };

    const newPoolCards = [...state.poolCards];
    if (fromPoolSlot >= 0) {
      newPoolCards[fromPoolSlot] = null;
    }

    const nextSelectingSeat = !allSelected && tributeStarted
      ? getNextSelectingSeat(tributeStarted.receivers, newSelectedCards)
      : null;

    set({
      step: allSelected ? 'returning' : 'selecting',
      selectedCards: newSelectedCards,
      poolCards: newPoolCards,
      messages,
      currentSelectingSeat: nextSelectingSeat,
      flyingCards: [...state.flyingCards, flyingCard],
    });
  },

  handleCardReturned: (actorSeat: number, payload: TributeCardReturnedPayload) => {
    const state = get();
    if (!payload.returnedCard) return;

    const returnedInfo: ReturnedCardInfo = {
      fromSeat: actorSeat,
      toSeat: payload.targetPlayer,
      card: payload.returnedCard as Card,
    };

    const messages = [
      ...state.messages,
      `座位 ${actorSeat} 向座位 ${payload.targetPlayer} 还贡`,
    ];

    const flyingCard: FlyingCard = {
      id: `return-${actorSeat}-${Date.now()}`,
      card: payload.returnedCard as Card,
      fromSeat: actorSeat,
      toSeat: payload.targetPlayer,
    };

    set({
      returnedCards: [...state.returnedCards, returnedInfo],
      messages,
      flyingCards: [...state.flyingCards, flyingCard],
    });
  },

  handleCompleted: () => {
    set({ step: 'completed' });
  },

  removeFlyingCard: (id: string) => {
    const state = get();
    const flyingCard = state.flyingCards.find((fc) => fc.id === id);
    
    if (flyingCard && flyingCard.toSeat === 'pool' && flyingCard.toPoolSlot !== undefined) {
      const newPoolCards = [...state.poolCards];
      if (flyingCard.toPoolSlot < newPoolCards.length) {
        newPoolCards[flyingCard.toPoolSlot] = flyingCard.card;
      } else {
        newPoolCards.push(flyingCard.card);
      }
      set({
        poolCards: newPoolCards,
        flyingCards: state.flyingCards.filter((fc) => fc.id !== id),
      });
    } else {
      set({
        flyingCards: state.flyingCards.filter((fc) => fc.id !== id),
      });
    }
  },

  reset: () => {
    set(initialState);
  },
}));
