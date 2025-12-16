import * as React from "react"
import * as LabelPrimitive from "@radix-ui/react-label"
import { cva, type VariantProps } from "class-variance-authority"
import { cn } from "@/lib/utils"

const inputVariants = cva(
  [
    // 基础布局
    "flex w-full",
    "text-fg-primary placeholder:text-fg-secondary/60",
    // 形状
    "rounded-md",
    // 背景与边框
    "bg-surface-elevated border border-stroke",
    // 阴影 - 内凹效果
    "shadow-[inset_0_1px_3px_rgba(0,0,0,0.06),var(--shadow-relief)]",
    // 过渡动画
    "transition-all duration-normal ease-smooth",
    // 聚焦状态 - 优雅高亮
    "focus-visible:outline-none",
    "focus-visible:border-state-active",
    "focus-visible:shadow-[inset_0_1px_3px_rgba(0,0,0,0.06),var(--shadow-relief),var(--shadow-glow-sm)]",
    "focus-visible:bg-surface-base",
    // 禁用状态
    "disabled:cursor-not-allowed disabled:bg-state-disabled/50",
    "disabled:text-fg-secondary disabled:border-stroke/50",
    // 文件输入特殊处理
    "file:border-0 file:bg-transparent file:text-sm file:font-medium",
  ],
  {
    variants: {
      inputSize: {
        sm: "h-9 px-3 py-1.5 text-sm",
        md: "h-11 px-4 py-2.5 text-base",
        lg: "h-13 px-5 py-3 text-lg",
      },
    },
    defaultVariants: {
      inputSize: "md",
    },
  }
)

export interface InputProps
  extends Omit<React.InputHTMLAttributes<HTMLInputElement>, 'size'>,
    VariantProps<typeof inputVariants> {
  error?: boolean
  icon?: React.ReactNode
  iconPosition?: 'left' | 'right'
}

const Input = React.forwardRef<HTMLInputElement, InputProps>(
  ({ className, type, inputSize, error, icon, iconPosition = 'left', ...props }, ref) => {
    if (icon) {
      return (
        <div className="relative">
          {iconPosition === 'left' && (
            <span className="absolute left-3 top-1/2 -translate-y-1/2 text-fg-secondary pointer-events-none">
              {icon}
            </span>
          )}
          <input
            type={type}
            className={cn(
              inputVariants({ inputSize }),
              iconPosition === 'left' && "pl-10",
              iconPosition === 'right' && "pr-10",
              error && [
                "border-action-danger",
                "focus-visible:border-action-danger",
                "focus-visible:shadow-[inset_0_1px_3px_rgba(0,0,0,0.06),var(--shadow-relief),var(--shadow-ruby)]",
              ],
              className
            )}
            ref={ref}
            {...props}
          />
          {iconPosition === 'right' && (
            <span className="absolute right-3 top-1/2 -translate-y-1/2 text-fg-secondary pointer-events-none">
              {icon}
            </span>
          )}
        </div>
      )
    }

    return (
      <input
        type={type}
        className={cn(
          inputVariants({ inputSize }),
          error && [
            "border-action-danger",
            "focus-visible:border-action-danger",
            "focus-visible:shadow-[inset_0_1px_3px_rgba(0,0,0,0.06),var(--shadow-relief),var(--shadow-ruby)]",
          ],
          className
        )}
        ref={ref}
        {...props}
      />
    )
  }
)
Input.displayName = "Input"

// 标签组件
const Label = React.forwardRef<
  React.ComponentRef<typeof LabelPrimitive.Root>,
  React.ComponentPropsWithoutRef<typeof LabelPrimitive.Root> & {
    required?: boolean
  }
>(({ className, required, children, ...props }, ref) => (
  <LabelPrimitive.Root
    ref={ref}
    className={cn(
      "inline-flex items-center gap-1",
      "text-sm font-medium leading-none text-fg-primary",
      "mb-2",
      "peer-disabled:cursor-not-allowed peer-disabled:opacity-70",
      className
    )}
    {...props}
  >
    {children}
    {required && (
      <span className="text-action-danger text-xs">*</span>
    )}
  </LabelPrimitive.Root>
))
Label.displayName = LabelPrimitive.Root.displayName

// Textarea 组件
const Textarea = React.forwardRef<
  HTMLTextAreaElement,
  React.TextareaHTMLAttributes<HTMLTextAreaElement> & {
    error?: boolean
  }
>(({ className, error, ...props }, ref) => (
  <textarea
    className={cn(
      // 基础布局
      "flex min-h-[100px] w-full resize-y",
      "text-fg-primary placeholder:text-fg-secondary/60",
      // 形状
      "rounded-md",
      // 背景与边框
      "bg-surface-elevated border border-stroke",
      // 阴影
      "shadow-[inset_0_1px_3px_rgba(0,0,0,0.06),var(--shadow-relief)]",
      // 内边距
      "px-4 py-3 text-base",
      // 过渡动画
      "transition-all duration-normal ease-smooth",
      // 聚焦状态
      "focus-visible:outline-none",
      "focus-visible:border-state-active",
      "focus-visible:shadow-[inset_0_1px_3px_rgba(0,0,0,0.06),var(--shadow-relief),var(--shadow-glow-sm)]",
      "focus-visible:bg-surface-base",
      // 禁用状态
      "disabled:cursor-not-allowed disabled:bg-state-disabled/50",
      "disabled:text-fg-secondary disabled:border-stroke/50",
      // 错误状态
      error && [
        "border-action-danger",
        "focus-visible:border-action-danger",
        "focus-visible:shadow-[inset_0_1px_3px_rgba(0,0,0,0.06),var(--shadow-relief),var(--shadow-ruby)]",
      ],
      className
    )}
    ref={ref}
    {...props}
  />
))
Textarea.displayName = "Textarea"

// 输入框组容器
const InputGroup = React.forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement>
>(({ className, ...props }, ref) => (
  <div
    ref={ref}
    className={cn("space-y-1.5", className)}
    {...props}
  />
))
InputGroup.displayName = "InputGroup"

// 帮助文本
const InputHelper = React.forwardRef<
  HTMLParagraphElement,
  React.HTMLAttributes<HTMLParagraphElement> & {
    error?: boolean
  }
>(({ className, error, ...props }, ref) => (
  <p
    ref={ref}
    className={cn(
      "text-xs mt-1.5",
      error ? "text-action-danger" : "text-fg-secondary",
      className
    )}
    {...props}
  />
))
InputHelper.displayName = "InputHelper"

export { Input, Label, Textarea, InputGroup, InputHelper, inputVariants }
