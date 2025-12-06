import React from 'react';

export interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'warning' | 'danger' | 'ghost';
  size?: 'sm' | 'md' | 'lg';
  loading?: boolean;
  fullWidth?: boolean;
}

const variantStyles: Record<NonNullable<ButtonProps['variant']>, string> = {
  primary: `
    bg-gradient-to-b from-btn-primary-from to-btn-primary-to 
    text-white border-btn-primary-from/30 
    hover:brightness-110
    active:brightness-90
  `,
  secondary: `
    bg-gradient-to-b from-btn-secondary-from to-btn-secondary-to 
    text-white border-btn-secondary-from/30 
    hover:brightness-110
    active:brightness-90
  `,
  warning: `
    bg-gradient-to-b from-btn-warning-from to-btn-warning-to 
    text-white border-btn-warning-from/30 
    hover:brightness-110
    active:brightness-90
  `,
  danger: `
    bg-gradient-to-b from-btn-danger-from to-btn-danger-to 
    text-white border-btn-danger-from/30 
    hover:brightness-110
    active:brightness-90
  `,
  ghost: `
    bg-white/70 backdrop-blur-sm 
    text-table-400 border-table-300 
    hover:bg-white/90
    active:bg-white/60
  `,
};

const sizeStyles: Record<NonNullable<ButtonProps['size']>, string> = {
  sm: 'py-1.5 px-4 text-sm rounded-card',
  md: 'py-2 px-6 text-base rounded-panel',
  lg: 'py-3 px-8 text-lg rounded-panel',
};

const disabledStyles = 'bg-disabled-bg text-disabled-text border-disabled-border cursor-not-allowed';

export const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  (
    {
      variant = 'primary',
      size = 'md',
      loading = false,
      fullWidth = false,
      disabled,
      className = '',
      children,
      ...props
    },
    ref
  ) => {
    const isDisabled = disabled || loading;

    const baseStyles = `
      inline-flex items-center justify-center
      font-semibold border shadow-card
      transition-all duration-200
      hover:scale-105 hover:shadow-card-hover
      active:scale-95
      focus:outline-none focus:ring-2 focus:ring-focus-ring focus:ring-offset-2
    `;

    const computedClassName = [
      baseStyles,
      sizeStyles[size],
      isDisabled ? disabledStyles : variantStyles[variant],
      fullWidth ? 'w-full' : '',
      className,
    ]
      .join(' ')
      .replace(/\s+/g, ' ')
      .trim();

    return (
      <button
        ref={ref}
        disabled={isDisabled}
        className={computedClassName}
        {...props}
      >
        {loading && (
          <svg
            className="animate-spin -ml-1 mr-2 h-4 w-4"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
          >
            <circle
              className="opacity-25"
              cx="12"
              cy="12"
              r="10"
              stroke="currentColor"
              strokeWidth="4"
            />
            <path
              className="opacity-75"
              fill="currentColor"
              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
            />
          </svg>
        )}
        {children}
      </button>
    );
  }
);

Button.displayName = 'Button';
