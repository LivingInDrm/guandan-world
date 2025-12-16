import React, { useMemo } from 'react';
import type { Card } from '../../types/proto';
import { isJoker } from '../../utils/cardUtils';
import { cn } from '@/lib/utils';
import {
  SUIT_SYMBOLS,
  JOKER_CONFIG,
  ANIMATIONS,
  SELECTED_COLORS,
  LEVEL_BADGE_CONFIG,
  CARD_BG,
  type CardSize,
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

const CardCorner: React.FC<{
  rank: number;
  suit: number;
  size: CardSize;
}> = ({ rank, suit, size }) => {
  const colorClass = getSuitColorClass(suit);
  const isSmall = size === 'small';
  
  return (
    <div 
      className={cn(
        'absolute flex leading-none',
        isSmall ? 'flex-col items-center' : 'flex-row items-center',
        colorClass
      )}
      style={{ 
        top: `var(${isSmall ? '--card-corner-top-small' : '--card-corner-top'})`,
        left: `var(${isSmall ? '--card-padding-small' : '--card-padding'})`,
        gap: `var(${isSmall ? '--card-corner-gap-small' : '--card-corner-gap'})`,
      }}
    >
      <span style={{ fontSize: `var(${isSmall ? '--card-font-size-small' : '--card-font-size'})`, fontWeight: 700, letterSpacing: '-0.05em' }}>
        {getRankText(rank)}
      </span>
      <span style={{ fontSize: `var(${isSmall ? '--card-icon-size-small' : '--card-icon-size'})` }}>
        {SUIT_SYMBOLS[suit]}
      </span>
    </div>
  );
};

const CardSuitLarge: React.FC<{
  suit: number;
  size: CardSize;
}> = ({ suit, size }) => {
  const isSmall = size === 'small';
  const scale = isSmall ? 1.5 : 1.2;
  
  return (
    <img 
      src={SUIT_IMAGES[suit]}
      alt={SUIT_SYMBOLS[suit]}
      draggable="false"
      className="absolute"
      style={{ 
        width: `calc(var(${isSmall ? '--card-center-icon-small' : '--card-center-icon'}) * ${scale})`,
        height: `calc(var(${isSmall ? '--card-center-icon-small' : '--card-center-icon'}) * ${scale})`,
        objectFit: 'contain',
        bottom: `var(${isSmall ? '--card-suit-bottom-small' : '--card-suit-bottom'})`,
        right: `var(${isSmall ? '--card-suit-right-small' : '--card-suit-right'})`,
      }}
    />
  );
};

const LevelBadge: React.FC<{
  suit: number;
  size: CardSize;
}> = ({ suit, size }) => {
  const isHeart = suit === 1;
  const config = isHeart ? LEVEL_BADGE_CONFIG.wild : LEVEL_BADGE_CONFIG.normal;
  const isNormal = size === 'normal';

  const viewBoxSize = 24;
  const fontSize = 12;

  const positionClass = isNormal
    ? "absolute top-0 right-0"
    : "absolute bottom-0 left-0";

  const trianglePoints = isNormal
    ? `${viewBoxSize},0 ${viewBoxSize},${viewBoxSize} 0,0`
    : `0,${viewBoxSize} ${viewBoxSize},${viewBoxSize} 0,0`;

  const textX = viewBoxSize * (isNormal ? 0.68 : 0.32);
  const textY = viewBoxSize * (isNormal ? 0.42 : 0.69);

  const sizeVar = isNormal ? '--card-badge-size' : '--card-badge-size-small';

  return (
    <div
      className={cn(positionClass, 'pointer-events-none')}
      style={{
        width: `var(${sizeVar})`,
        height: `var(${sizeVar})`,
      }}
    >
      <svg
        width="100%"
        height="100%"
        viewBox={`0 0 ${viewBoxSize} ${viewBoxSize}`}
        preserveAspectRatio="none"
      >
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

const JokerContent: React.FC<{
  card: Card;
  size: CardSize;
}> = ({ card, size }) => {
  const isBigJoker = card.rank === 16;
  const config = isBigJoker ? JOKER_CONFIG.big : JOKER_CONFIG.small;
  const jokerImg = isBigJoker ? bigJokerImg : smallJokerImg;
  const isSmall = size === 'small';
  
  return (
    <div className="absolute inset-0">
      <div 
        className={cn(
          'absolute flex flex-col items-center font-bold',
          config.color
        )}
        style={{ 
          top: `var(${isSmall ? '--card-joker-top-small' : '--card-joker-top'})`,
          left: `var(${isSmall ? '--card-joker-text-left-small' : '--card-joker-text-left'})`,
          fontSize: `var(${isSmall ? '--card-joker-font-size-small' : '--card-joker-font-size'})`,
          lineHeight: `var(${isSmall ? '--card-joker-line-height-small' : '--card-joker-line-height'})`,
        }}
      >
        {'JOKER'.split('').map((letter, i) => (
          <span key={i}>{letter}</span>
        ))}
      </div>
      
      <img 
        src={jokerImg}
        alt={isBigJoker ? 'Big Joker' : 'Small Joker'}
        draggable="false"
        className="absolute"
        style={{ 
          right: `var(${isSmall ? '--card-joker-right-small' : '--card-joker-right'})`,
          bottom: `var(${isSmall ? '--card-joker-img-bottom-small' : '--card-joker-img-bottom'})`,
          width: `calc(var(${isSmall ? '--card-center-icon-small' : '--card-center-icon'}) * 1.6)`,
          height: `calc(var(${isSmall ? '--card-center-icon-small' : '--card-center-icon'}) * 1.6)`,
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
  const isSmall = size === 'small';

  const bgClass = useMemo(() => {
    if (isJoker(card)) {
      return card.rank === 16 ? JOKER_CONFIG.big.bgGradient : JOKER_CONFIG.small.bgGradient;
    }
    return CARD_BG.base;
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
        'shadow-[0_2px_4px_rgba(0,0,0,0.1),0_4px_8px_rgba(0,0,0,0.08),inset_0_1px_0_rgba(255,255,255,0.5)]',
        bgClass,
        ANIMATIONS.transition,
        !isSelected && onClick && ANIMATIONS.hover,
        isSelected && ANIMATIONS.selected,
        className
      )}
      style={{
        width: `var(${isSmall ? '--card-width-small' : '--card-width'})`,
        height: `var(${isSmall ? '--card-height-small' : '--card-height'})`,
        borderRadius: `var(${isSmall ? '--card-border-radius-small' : '--card-border-radius'})`,
        zIndex,
        marginLeft: stackDirection === 'horizontal' && stackIndex > 0
          ? `calc(-1 * var(${isSmall ? '--card-stack-offset-h-small' : '--card-stack-offset-h'}))`
          : '0',
        marginTop: stackDirection === 'vertical' && stackIndex > 0
          ? `calc(-1 * var(${isSmall ? '--card-stack-offset-v-small' : '--card-stack-offset-v'}))`
          : '0',
        willChange: 'transform',
        border: isSelected
          ? `2px solid ${SELECTED_COLORS.border}`
          : '1px solid hsla(40, 20%, 70%, 0.5)',
        boxShadow: isSelected
          ? `${SELECTED_COLORS.glow}, 0_2px_4px_rgba(0,0,0,0.1)`
          : undefined,
      }}
      onClick={onClick}
    >
      <div
        className="absolute inset-x-0 top-0 h-1/3 pointer-events-none rounded-t-[inherit]"
        style={{
          background: 'linear-gradient(to bottom, rgba(255,255,255,0.3) 0%, transparent 100%)',
        }}
      />

      {isSelected && (
        <div
          className="absolute inset-0 pointer-events-none rounded-[inherit]"
          style={{
            background: SELECTED_COLORS.overlay,
          }}
        />
      )}

      {isJoker(card) ? (
        <JokerContent card={card} size={size} />
      ) : (
        <>
          <CardCorner
            rank={card.rank}
            suit={card.suit}
            size={size}
          />

          <CardSuitLarge
            suit={card.suit}
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

export { getRankText };

export default React.memo(CardDisplay);
