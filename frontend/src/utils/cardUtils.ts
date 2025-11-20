import type { Card as ProtoCard, PlayAction as ProtoPlayAction } from '../types/proto';
import type { FrontendCard as FrontendCardType } from '../types/frontend';

// Re-export FrontendCard for convenience
export type { FrontendCard } from '../types/frontend';

// Convert proto Card format to frontend Card format
// Proto Card: { suit: number, rank: number, deckIndex: number }
// Frontend Card: extends ProtoCard with { id: string, isJoker: boolean }
export function convertProtoCardToFrontend(protoCard: ProtoCard): FrontendCardType {
  const suit = protoCard.suit ?? -1;
  const rank = protoCard.rank ?? 0;
  const deckIndex = protoCard.deckIndex ?? 0;
  
  // Convert suit number to color string
  let color = '';
  switch (suit) {
    case 0:
      color = 'Spade';
      break;
    case 1:
      color = 'Heart';
      break;
    case 2:
      color = 'Club';
      break;
    case 3:
      color = 'Diamond';
      break;
    case -1:
    default:
      color = 'Joker';
      break;
  }
  
  // Generate card ID in the format: Color_Number_DeckIndex
  const id = `${color}_${rank}_${deckIndex}`;
  
  // Check if it's a joker (rank 15 or 16, or suit -1)
  const isJoker = suit === -1 || rank === 15 || rank === 16;
  
  return {
    ...protoCard,
    id,
    isJoker
  };
}

// Convert array of proto cards to frontend cards
export function convertProtoCardsToFrontend(protoCards: ProtoCard[]): FrontendCardType[] {
  if (!protoCards || !Array.isArray(protoCards)) {
    return [];
  }
  return protoCards.map(convertProtoCardToFrontend);
}

// Note: PlayAction from proto now uses camelCase (playerSeat, isPass, timestampMs)
// We keep it as-is and add frontend Card conversion
export interface FrontendPlayAction {
  playerSeat: number;
  cards: FrontendCardType[];
  isPass: boolean;
  timestamp: string;
  compType?: number;
}

// Convert proto PlayAction to frontend PlayAction with FrontendCard
export function convertProtoPlayActionToFrontend(protoPlay: ProtoPlayAction): FrontendPlayAction {
  return {
    playerSeat: protoPlay.playerSeat,
    cards: convertProtoCardsToFrontend(protoPlay.cards),
    isPass: protoPlay.isPass,
    timestamp: protoPlay.timestampMs ? new Date(Number(protoPlay.timestampMs)).toISOString() : new Date().toISOString(),
    compType: protoPlay.compType
  };
}

// Convert array of proto PlayActions to frontend PlayActions
export function convertProtoPlaysToFrontend(protoPlays: ProtoPlayAction[]): FrontendPlayAction[] {
  if (!protoPlays || !Array.isArray(protoPlays)) {
    return [];
  }
  return protoPlays.map(convertProtoPlayActionToFrontend);
}
