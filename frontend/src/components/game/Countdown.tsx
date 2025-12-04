import React, { useState, useEffect, useRef } from 'react';
import clockImage from '../../assets/clock.png';

interface CountdownProps {
  deadlineAtMs: number;
  isActive?: boolean;
  size?: 'small' | 'medium' | 'large';
  onTimeout?: () => void;
}

const sizeConfig = {
  small: { width: 48, fontSize: 14, top: '65%' },
  medium: { width: 72, fontSize: 22, top: '65%' },
  large: { width: 96, fontSize: 30, top: '65%' },
};

const calculateRemaining = (deadlineAtMs: number): number => {
  const remaining = deadlineAtMs - Date.now();
  return Math.max(0, Math.ceil(remaining / 1000));
};

const Countdown: React.FC<CountdownProps> = ({
  deadlineAtMs,
  isActive = true,
  size = 'medium',
  onTimeout,
}) => {
  const [secondsLeft, setSecondsLeft] = useState<number>(() =>
    calculateRemaining(deadlineAtMs)
  );
  const onTimeoutRef = useRef(onTimeout);
  const hasTriggeredTimeout = useRef(false);

  useEffect(() => {
    onTimeoutRef.current = onTimeout;
  }, [onTimeout]);

  useEffect(() => {
    hasTriggeredTimeout.current = false;
  }, [deadlineAtMs]);

  useEffect(() => {
    if (!isActive || !deadlineAtMs) return;

    const remaining = calculateRemaining(deadlineAtMs);
    setSecondsLeft(remaining);

    if (remaining <= 0 && !hasTriggeredTimeout.current) {
      hasTriggeredTimeout.current = true;
      onTimeoutRef.current?.();
      return;
    }

    const timer = setInterval(() => {
      const newRemaining = calculateRemaining(deadlineAtMs);
      setSecondsLeft(newRemaining);

      if (newRemaining <= 0) {
        clearInterval(timer);
        if (!hasTriggeredTimeout.current) {
          hasTriggeredTimeout.current = true;
          onTimeoutRef.current?.();
        }
      }
    }, 1000);

    return () => clearInterval(timer);
  }, [deadlineAtMs, isActive]);

  if (!deadlineAtMs) {
    return null;
  }

  const config = sizeConfig[size];

  const getTextColor = (): string => {
    if (secondsLeft <= 5) return '#ef4444';
    if (secondsLeft <= 10) return '#f97316';
    return '#4b5563';
  };

  return (
    <div
      className="relative inline-block"
      style={{ width: config.width, height: config.width }}
    >
      <img
        src={clockImage}
        alt="countdown"
        className="w-full h-full object-contain"
      />
      <div
        className="absolute left-1/2 font-bold"
        style={{
          top: config.top,
          transform: 'translate(-50%, -50%)',
          fontSize: config.fontSize,
          color: getTextColor(),
          textShadow: '0 1px 3px rgba(0, 0, 0, 0.6)',
        }}
      >
        {secondsLeft}
      </div>
    </div>
  );
};

export default Countdown;
