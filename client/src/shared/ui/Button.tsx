interface ButtonProps {
    children: React.ReactNode;
    onClick?: () => void;
    variant?: 'primary' | 'secondary' | 'gradient' | 'accent';
    className?: string;
    type?: 'button' | 'submit' | 'reset';
}

export const Button = ({ children, onClick, variant = 'primary', className = '', type = 'button' }: ButtonProps) => {
    // Базовые стили: убрал outline, добавил scale при клике, увеличил скругление
    const baseStyles = 'px-8 py-3 rounded-xl font-bold text-sm tracking-wide transition-all duration-300 focus:outline-none active:scale-95 hover:-translate-y-0.5 disabled:opacity-50 disabled:cursor-not-allowed';

    const variants = {
        // Бирюзовый (основной)
        primary: 'bg-bright-turquoise text-deep-midnight shadow-[0_0_20px_rgba(38,208,206,0.3)] hover:shadow-[0_0_30px_rgba(38,208,206,0.5)] border border-transparent',

        // Стеклянный (второстепенный) - идеально для фона deep-midnight
        secondary: 'bg-white/5 backdrop-blur-md text-white border border-white/10 hover:bg-white/10 hover:border-bright-turquoise/30 hover:text-bright-turquoise',

        // Градиент (как на главной)
        gradient: 'bg-gradient-primary text-white shadow-[0_10px_30px_rgba(38,208,206,0.3)] hover:shadow-[0_15px_40px_rgba(38,208,206,0.5)] border border-white/20',

        // Акцентный (Мятный - для важных действий)
        accent: 'bg-gradient-accent text-white shadow-[0_10px_30px_rgba(122,248,196,0.3)] hover:shadow-[0_15px_40px_rgba(122,248,196,0.5)] border border-white/20'
    };

    return (
        <button
            type={type}
            onClick={onClick}
            className={`${baseStyles} ${variants[variant]} ${className}`}
        >
            {children}
        </button>
    );
};