export interface User {
  id: string;
  username: string;
  nickname?: string;
  avatar_key?: string;
  online: boolean;
}

export interface AuthToken {
  access_token: string;
  refresh_token: string;
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
  nickname?: string;
  avatar_key?: string;
  seat: number;
  online: boolean;
  auto_play: boolean;
}

export interface Room {
  id: string;
  room_code?: string;
  status: RoomStatus;
  players: (Player | null)[];
  owner: string;
  created_at: string;
}

export interface ApiResponse<T = unknown> {
  success: boolean;
  data?: T;
  error?: string;
  message?: string;
}

export interface WSMessage {
  type: string;
  data: unknown;
  timestamp: string;
  player_id?: string;
}

export const WS_MESSAGE_TYPES = {
  JOIN_ROOM: 'join_room',
  LEAVE_ROOM: 'leave_room',
  START_GAME: 'start_game',
  GAME_PREPARE: 'game_prepare',
  COUNTDOWN: 'countdown',
  GAME_BEGIN: 'game_begin',
  PLAY_DECISION: 'play_decision',
  TRIBUTE_SELECT: 'tribute_select',
  TRIBUTE_RETURN: 'tribute_return',
  GAME_EVENT: 'game_event',
  PLAYER_VIEW: 'player_view',
  TRIBUTE_VIEW: 'tribute_view',
  ROOM_UPDATE: 'room_update',
  GAME_ACTION: 'game_action',
  PLAYER_TIMEOUT: 'player_timeout',
  AUTO_PLAY: 'auto_play',
  RECONNECT: 'reconnect',
  TURN_DEADLINE: 'turn_deadline',
  SYNC_GAME_STATE: 'sync_game_state'
} as const;

export type WSMessageType = typeof WS_MESSAGE_TYPES[keyof typeof WS_MESSAGE_TYPES];

export enum DealStatus {
  DEAL_STATUS_UNSPECIFIED = 0,
  DEAL_STATUS_WAITING = 1,
  DEAL_STATUS_DEALING = 2,
  DEAL_STATUS_TRIBUTE = 3,
  DEAL_STATUS_PLAYING = 4,
  DEAL_STATUS_FINISHED = 5,
  UNRECOGNIZED = -1,
}

export enum TributeStatus {
  TRIBUTE_STATUS_UNSPECIFIED = 0,
  TRIBUTE_STATUS_WAITING = 1,
  TRIBUTE_STATUS_SELECTING = 2,
  TRIBUTE_STATUS_RETURNING = 3,
  TRIBUTE_STATUS_FINISHED = 4,
  UNRECOGNIZED = -1,
}

export enum TributeType {
  TRIBUTE_TYPE_UNSPECIFIED = 0,
  TRIBUTE_TYPE_DOUBLE_DOWN = 1,
  TRIBUTE_TYPE_SINGLE_LAST = 2,
  TRIBUTE_TYPE_PARTNER_LAST = 3,
  TRIBUTE_TYPE_NONE = 4,
  UNRECOGNIZED = -1,
}

export enum CompType {
  COMP_TYPE_UNSPECIFIED = 0,
  COMP_TYPE_FOLD = 1,
  COMP_TYPE_ILLEGAL = 2,
  COMP_TYPE_SINGLE = 3,
  COMP_TYPE_PAIR = 4,
  COMP_TYPE_TRIPLE = 5,
  COMP_TYPE_FULL_HOUSE = 6,
  COMP_TYPE_STRAIGHT = 7,
  COMP_TYPE_PLATE = 8,
  COMP_TYPE_TUBE = 9,
  COMP_TYPE_JOKER_BOMB = 10,
  COMP_TYPE_NAIVE_BOMB = 11,
  COMP_TYPE_STRAIGHT_FLUSH = 12,
  UNRECOGNIZED = -1,
}

export interface Card {
  suit: number;
  rank: number;
  deckIndex: number;
}

export interface PlayAction {
  playerSeat: number;
  cards: Card[];
  compType: CompType;
  timestampMs: number;
  isPass: boolean;
}

export interface PlayerView {
  matchId: string;
  dealIndex: number;
  trickIndex?: number | undefined;
  seq: number;
  updatedAtMs: number;
  playerSeat: number;
  playerCards: Card[];
  teamLevels: number[];
  dealLevel: number;
  dealStatus: DealStatus;
  currentTurn?: number | undefined;
  plays: PlayAction[];
  leader?: number | undefined;
  playStates: number[];
  dealResult?: DealEndedPayload | undefined;
  matchResult?: MatchEndedPayload | undefined;
}

export interface TributePair {
  giver: number;
  receiver: number;
  tributeCard?: Card | undefined;
  returnCard?: Card | undefined;
}

export interface TributeView {
  matchId: string;
  dealIndex: number;
  seq: number;
  updatedAtMs: number;
  status: TributeStatus;
  tributePairs: TributePair[];
  poolCards: Card[];
  isImmune: boolean;
  tributeType: TributeType;
  givers: number[];
  receivers: number[];
}

export interface DealEndedPayload {
  matchId: string;
  dealIndex: number;
  winnerSeats: number[];
  teamLevelChanges: number[];
}

export interface MatchEndedPayload {
  matchId: string;
  winnerTeam: number;
  finalLevels: number[];
}

export enum GameActionType {
  GAME_ACTION_TYPE_UNSPECIFIED = 0,
  GAME_ACTION_TYPE_PLAY_DECISION = 1,
  GAME_ACTION_TYPE_TRIBUTE_SELECTION = 2,
  GAME_ACTION_TYPE_RETURN_TRIBUTE = 3,
}

export interface GameActionData {
  actionType: GameActionType;
  playerSeat: number;
  hand?: Card[];
  options?: Card[];
  timeout?: number;
}

export interface TurnDeadlineData {
  playerSeat: number;
  actionType: GameActionType;
  deadlineAtMs: number;
}
