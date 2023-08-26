/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        primary: "#323437",
        secondary: "#F6F7F9",
        accent: "#0180C9",
        success: {
          full: "#21A54F",
          faint: "rgba(33, 165, 79, 0.1)"
        },
        failure: "#B9170B"
      },
      textColor: {
        primary: "#323437",
        success: "#21A54F"
      }
    },
  },
  plugins: [],
}

