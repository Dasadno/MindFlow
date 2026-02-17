import React from 'react';

interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
    label?: string;
    className?: string;
}

const Input = ({ label, className = '', ...props }: InputProps) => {
    return (
        <div className={`group relative w-full ${className}`}>
            {/* Стили для анимации внутренней линии (сканирования) */}
            <style>{`
                @keyframes scan-line {
                    0% { transform: translateX(-100%); }
                    100% { transform: translateX(100%); }
                }
                .group:focus-within .animate-scan {
                    animation: scan-line 3s linear infinite;
                }
            `}</style>

            {/* Верхний лейбл в стиле "системного лога" */}
            {label && (
                <label className="block text-[10px] uppercase tracking-[0.2em] text-bright-turquoise/60 mb-2 ml-1 font-mono group-focus-within:text-bright-turquoise transition-colors">
                    {`> ${label}`}
                </label>
            )}

            <div className="relative overflow-hidden rounded-2xl transition-all duration-500">
                {/* Фоновое свечение при фокусе */}
                <div className="absolute inset-0 bg-bright-turquoise/5 opacity-0 group-focus-within:opacity-100 transition-opacity duration-500" />
                
                {/* Анимированная линия границы (появляется при фокусе) */}
                <div className="absolute bottom-0 left-0 w-full h-[2px] bg-gradient-to-r from-transparent via-bright-turquoise to-transparent -translate-x-full group-focus-within:animate-scan opacity-50" />

                <input
                    className={`
                        w-full px-6 py-4 rounded-2xl font-medium outline-none
                        bg-white/5 border border-white/10
                        text-white placeholder:text-white/20
                        focus:border-bright-turquoise/40 focus:bg-white/[0.08]
                        backdrop-blur-md
                        transition-all duration-300
                        shadow-[0_4px_20px_rgba(0,0,0,0.1)]
                        group-hover:border-white/20
                    `}
                    {...props}
                />
            </div>

            {/* Декоративный элемент в углу */}
            <div className="absolute -right-1 -bottom-1 w-2 h-2 border-r border-b border-white/20 group-focus-within:border-bright-turquoise transition-colors" />
        </div>
    );
};

export default Input;