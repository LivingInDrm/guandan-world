import type { Card, PlayAction } from '../types';

// Convert proto Card format to frontend Card format
// Proto Card: { suit: number, rank: number, deck_index: number }
// Frontend Card: { id: string, suit: number, rank: number, is_joker: boolean }
export function convertProtoCardToFrontend(protoCard: any): Card {
  const suit = protoCard.suit ?? -1;
  const rank = protoCard.rank ?? 0;
  const deckIndex = protoCard.deck_index ?? protoCard.deckIndex ?? 0;
  
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
    id,
    suit,
    rank,
    is_joker: isJoker
  };
}

// Convert array of proto cards to frontend cards
export function convertProtoCardsToFrontend(protoCards: any[]): Card[] {
  if (!protoCards || !Array.isArray(protoCards)) {
    return [];
  }
  return protoCards.map(convertProtoCardToFrontend);
}

// Convert proto PlayAction to frontend PlayAction
export function convertProtoPlayActionToFrontend(protoPlay: any): PlayAction {
  return {
    player_seat: protoPlay.player_seat ?? protoPlay.playerSeat ?? 0,
    cards: convertProtoCardsToFrontend(protoPlay.cards || []),
    is_pass: protoPlay.is_pass ?? protoPlay.isPass ?? false,
    timestamp: protoPlay.timestamp || new Date().toISOString()
  };
}

// Convert array of proto PlayActions to frontend PlayActions
export function convertProtoPlaysToFrontend(protoPlays: any[]): PlayAction[] {
  if (!protoPlays || !Array.isArray(protoPlays)) {
    return [];
  }
  return protoPlays.map(convertProtoPlayActionToFrontend);
}
