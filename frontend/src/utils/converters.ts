// Type converters for proto messages to frontend types
import type { PlayerView as ProtoPlayerView, TributeView as ProtoTributeView, TributePair } from '../types/proto';
import type { PlayerView, PlayerGameState, TributePhase } from '../types';
import { convertProtoCardsToFrontend, convertProtoPlaysToFrontend, type FrontendCard } from './cardUtils';

// Convert proto PlayerView to frontend PlayerView format
export function convertProtoPlayerView(proto: ProtoPlayerView): PlayerView {
  return {
    player_seat: proto.playerSeat,
    player_cards: convertProtoCardsToFrontend(proto.playerCards),
    team_levels: proto.teamLevels as [number, number],
    deal_level: proto.dealLevel,
    deal_status: proto.dealStatus,  // proto 枚举值可以直接使用
    trick_id: proto.trickIndex?.toString(),
    current_turn: proto.currentTurn,
    plays: convertProtoPlaysToFrontend(proto.plays),
  };
}

// Convert proto PlayerView to frontend PlayerGameState
export function convertProtoPlayerViewToGameState(proto: ProtoPlayerView): PlayerGameState {
  return {
    team_levels: proto.teamLevels as [number, number],
    deal_level: proto.dealLevel,
    deal_status: proto.dealStatus,  // proto 枚举值可以直接使用
    trick_id: proto.trickIndex?.toString(),
    current_turn: proto.currentTurn,
    plays: convertProtoPlaysToFrontend(proto.plays),
  };
}

// Convert proto TributeView to frontend TributePhase
export function convertProtoTributeView(proto: ProtoTributeView): TributePhase {
  // Build tribute_map from tributePairs
  const tribute_map: { [giver: number]: number } = {};
  const tribute_cards: { [giver: number]: FrontendCard } = {};
  const return_cards: { [receiver: number]: FrontendCard } = {};
  const selection_results: { [receiver: number]: number } = {};
  
  proto.tributePairs.forEach((pair: TributePair) => {
    tribute_map[pair.giver] = pair.receiver;
    
    if (pair.tributeCard) {
      tribute_cards[pair.giver] = convertProtoCardsToFrontend([pair.tributeCard])[0];
    }
    
    if (pair.returnCard) {
      return_cards[pair.receiver] = convertProtoCardsToFrontend([pair.returnCard])[0];
      // Track which giver the return card came from
      selection_results[pair.receiver] = pair.giver;
    }
  });
  
  return {
    status: proto.status,  // proto 枚举值可以直接使用
    tribute_map,
    tribute_cards,
    return_cards,
    pool_cards: convertProtoCardsToFrontend(proto.poolCards),
    selecting_player: proto.selectingPlayer,
    select_timeout: new Date(Number(proto.updatedAtMs)).toISOString(),
    is_immune: proto.isImmune,
    selection_results,
  };
}
