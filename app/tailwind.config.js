/** @type {import('tailwindcss').Config} */
export default {
  content: [],
  purge: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors:{
        'primary-bg':'var(--primary-bg)',
        'secondary-bg':'var(--secondary-bg)',
        'primary':'var(--primary)',
        'invert-bg':'var(--invert-bg)',
        'invert-bg-50':'var(--invert-bg-50)'
      }
    },
  },
  plugins: [],
}

