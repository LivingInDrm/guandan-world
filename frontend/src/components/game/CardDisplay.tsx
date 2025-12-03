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
  stackIndex?: number;
  size?: CardSize;
  className?: string;
}

// 左上角标组件
const CardCorner: React.FC<{
  rank: number;
  suit: number;
  sizeConfig: SizeConfig;
}> = ({ rank, suit, sizeConfig }) => {
  const colorClass = getSuitColorClass(suit);
  
  return (
    <div 
      className={`absolute flex flex-col items-center leading-none top-0.5 ${colorClass}`}
      style={{ 
        padding: sizeConfig.padding,
        width: '1.2em',
        left: sizeConfig.cornerTextLeft,
      }}
    >
      <span style={{ fontSize: sizeConfig.fontSize, fontWeight: 700 }}>
        {getRankText(rank)}
      </span>
      <span style={{ fontSize: sizeConfig.iconSize, marginTop: '-2px' }}>
        {SUIT_SYMBOLS[suit]}
      </span>
    </div>
  );
};

// 右下大花色组件
const CardSuitLarge: React.FC<{
  suit: number;
  sizeConfig: SizeConfig;
}> = ({ suit, sizeConfig }) => {
  return (
    <img 
      src={SUIT_IMAGES[suit]}
      alt={SUIT_SYMBOLS[suit]}
      className="absolute bottom-2 right-2"
      style={{ 
        width: `calc(${sizeConfig.centerIconSize} * 1.5)`,
        height: `calc(${sizeConfig.centerIconSize} * 1.5)`,
        objectFit: 'contain'
      }}
    />
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
        className={`absolute top-1 flex flex-col items-center leading-tight font-bold ${config.color}`}
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
        marginLeft: stackIndex > 0 ? (size === 'small' ? '-12px' : '-24px') : '0',
        willChange: 'transform',
      }}
      onClick={onClick}
    >
      {/* 卡片内边框装饰 (可选) */}
      {!isJoker(card) && (
        <div className="absolute inset-px border border-gray-100 rounded opacity-50 pointer-events-none" />
      )}

      {isJoker(card) ? (
        <JokerContent card={card} sizeConfig={sizeConfig} />
      ) : (
        <>
          <CardCorner 
            rank={card.rank} 
            suit={card.suit} 
            sizeConfig={sizeConfig}
          />
          
          <CardSuitLarge 
            suit={card.suit} 
            sizeConfig={sizeConfig} 
          />
        </>
      )}
    </div>
  );
};

// 导出辅助函数以保持向后兼容
export { getRankText };

export default React.memo(CardDisplay);
