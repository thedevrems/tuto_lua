/** @type {import('tailwindcss').Config} */
// Design tokens mirror CHARTE-GRAPHIQUE.md (modern black & white design system).
export default {
  darkMode: 'class',
  content: ['./index.html', './src/**/*.{ts,tsx}'],
  theme: {
    extend: {
      colors: {
        // ---- Charte graphique (light, monochrome) ----
        black: { DEFAULT: '#0A0A0A', soft: '#1A1A1A', muted: '#2A2A2A' },
        white: { DEFAULT: '#FFFFFF', soft: '#FAFAFA' },
        gray: {
          100: '#F5F5F5', 200: '#E5E5E5', 300: '#D4D4D4', 400: '#A3A3A3',
          500: '#737373', 600: '#525252', 700: '#404040', 800: '#262626', 900: '#171717',
        },
        success: { DEFAULT: '#16A34A', bg: '#F0FDF4', border: '#BBF7D0' },
        danger: { DEFAULT: '#DC2626', bg: '#FEF2F2', border: '#FECACA' },
        warning: { DEFAULT: '#EA580C', bg: '#FFF7ED', border: '#FED7AA' },
        info: { DEFAULT: '#0891B2', bg: '#F0FDFA', border: '#CCFBF1' },

        // ---- Legacy dark scale (kept while the learning IDE migrates) ----
        ink: {
          950: '#0a0a0b', 900: '#101012', 850: '#161618', 800: '#1c1c20',
          700: '#26262b', 600: '#3a3a42', 500: '#5b5b66', 400: '#8b8b96',
          300: '#b4b4bd', 200: '#d6d6dc', 100: '#ededf0', 50: '#f8f8fa',
        },
      },
      fontFamily: {
        sans: ['Inter', 'ui-sans-serif', 'system-ui', '-apple-system', 'Segoe UI', 'sans-serif'],
        mono: ['"JetBrains Mono"', 'ui-monospace', 'SFMono-Regular', 'Menlo', 'Consolas', 'monospace'],
      },
      borderRadius: {
        sm: '0.25rem', md: '0.5rem', lg: '0.75rem', xl: '1rem', '2xl': '1.5rem', full: '9999px',
      },
      boxShadow: {
        xs: '0 1px 2px 0 rgba(0,0,0,0.05)',
        sm: '0 1px 3px 0 rgba(0,0,0,0.1), 0 1px 2px -1px rgba(0,0,0,0.1)',
        md: '0 4px 6px -1px rgba(0,0,0,0.1), 0 2px 4px -2px rgba(0,0,0,0.1)',
        lg: '0 10px 15px -3px rgba(0,0,0,0.1), 0 4px 6px -4px rgba(0,0,0,0.1)',
        xl: '0 20px 25px -5px rgba(0,0,0,0.1), 0 8px 10px -6px rgba(0,0,0,0.1)',
        '2xl': '0 25px 50px -12px rgba(0,0,0,0.25)',
      },
      maxWidth: { container: '1280px' },
      transitionTimingFunction: { smooth: 'cubic-bezier(0.4, 0, 0.2, 1)' },
      transitionDuration: { fast: '150ms', base: '250ms', slow: '350ms' },
      keyframes: {
        'fade-in': { '0%': { opacity: '0' }, '100%': { opacity: '1' } },
        'fade-up': { '0%': { opacity: '0', transform: 'translateY(12px)' }, '100%': { opacity: '1', transform: 'translateY(0)' } },
      },
      animation: {
        'fade-in': 'fade-in 250ms cubic-bezier(0.4, 0, 0.2, 1)',
        'fade-up': 'fade-up 350ms cubic-bezier(0.4, 0, 0.2, 1)',
      },
    },
  },
  plugins: [],
}
