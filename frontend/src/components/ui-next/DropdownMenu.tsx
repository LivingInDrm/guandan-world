import * as React from "react"
import * as DropdownMenuPrimitive from "@radix-ui/react-dropdown-menu"
import { cn } from "@/lib/utils"
import { Check, ChevronRight, Circle } from "lucide-react"

const DropdownMenu = DropdownMenuPrimitive.Root

const DropdownMenuTrigger = DropdownMenuPrimitive.Trigger

const DropdownMenuGroup = DropdownMenuPrimitive.Group

const DropdownMenuPortal = DropdownMenuPrimitive.Portal

const DropdownMenuSub = DropdownMenuPrimitive.Sub

const DropdownMenuRadioGroup = DropdownMenuPrimitive.RadioGroup

const DropdownMenuSubTrigger = React.forwardRef<
  React.ComponentRef<typeof DropdownMenuPrimitive.SubTrigger>,
  React.ComponentPropsWithoutRef<typeof DropdownMenuPrimitive.SubTrigger> & {
    inset?: boolean
  }
>(({ className, inset, children, ...props }, ref) => (
  <DropdownMenuPrimitive.SubTrigger
    ref={ref}
    className={cn(
      "flex cursor-default select-none items-center rounded-ds-md px-3 py-2 text-sm outline-none",
      "text-ds-text-primary border-l-2 border-transparent",
      "hover:bg-gradient-to-r hover:from-ds-surface-emphasis hover:to-transparent",
      "hover:border-ds-state-active",
      "hover:shadow-[0_0_8px_hsla(45,100%,51%,0.3)]",
      "hover:scale-[1.02]",
      "focus:bg-gradient-to-r focus:from-ds-surface-emphasis focus:to-transparent",
      "focus:border-ds-state-active",
      "data-[state=open]:bg-gradient-to-r data-[state=open]:from-ds-surface-emphasis data-[state=open]:to-transparent",
      "data-[state=open]:border-ds-state-active",
      "transition-all duration-ds-normal ease-ds-bounce",
      inset && "pl-8",
      className
    )}
    {...props}
  >
    {children}
    <ChevronRight className="ml-auto h-4 w-4 text-ds-state-active" />
  </DropdownMenuPrimitive.SubTrigger>
))
DropdownMenuSubTrigger.displayName = DropdownMenuPrimitive.SubTrigger.displayName

const DropdownMenuSubContent = React.forwardRef<
  React.ComponentRef<typeof DropdownMenuPrimitive.SubContent>,
  React.ComponentPropsWithoutRef<typeof DropdownMenuPrimitive.SubContent>
>(({ className, ...props }, ref) => (
  <DropdownMenuPrimitive.SubContent
    ref={ref}
    className={cn(
      "z-50 min-w-[8rem] overflow-hidden p-1.5",
      "bg-ds-surface-elevated backdrop-blur-sm",
      "shadow-[var(--ds-elevation-3),var(--ds-shadow-relief),0_0_12px_hsla(145,60%,50%,0.15)]",
      "border border-ds-border-emphasis rounded-ds-lg",
      "data-[state=open]:animate-in data-[state=closed]:animate-out",
      "data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0",
      "data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95",
      "data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2",
      "data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2",
      className
    )}
    {...props}
  />
))
DropdownMenuSubContent.displayName = DropdownMenuPrimitive.SubContent.displayName

const DropdownMenuContent = React.forwardRef<
  React.ComponentRef<typeof DropdownMenuPrimitive.Content>,
  React.ComponentPropsWithoutRef<typeof DropdownMenuPrimitive.Content>
>(({ className, sideOffset = 4, ...props }, ref) => (
  <DropdownMenuPrimitive.Portal>
    <DropdownMenuPrimitive.Content
      ref={ref}
      sideOffset={sideOffset}
      className={cn(
        "z-50 min-w-[8rem] overflow-hidden p-1.5",
        "bg-ds-surface-elevated backdrop-blur-sm",
        "shadow-[var(--ds-elevation-3),var(--ds-shadow-relief),0_0_12px_hsla(145,60%,50%,0.15)]",
        "border border-ds-border-emphasis rounded-ds-lg",
        "data-[state=open]:animate-in data-[state=closed]:animate-out",
        "data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0",
        "data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95",
        "data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2",
        "data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2",
        className
      )}
      {...props}
    />
  </DropdownMenuPrimitive.Portal>
))
DropdownMenuContent.displayName = DropdownMenuPrimitive.Content.displayName

const DropdownMenuItem = React.forwardRef<
  React.ComponentRef<typeof DropdownMenuPrimitive.Item>,
  React.ComponentPropsWithoutRef<typeof DropdownMenuPrimitive.Item> & {
    inset?: boolean
  }
>(({ className, inset, ...props }, ref) => (
  <DropdownMenuPrimitive.Item
    ref={ref}
    className={cn(
      "relative flex cursor-default select-none items-center rounded-ds-md px-3 py-2 text-sm outline-none",
      "text-ds-text-primary border-l-2 border-transparent",
      "hover:bg-gradient-to-r hover:from-ds-surface-emphasis hover:to-transparent",
      "hover:border-ds-state-active",
      "hover:shadow-[0_0_8px_hsla(45,100%,51%,0.3)]",
      "hover:scale-[1.02]",
      "focus:bg-gradient-to-r focus:from-ds-surface-emphasis focus:to-transparent",
      "focus:border-ds-state-active",
      "data-[disabled]:pointer-events-none data-[disabled]:opacity-50",
      "transition-all duration-ds-normal ease-ds-bounce",
      inset && "pl-8",
      className
    )}
    {...props}
  />
))
DropdownMenuItem.displayName = DropdownMenuPrimitive.Item.displayName

