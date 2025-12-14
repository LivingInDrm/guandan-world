import * as React from "react"
import { Slot } from "@radix-ui/react-slot"
import { cva, type VariantProps } from "class-variance-authority"
import { cn } from "@/lib/utils"

const buttonVariants = cva(
  [
    "inline-flex items-center justify-center gap-2",
    "font-bold whitespace-nowrap",
    "rounded-full",
    "shadow-[var(--elevation-2),var(--shadow-relief)]",
    "transition-all duration-fast ease-bounce",
    "hover:shadow-elevation-3 hover:scale-105 hover:brightness-110",
    "active:scale-95 active:brightness-90",
    "disabled:opacity-50 disabled:cursor-not-allowed disabled:scale-100",
  ],
  {
    variants: {
      intent: {
        primary: "bg-action-primary text-fg-inverse",
        secondary: "bg-action-secondary text-fg-inverse",
        tertiary: "bg-action-tertiary text-fg-inverse",
        neutral: "bg-action-neutral text-fg-inverse",
        danger: "bg-action-danger text-fg-inverse",
      },
      size: {
        sm: "px-6 py-2 text-base min-w-[90px]",
        md: "px-8 py-3 text-lg min-w-[100px]",
        lg: "px-10 py-4 text-xl min-w-[120px]",
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
}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, intent, size, asChild = false, ...props }, ref) => {
    const Comp = asChild ? Slot : "button"
    return (
      <Comp
        className={cn(buttonVariants({ intent, size, className }))}
        ref={ref}
        {...props}
      />
    )
  }
)
Button.displayName = "Button"

export { Button, buttonVariants }
