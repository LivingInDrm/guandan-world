import { create } from 'zustand';
import type { WSMessage, Card, DealResult, MatchResult, TributePhase } from '../types';

interface GameState {
  isInGame: boolean;
  gameState: any; // Will be typed more specifically later
  playerSeat: number | null;
  isMyTurn: boolean;
  countdown: number | null;
  lastMessage: WSMessage | null;
  isConnected: boolean;
  error: string | null;
  currentPhase: string;
  selectedCards: Card[];
  dealResult: DealResult | null;
  matchResult: MatchResult | null;
  tributeInfo: TributePhase | null;
}

interface GameActions {
  setInGame: (inGame: boolean) => void;
  setGameState: (state: any) => void;
  setPlayerSeat: (seat: number | null) => void;
  setMyTurn: (isMyTurn: boolean) => void;
  setCountdown: (countdown: number | null) => void;
  setLastMessage: (message: WSMessage | null) => void;
  setConnected: (connected: boolean) => void;
  setError: (error: string | null) => void;
  clearError: () => void;
  reset: () => void;
  setCurrentPhase: (phase: string) => void;
  setSelectedCards: (cards: Card[]) => void;
  setDealResult: (result: DealResult | null) => void;
  setMatchResult: (result: MatchResult | null) => void;
  setTributeInfo: (info: TributePhase | null) => void;
}

type GameStore = GameState & GameActions;

const initialState: GameState = {
  isInGame: false,
  gameState: null,
  playerSeat: null,
  isMyTurn: false,
  countdown: null,
  lastMessage: null,
  isConnected: false,
  error: null,
  currentPhase: 'waiting_players',
  selectedCards: [],
  dealResult: null,
  matchResult: null,
  tributeInfo: null
};

export const useGameStore = create<GameStore>((set) => ({
  ...initialState,

  // Actions
  setInGame: (inGame: boolean) => set({ isInGame: inGame }),
  
  setGameState: (gameState: any) => set({ gameState }),
  
  setPlayerSeat: (seat: number | null) => set({ playerSeat: seat }),
  
  setMyTurn: (isMyTurn: boolean) => set({ isMyTurn }),
  
  setCountdown: (countdown: number | null) => set({ countdown }),
  
  setLastMessage: (message: WSMessage | null) => set({ lastMessage: message }),
  
  setConnected: (connected: boolean) => set({ isConnected: connected }),
  
  setError: (error: string | null) => set({ error }),
  
  clearError: () => set({ error: null }),
  
  setCurrentPhase: (phase: string) => set({ currentPhase: phase }),
  
  setSelectedCards: (cards: Card[]) => set({ selectedCards: cards }),
  
  setDealResult: (result: DealResult | null) => set({ dealResult: result }),
  
  setMatchResult: (result: MatchResult | null) => set({ matchResult: result }),
  
  setTributeInfo: (info: TributePhase | null) => set({ tributeInfo: info }),
  
  reset: () => set(initialState)
}));