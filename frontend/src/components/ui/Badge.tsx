import * as React from "react"
import { cva, type VariantProps } from "class-variance-authority"
import { cn } from "@/lib/utils"

const badgeVariants = cva(
  [
    "inline-flex items-center justify-center",
    "font-bold backdrop-blur-sm",
    "rounded-sm",
    "shadow-elevation-1",
    "transition-all duration-fast",
  ],
  {
    variants: {
      variant: {
        landlord: "bg-state-active text-fg-inverse",
        farmer: "bg-team-us text-fg-inverse",
        teammate: "bg-team-us/80 text-fg-inverse border border-team-us",
        owner: "bg-state-active text-fg-inverse",
        neutral: "bg-surface-elevated text-fg-secondary",
      },
      size: {
        sm: "px-1.5 py-0.5 text-xs",
        md: "px-2 py-1 text-sm",
      }
    },
    defaultVariants: {
      variant: "teammate",
      size: "sm",
    }
  }
)

export interface BadgeProps
  extends React.HTMLAttributes<HTMLDivElement>,
    VariantProps<typeof badgeVariants> {}

const Badge = React.forwardRef<HTMLDivElement, BadgeProps>(
  ({ className, variant, size, ...props }, ref) => (
    <div
      ref={ref}
      className={cn(badgeVariants({ variant, size }), className)}
      {...props}
    />
  )
)
Badge.displayName = "Badge"

export { Badge, badgeVariants }
