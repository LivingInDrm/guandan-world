// ========== Proto 生成的类型 ==========
// 从 proto 文件自动生成的类型（数据传输格式）
export type {
  Card as ProtoCard,
  PlayAction as ProtoPlayAction,
  PlayerBasicInfo,
  CompType,
} from './proto';

export {
  DealStatus as ProtoDealStatus,
  TributeStatus as ProtoTributeStatus,
} from './generated/view';

// 为前端提供方便的常量（兼容旧代码）
import { 
  DealStatus as ProtoDealStatusEnum, 
  TributeStatus as ProtoTributeStatusEnum 
} from './generated/view';

export const DealStatus = {
  UNSPECIFIED: ProtoDealStatusEnum.DEAL_STATUS_UNSPECIFIED,
  WAITING: ProtoDealStatusEnum.DEAL_STATUS_WAITING,
  DEALING: ProtoDealStatusEnum.DEAL_STATUS_DEALING,
  TRIBUTE: ProtoDealStatusEnum.DEAL_STATUS_TRIBUTE,
  PLAYING: ProtoDealStatusEnum.DEAL_STATUS_PLAYING,
  FINISHED: ProtoDealStatusEnum.DEAL_STATUS_FINISHED,
} as const;

export const TributeStatus = {
  UNSPECIFIED: ProtoTributeStatusEnum.TRIBUTE_STATUS_UNSPECIFIED,
  WAITING: ProtoTributeStatusEnum.TRIBUTE_STATUS_WAITING,
  SELECTING: ProtoTributeStatusEnum.TRIBUTE_STATUS_SELECTING,
  RETURNING: ProtoTributeStatusEnum.TRIBUTE_STATUS_RETURNING,
  FINISHED: ProtoTributeStatusEnum.TRIBUTE_STATUS_FINISHED,
} as const;

export type {
  PlayerView as ProtoPlayerView,
  TributeView as ProtoTributeView,
  TributePair,
} from './proto';

export type {
  GameEvent,
  EventType,
  VictoryType as ProtoVictoryType,
} from './proto';

// ========== 类型别名 ==========
// 直接使用 Proto 类型，不再需要前端扩展类型
import type { Card as ProtoCardType, PlayAction as ProtoPlayActionType } from './proto';

export type Card = ProtoCardType;
export type PlayAction = ProtoPlayActionType;

// Base types for the application

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
  TRIBUTE_VIEW: 'tribute_view',
  ROOM_UPDATE: 'room_update',
  GAME_ACTION: 'game_action',
  
  // Timeout and auto-play
  PLAYER_TIMEOUT: 'player_timeout',
  AUTO_PLAY: 'auto_play',
  RECONNECT: 'reconnect',
  
  // Turn deadline broadcast
  TURN_DEADLINE: 'turn_deadline',

  // Client sync request
  SYNC_GAME_STATE: 'sync_game_state'
} as const;

export type WSMessageType = typeof WS_MESSAGE_TYPES[keyof typeof WS_MESSAGE_TYPES];

// ========== Game Card and Action Types (使用 Proto) ==========
// Card 和 PlayAction 现在从 proto 生成，见文件顶部导出

// Game related types

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

export interface GameActionData {
  action_type: 'play_decision_required' | 'tribute_selection_required' | 'return_tribute_required';
  player_seat: number;
  hand?: Card[];
  trick_info?: TrickInfo;
  timeout?: number;
  options?: Card[];
}

export interface TurnDeadlineData {
  player_seat: number;
  action_type: 'play_decision' | 'tribute_selection' | 'return_tribute';
  deadline_at_ms: number;
}

// ========== Deal and Tribute Status (使用 Proto) ==========
// DealStatus 和 TributeStatus 现在从 proto/view.ts 导出，见文件顶部

// Deal result related types
// Note: VictoryType 直接使用 Proto 枚举，与 DealStatus、TributeStatus 保持一致
// 直接使用 proto 类型，避免重复定义和转换
export type { DealEndedPayload, PlayerDealStats, MatchEndedPayload } from './generated/event';