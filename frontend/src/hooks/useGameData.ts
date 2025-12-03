import { useGameStore } from '../store/gameStore';

export function usePlayerViewData() {
  const playerView = useGameStore(s => s.playerView);
  if (!playerView) return null;
  
  return {
    teamLevels: playerView.teamLevels.length === 2
      ? playerView.teamLevels as [number, number]
      : [2, 2] as [number, number],
    dealLevel: playerView.dealLevel,
    currentTurn: playerView.currentTurn,
    plays: playerView.plays,
    playStates: playerView.playStates.length === 4 
      ? playerView.playStates as [number, number, number, number] 
      : undefined,
  };
}

export function usePlayerHandData() {
  const playerView = useGameStore(state => state.playerView);
  return {
    hand: playerView?.playerCards ?? [],
  };
}
