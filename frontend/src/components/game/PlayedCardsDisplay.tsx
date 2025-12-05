import React from 'react';
import type { PlayAction } from '../../types/proto';
import CardDisplay from './CardDisplay';

export interface PlayedCardsDisplayProps {
  play: PlayAction | null;
  position: 'bottom' | 'left' | 'top' | 'right';
}

const getPositionClasses = (position: PlayedCardsDisplayProps['position']): string => {
  const base = 'absolute flex items-center justify-center';
  switch (position) {
    case 'bottom':
      return `${base} bottom-4 left-1/2 transform -translate-x-1/2`;
    case 'top':
      return `${base} top-4 left-1/2 transform -translate-x-1/2`;
    case 'left':
      return `${base} left-4 top-1/2 transform -translate-y-1/2`;
    case 'right':
      return `${base} right-4 top-1/2 transform -translate-y-1/2`;
    default:
      return base;
  }
};

const PlayedCardsDisplay: React.FC<PlayedCardsDisplayProps> = ({ play, position }) => {
  if (!play) {
    return <div className={getPositionClasses(position)} />;
  }

  if (play.isPass) {
    return (
      <div className={getPositionClasses(position)}>
        <div className="bg-gray-200 px-3 py-1 rounded text-sm text-gray-600 font-medium">
          不出
        </div>
      </div>
    );
  }

  if (play.cards && play.cards.length > 0) {
    return (
      <div className={getPositionClasses(position)}>
        <div className="flex flex-row items-center gap-0.5">
          {play.cards.map((card, index) => (
            <CardDisplay
              key={card.deckIndex}
              card={card}
              size="small"
              stackIndex={index}
            />
          ))}
        </div>
      </div>
    );
  }

  return <div className={getPositionClasses(position)} />;
};

export default React.memo(PlayedCardsDisplay);