const DropdownMenuCheckboxItem = React.forwardRef<
  React.ComponentRef<typeof DropdownMenuPrimitive.CheckboxItem>,
  React.ComponentPropsWithoutRef<typeof DropdownMenuPrimitive.CheckboxItem>
>(({ className, children, checked, ...props }, ref) => (
  <DropdownMenuPrimitive.CheckboxItem
    ref={ref}
    className={cn(
      "relative flex cursor-default select-none items-center rounded-ds-md py-2 pl-9 pr-3 text-sm outline-none",
      "text-ds-text-primary border-l-2 border-transparent",
      "hover:bg-gradient-to-r hover:from-ds-surface-emphasis hover:to-transparent",
      "hover:border-ds-state-active",
      "hover:shadow-[0_0_8px_hsla(45,100%,51%,0.3)]",
      "hover:scale-[1.02]",
      "focus:bg-gradient-to-r focus:from-ds-surface-emphasis focus:to-transparent",
      "focus:border-ds-state-active",
      "data-[disabled]:pointer-events-none data-[disabled]:opacity-50",
      "transition-all duration-ds-normal ease-ds-bounce",
      className
    )}
    checked={checked}
    {...props}
  >
    <span className="absolute left-2 flex h-3.5 w-3.5 items-center justify-center">
      <DropdownMenuPrimitive.ItemIndicator>
        <Check className="h-4 w-4 text-ds-state-active" />
      </DropdownMenuPrimitive.ItemIndicator>
    </span>
    {children}
  </DropdownMenuPrimitive.CheckboxItem>
))
DropdownMenuCheckboxItem.displayName = DropdownMenuPrimitive.CheckboxItem.displayName

const DropdownMenuRadioItem = React.forwardRef<
  React.ComponentRef<typeof DropdownMenuPrimitive.RadioItem>,
  React.ComponentPropsWithoutRef<typeof DropdownMenuPrimitive.RadioItem>
>(({ className, children, ...props }, ref) => (
  <DropdownMenuPrimitive.RadioItem
    ref={ref}
    className={cn(
      "relative flex cursor-default select-none items-center rounded-ds-md py-2 pl-9 pr-3 text-sm outline-none",
      "text-ds-text-primary border-l-2 border-transparent",
      "hover:bg-gradient-to-r hover:from-ds-surface-emphasis hover:to-transparent",
      "hover:border-ds-state-active",
      "hover:shadow-[0_0_8px_hsla(45,100%,51%,0.3)]",
      "hover:scale-[1.02]",
      "focus:bg-gradient-to-r focus:from-ds-surface-emphasis focus:to-transparent",
      "focus:border-ds-state-active",
      "data-[disabled]:pointer-events-none data-[disabled]:opacity-50",
      "transition-all duration-ds-normal ease-ds-bounce",
      className
    )}
    {...props}
  >
    <span className="absolute left-2 flex h-3.5 w-3.5 items-center justify-center">
      <DropdownMenuPrimitive.ItemIndicator>
        <Circle className="h-2 w-2 fill-current text-ds-state-active" />
      </DropdownMenuPrimitive.ItemIndicator>
    </span>
    {children}
  </DropdownMenuPrimitive.RadioItem>
))
DropdownMenuRadioItem.displayName = DropdownMenuPrimitive.RadioItem.displayName

const DropdownMenuLabel = React.forwardRef<
  React.ComponentRef<typeof DropdownMenuPrimitive.Label>,
  React.ComponentPropsWithoutRef<typeof DropdownMenuPrimitive.Label> & {
    inset?: boolean
  }
>(({ className, inset, ...props }, ref) => (
  <DropdownMenuPrimitive.Label
    ref={ref}
    className={cn(
      "px-3 py-2 text-sm font-semibold",
      "bg-gradient-to-r from-ds-state-active to-ds-action-primary bg-clip-text text-transparent",
      inset && "pl-8",
      className
    )}
    {...props}
  />
))
DropdownMenuLabel.displayName = DropdownMenuPrimitive.Label.displayName

const DropdownMenuSeparator = React.forwardRef<
  React.ComponentRef<typeof DropdownMenuPrimitive.Separator>,
  React.ComponentPropsWithoutRef<typeof DropdownMenuPrimitive.Separator>
>(({ className, ...props }, ref) => (
  <DropdownMenuPrimitive.Separator
    ref={ref}
    className={cn("-mx-1 my-1.5 h-[1.5px] bg-gradient-to-r from-transparent via-ds-border to-transparent", className)}
    {...props}
  />
))
DropdownMenuSeparator.displayName = DropdownMenuPrimitive.Separator.displayName

const DropdownMenuShortcut = ({
  className,
  ...props
}: React.HTMLAttributes<HTMLSpanElement>) => {
  return (
    <span
      className={cn("ml-auto text-xs tracking-widest text-ds-state-active opacity-70", className)}
      {...props}
    />
  )
}
DropdownMenuShortcut.displayName = "DropdownMenuShortcut"

export {
  DropdownMenu,
  DropdownMenuTrigger,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuCheckboxItem,
  DropdownMenuRadioItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuShortcut,
  DropdownMenuGroup,
  DropdownMenuPortal,
  DropdownMenuSub,
  DropdownMenuSubContent,
  DropdownMenuSubTrigger,
  DropdownMenuRadioGroup,
}
