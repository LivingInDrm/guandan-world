import * as React from "react"
import { cva, type VariantProps } from "class-variance-authority"
import { cn } from "@/lib/utils"

const cardVariants = cva(
  [
    "relative rounded-lg",
    "transition-all duration-normal ease-bounce",
    "hover:scale-[1.02]",
  ],
  {
    variants: {
      variant: {
        base: [
          "bg-surface-base",
          "shadow-relief",
        ],
        elevated: [
          "bg-surface-elevated",
          "shadow-[var(--elevation-3),var(--shadow-relief)]",
        ],
        emphasis: [
          "bg-surface-emphasis backdrop-blur-xl",
          "border border-state-active",
          "shadow-[var(--shadow-relief),var(--shadow-glow-lg)]",
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
