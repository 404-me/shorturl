/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './*.html',       // 匹配根目录下的所有 HTML 文件
    './**/*.html',    // 匹配所有子目录下的 HTML 文件
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}