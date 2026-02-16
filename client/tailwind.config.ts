/** @type {import('tailwindcss').Config} */
export default {
    content: ['./index.html', './src/**/*.{js,ts,jsx,tsx}'],
    theme: {
        extend: {
            colors: {
                // 🔵 Глубокие синие (Фон и основные блоки)
                'deep-midnight': '#0B1E3B',
                'dark-ocean': '#1A3C5E',

                // 💎 Голубые и Циан (Акценты и активные элементы)
                'bright-turquoise': '#26D0CE',
                'sky-blue': '#5BC0EB',

                // 🟢 Зеленые и Мятные (Процессы и "Мысли" ИИ)
                'light-mint': '#7AF8C4',
                'soft-teal': '#50E3C2',

                // Дополнительные цвета для текста
                'text-primary': '#FFFFFF',
                'text-secondary': '#B0BEC5',
            },
            backgroundImage: {
                'gradient-primary': 'linear-gradient(135deg, #26D0CE 0%, #1A3C5E 100%)',
                'gradient-accent': 'linear-gradient(135deg, #7AF8C4 0%, #26D0CE 100%)',
            },
        },
    },
    plugins: [],
};
