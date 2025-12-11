import * as React from "react"
import * as SliderPrimitive from "@radix-ui/react-slider"
import { cn } from "@/lib/utils"

const Slider = React.forwardRef<
  React.ComponentRef<typeof SliderPrimitive.Root>,
  React.ComponentPropsWithoutRef<typeof SliderPrimitive.Root>
>(({ className, ...props }, ref) => (
  <SliderPrimitive.Root
    ref={ref}
    className={cn(
      "relative flex w-full touch-none select-none items-center",
      className
    )}
    {...props}
  >
    <SliderPrimitive.Track className="relative h-2 w-full grow overflow-hidden bg-ds-state-disabled rounded-ds-sm shadow-[inset_0_1px_2px_rgba(0,0,0,0.1),var(--ds-shadow-relief)] transition-all duration-ds-normal ease-ds-smooth">
      <SliderPrimitive.Range className="absolute h-full bg-ds-state-active transition-all duration-ds-normal ease-ds-smooth" />
    </SliderPrimitive.Track>
    <SliderPrimitive.Thumb className="block h-5 w-5 rounded-full border-2 border-ds-state-active bg-ds-surface-base shadow-[var(--ds-elevation-2),var(--ds-shadow-relief)] transition-all duration-ds-fast ease-ds-bounce hover:scale-110 hover:shadow-[var(--ds-elevation-3),var(--ds-shadow-glow-sm)] active:scale-125 active:shadow-[var(--ds-elevation-3),var(--ds-shadow-glow-md)] focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ds-state-active focus-visible:ring-offset-2 focus-visible:shadow-[var(--ds-elevation-3),var(--ds-shadow-glow-md)] disabled:pointer-events-none disabled:opacity-50 disabled:scale-100 cursor-grab active:cursor-grabbing" />
  </SliderPrimitive.Root>
))
Slider.displayName = SliderPrimitive.Root.displayName

export { Slider }
