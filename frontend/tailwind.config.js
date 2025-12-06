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
        table: {
          50: '#EAF4EF',
          100: '#DDEEE5',
          200: '#D2E8DD',
          300: '#8CAA96',
          400: '#506E5A',
        },
        accent: {
          DEFAULT: '#FFC107',
          light: 'rgba(255, 193, 7, 0.15)',
          glow: 'rgba(255, 193, 7, 0.6)',
        },
        'team-us': '#525E6B',
        'team-us-dark': '#3E4854',
        'team-them': '#9E3737',
        'team-them-dark': '#7B2A2A',
        suit: {
          red: '#dc2626',
          black: '#374151',
        },
        badge: {
          level: '#22c55e',
          wild: '#ef4444',
        },
        btn: {
          primary: { from: '#10b981', to: '#059669' },
          secondary: { from: '#64748b', to: '#475569' },
          warning: { from: '#fbbf24', to: '#f59e0b' },
          danger: { from: '#ef4444', to: '#dc2626' },
        },
        focus: {
          ring: '#10b981',
          border: '#059669',
        },
        disabled: {
          bg: 'rgba(148, 163, 184, 0.5)',
          text: '#94a3b8',
          border: 'rgba(148, 163, 184, 0.3)',
        },
      },
      borderRadius: {
        card: '8px',
        panel: '12px',
        modal: '16px',
        board: '20px',
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
      },
      spacing: {
        4.5: '1.125rem',
        18: '4.5rem',
        30: '7.5rem',
      },
      animation: {
        'card-enter': 'card-enter 0.3s cubic-bezier(0.34, 1.56, 0.64, 1)',
        'card-fly': 'card-fly 0.5s ease-out',
        'glow-pulse': 'glow-pulse 2s ease-in-out infinite',
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
      },
      fontSize: {
        'card-sm': ['12px', '1'],
        'card-md': ['16px', '1'],
      },
    },
  },
  plugins: [],
}