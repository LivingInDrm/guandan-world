import { useGameStore } from '../store/gameStore';
import { useRoomStore } from '../store/roomStore';
import type { Player, PlayerGameState } from '../types';

export const useDealResultData = () => {
  const dealResult = useGameStore(s => s.dealResult);
  const room = useRoomStore(s => s.currentRoom);
  const gameState = useGameStore(s => s.gameState);
  
  if (!dealResult || !room) return null;
  
  return {
    dealResult,
    players: room.players.filter(p => p !== null) as Player[],
    teamLevels: (gameState as PlayerGameState)?.team_levels || [2, 2] as [number, number]
  };
};

export const useMatchResultData = () => {
  const matchResult = useGameStore(s => s.matchResult);
  const room = useRoomStore(s => s.currentRoom);
  
  if (!matchResult || !room) return null;
  
  return {
    matchResult,
    players: room.players.filter(p => p !== null) as Player[]
  };
};
