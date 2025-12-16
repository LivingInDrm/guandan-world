import * as React from "react"
import * as DialogPrimitive from "@radix-ui/react-dialog"
import { cn } from "@/lib/utils"
import { X } from "lucide-react"

const Dialog = DialogPrimitive.Root

const DialogTrigger = DialogPrimitive.Trigger

const DialogPortal = DialogPrimitive.Portal

const DialogClose = DialogPrimitive.Close

const DialogOverlay = React.forwardRef<
  React.ComponentRef<typeof DialogPrimitive.Overlay>,
  React.ComponentPropsWithoutRef<typeof DialogPrimitive.Overlay>
>(({ className, ...props }, ref) => (
  <DialogPrimitive.Overlay
    ref={ref}
    className={cn(
      "fixed inset-0 z-50",
      // 背景 - 优雅的毛玻璃效果
      "bg-primitive-neutral-900/40 backdrop-blur-sm",
      // 动画
      "data-[state=open]:animate-in data-[state=closed]:animate-out",
      "data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0",
      "duration-300",
      className
    )}
    {...props}
  />
))
DialogOverlay.displayName = DialogPrimitive.Overlay.displayName

const DialogContent = React.forwardRef<
  React.ComponentRef<typeof DialogPrimitive.Content>,
  React.ComponentPropsWithoutRef<typeof DialogPrimitive.Content>
>(({ className, children, ...props }, ref) => (
  <DialogPortal>
    <DialogOverlay />
    <DialogPrimitive.Content
      ref={ref}
      className={cn(
        // 定位
        "fixed left-[50%] top-[50%] z-50 translate-x-[-50%] translate-y-[-50%]",
        // 尺寸
        "grid w-full max-w-lg gap-5 p-6",
        // 背景 - 多层渐变营造深度
        "bg-gradient-to-b from-surface-base to-surface-elevated",
        "backdrop-blur-xl",
        // 边框 - 金色渐变边框
        "border-2 border-transparent rounded-lg",
        "[background:linear-gradient(to_bottom,hsl(0_0%_100%),hsl(40_8%_95%))_padding-box,linear-gradient(135deg,hsl(42_95%_68%/0.6),hsl(42_95%_52%),hsl(42_95%_38%/0.6))_border-box]",
        // 阴影 - 多层次立体感
        "shadow-[var(--elevation-3),var(--shadow-relief),0_0_40px_hsla(42,95%,52%,0.15)]",
        // 动画
        "transition-all duration-300 ease-bounce",
        "data-[state=open]:animate-in data-[state=closed]:animate-out",
        "data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0",
        "data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95",
        "data-[state=closed]:slide-out-to-left-1/2 data-[state=closed]:slide-out-to-top-[48%]",
        "data-[state=open]:slide-in-from-left-1/2 data-[state=open]:slide-in-from-top-[48%]",
        className
      )}
      {...props}
    >
      {/* 装饰性顶部高光 */}
      <div className="absolute inset-x-0 top-0 h-px bg-gradient-to-r from-transparent via-white/80 to-transparent" />

      {children}

      {/* 关闭按钮 */}
      <DialogPrimitive.Close
        className={cn(
          "absolute right-4 top-4",
          "flex items-center justify-center w-8 h-8",
          "rounded-full",
          // 背景渐变
          "bg-gradient-to-b from-surface-elevated to-[hsl(35,6%,82%)]",
          "border border-stroke/50",
          // 阴影
          "shadow-[var(--elevation-1),var(--shadow-relief)]",
          // 图标颜色
          "text-fg-secondary",
          // 过渡动画
          "transition-all duration-normal ease-bounce",
          // 悬停效果
          "hover:text-state-active hover:scale-110",
          "hover:shadow-[var(--elevation-2),var(--shadow-relief),var(--shadow-glow-sm)]",
          "hover:rotate-90",
          // 聚焦状态
          "focus:outline-none focus:ring-2 focus:ring-state-active focus:ring-offset-2",
          // 按下效果
          "active:scale-95",
          // 禁用状态
          "disabled:pointer-events-none disabled:opacity-50"
        )}
      >
        <X className="h-4 w-4" />
        <span className="sr-only">Close</span>
      </DialogPrimitive.Close>
    </DialogPrimitive.Content>
  </DialogPortal>
))
DialogContent.displayName = DialogPrimitive.Content.displayName

const DialogHeader = ({
  className,
  ...props
}: React.HTMLAttributes<HTMLDivElement>) => (
  <div
    className={cn(
      "flex flex-col space-y-2 text-center sm:text-left",
      "pb-4 border-b border-stroke/30",
      className
    )}
    {...props}
  />
)
DialogHeader.displayName = "DialogHeader"

const DialogFooter = ({
  className,
  ...props
}: React.HTMLAttributes<HTMLDivElement>) => (
  <div
    className={cn(
      "flex flex-col-reverse gap-2 sm:flex-row sm:justify-end",
      "pt-4 border-t border-stroke/30",
      className
    )}
    {...props}
  />
)
DialogFooter.displayName = "DialogFooter"

const DialogTitle = React.forwardRef<
  React.ComponentRef<typeof DialogPrimitive.Title>,
  React.ComponentPropsWithoutRef<typeof DialogPrimitive.Title>
>(({ className, ...props }, ref) => (
  <DialogPrimitive.Title
    ref={ref}
    className={cn(
      "font-display text-xl font-semibold leading-tight tracking-wide",
      // 金色渐变文字
      "bg-gradient-to-r from-[hsl(42,95%,45%)] via-[hsl(42,95%,55%)] to-[hsl(42,95%,45%)]",
      "bg-clip-text text-transparent",
      // 文字阴影增强可读性
      "drop-shadow-sm",
      className
    )}
    {...props}
  />
))
DialogTitle.displayName = DialogPrimitive.Title.displayName

const DialogDescription = React.forwardRef<
  React.ComponentRef<typeof DialogPrimitive.Description>,
  React.ComponentPropsWithoutRef<typeof DialogPrimitive.Description>
>(({ className, ...props }, ref) => (
  <DialogPrimitive.Description
    ref={ref}
    className={cn(
      "text-sm text-fg-secondary leading-relaxed",
      className
    )}
    {...props}
  />
))
DialogDescription.displayName = DialogPrimitive.Description.displayName

export {
  Dialog,
  DialogPortal,
  DialogOverlay,
  DialogClose,
  DialogTrigger,
  DialogContent,
  DialogHeader,
  DialogFooter,
  DialogTitle,
  DialogDescription,
}
