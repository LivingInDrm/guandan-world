import React, { useMemo } from 'react';
import type { Card } from '../../types/proto';
import { isJoker } from '../../utils/cardUtils';
import { cn } from '@/lib/utils';
import { 
  SUIT_SYMBOLS, 
  JOKER_CONFIG, 
  ANIMATIONS,
  STACK_OVERLAP,
  SELECTED_COLORS,
  LEVEL_BADGE_CONFIG,
  type CardSize,
  type SizeConfig,
  getCardSizeStyle,
  getSuitColorClass,
  getRankText
} from './cardStyles';
import bigJokerImg from '../../assets/big_joker.png';
import smallJokerImg from '../../assets/small_joker.png';
import spadeImg from '../../assets/spade.png';
import heartImg from '../../assets/heart.png';
import clubImg from '../../assets/club.png';
import diamondImg from '../../assets/diamond.png';

const SUIT_IMAGES = [spadeImg, heartImg, clubImg, diamondImg];

interface CardDisplayProps {
  card: Card;
  isSelected?: boolean;
  onClick?: () => void;
  onPointerEnter?: () => void;
  deckIndex?: number;
  stackIndex?: number;
  stackDirection?: 'horizontal' | 'vertical';
  size?: CardSize;
  className?: string;
  currentLevel?: number;
}

// 左上角标组件
const CardCorner: React.FC<{
  rank: number;
  suit: number;
  sizeConfig: SizeConfig;
  size: CardSize;
}> = ({ rank, suit, sizeConfig, size }) => {
  const colorClass = getSuitColorClass(suit);
  const isNormal = size === 'normal';
  
  return (
    <div 
      className={cn(
        'absolute flex leading-none top-0.5',
        isNormal ? 'flex-row items-center' : 'flex-col items-center',
        colorClass
      )}
      style={{ 
        padding: sizeConfig.padding,
        width: isNormal ? 'auto' : '1.2em',
        left: sizeConfig.cornerTextLeft,
        gap: isNormal ? '2px' : undefined,
      }}
    >
      <span style={{ fontSize: sizeConfig.fontSize, fontWeight: 700 }}>
        {getRankText(rank)}
      </span>
      <span style={{ fontSize: sizeConfig.iconSize, marginTop: isNormal ? undefined : '-1px' }}>
        {SUIT_SYMBOLS[suit]}
      </span>
    </div>
  );
};

// 右下大花色组件
const CardSuitLarge: React.FC<{
  suit: number;
  sizeConfig: SizeConfig;
  size: CardSize;
}> = ({ suit, sizeConfig, size }) => {
  const scale = size === 'normal' ? 1.2 : 1.5;
  const isNormal = size === 'normal';
  return (
    <img 
      src={SUIT_IMAGES[suit]}
      alt={SUIT_SYMBOLS[suit]}
      className="absolute"
      style={{ 
        width: `calc(${sizeConfig.centerIconSize} * ${scale})`,
        height: `calc(${sizeConfig.centerIconSize} * ${scale})`,
        objectFit: 'contain',
        bottom: isNormal ? 'calc(var(--base-unit) * 5)' : 'calc(var(--base-unit) * 2)',
        right: isNormal ? 'calc(var(--base-unit) * 3.5)' : 'calc(var(--base-unit) * 2)',
      }}
    />
  );
};

// 级牌角标组件
const LevelBadge: React.FC<{
  suit: number;
  size: CardSize;
}> = ({ suit, size }) => {
  const isHeart = suit === 1;
  const config = isHeart ? LEVEL_BADGE_CONFIG.wild : LEVEL_BADGE_CONFIG.normal;
  const badgeSize = size === 'normal' ? 24 : 18;
  const fontSize = size === 'normal' ? 12 : 9;
  const isNormal = size === 'normal';
  
  const positionClass = isNormal 
    ? "absolute top-0 right-0" 
    : "absolute bottom-0 left-0";
  
  const trianglePoints = isNormal
    ? `${badgeSize},0 ${badgeSize},${badgeSize} 0,0`
    : `0,${badgeSize} ${badgeSize},${badgeSize} 0,0`;
  
  const textX = badgeSize * (isNormal ? 0.68 : 0.32);
  const textY = badgeSize * (isNormal ? 0.42 : 0.69);
  
  return (
    <div 
      className={cn(positionClass, 'pointer-events-none')}
      style={{
        width: badgeSize,
        height: badgeSize,
      }}
    >
      <svg width={badgeSize} height={badgeSize} viewBox={`0 0 ${badgeSize} ${badgeSize}`}>
        <polygon 
          points={trianglePoints}
          fill={config.bgColor}
        />
        <text
          x={textX}
          y={textY}
          fill="white"
          fontSize={fontSize}
          fontWeight="bold"
          textAnchor="middle"
          dominantBaseline="middle"
        >
          {config.text}
        </text>
      </svg>
    </div>
  );
};

