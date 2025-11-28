import React, { useEffect, useState } from 'react';
import type { Card } from '../../../types';
import CardDisplay from '../CardDisplay';

interface CardFlyAnimationProps {
  card: Card;
  fromPosition: { x: number; y: number };
  toPosition: { x: number; y: number };
  duration?: number;
  onComplete: () => void;
}

const CardFlyAnimation: React.FC<CardFlyAnimationProps> = ({
  card,
  fromPosition,
  toPosition,
  duration = 500,
  onComplete,
}) => {
  const [position, setPosition] = useState(fromPosition);
  const [isAnimating, setIsAnimating] = useState(false);

  useEffect(() => {
    requestAnimationFrame(() => {
      setIsAnimating(true);
      setPosition(toPosition);
    });

    const timer = setTimeout(() => {
      onComplete();
    }, duration);

    return () => clearTimeout(timer);
  }, [toPosition, duration, onComplete]);

  return (
    <div
      className="fixed z-50 pointer-events-none"
      style={{
        left: position.x,
        top: position.y,
        transform: 'translate(-50%, -50%)',
        transition: isAnimating ? `all ${duration}ms ease-out` : 'none',
      }}
    >
      <CardDisplay card={card} size="small" />
    </div>
  );
};

export default CardFlyAnimation;
