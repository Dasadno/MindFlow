import { useState, type KeyboardEvent } from 'react';
import { Send } from 'lucide-react';

interface MessageInputProps {
    onSendMessage: (content: string) => void;
    disabled?: boolean;
}

export const MessageInput = ({ onSendMessage, disabled }: MessageInputProps) => {
    const [input, setInput] = useState('');

    const handleSend = (e?: React.FormEvent) => {
        e?.preventDefault();
        if (!input.trim() || disabled) return;

        onSendMessage(input);
        setInput('');
    };

    const handleKeyDown = (e: KeyboardEvent<HTMLTextAreaElement>) => {
        if (e.key === 'Enter' && !e.shiftKey) {
            e.preventDefault(); 
            handleSend();
        }
    };

    return (
        <div className="
            /* Контейнер для поля ввода */
            border-t border-bright-turquoise/20
            bg-dark-ocean
            p-4
            /* Sticky footer на мобилке */
            sticky bottom-0
        ">
            {/* ФОРМА ВВОДА */}
            <form className="
                flex items-end gap-3
                max-w-4xl mx-auto
            " onSubmit={handleSend}>
                {/* ТЕКСТОВОЕ ПОЛЕ */}
                <div className="flex-1">
                    <textarea
                        placeholder={disabled ? "Select an agent to chat..." : "Type your message..."}
                        rows={1}
                        value={input}
                        onChange={(e) => setInput(e.target.value)}
                        onKeyDown={handleKeyDown}
                        disabled={disabled}
                        className="
                            /* Базовые стили */
                            w-full
                            px-4 py-3
                            bg-deep-midnight
                            border border-bright-turquoise/30
                            rounded-xl
                            text-text-primary
                            placeholder:text-text-secondary/50
                            
                            /* Фокус */
                            focus:outline-none
                            focus:ring-2
                            focus:ring-bright-turquoise
                            focus:border-transparent
                            
                            /* Адаптивность */
                            text-sm md:text-base
                            
                            /* Автоматическое изменение высоты */
                            resize-none
                            min-h-[44px]
                            max-h-[120px]
                            
                            /* Скроллбар */
                            scrollbar-thin scrollbar-thumb-bright-turquoise/30 scrollbar-track-transparent

                            /* Disabled */
                            disabled:opacity-50
                            disabled:cursor-not-allowed
                        "
                    />
                </div>

                {/* КНОПКА ОТПРАВКИ */}
                <button
                    type="submit"
                    disabled={disabled || !input.trim()}
                    className="
                        /* Базовые стили */
                        px-4 py-3
                        bg-gradient-primary
                        text-white
                        rounded-xl
                        
                        /* Hover эффект */
                        hover:shadow-lg
                        
                        /* Фокус */
                        focus:outline-none
                        focus:ring-2
                        focus:ring-bright-turquoise
                        focus:ring-offset-2
                        focus:ring-offset-deep-midnight
                        
                        /* Transition */
                        transition-all
                        
                        /* Disabled состояние (когда поле пустое) */
                        disabled:opacity-50
                        disabled:cursor-not-allowed
                        
                        /* Адаптивность */
                        flex items-center justify-center
                        min-w-[44px]
                    "
                >
                    {/* Иконка отправки */}
                    <Send className="w-5 h-5" />

                    {/* Текст кнопки (скрывается на мобилке) */}
                    <span className="ml-2 hidden md:inline font-semibold">
                        Send
                    </span>
                </button>
            </form>
        </div>
    );
};
