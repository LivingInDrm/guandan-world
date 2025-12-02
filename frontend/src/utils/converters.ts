// Type converters for proto messages to frontend types
import type { PlayerView as ProtoPlayerView } from '../types/proto';
import type { PlayerView, PlayerGameState } from '../types';
import { convertProtoCardsToFrontend, convertProtoPlaysToFrontend } from './cardUtils';

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
    play_states: proto.playStates?.length === 4 
      ? proto.playStates as [number, number, number, number] 
      : undefined,
  };
}
