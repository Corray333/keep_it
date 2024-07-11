/** @type {import('tailwindcss').Config} */
export default {
  content: [],
  darkMode: false,
  purge: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors:{
        dark: "var(--dark)",
        black: "var(--black)",
        accent: "var(--accent)",
        semidark: "var(--semi-dark)",
      }
    },
  },
  plugins: [],
}

