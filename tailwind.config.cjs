/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./index.html",
    "./src/**/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      fontFamily: {
        'brand': ['Inter', 'sans-serif'], // Назови как хочешь
      },
    },
  },
  plugins: [],
}
