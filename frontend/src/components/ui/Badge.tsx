import * as React from "react"
import { cva, type VariantProps } from "class-variance-authority"
import { cn } from "@/lib/utils"

const badgeVariants = cva(
  [
    "inline-flex items-center justify-center gap-1",
    "font-display font-medium tracking-wide",
    "rounded-md",
    // 阴影
    "shadow-[var(--elevation-1),var(--shadow-relief)]",
    // 过渡动画
    "transition-all duration-fast ease-bounce",
    // 边框
    "border",
  ],
  {
    variants: {
      variant: {
        // 庄家标签 - 金色高贵
        landlord: [
          "bg-gradient-to-b from-state-active to-[hsl(42,95%,42%)]",
          "text-fg-inverse text-shadow-sm",
          "border-[hsl(42,95%,35%)]",
          "shadow-[var(--elevation-1),var(--shadow-relief),var(--shadow-glow-sm)]",
        ],
        // 农民标签 - 翡翠清新
        farmer: [
          "bg-gradient-to-b from-team-us to-[hsl(158,55%,35%)]",
          "text-fg-inverse text-shadow-sm",
          "border-[hsl(158,55%,28%)]",
        ],
        // 队友标签
        teammate: [
          "bg-team-us/20 backdrop-blur-sm",
          "text-team-us",
          "border-team-us/50",
        ],
        // 房主标签
        owner: [
          "bg-gradient-to-b from-state-active to-[hsl(42,95%,42%)]",
          "text-fg-inverse text-shadow-sm",
          "border-[hsl(42,95%,35%)]",
        ],
        // 中性标签
        neutral: [
          "bg-gradient-to-b from-surface-elevated to-[hsl(35,6%,85%)]",
          "text-fg-secondary",
          "border-stroke/50",
        ],
        // 成功标签
        success: [
          "bg-gradient-to-b from-action-primary to-[hsl(152,60%,35%)]",
          "text-fg-inverse text-shadow-sm",
          "border-[hsl(152,60%,28%)]",
        ],
        // 警告标签
        warning: [
          "bg-gradient-to-b from-action-secondary to-[hsl(28,85%,44%)]",
          "text-fg-inverse text-shadow-sm",
          "border-[hsl(28,85%,38%)]",
        ],
        // 危险标签
        danger: [
          "bg-gradient-to-b from-action-danger to-[hsl(5,75%,42%)]",
          "text-fg-inverse text-shadow-sm",
          "border-[hsl(5,75%,35%)]",
        ],
        // 轮廓样式
        outline: [
          "bg-transparent",
          "text-fg-primary",
          "border-stroke",
          "shadow-none",
        ],
        // 级别标签 - 特殊样式
        level: [
          "bg-gradient-to-b from-badge-level to-[hsl(142,71%,38%)]",
          "text-fg-inverse text-shadow-sm",
          "border-[hsl(142,71%,32%)]",
          "shadow-[var(--elevation-1),var(--shadow-relief),var(--shadow-jade)]",
        ],
        // 癞子标签 - 特殊样式
        wild: [
          "bg-gradient-to-b from-badge-wild to-[hsl(0,84%,50%)]",
          "text-fg-inverse text-shadow-sm",
          "border-[hsl(0,84%,42%)]",
          "shadow-[var(--elevation-1),var(--shadow-relief),var(--shadow-ruby)]",
          "animate-pulse-glow",
        ],
      },
      size: {
        xs: "px-1.5 py-0.5 text-[10px] min-h-[18px] mobile-landscape:px-1 mobile-landscape:py-0 mobile-landscape:text-[8px] mobile-landscape:min-h-[14px]",
        sm: "px-2 py-0.5 text-xs min-h-[22px] mobile-landscape:px-1.5 mobile-landscape:py-0.5 mobile-landscape:text-[10px] mobile-landscape:min-h-[18px]",
        md: "px-2.5 py-1 text-sm min-h-[26px] mobile-landscape:px-2 mobile-landscape:py-0.5 mobile-landscape:text-xs mobile-landscape:min-h-[20px]",
        lg: "px-3 py-1.5 text-base min-h-[32px] mobile-landscape:px-2.5 mobile-landscape:py-1 mobile-landscape:text-sm mobile-landscape:min-h-[26px]",
      }
    },
    defaultVariants: {
      variant: "neutral",
      size: "sm",
    }
  }
)

export interface BadgeProps
  extends React.HTMLAttributes<HTMLDivElement>,
    VariantProps<typeof badgeVariants> {
  icon?: React.ReactNode
  removable?: boolean
  onRemove?: () => void
}

const Badge = React.forwardRef<HTMLDivElement, BadgeProps>(
  ({ className, variant, size, icon, removable, onRemove, children, ...props }, ref) => (
    <div
      ref={ref}
      className={cn(badgeVariants({ variant, size }), className)}
      {...props}
    >
      {icon && (
        <span className="flex-shrink-0">{icon}</span>
      )}
      <span>{children}</span>
      {removable && (
        <button
          type="button"
          onClick={onRemove}
          className={cn(
            "flex-shrink-0 ml-0.5 -mr-0.5",
            "w-4 h-4 rounded-full",
            "flex items-center justify-center",
            "hover:bg-black/10 active:bg-black/20",
            "transition-colors duration-fast",
          )}
        >
          <svg
            className="w-3 h-3"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M6 18L18 6M6 6l12 12"
            />
          </svg>
        </button>
      )}
    </div>
  )
)
Badge.displayName = "Badge"

// 带点的徽章 - 用于显示数量或状态
const BadgeDot = React.forwardRef<
  HTMLSpanElement,
  React.HTMLAttributes<HTMLSpanElement> & {
    color?: 'primary' | 'success' | 'warning' | 'danger' | 'neutral'
    pulse?: boolean
  }
>(({ className, color = 'primary', pulse = false, ...props }, ref) => {
  const colorClasses = {
    primary: "bg-action-primary",
    success: "bg-action-primary",
    warning: "bg-action-secondary",
    danger: "bg-action-danger",
    neutral: "bg-fg-secondary",
  }

  return (
    <span
      ref={ref}
      className={cn(
        "inline-block w-2 h-2 rounded-full",
        colorClasses[color],
        pulse && "animate-pulse",
        className
      )}
      {...props}
    />
  )
})
BadgeDot.displayName = "BadgeDot"

export { Badge, BadgeDot, badgeVariants }
