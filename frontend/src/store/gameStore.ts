import { create } from 'zustand';
import type { WSMessage } from '../types';
import type { PlayerView } from '../types/proto';

interface GameState {
  isInGame: boolean;
  playerView: PlayerView | null;
  playerSeat: number | null;
  isMyTurn: boolean;
  countdown: number | null;
  lastMessage: WSMessage | null;
  isConnected: boolean;
  error: string | null;
  currentPhase: string;
}

interface GameActions {
  setInGame: (inGame: boolean) => void;
  setPlayerView: (view: PlayerView | null) => void;
  setPlayerSeat: (seat: number | null) => void;
  setMyTurn: (isMyTurn: boolean) => void;
  setCountdown: (countdown: number | null) => void;
  setLastMessage: (message: WSMessage | null) => void;
  setConnected: (connected: boolean) => void;
  setError: (error: string | null) => void;
  clearError: () => void;
  reset: () => void;
  setCurrentPhase: (phase: string) => void;
}

type GameStore = GameState & GameActions;

const initialState: GameState = {
  isInGame: false,
  playerView: null,
  playerSeat: null,
  isMyTurn: false,
  countdown: null,
  lastMessage: null,
  isConnected: false,
  error: null,
  currentPhase: 'waiting_players'
};

export const useGameStore = create<GameStore>((set) => ({
  ...initialState,

  // Actions
  setInGame: (inGame: boolean) => set({ isInGame: inGame }),
  
  setPlayerView: (playerView: PlayerView | null) => set({ playerView }),
  
  setPlayerSeat: (seat: number | null) => set({ playerSeat: seat }),
  
  setMyTurn: (isMyTurn: boolean) => set({ isMyTurn }),
  
  setCountdown: (countdown: number | null) => set({ countdown }),
  
  setLastMessage: (message: WSMessage | null) => set({ lastMessage: message }),
  
  setConnected: (connected: boolean) => set({ isConnected: connected }),
  
  setError: (error: string | null) => set({ error }),
  
  clearError: () => set({ error: null }),
  
  setCurrentPhase: (phase: string) => set({ currentPhase: phase }),
  
  reset: () => set(initialState)
}));