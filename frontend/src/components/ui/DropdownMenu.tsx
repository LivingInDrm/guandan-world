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

// 共享的菜单内容样式 - 精致牌桌风格
const menuContentStyles = [
  "z-50 min-w-[200px] overflow-hidden",
  // 内边距
  "p-2",
  // 背景 - 奶油色毛玻璃效果，呼应牌桌质感
  "bg-gradient-to-b from-[hsl(40,20%,98%)/98%] via-[hsl(38,15%,96%)/97%] to-[hsl(35,12%,94%)/96%]",
  "backdrop-blur-xl backdrop-saturate-150",
  // 边框 - 精致双层边框
  "border border-[hsl(38,25%,88%)]/80",
  "rounded-2xl",
  // 顶部翡翠绿装饰线
  "before:absolute before:top-0 before:inset-x-3 before:h-[2px]",
  "before:bg-gradient-to-r before:from-transparent before:via-[hsl(158,55%,42%)]/60 before:to-transparent",
  "before:rounded-full",
  // 阴影 - 多层次立体感，带微妙绿色光晕
  "shadow-[0_12px_48px_-12px_rgba(45,106,79,0.15),0_8px_24px_-8px_rgba(0,0,0,0.12),0_0_0_1px_rgba(255,255,255,0.6)_inset]",
  // 底部微妙反光
  "after:absolute after:bottom-0 after:inset-x-4 after:h-[1px]",
  "after:bg-gradient-to-r after:from-transparent after:via-white/40 after:to-transparent",
  // 动画
  "data-[state=open]:animate-in data-[state=closed]:animate-out",
  "data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0",
  "data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95",
  "data-[side=bottom]:slide-in-from-top-2",
  "data-[side=left]:slide-in-from-right-2",
  "data-[side=right]:slide-in-from-left-2",
  "data-[side=top]:slide-in-from-bottom-2",
  // 相对定位用于装饰元素
  "relative",
]

// 共享的菜单项样式 - 优雅交互
const menuItemStyles = [
  "relative flex cursor-pointer select-none items-center",
  "rounded-xl px-3 py-2.5 text-sm",
  "outline-none",
  // 文字颜色和字体
  "text-fg-primary font-display font-medium",
  // 过渡动画
  "transition-all duration-200 ease-out",
  // 悬停效果 - 翡翠绿渐变
  "hover:bg-gradient-to-r hover:from-[hsl(158,45%,94%)] hover:to-[hsl(158,40%,96%)]/50",
  "hover:text-[hsl(158,50%,28%)]",
  "hover:shadow-[inset_0_0_0_1px_hsl(158,45%,85%)]",
  // 聚焦效果
  "focus:bg-gradient-to-r focus:from-[hsl(158,45%,94%)] focus:to-[hsl(158,40%,96%)]/50",
  "focus:text-[hsl(158,50%,28%)]",
  "focus:shadow-[inset_0_0_0_1px_hsl(158,45%,85%)]",
  // 禁用状态
  "data-[disabled]:pointer-events-none data-[disabled]:opacity-40",
]

const DropdownMenuSubTrigger = React.forwardRef<
  React.ComponentRef<typeof DropdownMenuPrimitive.SubTrigger>,
  React.ComponentPropsWithoutRef<typeof DropdownMenuPrimitive.SubTrigger> & {
    inset?: boolean
  }
>(({ className, inset, children, ...props }, ref) => (
  <DropdownMenuPrimitive.SubTrigger
    ref={ref}
    className={cn(
      ...menuItemStyles,
      "data-[state=open]:bg-gradient-to-r data-[state=open]:from-[hsl(158,45%,92%)] data-[state=open]:to-transparent",
      "data-[state=open]:text-[hsl(158,50%,28%)]",
      inset && "pl-8",
      className
    )}
    {...props}
  >
    {children}
    <ChevronRight className="ml-auto h-4 w-4 text-[hsl(158,55%,42%)] opacity-60 transition-all duration-200 group-data-[state=open]:rotate-90 group-data-[state=open]:opacity-100" />
  </DropdownMenuPrimitive.SubTrigger>
))
DropdownMenuSubTrigger.displayName = DropdownMenuPrimitive.SubTrigger.displayName

