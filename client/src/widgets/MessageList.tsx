import { useRef, useEffect } from 'react';
import type { Event } from '../shared/types';

interface MessageListProps {
    messages: Event[];
}

export const MessageList = ({ messages }: MessageListProps) => {
    const bottomRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        bottomRef.current?.scrollIntoView({ behavior: 'smooth' });
    }, [messages]);

    return (
        <div className="
            /* Контейнер для сообщений */
            flex-1
            overflow-y-auto
            p-4 md:p-6
            space-y-4
            /* Кастомный скроллбар */
            scrollbar-thin scrollbar-thumb-bright-turquoise/30 scrollbar-track-transparent
        ">
            {messages.length === 0 && (
                <div className="text-center text-text-secondary mt-10">
                    No messages yet. Waiting for activity...
                </div>
            )}

            {messages.map((message, index) => {
                // Determine if sender is 'user' or 'ai' based on message source/speaker
                const isUser = message.source === 'user' || message.speaker === 'User';
                const senderName = message.speaker || 'Unknown';
                const content = message.content || '';
                const timestamp = message.created_at ? new Date(message.created_at).toLocaleTimeString() : '';

                if (message.type === 'system') {
                    return (
                        <div key={index} className="text-center text-xs text-text-secondary/50 my-2">
                            {content}
                        </div>
                    );
                }

                return (
                    <div
                        key={index}
                        className={`
                        /* Выравнивание: AI слева, пользователь справа */
                        flex
                        ${isUser ? 'justify-end' : 'justify-start'}
                    `}
                    >
                        {/* КОНТЕЙНЕР СООБЩЕНИЯ */}
                        <div className={`
                        /* Максимальная ширина сообщения */
                        max-w-[85%] md:max-w-[70%]
                        flex flex-col
                        ${isUser ? 'items-end' : 'items-start'}
                    `}>
                            {/* ИМЯ ОТПРАВИТЕЛЯ */}
                            <div className="
                            flex items-center gap-2 mb-1
                            px-2
                        ">
                                {/* Индикатор онлайн (только для AI) */}
                                {!isUser && (
                                    <div className="w-2 h-2 bg-light-mint rounded-full" />
                                )}

                                <span className="text-text-secondary text-xs font-medium">
                                    {senderName}
                                </span>

                                <span className="text-text-secondary/50 text-xs">
                                    {timestamp}
                                </span>
                            </div>

                            {/* ТЕЛО СООБЩЕНИЯ */}
                            <div className={`
                            px-4 py-3
                            rounded-2xl
                            shadow-md
                            
                            /* Разные стили для AI и пользователя */
                            ${!isUser
                                    ? 'bg-dark-ocean text-text-primary rounded-tl-none'
                                    : 'bg-gradient-primary text-white rounded-tr-none'
                                }
                        `}>
                                <p className="text-sm md:text-base leading-relaxed whitespace-pre-wrap">
                                    {content}
                                </p>
                            </div>
                        </div>
                    </div>
                )
            })}
            <div ref={bottomRef} />
        </div>
    );
};
