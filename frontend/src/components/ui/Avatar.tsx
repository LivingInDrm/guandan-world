import * as React from "react"
import * as AvatarPrimitive from "@radix-ui/react-avatar"
import { cva, type VariantProps } from "class-variance-authority"
import { cn } from "@/lib/utils"

const avatarVariants = cva(
  [
    "relative flex shrink-0 overflow-hidden",
    // 阴影 - 精致立体感
    "shadow-[var(--elevation-2),var(--shadow-relief)]",
    // 过渡动画
    "transition-all duration-normal ease-bounce",
  ],
  {
    variants: {
      size: {
        xs: "w-8 h-8 text-xs",
        sm: "w-10 h-10 text-sm",
        md: "w-14 h-14 text-base",
        lg: "w-16 h-16 text-lg",
        xl: "w-20 h-20 text-xl",
        "2xl": "w-24 h-24 text-2xl",
      },
      shape: {
        circle: "rounded-full",
        rounded: "rounded-lg",
        square: "rounded-sm",
      },
      ringState: {
        none: "",
        normal: [
          "ring-[3px] ring-stroke ring-offset-2 ring-offset-surface-base",
        ],
        active: [
          "ring-[3px] ring-state-active ring-offset-2 ring-offset-surface-base",
          "shadow-[var(--elevation-2),var(--shadow-relief),var(--shadow-glow-md)]",
          "animate-pulse-ring",
        ],
        teamUs: [
          "ring-[3px] ring-team-us ring-offset-2 ring-offset-surface-base",
          "shadow-[var(--elevation-2),var(--shadow-relief),var(--shadow-jade)]",
        ],
        teamThem: [
          "ring-[3px] ring-team-them ring-offset-2 ring-offset-surface-base",
          "shadow-[var(--elevation-2),var(--shadow-relief),var(--shadow-ruby)]",
        ],
        premium: [
          "ring-[3px] ring-transparent ring-offset-2 ring-offset-surface-base",
          // 金色渐变边框
          "[--tw-ring-color:transparent]",
          "before:absolute before:inset-[-4px] before:rounded-full before:-z-10",
          "before:bg-gradient-to-br before:from-[hsl(42,95%,68%)] before:via-[hsl(42,95%,52%)] before:to-[hsl(42,95%,38%)]",
          "before:animate-spin-slow",
          "shadow-[var(--elevation-2),var(--shadow-relief),var(--shadow-glow-lg)]",
        ],
      }
    },
    defaultVariants: {
      size: "md",
      shape: "circle",
      ringState: "none",
    }
  }
)

export interface AvatarProps
  extends React.ComponentPropsWithoutRef<typeof AvatarPrimitive.Root>,
    VariantProps<typeof avatarVariants> {
  src?: string
  alt?: string
  fallback?: string
  status?: 'online' | 'offline' | 'away' | 'busy'
}

const statusColors = {
  online: "bg-action-primary",
  offline: "bg-fg-secondary",
  away: "bg-action-secondary",
  busy: "bg-action-danger",
}

const Avatar = React.forwardRef<
  React.ElementRef<typeof AvatarPrimitive.Root>,
  AvatarProps
>(({ className, size, shape, ringState, src, alt, fallback, status, ...props }, ref) => (
  <div className="relative inline-flex">
    <AvatarPrimitive.Root
      ref={ref}
      className={cn(avatarVariants({ size, shape, ringState }), className)}
      {...props}
    >
      <AvatarPrimitive.Image
        src={src}
        alt={alt}
        className={cn(
          "aspect-square h-full w-full object-cover",
          // 边框装饰
          "border border-white/30",
          shape === 'circle' && "rounded-full",
          shape === 'rounded' && "rounded-lg",
          shape === 'square' && "rounded-sm",
        )}
      />
      <AvatarPrimitive.Fallback
        className={cn(
          "flex h-full w-full items-center justify-center",
          // 渐变背景
          "bg-gradient-to-br from-surface-elevated to-[hsl(35,6%,78%)]",
          // 边框
          "border border-stroke/50",
          // 文字样式
          "font-display font-semibold text-fg-secondary",
          shape === 'circle' && "rounded-full",
          shape === 'rounded' && "rounded-lg",
          shape === 'square' && "rounded-sm",
        )}
      >
        {fallback?.slice(0, 2).toUpperCase()}
      </AvatarPrimitive.Fallback>
    </AvatarPrimitive.Root>

    {/* 在线状态指示器 */}
    {status && (
      <span
        className={cn(
          "absolute bottom-0 right-0",
          "block rounded-full",
          "border-2 border-surface-base",
          "shadow-elevation-1",
          statusColors[status],
          // 尺寸根据头像大小调整
          size === 'xs' && "w-2 h-2",
          size === 'sm' && "w-2.5 h-2.5",
          size === 'md' && "w-3 h-3",
          size === 'lg' && "w-3.5 h-3.5",
          size === 'xl' && "w-4 h-4",
          size === '2xl' && "w-5 h-5",
          // 在线状态脉冲动画
          status === 'online' && "animate-pulse",
        )}
      />
    )}
  </div>
))
Avatar.displayName = "Avatar"

// 头像组 - 用于显示多个头像
const AvatarGroup = React.forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement> & {
    max?: number
    size?: VariantProps<typeof avatarVariants>['size']
  }
>(({ className, children, max = 4, size = 'md', ...props }, ref) => {
  const childArray = React.Children.toArray(children)
  const visibleChildren = childArray.slice(0, max)
  const remainingCount = childArray.length - max

  return (
    <div
      ref={ref}
      className={cn("flex -space-x-3", className)}
      {...props}
    >
      {visibleChildren.map((child, index) => (
        <div
          key={index}
          className="relative"
          style={{ zIndex: visibleChildren.length - index }}
        >
          {React.isValidElement(child)
            ? React.cloneElement(child as React.ReactElement<AvatarProps>, { size })
            : child
          }
        </div>
      ))}
      {remainingCount > 0 && (
        <div
          className={cn(
            avatarVariants({ size, shape: 'circle' }),
            "flex items-center justify-center",
            "bg-gradient-to-br from-state-active to-[hsl(42,95%,42%)]",
            "text-fg-inverse font-display font-semibold",
            "border-2 border-surface-base",
          )}
          style={{ zIndex: 0 }}
        >
          +{remainingCount}
        </div>
      )}
    </div>
  )
})
AvatarGroup.displayName = "AvatarGroup"

export { Avatar, AvatarGroup, avatarVariants }
