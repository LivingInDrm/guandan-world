import React, { useMemo, useState } from 'react';
import type { Card } from '../../types/proto';
import { isJoker } from '../../utils/cardUtils';
import { 
  SUIT_SYMBOLS, 
  JOKER_CONFIG, 
  ANIMATIONS,
  type CardSize,
  type SizeConfig,
  getCardSizeStyle,
  getSuitColorClass,
  getSuitShadowClass,
  getRankText
} from './cardStyles';

interface CardDisplayProps {
  card: Card;
  isSelected?: boolean;
  onClick?: () => void;
  stackIndex?: number;
  size?: CardSize;
  className?: string;
}

// 角标组件
const CardCorner: React.FC<{
  rank: number;
  suit: number;
  position: 'top-left' | 'bottom-right';
  sizeConfig: SizeConfig;
  size: CardSize;
}> = ({ rank, suit, position, sizeConfig, size }) => {
  const isBottom = position === 'bottom-right';
  const colorClass = getSuitColorClass(suit);
  const showSuit = size !== 'xs';
  
  return (
    <div 
      className={`absolute flex flex-col items-center leading-none ${colorClass} ${isBottom ? 'rotate-180 bottom-0.5 right-0.5' : 'top-0.5 left-0.5'}`}
      style={{ 
        padding: sizeConfig.padding,
        width: '1.2em',
      }}
    >
      <span style={{ fontSize: sizeConfig.fontSize, fontWeight: 700 }}>
        {getRankText(rank)}
      </span>
      {showSuit && (
        <span style={{ fontSize: sizeConfig.iconSize, marginTop: '-2px' }}>
          {SUIT_SYMBOLS[suit]}
        </span>
      )}
    </div>
  );
};

// 中央花色组件
const CardCenter: React.FC<{
  rank: number;
  suit: number;
  sizeConfig: SizeConfig;
}> = ({ suit, sizeConfig }) => {
  const colorClass = getSuitColorClass(suit);
  const shadowClass = getSuitShadowClass(suit);
  
  return (
    <div 
      className={`absolute inset-0 flex items-center justify-center ${colorClass} ${shadowClass}`}
      style={{ fontSize: sizeConfig.centerIconSize }}
    >
      {SUIT_SYMBOLS[suit]}
    </div>
  );
};

// 大小王专属内容
const JokerContent: React.FC<{
  card: Card;
  sizeConfig: SizeConfig;
}> = ({ card, sizeConfig }) => {
  const config = card.rank === 16 ? JOKER_CONFIG.big : JOKER_CONFIG.small;
  
  return (
    <div className="absolute inset-0 flex flex-col items-center justify-center p-1">
      <div 
        className="absolute top-1 left-1"
        style={{ fontSize: sizeConfig.iconSize }}
      >
        {config.icon}
      </div>
      
      <div 
        className={`writing-vertical-rl font-bold tracking-widest ${config.color}`}
        style={{ 
          fontSize: `calc(${sizeConfig.fontSize} * 1.2)`,
          textShadow: '0 1px 1px rgba(0,0,0,0.1)'
        }}
      >
        {config.text}
      </div>
      
      <div 
        className="absolute bottom-1 right-1 rotate-180"
        style={{ fontSize: sizeConfig.iconSize }}
      >
        {config.icon}
      </div>
    </div>
  );
};

const CardDisplay: React.FC<CardDisplayProps> = ({ 
  card, 
  isSelected = false, 
  onClick, 
  stackIndex = 0,
  size = 'normal',
  className = ''
}) => {
  const [isHovered, setIsHovered] = useState(false);
  const sizeConfig = useMemo(() => getCardSizeStyle(size), [size]);
  
  const bgClass = useMemo(() => {
    if (isJoker(card)) {
      return card.rank === 16 ? JOKER_CONFIG.big.bgGradient : JOKER_CONFIG.small.bgGradient;
    }
    return 'bg-white';
  }, [card]);
  
  const zIndex = useMemo(() => {
    if (isSelected) return 100;
    if (isHovered) return 50;
    return stackIndex;
  }, [isSelected, isHovered, stackIndex]);

  return (
    <div
      role="button"
      aria-pressed={isSelected}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
      className={`
        relative select-none
        card-3d-shadow border border-gray-200
        ${bgClass}
        ${ANIMATIONS.transition}
        cursor-pointer
        ${isSelected ? ANIMATIONS.selected : onClick ? ANIMATIONS.hover : ''}
        ${className}
      `}
      style={{
        width: sizeConfig.width,
        height: sizeConfig.height,
        borderRadius: sizeConfig.borderRadius,
        zIndex,
        marginLeft: stackIndex > 0 ? (size === 'small' || size === 'xs' ? '-12px' : '-24px') : '0',
        willChange: 'transform',
      }}
      onClick={onClick}
    >
      {/* 卡片内边框装饰 (可选) */}
      {!isJoker(card) && (
        <div className="absolute inset-1 border border-gray-100 rounded opacity-50 pointer-events-none" />
      )}

      {isJoker(card) ? (
        <JokerContent card={card} sizeConfig={sizeConfig} />
      ) : (
        <>
          <CardCorner 
            rank={card.rank} 
            suit={card.suit} 
            position="top-left" 
            sizeConfig={sizeConfig}
            size={size}
          />
          
          <CardCenter 
            rank={card.rank} 
            suit={card.suit} 
            sizeConfig={sizeConfig} 
          />
          
          <CardCorner 
            rank={card.rank} 
            suit={card.suit} 
            position="bottom-right" 
            sizeConfig={sizeConfig}
            size={size}
          />
        </>
      )}
    </div>
  );
};

// 导出辅助函数以保持向后兼容
export { getRankText };

export default React.memo(CardDisplay);
