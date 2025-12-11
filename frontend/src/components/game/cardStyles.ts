export type CardSize = 'small' | 'normal';

export interface SizeConfig {
  width: string;
  height: string;
  fontSize: string;
  iconSize: string;
  centerIconSize: string;
  borderRadius: string;
  padding: string;
  jokerFontSize: string;
  jokerImgBottom: string;
  jokerTextLeft: string;
  cornerTextLeft: string;
}

export const CARD_SIZES: Record<CardSize, SizeConfig> = {
  small: {
    width: '48px',
    height: '68px',
    fontSize: '12px',
    iconSize: '12px',
    centerIconSize: '20px',
    borderRadius: '6px',
    padding: '4px',
    jokerFontSize: '10px',
    jokerImgBottom: '10px',
    jokerTextLeft: '2px',
    cornerTextLeft: '0',
  },
  normal: {
    width: '70px',
    height: '98px',
    fontSize: '16px',
    iconSize: '14px',
    centerIconSize: '32px',
    borderRadius: '8px',
    padding: '6px',
    jokerFontSize: '14px',
    jokerImgBottom: '16px',
    jokerTextLeft: '2px',
    cornerTextLeft: '0',
  },
};

export const SUIT_SYMBOLS = ['♠', '♥', '♣', '♦'];

export const SUIT_COLORS = {
  black: {
    text: 'text-suit-black',
    shadow: 'text-shadow-sm',
  },
  red: {
    text: 'text-suit-red',
    shadow: 'drop-shadow-sm',
  }
};

export const JOKER_CONFIG = {
  small: {
    rank: 15,
    text: 'JOKER',
    color: 'text-fg-secondary',
    bgGradient: 'bg-surface-elevated',
  },
  big: {
    rank: 16,
    text: 'JOKER',
    color: 'text-suit-red',
    bgGradient: 'bg-gradient-to-br from-yellow-50 via-yellow-100 to-red-100',
  }
};

export const SELECTED_COLORS = {
  border: 'hsl(var(--color-state-active))',
  overlay: 'hsl(var(--color-state-active) / 0.15)',
  glow: 'var(--shadow-glow-md)',
};

export const LEVEL_BADGE_CONFIG = {
  normal: {
    bgColor: 'hsl(var(--badge-level))',
    text: '级',
  },
  wild: {
    bgColor: 'hsl(var(--badge-wild))',
    text: '万',
  }
};

export const ANIMATIONS = {
  hover: 'hover:scale-105 hover:shadow-elevation-3',
  selected: '',
  transition: 'transition-[transform,box-shadow,border-color] duration-fast ease-bounce',
};

export const getCardSizeStyle = (size: CardSize) => CARD_SIZES[size];

export const getSuitColorClass = (suit: number) => {
  return suit === 1 || suit === 3 ? SUIT_COLORS.red.text : SUIT_COLORS.black.text;
};

export const getSuitShadowClass = (suit: number) => {
  return suit === 1 || suit === 3 ? SUIT_COLORS.red.shadow : SUIT_COLORS.black.shadow;
};

export const getRankText = (rank: number): string => {
  if (rank === 15) return 'S';
  if (rank === 16) return 'B';
  if (rank <= 10) return rank.toString();
  switch (rank) {
    case 11: return 'J';
    case 12: return 'Q';
    case 13: return 'K';
    case 14: return 'A';
    default: return rank.toString();
  }
};

export const STACK_OVERLAP = {
  vertical: 0.7,
  horizontal: 0.7,
};
