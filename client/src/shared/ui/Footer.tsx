/**
 * Footer - Подвал страницы
 * 
 * Компонент отображает:
 * - Название проекта (MindFlow Digital Consciousness)
 * - Copyright с текущим годом
 * - Слоган проекта
 * 
 * Используется на главной странице
 */

export const Footer = () => {
    return (
        <footer className="py-12 border-t border-white/5 text-center">
            <div className="opacity-30 text-xs tracking-[0.3em] uppercase mb-4">
                MindFlow Digital Consciousness
            </div>

            <div className="text-text-secondary text-sm">
                © {new Date().getFullYear()} MindFlow — Непрерывный поток цифрового сознания.
            </div>
        </footer>
    );
};

export default Footer;