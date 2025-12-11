import * as React from "react"
import * as LabelPrimitive from "@radix-ui/react-label"
import { cn } from "@/lib/utils"

export interface InputProps
  extends React.InputHTMLAttributes<HTMLInputElement> {}

const Input = React.forwardRef<HTMLInputElement, InputProps>(
  ({ className, type, ...props }, ref) => {
    return (
      <input
        type={type}
        className={cn(
          "flex h-10 w-full px-3 py-2",
          "border border-stroke bg-surface-elevated text-fg-primary",
          "rounded-md shadow-[var(--elevation-2),var(--shadow-relief)]",
          "placeholder:text-fg-secondary",
          "focus-visible:outline-none focus-visible:border-state-active",
          "focus-visible:shadow-[var(--elevation-3),var(--shadow-relief),0_0_12px_hsla(45,100%,51%,0.3)]",
          "focus-visible:scale-[1.01]",
          "disabled:cursor-not-allowed disabled:bg-state-disabled disabled:text-fg-secondary disabled:scale-100",
          "transition-all duration-normal ease-bounce",
          className
        )}
        ref={ref}
        {...props}
      />
    )
  }
)
Input.displayName = "Input"

const Label = React.forwardRef<
  React.ComponentRef<typeof LabelPrimitive.Root>,
  React.ComponentPropsWithoutRef<typeof LabelPrimitive.Root>
>(({ className, ...props }, ref) => (
  <LabelPrimitive.Root
    ref={ref}
    className={cn(
      "text-sm font-medium leading-none text-fg-primary",
      "peer-disabled:cursor-not-allowed peer-disabled:opacity-70",
      className
    )}
    {...props}
  />
))
Label.displayName = LabelPrimitive.Root.displayName

export { Input, Label }
