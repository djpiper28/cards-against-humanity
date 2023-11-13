import type { Config } from "tailwindcss";

const config: Config = {
  content: ["./src/**/*.{js,jsx,ts,tsx}"],
  theme: {
    extend: {
      colors: {
        "error-colour": "#FF3333",
      },
    },
  },
  plugins: [],
  fontFamily: {
    sans: ["Roboto", "sans-serif"],
    serif: ["serif"],
  },
};

export default config;
