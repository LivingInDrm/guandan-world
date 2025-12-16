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
      // 响应式断点 - 移动端横屏适配
      screens: {
        'xs': '480px',
        // 移动端横屏：高度受限 (iPhone 横屏等)
        'mobile-landscape': { 'raw': '(max-height: 500px) and (orientation: landscape)' },
        // 平板横屏：中等高度
        'tablet-landscape': { 'raw': '(min-width: 768px) and (max-height: 700px) and (orientation: landscape)' },
        // 基于高度的断点
        'short': { 'raw': '(max-height: 450px)' },
        'medium-height': { 'raw': '(min-height: 451px) and (max-height: 650px)' },
      },
      fontFamily: {
        sans: ['Noto Serif SC', 'Songti SC', 'STSong', 'Georgia', 'serif'],
        display: ['ZCOOL KuaiLe', 'Noto Serif SC', 'cursive'],
        mono: ['SF Mono', 'Fira Code', 'JetBrains Mono', 'monospace'],
      },
      colors: {
        action: {
          primary: 'hsl(var(--color-action-primary))',
          'primary-dark': 'hsl(var(--color-action-primary-dark))',
          'primary-darker': 'hsl(var(--color-action-primary-darker))',
          'primary-darkest': 'hsl(var(--color-action-primary-darkest))',
          secondary: 'hsl(var(--color-action-secondary))',
          tertiary: 'hsl(var(--color-action-tertiary))',
          neutral: 'hsl(var(--color-action-neutral))',
          danger: 'hsl(var(--color-action-danger))',
        },
        surface: {
          base: 'hsl(var(--color-surface-base))',
          elevated: 'hsl(var(--color-surface-elevated))',
          emphasis: 'hsl(var(--color-surface-emphasis))',
          muted: 'hsl(var(--color-surface-muted))',
          game: 'hsl(var(--color-surface-game))',
          'game-dark': 'hsl(var(--color-surface-game-dark))',
          'play-area': 'hsl(var(--color-surface-play-area))',
          'play-area-dark': 'hsl(var(--color-surface-play-area-dark))',
        },
        fg: {
          primary: 'hsl(var(--color-text-primary))',
          secondary: 'hsl(var(--color-text-secondary))',
          inverse: 'hsl(var(--color-text-inverse))',
        },
        state: {
          active: 'hsl(var(--color-state-active))',
          'active-dark': 'hsl(var(--color-state-active-dark))',
          'active-muted': 'hsl(var(--color-state-active-muted))',
          'active-light': 'hsl(var(--color-state-active-light))',
          'active-lighter': 'hsl(var(--color-state-active-lighter))',
          disabled: 'hsl(var(--color-state-disabled))',
        },
        stroke: {
          DEFAULT: 'hsl(var(--color-border-base))',
          emphasis: 'hsl(var(--color-border-emphasis))',
          game: 'hsl(var(--color-stroke-game))',
          'game-dark': 'hsl(var(--color-stroke-game-dark))',
          'play-area': 'hsl(var(--color-stroke-play-area))',
          header: 'hsl(var(--color-stroke-header))',
        },
        team: {
          us: 'hsl(var(--color-team-us))',
          them: 'hsl(var(--color-team-them))',
        },
        primitive: {
          neutral: {
            100: 'hsl(var(--primitive-neutral-100))',
            900: 'hsl(var(--primitive-neutral-900))',
          },
        },
        suit: {
          red: 'hsl(var(--suit-red))',
          black: 'hsl(var(--suit-black))',
        },
        badge: {
          level: 'hsl(var(--badge-level))',
          wild: 'hsl(var(--badge-wild))',
        },
        error: 'hsl(var(--color-error))',
      },
      borderRadius: {
        sm: 'var(--radius-sm)',
        md: 'var(--radius-md)',
        lg: 'var(--radius-lg)',
        full: 'var(--radius-full)',
      },
      boxShadow: {
        card: '0 2px 4px rgba(0,0,0,0.1)',
        'card-hover': '0 4px 12px rgba(0,0,0,0.15)',
        'card-3d': '0 2px 3px rgba(0,0,0,0.12), 0 4px 8px rgba(0,0,0,0.08), inset 0 1px 0 rgba(255,255,255,0.6)',
        'card-elevated': 'var(--elevation-3), var(--shadow-relief)',
        'card-interactive': 'var(--elevation-2), var(--shadow-relief)',
        panel: '0 2px 8px rgba(0,0,0,0.15)',
        modal: '0 8px 32px rgba(0,0,0,0.25)',
        'glow-accent': '0 0 10px rgba(255, 193, 7, 0.6)',
        'inset-soft': 'inset 0 0 12px rgba(0,0,0,0.08)',
        input: '0 1px 2px rgba(0,0,0,0.05)',
        'input-focus': '0 0 0 3px rgba(16, 185, 129, 0.2)',
        'elevation-0': 'var(--elevation-0)',
        'elevation-1': 'var(--elevation-1)',
        'elevation-2': 'var(--elevation-2)',
        'elevation-3': 'var(--elevation-3)',
        relief: 'var(--shadow-relief)',
        'glow-sm': 'var(--shadow-glow-sm)',
        'glow-md': 'var(--shadow-glow-md)',
        'glow-lg': 'var(--shadow-glow-lg)',
        jade: 'var(--shadow-jade)',
        ruby: 'var(--shadow-ruby)',
        'inner-light': 'var(--shadow-inner-light)',
        paper: 'var(--shadow-paper)',
      },
      backgroundImage: {
        'gradient-primary': 'linear-gradient(to right, hsl(var(--color-action-primary-dark)), hsl(var(--color-action-primary)))',
        'gradient-active': 'linear-gradient(to right, hsl(var(--color-state-active)), hsl(var(--color-state-active-dark)))',
        'gradient-table': 'linear-gradient(to bottom, hsl(var(--color-surface-game)), hsl(var(--color-surface-game-dark)))',
        'gradient-play-area': 'linear-gradient(to bottom, hsl(var(--color-surface-play-area)), hsl(var(--color-surface-play-area-dark)))',
        'gradient-header': 'linear-gradient(to right, hsl(var(--color-action-primary-dark)), hsl(var(--color-action-primary-darker)))',
        'gradient-header-bar': 'linear-gradient(to bottom, hsl(var(--color-action-primary-dark)), hsl(var(--color-action-primary-darkest)))',
        'gradient-gold': 'linear-gradient(to right, hsl(var(--color-state-active-lighter)), hsl(var(--color-state-active)), hsl(var(--color-state-active-lighter)))',
        'gradient-gold-text': 'linear-gradient(to right, hsl(var(--color-state-active-light)), hsl(var(--color-state-active)))',
      },
      spacing: {
        4.5: '1.125rem',
        18: 'var(--spacing-18)',
        22: 'var(--spacing-22)',
        30: 'var(--spacing-30)',
      },
      transitionDuration: {
        fast: 'var(--duration-fast)',
        normal: 'var(--duration-normal)',
        slow: 'var(--duration-slow)',
      },
      transitionTimingFunction: {
        smooth: 'var(--ease-smooth)',
        bounce: 'var(--ease-bounce)',
        snap: 'var(--ease-snap)',
      },
      animation: {
        'card-enter': 'card-enter 0.3s cubic-bezier(0.34, 1.56, 0.64, 1)',
        'card-fly': 'card-fly 0.5s ease-out',
        'glow-pulse': 'glow-pulse 2s ease-in-out infinite',
        'bounce-select': 'bounce-select 0.4s cubic-bezier(0.34, 1.56, 0.64, 1)',
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
        'bounce-select': {
          '0%': { transform: 'translateY(0) scale(1)' },
          '50%': { transform: 'translateY(-20px) scale(1.05)' },
          '100%': { transform: 'translateY(-12px) scale(1.05)' },
        },
      },
    },
  },
  plugins: [tailwindcssAnimate],
}
