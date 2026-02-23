import type { AgentSummary } from '../shared/types';
import { Plus, Radio } from 'lucide-react';
import { useState } from 'react';
import { NewAgentPopUp } from './New-Agent';

interface ChatSidebarProps {
    agents: AgentSummary[];
    selectedAgentId: string | null;
    onSelectAgent: (id: string) => void;
    className?: string;
}

export const ChatSidebar = ({ agents, selectedAgentId, onSelectAgent, className }: ChatSidebarProps) => {
    const [isNewAgentModalOpen, setIsNewAgentModalOpen] = useState(false);
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
                    onClick={() => setIsNewAgentModalOpen(true)}
                    className="w-full h-14 bg-gradient-to-r from-bright-turquoise to-soft-teal rounded-2xl flex items-center justify-center gap-3 text-deep-midnight font-bold shadow-lg shadow-bright-turquoise/20 hover:scale-[1.02] active:scale-[0.98] transition-all duration-300 group"
                >
                    <Plus className="w-5 h-5 group-hover:rotate-90 transition-transform duration-300" />
                    <span>Новый агент</span>
                </button>
            </div>

            <NewAgentPopUp
                isOpen={isNewAgentModalOpen}
                onClose={() => setIsNewAgentModalOpen(false)}
            />
        </aside>
    );
};