import * as React from "react"
import { Slot } from "@radix-ui/react-slot"
import { cva, type VariantProps } from "class-variance-authority"
import { cn } from "@/lib/utils"

const buttonVariants = cva(
  [
    // 基础布局
    "relative inline-flex items-center justify-center gap-2.5",
    "font-display font-medium tracking-wide whitespace-nowrap",
    "overflow-hidden isolate",
    // 形状与尺寸
    "rounded-lg",
    // 精致阴影 - 多层次立体感
    "shadow-[var(--elevation-2),var(--shadow-relief)]",
    // 过渡动画
    "transition-all duration-normal ease-bounce",
    // 悬停效果 - 优雅浮起
    "hover:shadow-[var(--elevation-3),var(--shadow-relief),var(--shadow-glow-sm)]",
    "hover:-translate-y-0.5 hover:brightness-105",
    // 激活效果 - 按下反馈
    "active:translate-y-0 active:scale-[0.98] active:brightness-95",
    "active:shadow-[var(--elevation-1),var(--shadow-relief)]",
    // 禁用状态
    "disabled:opacity-50 disabled:cursor-not-allowed",
    "disabled:translate-y-0 disabled:scale-100 disabled:shadow-none",
    // 聚焦状态
    "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-state-active focus-visible:ring-offset-2",
    // 内部光泽装饰
    "before:absolute before:inset-0 before:rounded-lg",
    "before:bg-gradient-to-b before:from-white/25 before:to-transparent before:to-50%",
    "before:pointer-events-none before:transition-opacity before:duration-normal",
    "hover:before:opacity-80",
  ],
  {
    variants: {
      intent: {
        primary: [
          "bg-gradient-to-b from-action-primary to-[hsl(152,60%,36%)]",
          "text-fg-inverse text-shadow-sm",
          "border border-[hsl(152,60%,30%)]",
          "hover:shadow-[var(--elevation-3),var(--shadow-relief),var(--shadow-jade)]",
        ],
        secondary: [
          "bg-gradient-to-b from-action-secondary to-[hsl(28,85%,44%)]",
          "text-fg-inverse text-shadow-sm",
          "border border-[hsl(28,85%,38%)]",
        ],
        tertiary: [
          "bg-gradient-to-b from-action-tertiary to-[hsl(42,95%,44%)]",
          "text-fg-inverse text-shadow-sm",
          "border border-[hsl(42,95%,38%)]",
          "hover:shadow-[var(--elevation-3),var(--shadow-relief),var(--shadow-glow-md)]",
        ],
        neutral: [
          "bg-gradient-to-b from-surface-elevated to-[hsl(35,6%,78%)]",
          "text-fg-primary",
          "border border-stroke",
        ],
        danger: [
          "bg-gradient-to-b from-action-danger to-[hsl(5,75%,42%)]",
          "text-fg-inverse text-shadow-sm",
          "border border-[hsl(5,75%,35%)]",
          "hover:shadow-[var(--elevation-3),var(--shadow-relief),var(--shadow-ruby)]",
        ],
        ghost: [
          "bg-transparent shadow-none",
          "text-fg-primary",
          "border border-transparent",
          "hover:bg-surface-emphasis hover:border-stroke-emphasis",
          "hover:shadow-none",
          "before:hidden",
        ],
        outline: [
          "bg-transparent",
          "text-state-active",
          "border-2 border-state-active",
          "shadow-none",
          "hover:bg-state-active/10",
          "hover:shadow-[var(--shadow-glow-sm)]",
          "before:hidden",
        ],
      },
      size: {
        xs: "px-4 py-1.5 text-sm min-w-[72px]",
        sm: "px-5 py-2 text-base min-w-[88px]",
        md: "px-7 py-2.5 text-lg min-w-[100px]",
        lg: "px-9 py-3 text-xl min-w-[120px]",
        xl: "px-11 py-4 text-2xl min-w-[140px]",
      },
    },
    defaultVariants: {
      intent: "primary",
      size: "md",
    }
  }
)

export interface ButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement>,
    VariantProps<typeof buttonVariants> {
  asChild?: boolean
  loading?: boolean
}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, intent, size, asChild = false, loading = false, disabled, children, ...props }, ref) => {
    const Comp = asChild ? Slot : "button"
    return (
      <Comp
        className={cn(buttonVariants({ intent, size, className }))}
        ref={ref}
        disabled={disabled || loading}
        {...props}
      >
        {loading ? (
          <>
            <span className="absolute inset-0 flex items-center justify-center">
              <span className="w-5 h-5 border-2 border-current border-t-transparent rounded-full animate-spin" />
            </span>
            <span className="opacity-0">{children}</span>
          </>
        ) : (
          children
        )}
        {/* 点击涟漪效果容器 */}
        <span
          className="absolute inset-0 pointer-events-none"
          aria-hidden="true"
        />
      </Comp>
    )
  }
)
Button.displayName = "Button"

export { Button, buttonVariants }
