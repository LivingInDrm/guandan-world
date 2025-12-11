import tailwindcssAnimate from 'tailwindcss-animate';

/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
    "./**/*.stories.tsx",
  ],
  theme: {
    extend: {
      colors: {
        background: 'hsl(var(--background))',
        foreground: 'hsl(var(--foreground))',
        card: {
          DEFAULT: 'hsl(var(--card))',
          foreground: 'hsl(var(--card-foreground))',
        },
        popover: {
          DEFAULT: 'hsl(var(--popover))',
          foreground: 'hsl(var(--popover-foreground))',
        },
        primary: {
          DEFAULT: 'hsl(var(--primary))',
          foreground: 'hsl(var(--primary-foreground))',
        },
        secondary: {
          DEFAULT: 'hsl(var(--secondary))',
          foreground: 'hsl(var(--secondary-foreground))',
        },
        muted: {
          DEFAULT: 'hsl(var(--muted))',
          foreground: 'hsl(var(--muted-foreground))',
        },
        destructive: {
          DEFAULT: 'hsl(var(--destructive))',
          foreground: 'hsl(var(--destructive-foreground))',
        },
        border: 'hsl(var(--border))',
        input: 'hsl(var(--input))',
        ring: 'hsl(var(--ring))',
        accent: {
          DEFAULT: 'hsl(var(--accent))',
          foreground: 'hsl(var(--accent-foreground))',
          light: 'rgba(255, 193, 7, 0.15)',
          glow: 'rgba(255, 193, 7, 0.6)',
        },
        table: {
          50: 'hsl(var(--table-50))',
          100: 'hsl(var(--table-100))',
          200: 'hsl(var(--table-200))',
          300: 'hsl(var(--table-300))',
          400: 'hsl(var(--table-400))',
          800: 'hsl(var(--table-800))',
          900: 'hsl(var(--table-900))',
        },
        'team-us': 'hsl(var(--team-us))',
        'team-us-dark': 'hsl(var(--team-us-dark))',
        'team-them': 'hsl(var(--team-them))',
        'team-them-dark': 'hsl(var(--team-them-dark))',
        suit: {
          red: 'hsl(var(--suit-red))',
          black: 'hsl(var(--suit-black))',
        },
        badge: {
          level: 'hsl(var(--badge-level))',
          wild: 'hsl(var(--badge-wild))',
        },
        disabled: {
          bg: 'rgba(148, 163, 184, 0.5)',
          text: '#94a3b8',
          border: 'rgba(148, 163, 184, 0.3)',
        },
        'ds-action': {
          primary: 'hsl(var(--ds-color-action-primary))',
          secondary: 'hsl(var(--ds-color-action-secondary))',
          neutral: 'hsl(var(--ds-color-action-neutral))',
          danger: 'hsl(var(--ds-color-action-danger))',
        },
        'ds-surface': {
          base: 'hsl(var(--ds-color-surface-base))',
          elevated: 'hsl(var(--ds-color-surface-elevated))',
          emphasis: 'hsl(var(--ds-color-surface-emphasis))',
        },
        'ds-team': {
          us: 'hsl(var(--ds-color-team-us))',
          them: 'hsl(var(--ds-color-team-them))',
        },
        'ds-state': {
          active: 'hsl(var(--ds-color-state-active))',
          disabled: 'hsl(var(--ds-color-state-disabled))',
        },
        'ds-text': {
          primary: 'hsl(var(--ds-color-text-primary))',
          secondary: 'hsl(var(--ds-color-text-secondary))',
          inverse: 'hsl(var(--ds-color-text-inverse))',
        },
        'ds-border': {
          DEFAULT: 'hsl(var(--ds-color-border-base))',
          emphasis: 'hsl(var(--ds-color-border-emphasis))',
        },
        'ds-error': 'hsl(var(--ds-color-error))',
      },
      backgroundImage: {
        'gradient-primary': 'var(--gradient-primary)',
        'gradient-secondary': 'var(--gradient-secondary)',
        'gradient-warning': 'var(--gradient-warning)',
        'gradient-danger': 'var(--gradient-danger)',
      },
      borderRadius: {
        sm: 'calc(var(--radius) - 4px)',
        md: 'calc(var(--radius) - 2px)',
        lg: 'var(--radius)',
        xl: 'calc(var(--radius) + 4px)',
        '2xl': 'calc(var(--radius) + 8px)',
        'ds-sm': 'var(--ds-radius-sm)',
        'ds-md': 'var(--ds-radius-md)',
        'ds-lg': 'var(--ds-radius-lg)',
        'ds-full': 'var(--ds-radius-full)',
      },
      boxShadow: {
        card: '0 2px 4px rgba(0,0,0,0.1)',
        'card-hover': '0 4px 12px rgba(0,0,0,0.15)',
        'card-3d':
          '0 2px 3px rgba(0,0,0,0.12), 0 4px 8px rgba(0,0,0,0.08), inset 0 1px 0 rgba(255,255,255,0.6)',
        panel: '0 2px 8px rgba(0,0,0,0.15)',
        modal: '0 8px 32px rgba(0,0,0,0.25)',
        'glow-accent': '0 0 10px rgba(255, 193, 7, 0.6)',
        'inset-soft': 'inset 0 0 12px rgba(0,0,0,0.08)',
        input: '0 1px 2px rgba(0,0,0,0.05)',
        'input-focus': '0 0 0 3px rgba(16, 185, 129, 0.2)',
        'ds-elevation-0': 'var(--ds-elevation-0)',
        'ds-elevation-1': 'var(--ds-elevation-1)',
        'ds-elevation-2': 'var(--ds-elevation-2)',
        'ds-elevation-3': 'var(--ds-elevation-3)',
        'ds-relief': 'var(--ds-shadow-relief)',
        'ds-glow-sm': 'var(--ds-shadow-glow-sm)',
        'ds-glow-md': 'var(--ds-shadow-glow-md)',
        'ds-glow-lg': 'var(--ds-shadow-glow-lg)',
      },
      spacing: {
        4.5: '1.125rem',
        18: '4.5rem',
        30: '7.5rem',
        'ds-18': 'var(--ds-spacing-18)',
        'ds-22': 'var(--ds-spacing-22)',
        'ds-30': 'var(--ds-spacing-30)',
      },
      transitionDuration: {
        'ds-fast': 'var(--ds-duration-fast)',
        'ds-normal': 'var(--ds-duration-normal)',
        'ds-slow': 'var(--ds-duration-slow)',
      },
      transitionTimingFunction: {
        'ds-smooth': 'var(--ds-ease-smooth)',
        'ds-bounce': 'var(--ds-ease-bounce)',
        'ds-snap': 'var(--ds-ease-snap)',
      },
      animation: {
        'card-enter': 'card-enter 0.3s cubic-bezier(0.34, 1.56, 0.64, 1)',
        'card-fly': 'card-fly 0.5s ease-out',
        'glow-pulse': 'glow-pulse 2s ease-in-out infinite',
        'ds-pulse-ring': 'ds-pulse-ring 2s cubic-bezier(0.4, 0, 0.6, 1) infinite',
        'ds-bounce-select': 'ds-bounce-select 0.4s cubic-bezier(0.34, 1.56, 0.64, 1)',
      },
      keyframes: {
        'card-enter': {
          '0%': { opacity: '0', transform: 'scale(0.8) translateY(10px)' },
          '100%': { opacity: '1', transform: 'scale(1) translateY(0)' },
        },
        'card-fly': {
          '0%': { opacity: '0.8' },
          '100%': { opacity: '1' },
        },
        'glow-pulse': {
          '0%, 100%': { boxShadow: '0 0 10px rgba(255, 193, 7, 0.4)' },
          '50%': { boxShadow: '0 0 20px rgba(255, 193, 7, 0.8)' },
        },
        'ds-pulse-ring': {
          '0%': { boxShadow: '0 0 0 0 rgba(250, 204, 21, 0.7)' },
          '70%': { boxShadow: '0 0 0 12px rgba(250, 204, 21, 0)' },
          '100%': { boxShadow: '0 0 0 0 rgba(250, 204, 21, 0)' },
        },
        'ds-bounce-select': {
          '0%': { transform: 'translateY(0) scale(1)' },
          '50%': { transform: 'translateY(-20px) scale(1.05)' },
          '100%': { transform: 'translateY(-12px) scale(1.05)' },
        },
      },
    },
  },
  plugins: [tailwindcssAnimate],
}