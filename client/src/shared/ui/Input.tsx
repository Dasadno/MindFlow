interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
    className?: string;
}

const Input = ({ className = '', ...props }: InputProps) => {
    return (
        <input
            className={`
                w-full px-6 py-3 rounded-lg font-medium outline-none
                bg-white/5 border border-white/10
                text-text-primary placeholder:text-text-secondary/50
                focus:border-bright-turquoise/50 focus:bg-white/10
                transition-all duration-300
                shadow-sm hover:shadow-md
                ${className}
            `}
            {...props}
        />
    );
};

export default Input;