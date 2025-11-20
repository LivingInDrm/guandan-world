// 前端扩展的类型定义
// 在 proto 类型基础上添加 UI 需要的额外字段

import type { Card as ProtoCard } from './proto';

// 前端扩展的 Card 类型（包含 UI 需要的 id 字段）
// id 格式: Color_Rank_DeckIndex (例如: "Spade_13_12")
export interface FrontendCard extends ProtoCard {
  id: string;
  isJoker: boolean;
}
