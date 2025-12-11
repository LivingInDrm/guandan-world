import * as React from "react"
import { cva, type VariantProps } from "class-variance-authority"
import { cn } from "@/lib/utils"

const badgeVariants = cva(
  [
    "inline-flex items-center justify-center",
    "font-bold backdrop-blur-sm",
    "rounded-ds-sm",
    "shadow-ds-elevation-1",
    "transition-all duration-ds-fast",
  ],
  {
    variants: {
      variant: {
        landlord: "bg-ds-state-active text-ds-text-inverse",
        farmer: "bg-ds-team-us text-ds-text-inverse",
        teammate: "bg-ds-team-us/80 text-ds-text-inverse border border-ds-team-us",
        owner: "bg-ds-state-active text-ds-text-inverse",
        neutral: "bg-ds-surface-elevated text-ds-text-secondary",
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
