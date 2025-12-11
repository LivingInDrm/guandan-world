import * as React from "react"
import { cva, type VariantProps } from "class-variance-authority"
import { cn } from "@/lib/utils"

const cardVariants = cva(
  [
    "relative rounded-ds-lg",
    "transition-all duration-ds-normal ease-ds-bounce",
    "hover:scale-[1.02]",
  ],
  {
    variants: {
      variant: {
        base: [
          "bg-ds-surface-base",
          "shadow-ds-relief",
        ],
        elevated: [
          "bg-ds-surface-elevated",
          "shadow-[var(--ds-elevation-3),var(--ds-shadow-relief)]",
        ],
        emphasis: [
          "bg-ds-surface-emphasis backdrop-blur-xl",
          "border border-ds-state-active",
          "shadow-[var(--ds-shadow-relief),var(--ds-shadow-glow-lg)]",
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
  interactive?: boolean;
}

const Card = React.forwardRef<HTMLDivElement, CardProps>(
  ({ className, variant, interactive = true, children, ...props }, ref) => (
    <div
      ref={ref}
      className={cn(
        cardVariants({ variant }),
        !interactive && "hover:scale-100",
        className
      )}
      {...props}
    >
      {children}
    </div>
  )
)
Card.displayName = "Card"

export { Card, cardVariants }
