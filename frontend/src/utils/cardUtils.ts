import type { Card as ProtoCard } from '../types/proto';

export function isJoker(card: ProtoCard): boolean {
  return card.suit === -1 || card.rank === 15 || card.rank === 16;
}
