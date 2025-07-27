import { create } from 'zustand';
import type { WSMessage } from '../types';

interface GameState {
  isInGame: boolean;
  gameState: any; // Will be typed more specifically later
  playerSeat: number | null;
  isMyTurn: boolean;
  countdown: number | null;
  lastMessage: WSMessage | null;
  isConnected: boolean;
  error: string | null;
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
  error: null
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
  
  reset: () => set(initialState)
}));