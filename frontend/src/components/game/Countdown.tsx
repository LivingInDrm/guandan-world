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

  const config = sizeConfig[size];
  const circumference = 2 * Math.PI * config.radius;
  const totalSeconds = totalSecondsRef.current || 1;
  const progress = Math.max(0, Math.min(secondsLeft / totalSeconds, 1));
  const offset = circumference * (1 - progress);

  const getRingColor = (): string => {
    if (secondsLeft <= 5) return '#ef4444';
    if (secondsLeft <= 10) return '#f97316';
    return '#fbbf24';
  };

  const getTextColor = (): string => {
    if (secondsLeft <= 5) return '#ef4444';
    if (secondsLeft <= 10) return '#f97316';
    return '#4b5563';
  };

  const viewBoxSize = config.radius * 2 + config.strokeWidth * 2;
  const center = viewBoxSize / 2;

return (
  <div
    className="relative inline-flex items-center justify-center"
    style={{ width: config.width, height: config.width }}
  >
    {/* 背景玻璃圆盘 */}
    <div
      className="
        absolute inset-0 rounded-full
        bg-gradient-to-b from-white/90 to-slate-100/85
        backdrop-blur-sm
        border border-white/70
        shadow-[0_4px_12px_rgba(0,0,0,0.2)]
      "
    />

    {/* 外环进度（SVG） */}
    <svg
      width={config.width}
      height={config.width}
      viewBox={`0 0 ${viewBoxSize} ${viewBoxSize}`}
      className="absolute"
    >
      {/* 底环 */}
      <circle
        cx={center}
        cy={center}
        r={config.radius}
        stroke="rgba(148,163,184,0.45)"   // slate-400 的透明版
        strokeWidth={config.strokeWidth}
        fill="none"
      />

      {/* 外环阴影 */}
      <circle
        cx={center}
        cy={center}
        r={config.radius}
        stroke="rgba(0,0,0,0.15)"
        strokeWidth={config.strokeWidth + 2}
        fill="none"
        className="blur-[3px]"
      />

      {/* 倒计时环（主色渐变 + 抖动发光） */}
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
        style={{
          filter: secondsLeft <= 5 ? 'drop-shadow(0 0 6px rgba(239,68,68,0.7))'
                                   : 'drop-shadow(0 0 4px rgba(251,191,36,0.6))',
          transition: 'stroke 0.2s ease-out',
        }}
      />
    </svg>

    {/* 中心数字 */}
    <div
      className="relative font-bold leading-none"
      style={{
        fontSize: config.fontSize,
        color: getTextColor(),
        textShadow: '0 1px 2px rgba(0,0,0,0.4)',
      }}
    >
      {secondsLeft}
    </div>
  </div>
  );
};

export default Countdown;
