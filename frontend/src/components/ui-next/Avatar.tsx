import * as React from "react"
import * as AvatarPrimitive from "@radix-ui/react-avatar"
import { cva, type VariantProps } from "class-variance-authority"
import { cn } from "@/lib/utils"

const avatarVariants = cva(
  [
    "relative flex shrink-0 rounded-ds-lg overflow-hidden",
    "transition-all duration-ds-normal",
  ],
  {
    variants: {
      size: {
        sm: "w-10 h-10",
        md: "w-14 h-14",
        lg: "w-16 h-16",
        xl: "w-20 h-20",
        "2xl": "w-24 h-24",
      },
      ringState: {
        none: "",
        normal: "ring-4 ring-ds-border ring-offset-2",
        active: [
          "ring-4 ring-ds-state-active ring-offset-2",
          "animate-ds-pulse-ring",
        ],
        teamUs: "ring-4 ring-ds-team-us ring-offset-2",
        teamThem: "ring-4 ring-ds-team-them ring-offset-2",
      }
    },
    defaultVariants: {
      size: "md",
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
}

const Avatar = React.forwardRef<
  React.ElementRef<typeof AvatarPrimitive.Root>,
  AvatarProps
>(({ className, size, ringState, src, alt, fallback, ...props }, ref) => (
  <AvatarPrimitive.Root
    ref={ref}
    className={cn(avatarVariants({ size, ringState }), className)}
    {...props}
  >
    <AvatarPrimitive.Image
      src={src}
      alt={alt}
      className="aspect-square h-full w-full object-cover"
    />
    <AvatarPrimitive.Fallback className="flex h-full w-full items-center justify-center bg-ds-surface-elevated text-ds-text-secondary">
      {fallback?.slice(0, 2).toUpperCase()}
    </AvatarPrimitive.Fallback>
  </AvatarPrimitive.Root>
))
Avatar.displayName = "Avatar"

export { Avatar, avatarVariants }
