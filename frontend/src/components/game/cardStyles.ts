import type { Card } from '../../types';

// 卡片尺寸定义
export type CardSize = 'xs' | 'small' | 'normal' | 'large';

export interface SizeConfig {
  width: string;
  height: string;
  fontSize: string;
  iconSize: string;
  centerIconSize: string;
  borderRadius: string;
  padding: string;
}

export const CARD_SIZES: Record<CardSize, SizeConfig> = {
  xs: {
    width: '32px',
    height: '48px',
    fontSize: '10px',
    iconSize: '10px',
    centerIconSize: '14px',
    borderRadius: '4px',
    padding: '2px',
  },
  small: {
    width: '48px',
    height: '68px',
    fontSize: '12px',
    iconSize: '12px',
    centerIconSize: '20px',
    borderRadius: '6px',
    padding: '4px',
  },
  normal: {
    width: '70px',
    height: '98px',
    fontSize: '16px',
    iconSize: '14px',
    centerIconSize: '32px',
    borderRadius: '8px',
    padding: '6px',
  },
  large: {
    width: '90px',
    height: '126px',
    fontSize: '20px',
    iconSize: '18px',
    centerIconSize: '40px',
    borderRadius: '10px',
    padding: '8px',
  }
};

// 花色配置
export const SUIT_SYMBOLS = ['♠', '♥', '♣', '♦'];

export const SUIT_COLORS = {
  black: {
    text: 'text-gray-900',
    shadow: 'text-shadow-sm',
  },
  red: {
    text: 'text-transparent bg-clip-text bg-gradient-to-br from-red-600 to-red-500',
    shadow: 'drop-shadow-sm', // 使用 drop-shadow 因为 text-shadow 对透明文字无效
  }
};

// 特殊牌型配置
export const JOKER_CONFIG = {
  small: {
    rank: 15,
    text: '小王',
    color: 'text-gray-900',
    bgGradient: 'bg-gradient-to-br from-gray-50 to-gray-200',
    icon: '⭐',
  },
  big: {
    rank: 16,
    text: '大王',
    color: 'text-red-600',
    bgGradient: 'bg-gradient-to-br from-yellow-50 via-yellow-100 to-red-100',
    icon: '👑',
  }
};

// 动画配置
export const ANIMATIONS = {
  hover: 'hover:scale-105 hover:shadow-lg',
  selected: 'transform -translate-y-5 shadow-xl',
  transition: 'transition-[transform,box-shadow] duration-200 ease-[cubic-bezier(0.34,1.56,0.64,1)]',
};

// 工具函数
export const getCardSizeStyle = (size: CardSize) => CARD_SIZES[size];

export const getSuitColorClass = (suit: number) => {
  return suit === 1 || suit === 3 ? SUIT_COLORS.red.text : SUIT_COLORS.black.text;
};

export const getSuitShadowClass = (suit: number) => {
  return suit === 1 || suit === 3 ? SUIT_COLORS.red.shadow : SUIT_COLORS.black.shadow;
};

export const getRankText = (rank: number): string => {
  if (rank === 15) return 'S'; // Small Joker 简写，完整显示在中央
  if (rank === 16) return 'B'; // Big Joker 简写
  if (rank <= 10) return rank.toString();
  switch (rank) {
    case 11: return 'J';
    case 12: return 'Q';
    case 13: return 'K';
    case 14: return 'A';
    default: return rank.toString();
  }
};
