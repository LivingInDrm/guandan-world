import * as React from "react"
import { Slot } from "@radix-ui/react-slot"
import { cva, type VariantProps } from "class-variance-authority"
import { cn } from "@/lib/utils"

const buttonVariants = cva(
  "inline-flex items-center justify-center gap-2 whitespace-nowrap font-semibold border shadow-card transition-all duration-200 hover:scale-105 hover:shadow-card-hover active:scale-95 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:bg-disabled-bg disabled:text-disabled-text disabled:border-disabled-border disabled:scale-100 disabled:shadow-none [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0",
  {
    variants: {
      variant: {
        primary:
          "bg-gradient-primary text-white border-white/30 hover:brightness-110 active:brightness-90",
        secondary:
          "bg-gradient-secondary text-white border-white/30 hover:brightness-110 active:brightness-90",
        warning:
          "bg-gradient-warning text-white border-white/30 hover:brightness-110 active:brightness-90",
        danger:
          "bg-gradient-danger text-white border-white/30 hover:brightness-110 active:brightness-90",
        ghost:
          "bg-white/70 backdrop-blur-sm text-table-400 border-table-300 hover:bg-white/90 active:bg-white/60 shadow-none",
        outline:
          "border border-input bg-background hover:bg-accent hover:text-accent-foreground",
        link: "text-primary underline-offset-4 hover:underline shadow-none border-none",
      },
      size: {
        sm: "py-1.5 px-4 text-sm rounded-sm",
        md: "py-2 px-6 text-base rounded-lg",
        lg: "py-3 px-8 text-lg rounded-lg",
        icon: "h-10 w-10 rounded-md",
      },
    },
    defaultVariants: {
      variant: "primary",
      size: "md",
    },
  }
)

export interface ButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement>,
    VariantProps<typeof buttonVariants> {
  asChild?: boolean
}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant, size, asChild = false, ...props }, ref) => {
    const Comp = asChild ? Slot : "button"
    return (
      <Comp
        className={cn(buttonVariants({ variant, size, className }))}
        ref={ref}
        {...props}
      />
    )
  }
)
Button.displayName = "Button"

export { Button, buttonVariants }
