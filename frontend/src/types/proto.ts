// Proto 生成的类型（数据传输格式）
// 从 generated 目录选择性导出需要的类型，避免重复导出辅助类型

// 从 common.proto 导出
export type { Card, PlayAction, PlayerBasicInfo } from './generated/common';
export { CompType } from './generated/common';

// 从 view.proto 导出
export type { PlayerView, TributeView, TributePair } from './generated/view';
export { DealStatus, TributeStatus } from './generated/view';

// 从 event.proto 导出
export type { 
  GameEvent,
  MatchStartedPayload,
  MatchEndedPayload,
  DealStartedPayload,
  CardsDealtPayload,
  DealEndedPayload,
  TributeStartedPayload,
  TributeExemptedPayload,
  TributeCardSubmittedPayload,
  TributeCardSelectedPayload,
  TributeCardReturnedPayload,
  TributeCompletedPayload,
  TrickStartedPayload,
  TrickEndedPayload,
  PlayerPlayedPayload,
  PlayerPassedPayload,
  PlayerTimeoutPayload,
  PlayerDisconnectPayload,
  PlayerReconnectPayload,
} from './generated/event';

export { 
  EventType, 
  VictoryType,
  TributeType,
  PlayerTimeoutActionType,
} from './generated/event';

// 从 action.proto 导出
export type { GameAction, TurnDeadline } from './generated/action';
export { GameActionType } from './generated/action';