// 大小王专属内容
const JokerContent: React.FC<{
  card: Card;
  sizeConfig: SizeConfig;
}> = ({ card, sizeConfig }) => {
  const isBigJoker = card.rank === 16;
  const config = isBigJoker ? JOKER_CONFIG.big : JOKER_CONFIG.small;
  const jokerImg = isBigJoker ? bigJokerImg : smallJokerImg;
  
  return (
    <div className="absolute inset-0">
      <div 
        className={cn(
          'absolute top-1 flex flex-col items-center leading-tight font-bold',
          config.color
        )}
        style={{ 
          fontSize: sizeConfig.jokerFontSize,
          left: sizeConfig.jokerTextLeft,
        }}
      >
        {'JOKER'.split('').map((letter, i) => (
          <span key={i}>{letter}</span>
        ))}
      </div>
      
      <img 
        src={jokerImg}
        alt={isBigJoker ? 'Big Joker' : 'Small Joker'}
        className="absolute right-1"
        style={{ 
          bottom: sizeConfig.jokerImgBottom,
          width: `calc(${sizeConfig.centerIconSize} * 1.6)`,
          height: `calc(${sizeConfig.centerIconSize} * 1.6)`,
          objectFit: 'contain'
        }}
      />
    </div>
  );
};

const CardDisplay: React.FC<CardDisplayProps> = ({ 
  card, 
  isSelected = false, 
  onClick,
  onPointerEnter,
  deckIndex,
  stackIndex = 0,
  stackDirection = 'horizontal',
  size = 'normal',
  className = '',
  currentLevel
}) => {
  const sizeConfig = useMemo(() => getCardSizeStyle(size), [size]);
  
  const bgClass = useMemo(() => {
    if (isJoker(card)) {
      return card.rank === 16 ? JOKER_CONFIG.big.bgGradient : JOKER_CONFIG.small.bgGradient;
    }
    return 'bg-surface-base';
  }, [card]);
  
  const zIndex = stackIndex;

  return (
    <div
      role="button"
      aria-pressed={isSelected}
      data-deck-index={deckIndex}
      onPointerEnter={onPointerEnter}
      className={cn(
        'relative select-none cursor-pointer',
        'shadow-elevation-2',
        bgClass,
        ANIMATIONS.transition,
        !isSelected && onClick && ANIMATIONS.hover,
        className
      )}
      style={{
        width: sizeConfig.width,
        height: sizeConfig.height,
        borderRadius: sizeConfig.borderRadius,
        zIndex,
        marginLeft: stackDirection === 'horizontal' && stackIndex > 0 ? `${-parseFloat(sizeConfig.width) * STACK_OVERLAP.horizontal}px` : '0',
        marginTop: stackDirection === 'vertical' && stackIndex > 0 ? `${-parseFloat(sizeConfig.height) * STACK_OVERLAP.vertical}px` : '0',
        willChange: 'transform',
        border: '1px solid hsl(var(--color-border-base))',
        outline: isSelected ? `2px solid ${SELECTED_COLORS.border}` : 'none',
        boxShadow: isSelected 
          ? 'var(--shadow-glow-md), var(--elevation-2)' 
          : undefined,
      }}
      onClick={onClick}
    >
      {isSelected && (
        <div 
          className="absolute inset-0 pointer-events-none"
          style={{ 
            background: SELECTED_COLORS.overlay,
            borderRadius: `calc(${sizeConfig.borderRadius} - 2px)`,
          }}
        />
      )}
      {!isJoker(card) && (
        <div className="absolute inset-px border border-stroke/30 rounded opacity-50 pointer-events-none" />
      )}

      {isJoker(card) ? (
        <JokerContent card={card} sizeConfig={sizeConfig} />
      ) : (
        <>
          <CardCorner 
            rank={card.rank} 
            suit={card.suit} 
            sizeConfig={sizeConfig}
            size={size}
          />
          
          <CardSuitLarge 
            suit={card.suit} 
            sizeConfig={sizeConfig}
            size={size}
          />
          
          {currentLevel !== undefined && card.rank === currentLevel && (
            <LevelBadge suit={card.suit} size={size} />
          )}
        </>
      )}
    </div>
  );
};

// 导出辅助函数以保持向后兼容
export { getRankText };

export default React.memo(CardDisplay);
