import React, { useState, useEffect, useRef } from 'react';

interface CountdownProps {
  deadlineAtMs: number;
  isActive?: boolean;
  size?: 'small' | 'medium' | 'large';
  onTimeout?: () => void;
}

const sizeConfig = {
  small: { width: 48, fontSize: 18, radius: 18, strokeWidth: 4 },
  medium: { width: 72, fontSize: 26, radius: 28, strokeWidth: 6 },
  large: { width: 96, fontSize: 36, radius: 38, strokeWidth: 8 },
};

const calculateRemaining = (deadlineAtMs: number): number => {
  const remaining = deadlineAtMs - Date.now();
  return Math.max(0, Math.ceil(remaining / 1000));
};

const DEFAULT_SECONDS = 30;

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
  const totalSecondsRef = useRef<number>(0);

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
    totalSecondsRef.current = remaining;

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

  const displaySeconds = isActive ? secondsLeft : DEFAULT_SECONDS;
  const displayProgress = isActive
    ? Math.max(0, Math.min(secondsLeft / (totalSecondsRef.current || 1), 1))
    : 1;

  const config = sizeConfig[size];
  const circumference = 2 * Math.PI * config.radius;
  const offset = circumference * (1 - displayProgress);

  const getRingColor = (): string => {
    if (!isActive) return 'hsl(0, 0%, 50%)';
    if (secondsLeft <= 5) return 'hsl(0, 70%, 50%)'; // 朱砂红
    if (secondsLeft <= 10) return 'hsl(42, 95%, 52%)'; // 帝王金
    return 'hsl(158, 55%, 42%)'; // 翡翠绿
  };

  const getTextColor = (): string => {
    if (!isActive) return 'hsla(0, 0%, 100%, 0.5)';
    if (secondsLeft <= 5) return 'hsl(0, 70%, 60%)';
    if (secondsLeft <= 10) return 'hsl(42, 95%, 55%)';
    return 'hsla(0, 0%, 100%, 0.9)';
  };

  const getRingFilter = (): string => {
    if (!isActive) return 'none';
    if (secondsLeft <= 5) return 'drop-shadow(0 0 8px hsla(0, 70%, 50%, 0.7))';
    if (secondsLeft <= 10) return 'drop-shadow(0 0 6px hsla(42, 95%, 52%, 0.6))';
    return 'drop-shadow(0 0 6px hsla(158, 55%, 42%, 0.5))';
  };

  const viewBoxSize = config.radius * 2 + config.strokeWidth * 2;
  const center = viewBoxSize / 2;

  return (
    <div
      className="relative inline-flex items-center justify-center mobile-landscape:scale-75"
      style={{ width: config.width, height: config.width }}
    >
      <div
        className="
          absolute inset-0 rounded-full
          bg-gradient-to-b from-black/50 to-black/70
          backdrop-blur-sm
          border border-white/20
          shadow-[0_2px_8px_rgba(0,0,0,0.3),inset_0_1px_0_hsla(0,0%,100%,0.1)]
        "
      />

      <svg
        width={config.width}
        height={config.width}
        viewBox={`0 0 ${viewBoxSize} ${viewBoxSize}`}
        className="absolute"
      >
        {/* 底色轨道 */}
        <circle
          cx={center}
          cy={center}
          r={config.radius}
          stroke="hsla(0, 0%, 100%, 0.15)"
          strokeWidth={config.strokeWidth}
          fill="none"
        />

        {/* 内阴影效果 */}
        <circle
          cx={center}
          cy={center}
          r={config.radius}
          stroke="hsla(0, 0%, 0%, 0.3)"
          strokeWidth={config.strokeWidth + 2}
          fill="none"
          className="blur-[3px]"
        />

        <circle
          cx={center}
          cy={center}
          r={config.radius}
          stroke={getRingColor()}
          strokeWidth={config.strokeWidth}
          fill="none"
          strokeDasharray={circumference}
          strokeDashoffset={offset}
          strokeLinecap="round"
          transform={`rotate(-90 ${center} ${center})`}
          className="transition-all duration-fast ease-smooth"
          style={{ filter: getRingFilter() }}
        />
      </svg>

      <div
        className="relative font-bold leading-none font-display"
        style={{
          fontSize: config.fontSize,
          color: getTextColor(),
          textShadow: '0 2px 4px rgba(0, 0, 0, 0.5)',
        }}
      >
        {displaySeconds}
      </div>
    </div>
  );
};

export default Countdown;
