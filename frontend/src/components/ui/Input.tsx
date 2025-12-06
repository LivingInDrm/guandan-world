import React from 'react';

export interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
  helperText?: string;
  fullWidth?: boolean;
}

export const Input = React.forwardRef<HTMLInputElement, InputProps>(
  ({ label, error, helperText, fullWidth = true, className = '', id, ...props }, ref) => {
    const inputId = id || (label ? label.toLowerCase().replace(/\s+/g, '-') : undefined);

    const baseStyles = `
      px-3 py-2 
      rounded-card border shadow-input
      transition-all duration-200
      focus:outline-none focus:border-focus-border focus:shadow-input-focus
    `;

    const stateStyles = error
      ? 'border-red-500'
      : 'border-table-300';

    const disabledStyles = props.disabled
      ? 'bg-disabled-bg text-disabled-text cursor-not-allowed'
      : 'bg-white';

    const computedClassName = [
      baseStyles,
      stateStyles,
      disabledStyles,
      fullWidth ? 'w-full' : '',
      className,
    ]
      .join(' ')
      .replace(/\s+/g, ' ')
      .trim();

    return (
      <div className={fullWidth ? 'w-full' : ''}>
        {label && (
          <label
            htmlFor={inputId}
            className="block text-sm font-medium text-table-400 mb-1"
          >
            {label}
          </label>
        )}
        <input
          ref={ref}
          id={inputId}
          className={computedClassName}
          {...props}
        />
        {error && (
          <p className="mt-1 text-sm text-red-500">{error}</p>
        )}
        {!error && helperText && (
          <p className="mt-1 text-xs text-table-300">{helperText}</p>
        )}
      </div>
    );
  }
);

Input.displayName = 'Input';
