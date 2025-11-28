import type { Player, Card } from '../../../types';
import type { TributeStep, FlyingCard } from '../../../store/tributeStore';
import type { TributeStartedPayload, TributeExemptedPayload } from '../../../types/generated/event';

export interface TributeData {
  step: TributeStep;
  tributeStarted: TributeStartedPayload | null;
  tributeExempted: TributeExemptedPayload | null;
  submittedCards: { [giverSeat: number]: Card };
  poolCards: (Card | null)[];
  selectedCards: { [receiverSeat: number]: Card };
  returnedCards: Array<{
    fromSeat: number;
    toSeat: number;
    card: Card;
  }>;
  messages: string[];
  currentSelectingSeat: number | null;
  flyingCards: FlyingCard[];
  players: Player[];
  playerSeat: number | null;
}

export interface TributeBoardProps {
  tributeData: TributeData;
  players: (Player | null)[];
  currentPlayerSeat: number;
  playerHand: Card[];
  selectedCards: Card[];
  onCardSelect: (cards: Card[]) => void;
  onSelectTribute: (deckIndex: number) => void;
  onReturnTribute: (deckIndex: number) => void;
}
