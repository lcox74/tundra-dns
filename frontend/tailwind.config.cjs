/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          full: "#323437",
          faint: "rgba(50, 52, 55, 0.1)"
        },
        secondary: "#F6F7F9",
        accent: "#0180C9",
        success: {
          full: "#21A54F",
          faint: "rgba(33, 165, 79, 0.1)"
        },
        failure: "#B9170B",

        a: "#B9170B",
        mx: "#0180C9",
        cname: "#21A54F",
        soa: "#FF9933",
        txt: "#66CCCC",
        ns: "#008B8B"
      },
      textColor: {
        primary: "#323437",
        success: "#21A54F",
        faint: "rgba(50, 52, 55, 0.8)"
      },
      fontFamily: {
        sans: ["Inter", "sans-serif"],
        mono: ["Fira Code", "monospace"]
      },
    },
  },
  plugins: [],
}

