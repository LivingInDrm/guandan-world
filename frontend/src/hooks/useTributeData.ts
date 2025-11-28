import { useTributeStore } from '../store/tributeStore';
import { useRoomStore } from '../store/roomStore';
import { useGameStore } from '../store/gameStore';
import type { Player } from '../types';

export const useTributeData = () => {
  const step = useTributeStore((s) => s.step);
  const tributeStarted = useTributeStore((s) => s.tributeStarted);
  const tributeExempted = useTributeStore((s) => s.tributeExempted);
  const submittedCards = useTributeStore((s) => s.submittedCards);
  const poolCards = useTributeStore((s) => s.poolCards);
  const selectedCards = useTributeStore((s) => s.selectedCards);
  const returnedCards = useTributeStore((s) => s.returnedCards);
  const messages = useTributeStore((s) => s.messages);
  const currentSelectingSeat = useTributeStore((s) => s.currentSelectingSeat);
  const flyingCards = useTributeStore((s) => s.flyingCards);

  const room = useRoomStore((s) => s.currentRoom);
  const playerSeat = useGameStore((s) => s.playerSeat);

  if (step === 'idle' || !room) return null;

  return {
    step,
    tributeStarted,
    tributeExempted,
    submittedCards,
    poolCards,
    selectedCards,
    returnedCards,
    messages,
    currentSelectingSeat,
    flyingCards,
    players: room.players.filter((p) => p !== null) as Player[],
    playerSeat,
  };
};
