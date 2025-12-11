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
    if (!isActive) return 'hsl(var(--ds-primitive-neutral-500))';
    if (secondsLeft <= 5) return 'hsl(var(--ds-primitive-danger-500))';
    if (secondsLeft <= 10) return 'hsl(var(--ds-primitive-warning-500))';
    return 'hsl(var(--ds-primitive-accent-500))';
  };

  const getTextColor = (): string => {
    if (!isActive) return 'hsl(var(--ds-primitive-neutral-500))';
    if (secondsLeft <= 5) return 'hsl(var(--ds-primitive-danger-500))';
    if (secondsLeft <= 10) return 'hsl(var(--ds-primitive-warning-500))';
    return 'hsl(var(--ds-primitive-neutral-700))';
  };

  const getRingFilter = (): string => {
    if (!isActive) return 'none';
    if (secondsLeft <= 5) return 'drop-shadow(0 0 6px hsl(var(--ds-primitive-danger-500) / 0.7))';
    return 'drop-shadow(0 0 4px hsl(var(--ds-primitive-accent-500) / 0.6))';
  };

  const viewBoxSize = config.radius * 2 + config.strokeWidth * 2;
  const center = viewBoxSize / 2;

  return (
    <div
      className="relative inline-flex items-center justify-center"
      style={{ width: config.width, height: config.width }}
    >
      <div
        className="
          absolute inset-0 rounded-ds-full
          bg-gradient-to-b from-white/90 to-ds-surface-elevated/85
          backdrop-blur-sm
          border border-white/70
          shadow-ds-elevation-2
        "
      />

      <svg
        width={config.width}
        height={config.width}
        viewBox={`0 0 ${viewBoxSize} ${viewBoxSize}`}
        className="absolute"
      >
        <circle
          cx={center}
          cy={center}
          r={config.radius}
          stroke="hsl(var(--ds-primitive-neutral-500) / 0.45)"
          strokeWidth={config.strokeWidth}
          fill="none"
        />

        <circle
          cx={center}
          cy={center}
          r={config.radius}
          stroke="hsl(var(--ds-primitive-neutral-900) / 0.15)"
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
          className="transition-all duration-ds-fast ease-ds-smooth"
          style={{ filter: getRingFilter() }}
        />
      </svg>

      <div
        className="relative font-bold leading-none"
        style={{
          fontSize: config.fontSize,
          color: getTextColor(),
          textShadow: '0 1px 2px hsl(var(--ds-primitive-neutral-900) / 0.4)',
        }}
      >
        {displaySeconds}
      </div>
    </div>
  );
};

export default Countdown;
