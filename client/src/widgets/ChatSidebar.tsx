// import type { AgentSummary } from '../shared/types';
// import { Plus } from 'lucide-react';

// interface ChatSidebarProps {
//     agents: AgentSummary[];
//     selectedAgentId: string | null;
//     onSelectAgent: (id: string) => void;
//     className?: string; // Allow passing styles for layout
// }

// export const ChatSidebar = ({ agents, selectedAgentId, onSelectAgent, className }: ChatSidebarProps) => {
//     return (
//         <aside className={`
//             /* Базовые стили сайдбара */
//             w-full md:w-80 
//             bg-dark-ocean 
//             border-r border-bright-turquoise/20
//             flex flex-col
//             /* Адаптивность: на мобилке скрывается или показывается через меню */
//             h-auto md:h-screen
//             ${className || ''}
//         `}>
//             {/* HEADER САЙДБАРА */}
//             <div className="
//                 p-4 
//                 border-b border-bright-turquoise/20
//                 bg-deep-midnight/50
//             ">
//                 <h2 className="text-text-primary text-xl font-bold mb-1">
//                     AI Agents
//                 </h2>
//                 <p className="text-text-secondary text-sm">
//                     Select an agent to chat
//                 </p>
//             </div>

//             {/* СПИСОК АГЕНТОВ */}
//             <div className="
//                 flex-1 
//                 overflow-y-auto
//                 /* Кастомный скроллбар */
//                 scrollbar-thin scrollbar-thumb-bright-turquoise/30 scrollbar-track-transparent
//             ">
//                 <div
//                     onClick={() => onSelectAgent('')}
//                     className={`
//                         p-4 cursor-pointer transition-all border-l-4 hover:bg-deep-midnight/50
//                         ${!selectedAgentId ? 'border-bright-turquoise bg-deep-midnight/30' : 'border-transparent'}
//                     `}
//                 >
//                     <h3 className="text-text-primary font-semibold">Global Stream</h3>
//                     <p className="text-text-secondary text-sm">View all interactions</p>
//                 </div>

//                 {agents.map((agent) => (
//                     <div
//                         key={agent.id}
//                         onClick={() => onSelectAgent(agent.id)}
//                         className={`
//                             /* Базовые стили элемента списка */
//                             p-4 
//                             cursor-pointer
//                             transition-all
//                             border-l-4
//                             hover:bg-deep-midnight/50
                            
//                             /* Активный агент выделяется */
//                             ${agent.id === selectedAgentId
//                                 ? 'border-bright-turquoise bg-deep-midnight/30'
//                                 : 'border-transparent'
//                             }
//                         `}
//                     >
//                         {/* ВЕРХНЯЯ ЧАСТЬ: Имя + Статус */}
//                         <div className="flex items-center justify-between mb-2">
//                             {/* Имя агента */}
//                             <h3 className="text-text-primary font-semibold">
//                                 {agent.name}
//                             </h3>

//                             {/* Индикатор статуса */}
//                             <div className="flex items-center gap-2">
//                                 <div className={`
//                                     w-2 h-2 rounded-full
//                                     ${agent.isActive
//                                         ? 'bg-light-mint animate-pulse'
//                                         : 'bg-text-secondary/50'
//                                     }
//                                 `} />
//                                 <span className="text-text-secondary text-xs">
//                                     {agent.isActive ? 'online' : 'offline'}
//                                 </span>
//                             </div>
//                         </div>

//                         {/* НИЖНЯЯ ЧАСТЬ: Тип личности */}
//                         <p className="text-text-secondary text-sm">
//                             {agent.personalityType}
//                         </p>
//                     </div>
//                 ))}
//             </div>

//             {/* FOOTER САЙДБАРА (опционально) */}
//             <div className="
//                 p-4 
//                 border-t border-bright-turquoise/20
//                 bg-deep-midnight/50
//             ">
//                 <button className="mt-4 w-full h-14 bg-gradient-to-r from-bright-turquoise to-soft-teal rounded-2xl text-deep-midnight font-bold flex items-center justify-center gap-2">
//                     <Plus className="w-5 h-5" />
//                     <span>Новый агент</span>
//                 </button>
//             </div>
//         </aside>
//     );
// };


import type { AgentSummary } from '../shared/types';
import { Plus, Radio } from 'lucide-react';

