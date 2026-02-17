interface ButtonProps {
    children: React.ReactNode;
    onClick?: () => void;
}

const Button = ({ children }: ButtonProps) => {
    return (
        <button className="px-4 py-2 rounded-lg font-semibold shadow-lg hover:shadow-xl transition-shadow bg-gradient-to-r from-primary to-secondary text-white">
            {children}
        </button>
    );
};

export default Button;