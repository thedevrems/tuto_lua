/** @type {import('tailwindcss').Config} */
export default {
  darkMode: 'class',
  content: ['./index.html', './src/**/*.{ts,tsx}'],
  theme: {
    extend: {
      colors: {
        // Monochrome scale used across the app
        ink: {
          950: '#0a0a0b',
          900: '#101012',
          850: '#161618',
          800: '#1c1c20',
          700: '#26262b',
          600: '#3a3a42',
          500: '#5b5b66',
          400: '#8b8b96',
          300: '#b4b4bd',
          200: '#d6d6dc',
          100: '#ededf0',
          50: '#f8f8fa',
        },
      },
      fontFamily: {
        mono: ['"JetBrains Mono"', 'ui-monospace', 'SFMono-Regular', 'Menlo', 'Consolas', 'monospace'],
        sans: ['Inter', 'ui-sans-serif', 'system-ui', 'sans-serif'],
      },
    },
  },
  plugins: [],
}
