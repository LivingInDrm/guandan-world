export type CardSize = 'small' | 'normal';

export const SUIT_SYMBOLS = ['♠', '♥', '♣', '♦'];

export const SUIT_COLORS = {
  black: {
    text: 'text-suit-black',
    shadow: 'drop-shadow-[0_1px_1px_rgba(0,0,0,0.2)]',
  },
  red: {
    text: 'text-suit-red',
    shadow: 'drop-shadow-[0_1px_1px_rgba(180,0,0,0.2)]',
  }
};

export const JOKER_CONFIG = {
  small: {
    rank: 15,
    text: 'JOKER',
    color: 'text-[hsl(158,55%,35%)]',
    bgGradient: 'bg-gradient-to-br from-[hsl(158,20%,95%)] via-[hsl(158,25%,90%)] to-[hsl(158,30%,85%)]',
  },
  big: {
    rank: 16,
    text: 'JOKER',
    color: 'text-suit-red',
    bgGradient: 'bg-gradient-to-br from-[hsl(42,60%,95%)] via-[hsl(42,70%,88%)] to-[hsl(0,60%,90%)]',
  }
};

export const SELECTED_COLORS = {
  border: 'hsl(42, 95%, 52%)',
  overlay: 'hsla(42, 95%, 52%, 0.12)',
  glow: '0 0 12px hsla(42, 95%, 52%, 0.5), 0 0 4px hsla(42, 95%, 52%, 0.3)',
};

export const LEVEL_BADGE_CONFIG = {
  normal: {
    bgColor: 'hsl(158, 55%, 38%)',
    text: '级',
  },
  wild: {
    bgColor: 'hsl(0, 70%, 50%)',
    text: '万',
  }
};

export const ANIMATIONS = {
  hover: 'hover:scale-105 hover:-translate-y-1 hover:shadow-[var(--elevation-3),var(--shadow-relief)]',
  selected: 'scale-105',
  transition: 'transition-all duration-150 ease-[cubic-bezier(0.34,1.56,0.64,1)]',
};

export const CARD_BORDER = {
  normal: 'border border-stroke/40',
  hover: 'hover:border-stroke/60',
  selected: 'border-2 border-state-active',
};

export const CARD_BG = {
  base: 'bg-gradient-to-br from-[hsl(40,15%,98%)] via-[hsl(40,10%,96%)] to-[hsl(40,8%,93%)]',
  shine: 'before:absolute before:inset-0 before:bg-gradient-to-br before:from-white/40 before:via-transparent before:to-transparent before:rounded-[inherit] before:pointer-events-none',
};

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
