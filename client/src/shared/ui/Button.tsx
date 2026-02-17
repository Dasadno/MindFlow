interface ButtonProps {
    children: React.ReactNode;
    onClick?: () => void;
    variant?: 'primary' | 'secondary' | 'gradient';
    className?: string;
}

export const Button = ({ children, onClick, variant = 'primary', className = '' }: ButtonProps) => {
    const baseStyles = 'px-6 py-3 rounded-lg font-semibold transition-all focus:outline-none focus:ring-2 focus:ring-bright-turquoise focus:ring-offset-2 focus:ring-offset-deep-midnight';

    const variants = {
        primary: 'bg-bright-turquoise hover:bg-sky-blue text-white shadow-md hover:shadow-lg',
        secondary: 'bg-dark-ocean hover:bg-dark-ocean/80 text-text-primary border border-bright-turquoise/30',
        gradient: 'bg-gradient-primary text-white shadow-lg hover:shadow-xl',
    };

    return (
        <button
            onClick={onClick}
            className={`${baseStyles} ${variants[variant]} ${className}`}
        >
            {children}
        </button>
    );
};
