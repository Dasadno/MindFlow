import type { LucideIcon } from 'lucide-react';


interface IconButtonProps {
    icon: LucideIcon;
    onClick?: () => void;
    className?: string;
}

export const IconButton = ({ icon: Icon, onClick, className = '' }: IconButtonProps) => {
    return (
        <button
            onClick={onClick}
            type="button"
            className={`
                relative group p-3 rounded-xl
                /* Базовый слой: стекло и глубокий синий */
                bg-white/5 backdrop-blur-md
                border border-white/10
                text-text-secondary
                
                /* Анимация и переходы */
                transition-all duration-300 ease-out
                
                /* Эффекты при наведении */
                hover:border-bright-turquoise/50
                hover:text-bright-turquoise
                hover:bg-dark-ocean/30
                hover:-translate-y-0.5
                
                /* Тень-свечение (как у твоих карточек) */
                hover:shadow-[0_0_20px_rgba(38,208,206,0.25)]
                
                /* Эффект при клике */
                active:scale-95 active:translate-y-0
                
                /* Дополнительные классы пользователя */
                ${className}
            `}
        >
            {/* Декоративный внутренний блик при наведении (как в Bento) */}
            <div className="absolute inset-0 rounded-xl bg-gradient-to-br from-bright-turquoise/10 to-transparent opacity-0 group-hover:opacity-100 transition-opacity" />

            <Icon className="relative z-10 w-5 h-5 transition-transform duration-300 group-hover:scale-110" />
        </button>
    );
};