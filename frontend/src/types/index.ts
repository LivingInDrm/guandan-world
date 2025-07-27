// Base types for the application

export interface User {
  id: string;
  username: string;
  online: boolean;
}

export interface AuthToken {
  token: string;
  expires_at: string;
  user_id: string;
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface RegisterRequest {
  username: string;
  password: string;
}

export interface AuthResponse {
  user: User;
  token: AuthToken;
}

// Room related types
export const RoomStatus = {
  WAITING: 0,
  READY: 1,
  PLAYING: 2,
  CLOSED: 3
} as const;

export type RoomStatus = typeof RoomStatus[keyof typeof RoomStatus];

export interface Player {
  id: string;
  username: string;
  seat: number;
  online: boolean;
  auto_play: boolean;
}

export interface Room {
  id: string;
  status: RoomStatus;
  players: (Player | null)[];
  owner: string;
  created_at: string;
}

export interface RoomInfo {
  id: string;
  status: RoomStatus;
  player_count: number;
  players: Player[];
  owner: string;
  can_join: boolean;
}

export interface RoomListResponse {
  rooms: RoomInfo[];
  total_count: number;
  page: number;
  limit: number;
}

export interface CreateRoomRequest {
  // Room creation parameters if needed
}

export interface JoinRoomRequest {
  room_id: string;
}

// API Response wrapper
export interface ApiResponse<T = any> {
  success: boolean;
  data?: T;
  error?: string;
  message?: string;
}

// WebSocket message types
export interface WSMessage {
  type: string;
  data: any;
  timestamp: string;
  player_id?: string;
}

// WebSocket message type constants
export const WS_MESSAGE_TYPES = {
  // Room management
  JOIN_ROOM: 'join_room',
  LEAVE_ROOM: 'leave_room',
  START_GAME: 'start_game',
  
  // Game flow
  GAME_PREPARE: 'game_prepare',
  COUNTDOWN: 'countdown',
  GAME_BEGIN: 'game_begin',
  
  // Game operations
  PLAY_DECISION: 'play_decision',
  TRIBUTE_SELECT: 'tribute_select',
  TRIBUTE_RETURN: 'tribute_return',
  
  // State sync
  GAME_EVENT: 'game_event',
  PLAYER_VIEW: 'player_view',
  ROOM_UPDATE: 'room_update',
  
  // Timeout and auto-play
  PLAYER_TIMEOUT: 'player_timeout',
  AUTO_PLAY: 'auto_play',
  RECONNECT: 'reconnect'
} as const;

export type WSMessageType = typeof WS_MESSAGE_TYPES[keyof typeof WS_MESSAGE_TYPES];

// Game related types
export interface Card {
  id: string;
  suit: number; // 0=spades, 1=hearts, 2=clubs, 3=diamonds
  rank: number; // 2-14 (11=J, 12=Q, 13=K, 14=A), 15=small joker, 16=big joker
  is_joker: boolean;
}

export interface PlayAction {
  player_seat: number;
  cards: Card[];
  is_pass: boolean;
  timestamp: string;
}

export interface TrickInfo {
  id: string;
  leader: number;
  current_turn: number;
  plays: PlayAction[];
  winner: number;
  lead_comp: any; // Card composition
  status: TrickStatus;
  start_time: string;
  turn_timeout: string;
  next_leader: number;
}

export const TrickStatus = {
  WAITING: 0,
  PLAYING: 1,
  FINISHED: 2
} as const;

export type TrickStatus = typeof TrickStatus[keyof typeof TrickStatus];

export const PlayerStatus = {
  WAITING: 'waiting',
  PLAYING: 'playing',
  PLAYED: 'played',
  PASSED: 'passed',
  FINISHED: 'finished'
} as const;

export type PlayerStatus = typeof PlayerStatus[keyof typeof PlayerStatus];

export interface GameState {
  match_id: string;
  current_deal: DealState;
  team_levels: [number, number]; // Team 0 and Team 1 levels
  current_player: number;
  status: string;
}

export interface DealState {
  id: string;
  level: number;
  status: string;
  current_trick: TrickInfo | null;
  player_cards: Card[][];
  rankings: number[];
  start_time: string;
}

export interface PlayerGameState {
  player_seat: number;
  hand: Card[];
  game_state: GameState;
  can_play: boolean;
  turn_timeout: string | null;
}

// Tribute phase related types
export const TributeStatus = {
  WAITING: 'waiting',
  SELECTING: 'selecting',
  RETURNING: 'returning',
  FINISHED: 'finished'
} as const;

export type TributeStatus = typeof TributeStatus[keyof typeof TributeStatus];

export const TributeActionType = {
  NONE: 'none',
  SELECT: 'select',
  RETURN: 'return'
} as const;

export type TributeActionType = typeof TributeActionType[keyof typeof TributeActionType];

export interface TributeAction {
  type: TributeActionType;
  player_id: number;
  options: Card[];
  target_player?: number;
  timeout: number; // seconds
}

export interface TributePhase {
  status: TributeStatus;
  tribute_map: { [giver: number]: number }; // giver -> receiver (-1 for pool)
  tribute_cards: { [giver: number]: Card };
  return_cards: { [receiver: number]: Card };
  pool_cards: Card[];
  selecting_player: number;
  select_timeout: string;
  is_immune: boolean;
  selection_results: { [receiver: number]: number }; // receiver -> original_giver
}

export interface TributeStatusInfo {
  phase: TributeStatus;
  tribute_cards: { [giver: number]: Card };
  return_cards: { [receiver: number]: Card };
  tribute_map: { [giver: number]: number };
  pending_actions: TributeAction[];
  is_immune: boolean;
}

// Deal result related types
export const VictoryType = {
  DOUBLE_DOWN: 'double_down',   // rank1, rank2 same team (+3 levels)
  SINGLE_LAST: 'single_last',   // rank1, rank3 same team (+2 levels)
  PARTNER_LAST: 'partner_last'  // rank1, rank4 same team (+1 level)
} as const;

export type VictoryType = typeof VictoryType[keyof typeof VictoryType];

export interface PlayerDealStats {
  player_seat: number;
  cards_played: number;
  tricks_won: number;
  pass_count: number;
  timeout_count: number;
  finish_rank: number; // 1-4, with 1 being first to finish
}

export interface TributeInfo {
  has_tribute: boolean;
  tribute_map: { [giver: number]: number }; // giver -> receiver
  tribute_cards: { [giver: number]: Card }; // giver -> card
  return_cards: { [receiver: number]: Card }; // receiver -> card
}

export interface DealStatistics {
  total_tricks: number;
  player_stats: PlayerDealStats[];
  tribute_info: TributeInfo;
}

export interface DealResult {
  rankings: number[];        // Order of finishing (seat numbers)
  winning_team: number;      // 0 or 1
  victory_type: VictoryType;
  upgrades: [number, number]; // Level upgrades for each team
  duration: number;          // Duration in milliseconds
  trick_count: number;
  statistics: DealStatistics;
}

export interface MatchStatistics {
  total_deals: number;
  total_duration: number;    // Duration in milliseconds
  final_levels: [number, number]; // Final levels for each team
  winner: number;            // Winning team (0 or 1)
}

export interface MatchResult {
  match_id: string;
  winner: number;           // Winning team (0 or 1)
  final_levels: [number, number];
  statistics: MatchStatistics;
  players: Player[];
}