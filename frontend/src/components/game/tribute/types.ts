import type { ReactNode } from 'react';
import type { TributePhase, Player, Card } from '../../../types';

export type UIPhase = 
  | 'START' 
  | 'IMMUNITY_CHECK' 
  | 'SUBMITTING' 
  | 'SELECTING' 
  | 'RETURNING' 
  | 'FINISHED';

export interface PhaseConfig {
  title: string;
  icon?: string;
  duration?: number;
  renderContent: (props: TributeFlowProps) => ReactNode;
}

export interface ReturnTask {
  receiver: number; // The player who received the tribute (and needs to return)
  giver: number;    // The player who gave the tribute
  done: boolean;
}

export interface TributeFlowProps {
  tributePhase: TributePhase;
  players: (Player | null)[];
  currentPlayerSeat: number;
  playerHand: Card[];
  selectedCards: Card[];
  onCardSelect: (cards: Card[]) => void;
  onSelectTribute: (cardIndex: number) => void;
  onReturnTribute: (cardIndex: number) => void;
}

export interface PhaseContentProps extends TributeFlowProps {
  phase: UIPhase;
}
