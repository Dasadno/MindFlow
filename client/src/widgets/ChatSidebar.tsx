import type { AgentSummary } from '../shared/types';

interface ChatSidebarProps {
    agents: AgentSummary[];
    selectedAgentId: string | null;
    onSelectAgent: (id: string) => void;
    className?: string; // Allow passing styles for layout
}

export const ChatSidebar = ({ agents, selectedAgentId, onSelectAgent, className }: ChatSidebarProps) => {
    return (
        <aside className={`
            /* Базовые стили сайдбара */
            w-full md:w-80 
            bg-dark-ocean 
            border-r border-bright-turquoise/20
            flex flex-col
            /* Адаптивность: на мобилке скрывается или показывается через меню */
            h-auto md:h-screen
            ${className || ''}
        `}>
            {/* HEADER САЙДБАРА */}
            <div className="
                p-4 
                border-b border-bright-turquoise/20
                bg-deep-midnight/50
            ">
                <h2 className="text-text-primary text-xl font-bold mb-1">
                    AI Agents
                </h2>
                <p className="text-text-secondary text-sm">
                    Select an agent to chat
                </p>
            </div>

            {/* СПИСОК АГЕНТОВ */}
            <div className="
                flex-1 
                overflow-y-auto
                /* Кастомный скроллбар */
                scrollbar-thin scrollbar-thumb-bright-turquoise/30 scrollbar-track-transparent
            ">
                <div
                    onClick={() => onSelectAgent('')}
                    className={`
                        p-4 cursor-pointer transition-all border-l-4 hover:bg-deep-midnight/50
                        ${!selectedAgentId ? 'border-bright-turquoise bg-deep-midnight/30' : 'border-transparent'}
                    `}
                >
                    <h3 className="text-text-primary font-semibold">Global Stream</h3>
                    <p className="text-text-secondary text-sm">View all interactions</p>
                </div>

                {agents.map((agent) => (
                    <div
                        key={agent.id}
                        onClick={() => onSelectAgent(agent.id)}
                        className={`
                            /* Базовые стили элемента списка */
                            p-4 
                            cursor-pointer
                            transition-all
                            border-l-4
                            hover:bg-deep-midnight/50
                            
                            /* Активный агент выделяется */
                            ${agent.id === selectedAgentId
                                ? 'border-bright-turquoise bg-deep-midnight/30'
                                : 'border-transparent'
                            }
                        `}
                    >
                        {/* ВЕРХНЯЯ ЧАСТЬ: Имя + Статус */}
                        <div className="flex items-center justify-between mb-2">
                            {/* Имя агента */}
                            <h3 className="text-text-primary font-semibold">
                                {agent.name}
                            </h3>

                            {/* Индикатор статуса */}
                            <div className="flex items-center gap-2">
                                <div className={`
                                    w-2 h-2 rounded-full
                                    ${agent.isActive
                                        ? 'bg-light-mint animate-pulse'
                                        : 'bg-text-secondary/50'
                                    }
                                `} />
                                <span className="text-text-secondary text-xs">
                                    {agent.isActive ? 'online' : 'offline'}
                                </span>
                            </div>
                        </div>

                        {/* НИЖНЯЯ ЧАСТЬ: Тип личности */}
                        <p className="text-text-secondary text-sm">
                            {agent.personalityType}
                        </p>
                    </div>
                ))}
            </div>

            {/* FOOTER САЙДБАРА (опционально) */}
            <div className="
                p-4 
                border-t border-bright-turquoise/20
                bg-deep-midnight/50
            ">
                <button className="
                    w-full
                    px-4 py-2
                    bg-gradient-primary
                    text-white
                    rounded-lg
                    font-semibold
                    hover:shadow-lg
                    transition-shadow
                ">
                    + New Agent
                </button>
            </div>
        </aside>
    );
};
