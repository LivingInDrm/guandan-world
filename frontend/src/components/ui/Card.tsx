import * as React from "react"
import { cva, type VariantProps } from "class-variance-authority"
import { cn } from "@/lib/utils"

const cardVariants = cva(
  [
    // 基础样式
    "relative rounded-lg overflow-hidden",
    // 过渡动画
    "transition-all duration-normal ease-bounce",
  ],
  {
    variants: {
      variant: {
        base: [
          "bg-surface-base",
          "shadow-paper",
          "border border-stroke/50",
        ],
        elevated: [
          "bg-gradient-to-b from-surface-elevated to-[hsl(40,8%,93%)]",
          "shadow-[var(--elevation-2),var(--shadow-relief)]",
          "border border-stroke/30",
          // 顶部高光边
          "before:absolute before:inset-x-0 before:top-0 before:h-px",
          "before:bg-gradient-to-r before:from-transparent before:via-white/60 before:to-transparent",
        ],
        emphasis: [
          "bg-surface-emphasis backdrop-blur-xl",
          "shadow-[var(--elevation-3),var(--shadow-relief),var(--shadow-glow-sm)]",
          "border-2 border-state-active/30",
          // 渐变边框效果
          "before:absolute before:inset-0 before:rounded-lg before:p-[2px]",
          "before:bg-gradient-to-br before:from-state-active/40 before:via-transparent before:to-action-primary/30",
          "before:-z-10",
        ],
        glass: [
          "bg-white/60 backdrop-blur-lg",
          "shadow-[var(--elevation-2),var(--shadow-inner-light)]",
          "border border-white/50",
        ],
        outlined: [
          "bg-transparent",
          "shadow-none",
          "border-2 border-stroke",
        ],
        premium: [
          // 高端金色边框卡片
          "bg-gradient-to-br from-[hsl(40,8%,97%)] to-[hsl(40,8%,92%)]",
          "shadow-[var(--elevation-3),var(--shadow-relief)]",
          "border-2 border-transparent",
          // 金色渐变边框
          "[background:linear-gradient(to_bottom_right,hsl(40_8%_97%),hsl(40_8%_92%))_padding-box,linear-gradient(135deg,hsl(42_95%_68%),hsl(42_95%_52%),hsl(42_95%_38%))_border-box]",
          "hover:shadow-[var(--elevation-3),var(--shadow-relief),var(--shadow-glow-md)]",
        ],
      },
    },
    defaultVariants: {
      variant: "base",
    }
  }
)

export interface CardProps
  extends React.HTMLAttributes<HTMLDivElement>,
    VariantProps<typeof cardVariants> {
  interactive?: boolean
  hoverable?: boolean
}

const Card = React.forwardRef<HTMLDivElement, CardProps>(
  ({ className, variant, interactive = false, hoverable = true, children, ...props }, ref) => (
    <div
      ref={ref}
      className={cn(
        cardVariants({ variant }),
        hoverable && [
          "hover:shadow-[var(--elevation-3),var(--shadow-relief)]",
          "hover:-translate-y-0.5",
        ],
        interactive && [
          "cursor-pointer",
          "active:translate-y-0 active:scale-[0.99]",
          "active:shadow-[var(--elevation-1)]",
        ],
        !hoverable && "hover:translate-y-0",
        className
      )}
      {...props}
    >
      {children}
    </div>
  )
)
Card.displayName = "Card"

// 卡片头部
const CardHeader = React.forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement>
>(({ className, ...props }, ref) => (
  <div
    ref={ref}
    className={cn(
      "flex flex-col space-y-1.5 p-5",
      "border-b border-stroke/30",
      className
    )}
    {...props}
  />
))
CardHeader.displayName = "CardHeader"

// 卡片标题
const CardTitle = React.forwardRef<
  HTMLHeadingElement,
  React.HTMLAttributes<HTMLHeadingElement>
>(({ className, ...props }, ref) => (
  <h3
    ref={ref}
    className={cn(
      "font-display text-xl font-semibold leading-tight tracking-wide",
      "text-fg-primary",
      className
    )}
    {...props}
  />
))
CardTitle.displayName = "CardTitle"

// 卡片描述
const CardDescription = React.forwardRef<
  HTMLParagraphElement,
  React.HTMLAttributes<HTMLParagraphElement>
>(({ className, ...props }, ref) => (
  <p
    ref={ref}
    className={cn(
      "text-sm text-fg-secondary leading-relaxed",
      className
    )}
    {...props}
  />
))
CardDescription.displayName = "CardDescription"

// 卡片内容
const CardContent = React.forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement>
>(({ className, ...props }, ref) => (
  <div
    ref={ref}
    className={cn("p-5", className)}
    {...props}
  />
))
CardContent.displayName = "CardContent"

// 卡片底部
const CardFooter = React.forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement>
>(({ className, ...props }, ref) => (
  <div
    ref={ref}
    className={cn(
      "flex items-center p-5 pt-0",
      className
    )}
    {...props}
  />
))
CardFooter.displayName = "CardFooter"

export {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
  CardFooter,
  cardVariants
}