const DropdownMenuSubContent = React.forwardRef<
  React.ComponentRef<typeof DropdownMenuPrimitive.SubContent>,
  React.ComponentPropsWithoutRef<typeof DropdownMenuPrimitive.SubContent>
>(({ className, ...props }, ref) => (
  <DropdownMenuPrimitive.SubContent
    ref={ref}
    className={cn(...menuContentStyles, className)}
    {...props}
  />
))
DropdownMenuSubContent.displayName = DropdownMenuPrimitive.SubContent.displayName

const DropdownMenuContent = React.forwardRef<
  React.ComponentRef<typeof DropdownMenuPrimitive.Content>,
  React.ComponentPropsWithoutRef<typeof DropdownMenuPrimitive.Content>
>(({ className, sideOffset = 8, ...props }, ref) => (
  <DropdownMenuPrimitive.Portal>
    <DropdownMenuPrimitive.Content
      ref={ref}
      sideOffset={sideOffset}
      className={cn(...menuContentStyles, className)}
      {...props}
    />
  </DropdownMenuPrimitive.Portal>
))
DropdownMenuContent.displayName = DropdownMenuPrimitive.Content.displayName

const DropdownMenuItem = React.forwardRef<
  React.ComponentRef<typeof DropdownMenuPrimitive.Item>,
  React.ComponentPropsWithoutRef<typeof DropdownMenuPrimitive.Item> & {
    inset?: boolean
    destructive?: boolean
  }
>(({ className, inset, destructive, ...props }, ref) => (
  <DropdownMenuPrimitive.Item
    ref={ref}
    className={cn(
      ...menuItemStyles,
      destructive && [
        "text-action-danger",
        "hover:bg-gradient-to-r hover:from-action-danger/10 hover:to-action-danger/5",
        "hover:text-action-danger",
        "hover:shadow-[inset_0_0_0_1px_hsl(0,70%,90%)]",
        "focus:bg-gradient-to-r focus:from-action-danger/10 focus:to-action-danger/5",
        "focus:text-action-danger",
        "focus:shadow-[inset_0_0_0_1px_hsl(0,70%,90%)]",
      ],
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
      ...menuItemStyles,
      "pl-9",
      className
    )}
    checked={checked}
    {...props}
  >
    <span className="absolute left-3 flex h-4 w-4 items-center justify-center">
      <DropdownMenuPrimitive.ItemIndicator>
        <Check className="h-4 w-4 text-[hsl(158,55%,42%)]" />
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
      ...menuItemStyles,
      "pl-9",
      className
    )}
    {...props}
  >
    <span className="absolute left-3 flex h-4 w-4 items-center justify-center">
      <DropdownMenuPrimitive.ItemIndicator>
        <Circle className="h-2.5 w-2.5 fill-[hsl(158,55%,42%)] text-[hsl(158,55%,42%)]" />
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
      "px-3 py-2 text-xs font-display font-bold tracking-wider uppercase",
      // 翡翠绿渐变文字
      "bg-gradient-to-r from-[hsl(158,55%,35%)] via-[hsl(158,50%,42%)] to-[hsl(158,55%,35%)]",
      "bg-clip-text text-transparent",
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
    className={cn(
      "-mx-1 my-2 h-px",
      "bg-gradient-to-r from-transparent via-[hsl(158,30%,80%)] to-transparent",
      className
    )}
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
      className={cn(
        "ml-auto text-[10px] tracking-widest uppercase",
        "text-[hsl(158,40%,55%)]",
        "font-display font-medium",
        "px-1.5 py-0.5 rounded-md",
        "bg-[hsl(158,30%,95%)]",
        className
      )}
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
