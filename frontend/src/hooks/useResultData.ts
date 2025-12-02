import { useGameStore } from '../store/gameStore';
import { useRoomStore } from '../store/roomStore';
import type { Player } from '../types';

export const useDealResultData = () => {
  const dealResult = useGameStore(s => s.dealResult);
  const room = useRoomStore(s => s.currentRoom);
  const playerView = useGameStore(s => s.playerView);
  
  if (!dealResult || !room) return null;
  
  return {
    dealResult,
    players: room.players.filter(p => p !== null) as Player[],
    teamLevels: playerView?.teamLevels?.length === 2
      ? playerView.teamLevels as [number, number]
      : [2, 2] as [number, number]
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
