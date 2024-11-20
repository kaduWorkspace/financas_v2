/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./resources/views/**/*.{html,tmpl,js}"],
  theme: {
    extend: {
        textShadow: {
            sm: '1px 1px 2px rgba(0, 0, 0, 0.5)',
            DEFAULT: '2px 2px 4px rgba(0, 0, 0, 0.7)',
            lg: '3px 3px 6px rgba(0, 0, 0, 0.8)',
        },
        colors: {
            preto: "#0e1110",
            pretoTransparente: 'rgba(5, 5, 5, 0.90)',
            pretoMaisTransparente: 'rgba(5, 5, 5, 0.50)',
            azul: "#00b8ff"
        },
        borderRadius: {
            app: "20px",
            appMajor: "40px",
            appMinor: "18px"
        },
    },
  },
  plugins: [],
}