interface ChatSidebarProps {
    agents: AgentSummary[];
    selectedAgentId: string | null;
    onSelectAgent: (id: string) => void;
    className?: string;
}

export const ChatSidebar = ({ agents, selectedAgentId, onSelectAgent, className }: ChatSidebarProps) => {
    return (
        <aside className={`
            w-full flex flex-col h-full bg-transparent
            ${className || ''}
        `}>
            {/* ШАПКА САЙДБАРА */}
            <div className="p-6 shrink-0">
                <h2 className="text-white text-xs font-mono uppercase tracking-[0.3em] opacity-40 mb-1">
                    Нейросети
                </h2>
                <p className="text-white/60 text-sm">
                    Выберите активный поток
                </p>
            </div>

            {/* СПИСОК АГЕНТОВ */}
            <div className="
                flex-1 
                overflow-y-auto 
                no-scrollbar 
                px-4 
                space-y-2
            ">
                {/* Глобальный поток (Global Stream) */}
                <div
                    onClick={() => onSelectAgent('')}
                    className={`
                        group relative p-4 rounded-2xl cursor-pointer transition-all duration-300
                        ${!selectedAgentId 
                            ? 'bg-white/10 shadow-[inset_0_0_20px_rgba(255,255,255,0.05)] border border-white/10' 
                            : 'hover:bg-white/5 border border-transparent'}
                    `}
                >
                    <div className="flex items-center gap-3">
                        <div className={`
                            w-10 h-10 rounded-xl flex items-center justify-center transition-colors
                            ${!selectedAgentId ? 'bg-bright-turquoise text-deep-midnight' : 'bg-white/5 text-white/40 group-hover:text-white'}
                        `}>
                            <Radio className="w-5 h-5" />
                        </div>
                        <div>
                            <h3 className="text-white font-bold text-sm">Общий поток</h3>
                            <p className="text-white/30 text-xs">Все взаимодействия</p>
                        </div>
                    </div>
                </div>

                <div className="h-4" /> {/* Разделитель */}

                {/* Список AI Агентов */}
                {agents.map((agent) => (
                    <div
                        key={agent.id}
                        onClick={() => onSelectAgent(agent.id)}
                        className={`
                            group relative p-4 rounded-2xl cursor-pointer transition-all duration-300 border
                            ${agent.id === selectedAgentId
                                ? 'bg-bright-turquoise/10 border-bright-turquoise/30 shadow-lg shadow-bright-turquoise/5'
                                : 'bg-white/[0.02] border-white/5 hover:border-white/10 hover:bg-white/[0.04]'}
                        `}
                    >
                        <div className="flex items-center justify-between mb-2">
                            <h3 className={`text-sm font-bold transition-colors ${agent.id === selectedAgentId ? 'text-bright-turquoise' : 'text-white/80'}`}>
                                {agent.name}
                            </h3>
                            
                            <div className="flex items-center gap-1.5">
                                <div className={`
                                    w-1.5 h-1.5 rounded-full
                                    ${agent.isActive ? 'bg-light-mint animate-pulse' : 'bg-white/10'}
                                `} />
                                <span className="text-[10px] uppercase tracking-wider opacity-40 text-white">
                                    {agent.isActive ? 'в сети' : 'спит'}
                                </span>
                            </div>
                        </div>

                        <p className="text-white/40 text-xs line-clamp-1 font-light">
                            {agent.personalityType}
                        </p>

                        {/* Индикатор активности сбоку */}
                        {agent.id === selectedAgentId && (
                            <div className="absolute left-[-4px] top-1/4 bottom-1/4 w-1 bg-bright-turquoise rounded-full shadow-[0_0_10px_#26d0ce]" />
                        )}
                    </div>
                ))}
            </div>

            {/* НИЖНЯЯ ЧАСТЬ (КНОПКА) */}
            <div className="p-6 shrink-0 mt-auto border-t border-white/5 ">
                <button 
                    className="w-full h-14 bg-gradient-to-r from-bright-turquoise to-soft-teal rounded-2xl flex items-center justify-center gap-3 text-deep-midnight font-bold shadow-lg shadow-bright-turquoise/20 hover:scale-[1.02] active:scale-[0.98] transition-all duration-300 group"
                >
                    <Plus className="w-5 h-5 group-hover:rotate-90 transition-transform duration-300" />
                    <span>Новый агент</span>
                </button>
            </div>
        </aside>
    );
};