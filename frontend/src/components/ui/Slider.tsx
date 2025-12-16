import * as React from "react"
import * as SliderPrimitive from "@radix-ui/react-slider"
import { cn } from "@/lib/utils"

export interface SliderProps
  extends React.ComponentPropsWithoutRef<typeof SliderPrimitive.Root> {
  showTooltip?: boolean
  variant?: 'default' | 'gold' | 'jade'
}

const Slider = React.forwardRef<
  React.ComponentRef<typeof SliderPrimitive.Root>,
  SliderProps
>(({ className, showTooltip = false, variant = 'default', ...props }, ref) => {
  const [showValue, setShowValue] = React.useState(false)

  const variantStyles = {
    default: {
      track: "bg-state-disabled",
      range: "bg-gradient-to-r from-state-active to-[hsl(42,95%,45%)]",
      thumb: "border-state-active hover:shadow-[var(--elevation-3),var(--shadow-relief),var(--shadow-glow-md)]",
      glow: "var(--shadow-glow-sm)",
    },
    gold: {
      track: "bg-[hsl(42,20%,85%)]",
      range: "bg-gradient-to-r from-[hsl(42,95%,58%)] via-[hsl(42,95%,52%)] to-[hsl(42,95%,45%)]",
      thumb: "border-[hsl(42,95%,45%)] hover:shadow-[var(--elevation-3),var(--shadow-relief),var(--shadow-glow-lg)]",
      glow: "var(--shadow-glow-md)",
    },
    jade: {
      track: "bg-[hsl(158,20%,85%)]",
      range: "bg-gradient-to-r from-[hsl(158,55%,55%)] via-[hsl(158,55%,42%)] to-[hsl(158,55%,35%)]",
      thumb: "border-[hsl(158,55%,42%)] hover:shadow-[var(--elevation-3),var(--shadow-relief),var(--shadow-jade)]",
      glow: "var(--shadow-jade)",
    },
  }

  const styles = variantStyles[variant]

  return (
    <SliderPrimitive.Root
      ref={ref}
      className={cn(
        "relative flex w-full touch-none select-none items-center",
        "py-2",
        className
      )}
      onPointerDown={() => setShowValue(true)}
      onPointerUp={() => setShowValue(false)}
      {...props}
    >
      {/* 轨道 */}
      <SliderPrimitive.Track
        className={cn(
          "relative h-2.5 w-full grow overflow-hidden rounded-full",
          styles.track,
          // 内凹阴影
          "shadow-[inset_0_1px_3px_rgba(0,0,0,0.15),inset_0_-1px_1px_rgba(255,255,255,0.5)]",
          "transition-all duration-normal ease-smooth"
        )}
      >
        {/* 已选范围 */}
        <SliderPrimitive.Range
          className={cn(
            "absolute h-full",
            styles.range,
            // 光泽效果
            "shadow-[inset_0_1px_0_rgba(255,255,255,0.4)]",
            "transition-all duration-normal ease-smooth"
          )}
        />
      </SliderPrimitive.Track>

      {/* 滑块 */}
      <SliderPrimitive.Thumb
        className={cn(
          "block h-6 w-6 rounded-full",
          // 背景
          "bg-gradient-to-b from-surface-base to-surface-elevated",
          // 边框
          "border-[3px]",
          styles.thumb,
          // 阴影
          "shadow-[var(--elevation-2),var(--shadow-relief)]",
          // 过渡动画
          "transition-all duration-fast ease-bounce",
          // 悬停效果
          "hover:scale-110",
          // 激活效果
          "active:scale-125",
          `active:shadow-[var(--elevation-3),var(--shadow-relief),${styles.glow}]`,
          // 聚焦状态
          "focus-visible:outline-none",
          "focus-visible:ring-2 focus-visible:ring-state-active focus-visible:ring-offset-2",
          `focus-visible:shadow-[var(--elevation-3),var(--shadow-relief),${styles.glow}]`,
          // 禁用状态
          "disabled:pointer-events-none disabled:opacity-50 disabled:scale-100",
          // 光标
          "cursor-grab active:cursor-grabbing"
        )}
      >
        {/* 滑块内部装饰 */}
        <span className="absolute inset-1 rounded-full bg-gradient-to-b from-white/50 to-transparent" />
      </SliderPrimitive.Thumb>

      {/* 数值提示 (可选) */}
      {showTooltip && showValue && props.value && (
        <div
          className={cn(
            "absolute -top-8 left-1/2 -translate-x-1/2",
            "px-2 py-1 rounded-md",
            "bg-primitive-neutral-900/90 text-fg-inverse",
            "text-xs font-mono",
            "shadow-elevation-2",
            "animate-in fade-in-0 zoom-in-95 duration-fast"
          )}
        >
          {Array.isArray(props.value) ? props.value[0] : props.value}
        </div>
      )}
    </SliderPrimitive.Root>
  )
})
Slider.displayName = SliderPrimitive.Root.displayName

// 带刻度的滑块
const SliderWithMarks = React.forwardRef<
  React.ComponentRef<typeof SliderPrimitive.Root>,
  SliderProps & {
    marks?: { value: number; label?: string }[]
  }
>(({ className, marks = [], variant = 'default', min = 0, max = 100, ...props }, ref) => {
  return (
    <div className="relative pt-2 pb-6">
      <Slider
        ref={ref}
        className={className}
        variant={variant}
        min={min}
        max={max}
        {...props}
      />
      {/* 刻度标记 */}
      <div className="absolute inset-x-0 top-[22px] flex justify-between px-3">
        {marks.map((mark) => {
          const position = ((mark.value - min) / (max - min)) * 100
          return (
            <div
              key={mark.value}
              className="absolute flex flex-col items-center"
              style={{ left: `${position}%`, transform: 'translateX(-50%)' }}
            >
              <span className="w-px h-2 bg-stroke" />
              {mark.label && (
                <span className="mt-1 text-xs text-fg-secondary whitespace-nowrap">
                  {mark.label}
                </span>
              )}
            </div>
          )
        })}
      </div>
    </div>
  )
})
SliderWithMarks.displayName = "SliderWithMarks"

export { Slider, SliderWithMarks }
