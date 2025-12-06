import React from 'react';

export interface CardProps extends React.HTMLAttributes<HTMLDivElement> {
  variant?: 'default' | 'glass' | 'gradient';
  padding?: 'none' | 'sm' | 'md' | 'lg';
}

const variantStyles: Record<NonNullable<CardProps['variant']>, string> = {
  default: 'bg-white rounded-panel shadow-panel',
  glass: 'bg-white/70 backdrop-blur-sm rounded-panel shadow-panel border border-white/50',
  gradient: 'bg-gradient-to-b from-table-50 via-table-100 to-table-200 rounded-panel shadow-panel border border-white/30',
};

const paddingStyles: Record<NonNullable<CardProps['padding']>, string> = {
  none: '',
  sm: 'p-3',
  md: 'p-5',
  lg: 'p-8',
};

export const Card = React.forwardRef<HTMLDivElement, CardProps>(
  (
    {
      variant = 'default',
      padding = 'md',
      className = '',
      children,
      ...props
    },
    ref
  ) => {
    const computedClassName = [
      variantStyles[variant],
      paddingStyles[padding],
      className,
    ]
      .join(' ')
      .replace(/\s+/g, ' ')
      .trim();

    return (
      <div ref={ref} className={computedClassName} {...props}>
        {children}
      </div>
    );
  }
);

Card.displayName = 'Card';
