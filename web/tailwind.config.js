/** @type {import('tailwindcss').Config} */
export default {
  content: ["./src/**/*.{html,js,svelte,ts}"],
  theme: {
    extend: {
      screens: {
        "min-sm": "340px",
      },
    },
  },
  plugins: [],
};
