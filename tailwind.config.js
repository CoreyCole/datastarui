/*eslint-env node*/

/**
 * @type {import('tailwindcss').Config}
 */
module.exports = {
  darkMode: "class",
  content: [
    "./components/**/*.templ",
    "./pages/**/*.templ",
    "./layouts/**/*.templ",
  ],
}