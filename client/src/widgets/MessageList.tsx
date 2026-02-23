import { useRef, useEffect } from 'react';
import type { Event } from '../shared/types';

interface MessageListProps {
    messages: (Event & { color?: string })[]; 
}

export const MessageList = ({ messages }: MessageListProps) => {
    const bottomRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        bottomRef.current?.scrollIntoView({ behavior: 'smooth' });
    }, [messages]);

    return (
        <div className="flex-1 overflow-y-auto p-4 md:p-6 space-y-6 no-scrollbar">
            {messages.length === 0 && (
                <div className="flex flex-col items-center justify-center mt-20 opacity-30">
                    <div className="w-16 h-16 border-2 border-bright-turquoise/20 rounded-full animate-ping mb-4" />
                    <p className="text-white font-mono text-xs tracking-widest">AWAITING NEURAL LINK...</p>
                </div>
            )}

            {messages.map((message, index) => {
                const isUser = message.source === 'user' || message.speaker === 'User';
                const senderName = message.speaker || 'Unknown';
                const content = message.content || '';
                const timestamp = message.created_at ? new Date(message.created_at).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }) : '';
                
                const agentColor = message.color || '#26d0ce';

                if (message.type === 'system') {
                    return (
                        <div key={index} className="text-center">
                            <span className="px-3 py-1 bg-white/5 rounded-full text-[10px] text-white/30 uppercase tracking-[0.2em]">
                                {content}
                            </span>
                        </div>
                    );
                }

                return (
                    <div
                        key={index}
                        className={`flex w-full animate-slide-up ${isUser ? 'justify-end' : 'justify-start'}`}
                    >
                        <div className={`flex flex-col max-w-[85%] md:max-w-[75%] ${isUser ? 'items-end' : 'items-start'}`}>
                            
                            {/* Заголовок сообщения */}
                            <div className="flex items-center gap-2 mb-1.5 px-1">
                                {!isUser && (
                                    <div 
                                        className="w-1.5 h-1.5 rounded-full shadow-[0_0_8px_currentcolor]" 
                                        style={{ backgroundColor: agentColor, color: agentColor }}
                                    />
                                )}
                                <span 
                                    className="text-[11px] font-bold uppercase tracking-wider transition-colors"
                                    style={{ color: isUser ? '#ffffff' : agentColor }}
                                >
                                    {senderName}
                                </span>
                                <span className="text-[10px] text-white/20 font-mono">{timestamp}</span>
                            </div>

                            {/* Тело сообщения */}
                            <div 
                                className={`
                                    relative px-4 py-3 rounded-2xl text-sm md:text-base leading-relaxed
                                    transition-all duration-300
                                    ${isUser 
                                        ? 'bg-gradient-to-br from-bright-turquoise to-deep-midnight text-white rounded-tr-none border border-white/10' 
                                        : 'bg-white/[0.03] text-white/90 rounded-tl-none border-l-2'
                                    }
                                `}
                                style={!isUser ? { 
                                    borderLeftColor: agentColor,
                                    backgroundColor: `${agentColor}08`, // Добавляем 5% прозрачности цвета к фону (HEX + 08)
                                    boxShadow: `inset 0 0 20px ${agentColor}05` // Легкое внутреннее свечение
                                } : {}}
                            >
                                <p className="whitespace-pre-wrap relative z-10">{content}</p>
                                
                                {/* Декоративный эффект для агентов */}
                                {!isUser && (
                                    <div 
                                        className="absolute inset-0 opacity-[0.03] pointer-events-none rounded-2xl"
                                        style={{ background: `linear-gradient(135deg, ${agentColor} 0%, transparent 100%)` }}
                                    />
                                )}
                            </div>
                        </div>
                    </div>
                );
            })}
            <div ref={bottomRef} className="h-4" />
        </div>
    );
};