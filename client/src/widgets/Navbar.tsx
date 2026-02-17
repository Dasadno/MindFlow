/**
 * Navbar - Навигационная панель для главной страницы
 * 
 * Компонент отображает:
 * - Логотип MindFlow
 * - Навигационные ссылки (Кто мы?, Миссия, GitHub)
 * - Кнопку "Запустить поток"
 * 
 * Особенности:
 * - Фиксированная позиция вверху страницы
 * - Backdrop blur эффект
 * - Адаптивная верстка (навигация скрывается на мобилке)
 * - Градиентные hover эффекты на ссылках
 */

export const Navbar = () => {
    return (
        <nav className="fixed top-6 left-1/2 -translate-x-1/2 w-[90%] max-w-5xl z-50 backdrop-blur-xl bg-deep-midnight/40 border border-white/10 rounded-2xl px-8 py-4 flex justify-between items-center shadow-2xl">
            {/* ЛОГОТИП */}
            <div className="flex items-center gap-4 cursor-pointer select-none group">
                <img
                    src="/cover2.png"
                    alt="MindFlow Logo"
                    className="w-10 h-10 rounded-full object-cover border border-white/10 shadow-[0_0_10px_rgba(38,208,206,0.5)] group-hover:shadow-[0_0_20px_rgba(38,208,206,0.8)] transition-all duration-300"
                />

                <div className="text-xl font-bold bg-gradient-accent bg-clip-text text-transparent tracking-tighter uppercase text-white">
                    MindFlow
                </div>
            </div>

            {/* НАВИГАЦИОННЫЕ ССЫЛКИ (скрываются на мобилке) */}
            <div className="hidden md:flex space-x-1">
                {/* Ссылка: Кто мы? */}
                <a href="#about" className="group relative px-4 py-2 overflow-hidden rounded-lg transition-all hover:bg-white/5">
                    <span className="relative z-10 transition-all duration-300 group-hover:text-bright-turquoise group-hover:drop-shadow-[0_0_10px_rgba(38,208,206,0.8)]">
                        Кто мы?
                    </span>
                    {/* Градиентная линия снизу */}
                    <span className="absolute bottom-0 left-0 w-full h-[2px] bg-gradient-to-r from-transparent via-bright-turquoise to-transparent transform scale-x-0 transition-transform duration-300 group-hover:scale-x-100" />
                </a>

                {/* Ссылка: Миссия */}
                <a href="#mission" className="group relative px-4 py-2 overflow-hidden rounded-lg transition-all hover:bg-white/5">
                    <span className="relative z-10 transition-all duration-300 group-hover:text-bright-turquoise group-hover:drop-shadow-[0_0_10px_rgba(38,208,206,0.8)]">
                        Миссия
                    </span>
                    <span className="absolute bottom-0 left-0 w-full h-[2px] bg-gradient-to-r from-transparent via-bright-turquoise to-transparent transform scale-x-0 transition-transform duration-300 group-hover:scale-x-100" />
                </a>

                {/* Ссылка: GitHub */}
                <a href="https://github.com/Dasadno/Milk-IslandAI" target="_blank" rel="noopener noreferrer" className="group relative px-4 py-2 overflow-hidden rounded-lg transition-all hover:bg-white/5">
                    <span className="relative z-10 transition-all duration-300 group-hover:text-bright-turquoise group-hover:drop-shadow-[0_0_10px_rgba(38,208,206,0.8)]">
                        GitHub
                    </span>
                    <span className="absolute bottom-0 left-0 w-full h-[2px] bg-gradient-to-r from-transparent via-bright-turquoise to-transparent transform scale-x-0 transition-transform duration-300 group-hover:scale-x-100" />
                </a>
            </div>

            {/* КНОПКА CTA */}
            <button className="bg-gradient-accent text-white px-6 py-2.5 rounded-xl text-sm font-black hover:shadow-[0_0_25px_rgba(122,248,196,0.6)] hover:scale-105 transition-all active:scale-95 border border-white/20">
                Запустить поток
            </button>
        </nav>
    );
};