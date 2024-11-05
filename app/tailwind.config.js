/** @type {import('tailwindcss').Config} */
export default {
  content: [],
  purge: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors:{
        'primary-bg':'var(--dark)',
        'secondary-bg':'var(--half-dark)',
        'primary':'var(--primary)',
        'invert-bg':'var(--light)',
      }
    },
  },
  plugins: [],
}

