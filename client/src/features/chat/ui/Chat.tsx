// import { useState, useEffect } from 'react';
// import { Menu, X, Send, MoreVertical, Plus } from 'lucide-react';
// import { useChatStore } from '../model/store';
// import { ChatSidebar, MessageList } from '../../../widgets';
// import { Button } from '@/shared/ui/Button';
// import Input from '@/shared/ui/Input';

// const getAgentColor = (id: string) => {
//     if (!id) return '#26d0ce'; // Дефолтный бирюзовый для общего потока
//     let hash = 0;
//     for (let i = 0; i < id.length; i++) {
//         hash = id.charCodeAt(i) + ((hash << 5) - hash);
//     }
//     const hue = Math.abs(hash % 360);
//     // Возвращаем HSL: насыщенность 75%, яркость 65% для неона
//     return `hsl(${hue}, 75%, 65%)`;
// };

// export const Chat = () => {
//     const {
//         agents,
//         messages,
//         selectedAgentId,
//         fetchAgents,
//         selectAgent,
//         sendMessage,
//         initConnection,
//         disconnect
//     } = useChatStore();

//     const [isSidebarOpen, setIsSidebarOpen] = useState(false);
//     const [inputValue, setInputValue] = useState('');

//     useEffect(() => {
//         fetchAgents();
//         initConnection();
//         return () => disconnect();
//     }, [fetchAgents, initConnection, disconnect]);

//     const currentAgent = selectedAgentId ? agents.find(a => a.id === selectedAgentId) : null;

//     const handleSend = (e?: React.FormEvent) => {
//         e?.preventDefault();
//         if (inputValue.trim() && selectedAgentId) {
//             sendMessage(inputValue);
//             setInputValue('');
//         }
//     };

//     return (
//         <div className="flex h-screen overflow-hidden bg-deep-midnight relative selection:bg-bright-turquoise/30 selection:text-white font-sans">

//             {/* --- Глобальные стили для скрытия скроллбаров --- */}
//             <style>{`
//                 @keyframes float {
//                     0%, 100% { transform: translateY(0px); opacity: 0.3; }
//                     50% { transform: translateY(-20px); opacity: 0.6; }
//                 }
//                 @keyframes slideUp {
//                     from { opacity: 0; transform: translateY(10px); }
//                     to { opacity: 1; transform: translateY(0); }
//                 }
//                 .animate-slide-up { animation: slideUp 0.4s ease-out forwards; }

//                 /* Скрываем скроллбары для всех браузеров */
//                 .no-scrollbar::-webkit-scrollbar { display: none; }
//                 .no-scrollbar { -ms-overflow-style: none; scrollbar-width: none; }

//                 .glass-panel {
//                     background: linear-gradient(180deg, rgba(255,255,255,0.05) 0%, rgba(255,255,255,0) 100%);
//                     backdrop-filter: blur(20px);
//                 }
//             `}</style>

//             {/* Фоновые градиенты */}
//             <div className="fixed inset-0 overflow-hidden pointer-events-none">
//                 <div className="absolute top-[-10%] left-[-5%] w-[50%] h-[50%] bg-bright-turquoise/10 blur-[140px] rounded-full animate-pulse" />
//                 <div className="absolute bottom-[-5%] right-[-5%] w-[40%] h-[40%] bg-soft-teal/10 blur-[120px] rounded-full" style={{ animation: 'float 10s infinite' }} />
//             </div>

//             {/* --- САЙДБАР (Desktop) --- */}
//             <aside className="hidden md:flex flex-col w-80 bg-white/[0.01] border-r border-white/5 backdrop-blur-3xl relative z-20">
//                 <div className="p-6">
//                     <div className="flex items-center gap-3 group transition-transform duration-300 hover:scale-[1.02]">
//                         <img 
//                             src="/cover2.png" 
//                             alt="Логотип" 
//                             className="w-10 h-10 rounded-2xl object-cover shadow-[0_0_20px_rgba(38,208,206,0.3)]" 
//                         />
//                         <div className="text-xl font-bold bg-gradient-accent bg-clip-text text-transparent tracking-tighter uppercase text-white">
//                             MindFlow
//                         </div>
//                     </div>
//                 </div>

//                 <div className="flex-1 overflow-y-auto no-scrollbar px-4 space-y-2 pb-24">
//                     <ChatSidebar
//                         agents={agents}
//                         selectedAgentId={selectedAgentId}
//                         onSelectAgent={selectAgent}
//                     />
//                 </div>
//             </aside>

//             {/* --- ОСНОВНАЯ ОБЛАСТЬ ЧАТА --- */}
//             <main className="flex-1 flex flex-col overflow-hidden relative z-10">

//                 {/* ХЕДЕР */}
//                 <header className="h-20 flex items-center justify-between px-8 border-b border-white/5 glass-panel shrink-0">
//                     <div className="flex items-center gap-5">
//                         <button
//                             onClick={() => setIsSidebarOpen(true)}
//                             className="md:hidden p-2 text-white/70 hover:text-bright-turquoise transition-colors"
//                         >
//                             <Menu className="w-6 h-6" />
//                         </button>

//                         <div className="flex items-center gap-4 animate-slide-up">
//                             <div className="relative group">
//                                 <div className="w-12 h-12 bg-gradient-primary rounded-2xl flex items-center justify-center text-white font-bold text-xl shadow-lg border border-white/10 group-hover:shadow-bright-turquoise/20 transition-all">
//                                     {currentAgent ? currentAgent.name[0] : 'М'}
//                                 </div>
//                                 {currentAgent?.isActive && (
//                                     <div className="absolute -bottom-1 -right-1 w-4 h-4 bg-light-mint rounded-full border-2 border-deep-midnight animate-pulse" />
//                                 )}
//                             </div>

//                             <div>
//                                 <h1 className="text-white font-semibold tracking-tight text-lg leading-none mb-1">
//                                     {currentAgent ? currentAgent.name : 'Системный поток'}
//                                 </h1>
//                                 <span className="text-[10px] uppercase tracking-[0.2em] text-bright-turquoise/60 font-medium">
//                                     {currentAgent?.isActive ? 'Соединение активно' : 'Ожидание выбора агента'}
//                                 </span>
//                             </div>
//                         </div>
//                     </div>

//                     <button className="w-11 h-11 flex items-center justify-center rounded-2xl text-white/40 hover:text-white hover:bg-white/10 hover:scale-110 active:scale-90 transition-all duration-200 border border-transparent hover:border-white/10">
//                         <MoreVertical className="w-5 h-5" />
//                     </button>
//                 </header>

//                 {/* ОБЛАСТЬ СООБЩЕНИЙ */}
//                 <div className="flex-1 overflow-y-auto no-scrollbar scroll-smooth px-6 py-10 flex flex-col relative">
//                     {messages.length > 0 ? (
//                         <div className="flex flex-col gap-6">
//                             <MessageList messages={messages} />
//                         </div>
//                     ) : (
//                         /* Пустое состояние зафиксировано по центру и не мешает скроллу */
//                         <div className="absolute inset-0 flex flex-col items-center justify-center pointer-events-none">
//                             <div className="group relative">
//                                 <div className="absolute inset-0 bg-bright-turquoise/10 blur-3xl rounded-full scale-150 transition-transform duration-700" />
//                                 <div className="w-32 h-32 rounded-[40px] bg-white/[0.03] flex items-center justify-center mb-6 border border-white/5 overflow-hidden relative z-10">
//                                     <img 
//                                         src="/cover2.png" 
//                                         alt="Логотип" 
//                                         className="w-full h-full object-cover opacity-70 tracking-[0.4em] animate-pulse" 
//                                     />
//                                 </div>
//                             </div>
//                             <p className="text-white/80 font-mono text-[11px] tracking-[0.4em] uppercase animate-pulse text-center">
//                                 Сообщений пока нет<br/>Ожидание активности системы
//                             </p>
//                         </div>
//                     )}
//                 </div>

//                 {/* ОБЛАСТЬ ВВОДА */}
//                 <div className="p-6 md:p-10 bg-gradient-to-t from-deep-midnight via-deep-midnight/90 to-transparent">
//                     <form 
//                         onSubmit={handleSend} 

//                         className="max-w-4xl mx-auto flex items-center gap-3 p-2 rounded-[32px] bg-white/[0.03] border border-white/5 backdrop-blur-3xl focus-within:border-bright-turquoise/30 transition-all duration-500 shadow-2xl"
//                     >
//                         <div className="flex-1">
//                             <Input
//                                 placeholder={selectedAgentId ? "Введите ваше сообщение..." : "Выберите агента для начала связи..."}
//                                 value={inputValue}
//                                 onChange={(e) => setInputValue(e.target.value)}
//                                 disabled={!selectedAgentId}
//                                 className="bg-transparent border-none shadow-none focus:ring-0 text-white placeholder:text-white/10 py-5 px-6 text-lg"
//                             />
//                         </div>

//                         <Button
//                             variant="gradient"
//                             type="submit"
//                             className="h-14 px-8 rounded-3xl flex items-center justify-center gap-3 group overflow-hidden relative active:scale-95 transition-transform shrink-0"
//                         >
//                             <span className="hidden md:inline font-bold tracking-tighter uppercase text-sm">Отправить</span>
//                             <Send className="w-5 h-5 group-hover:translate-x-1 group-hover:-translate-y-1 transition-transform" />
//                         </Button>
//                     </form>

//                     <div className="mt-6 flex justify-center items-center gap-6">
//                         <div className="h-[1px] w-12 bg-white/5" />
//                         <span className="text-[10px] font-mono text-white/80 tracking-[0.5em] uppercase">
//                             Нейронная связь MindFlow v1.0
//                         </span>
//                         <div className="h-[1px] w-12 bg-white/5" />
//                     </div>
//                 </div>

//             </main>
//         </div>
//     );
// };


import { useState, useEffect } from 'react';
import { Menu, X, Send, MoreVertical } from 'lucide-react';
import { useChatStore } from '../model/store';
import { ChatSidebar, MessageList } from '../../../widgets';
import { Button } from '@/shared/ui/Button';
import Input from '@/shared/ui/Input';

export const Chat = () => {
    const {
        agents,
        messages,
        selectedAgentId,
        fetchAgents,
        selectAgent,
        sendMessage,
        initConnection,
        disconnect
    } = useChatStore();

    const [isSidebarOpen, setIsSidebarOpen] = useState(false);
    const [inputValue, setInputValue] = useState('');

    useEffect(() => {
        fetchAgents();
        initConnection();
        return () => disconnect();
    }, [fetchAgents, initConnection, disconnect]);

    const currentAgent = selectedAgentId ? agents.find(a => a.id === selectedAgentId) : null;

    const handleSend = (e?: React.FormEvent) => {
        e?.preventDefault();
        if (inputValue.trim() && selectedAgentId) {
            sendMessage(inputValue);
            setInputValue('');
        }
    };

    // Закрытие сайдбара при выборе агента на мобилке
    const handleSelectAgent = (id: string) => {
        selectAgent(id);
        setIsSidebarOpen(false);
    };

    return (
        <div className="flex h-screen overflow-hidden bg-deep-midnight relative selection:bg-bright-turquoise/30 selection:text-white font-sans">

            <style>{`
                @keyframes float {
                    0%, 100% { transform: translateY(0px); opacity: 0.3; }
                    50% { transform: translateY(-20px); opacity: 0.6; }
                }
                @keyframes slideUp {
                    from { opacity: 0; transform: translateY(10px); }
                    to { opacity: 1; transform: translateY(0); }
                }
                .animate-slide-up { animation: slideUp 0.4s ease-out forwards; }
                .no-scrollbar::-webkit-scrollbar { display: none; }
                .no-scrollbar { -ms-overflow-style: none; scrollbar-width: none; }
                .glass-panel {
                    background: linear-gradient(180deg, rgba(255,255,255,0.05) 0%, rgba(255,255,255,0) 100%);
                    backdrop-filter: blur(20px);
                }
            `}</style>

            {/* Фоновые градиенты */}
            <div className="fixed inset-0 overflow-hidden pointer-events-none">
                <div className="absolute top-[-10%] left-[-5%] w-[50%] h-[50%] bg-bright-turquoise/10 blur-[140px] rounded-full animate-pulse" />
                <div className="absolute bottom-[-5%] right-[-5%] w-[40%] h-[40%] bg-soft-teal/10 blur-[120px] rounded-full" style={{ animation: 'float 10s infinite' }} />
            </div>

            {/* --- МОБИЛЬНЫЙ OVERLAY --- */}
            {isSidebarOpen && (
                <div
                    className="fixed inset-0 bg-black/60 backdrop-blur-sm z-40 md:hidden transition-opacity duration-300"
                    onClick={() => setIsSidebarOpen(false)}
                />
            )}

            {/* --- САЙДБАР (Desktop + Mobile) --- */}
            <aside className={`
                fixed inset-y-0 left-0 z-50 w-72 bg-[#0a0f1a] border-r border-white/5 backdrop-blur-3xl transition-transform duration-300 ease-in-out
                md:relative md:translate-x-0 md:flex md:flex-col md:w-80 md:bg-white/[0.01]
                ${isSidebarOpen ? 'translate-x-0' : '-translate-x-full'}
            `}>
                <div className="p-6 flex items-center justify-between">
                    <div className="flex items-center gap-3 group">
                        <img src="/cover2.png" alt="Лого" className="w-10 h-10 rounded-2xl object-cover shadow-[0_0_20px_rgba(38,208,206,0.3)]" />
                        <div className="text-xl font-bold bg-gradient-accent bg-clip-text text-transparent uppercase text-white">MindFlow</div>
                    </div>
                    {/* Кнопка закрытия для мобилок */}
                    <button onClick={() => setIsSidebarOpen(false)} className="md:hidden p-2 text-white/50 hover:text-white">
                        <X className="w-6 h-6" />
                    </button>
                </div>

                <div className="flex-1 overflow-y-auto no-scrollbar px-4 space-y-2 pb-24">
                    <ChatSidebar
                        agents={agents}
                        selectedAgentId={selectedAgentId}
                        onSelectAgent={handleSelectAgent}
                    />
                </div>
            </aside>

            {/* --- ОСНОВНАЯ ОБЛАСТЬ ЧАТА --- */}
            <main className="flex-1 flex flex-col overflow-hidden relative z-10">

                {/* ХЕДЕР */}
                <header className="h-20 flex items-center justify-between px-4 md:px-8 border-b border-white/5 glass-panel shrink-0">
                    <div className="flex items-center gap-3 md:gap-5">
                        <button
                            onClick={() => setIsSidebarOpen(true)}
                            className="md:hidden p-2 text-white/70 hover:text-bright-turquoise transition-colors"
                        >
                            <Menu className="w-6 h-6" />
                        </button>

                        <div className="flex items-center gap-3 md:gap-4 animate-slide-up">
                            <div className="relative group">
                                <div className="w-10 h-10 md:w-12 md:h-12 bg-gradient-primary rounded-xl md:rounded-2xl flex items-center justify-center text-white font-bold text-lg md:text-xl shadow-lg border border-white/10">
                                    {currentAgent ? currentAgent.name[0] : 'М'}
                                </div>
                                {currentAgent?.isActive && (
                                    <div className="absolute -bottom-1 -right-1 w-3 h-3 md:w-4 md:h-4 bg-light-mint rounded-full border-2 border-deep-midnight animate-pulse" />
                                )}
                            </div>

                            <div className="max-w-[150px] md:max-w-none">
                                <h1 className="text-white font-semibold tracking-tight text-sm md:text-lg leading-none mb-1 truncate">
                                    {currentAgent ? currentAgent.name : 'Системный поток'}
                                </h1>
                                <span className="text-[9px] md:text-[10px] uppercase tracking-[0.1em] md:tracking-[0.2em] text-bright-turquoise/60 font-medium block truncate">
                                    {currentAgent?.isActive ? 'Соединение активно' : 'Ожидание выбора'}
                                </span>
                            </div>
                        </div>
                    </div>

                    <button className="w-10 h-10 md:w-11 md:h-11 flex items-center justify-center rounded-xl md:rounded-2xl text-white/40 hover:text-white hover:bg-white/10 transition-all border border-transparent hover:border-white/10">
                        <MoreVertical className="w-5 h-5" />
                    </button>
                </header>

                {/* ОБЛАСТЬ СООБЩЕНИЙ */}
                <div className="flex-1 overflow-y-auto no-scrollbar scroll-smooth px-4 md:px-6 py-6 md:py-10 flex flex-col relative">
                    {messages.length > 0 ? (
                        <div className="flex flex-col gap-6">
                            <MessageList messages={messages} />
                        </div>
                    ) : (
                        <div className="absolute inset-0 flex flex-col items-center justify-center pointer-events-none p-4">
                            <div className="group relative">
                                <div className="absolute inset-0 bg-bright-turquoise/10 blur-3xl rounded-full scale-150" />
                                <div className="w-24 h-24 md:w-32 md:h-32 rounded-[32px] md:rounded-[40px] bg-white/[0.03] flex items-center justify-center mb-6 border border-white/5 overflow-hidden relative z-10">
                                    <img src="/cover2.png" alt="Логотип" className="w-full h-full object-cover opacity-70 animate-pulse" />
                                </div>
                            </div>
                            <p className="text-white/80 font-mono text-[9px] md:text-[11px] tracking-[0.3em] md:tracking-[0.4em] uppercase animate-pulse text-center leading-relaxed">
                                Сообщений пока нет<br />Ожидание активности системы
                            </p>
                        </div>
                    )}
                </div>

                {/* ОБЛАСТЬ ВВОДА */}
                <div className="p-4 md:p-10 bg-gradient-to-t from-deep-midnight via-deep-midnight/90 to-transparent">
                    <form
                        onSubmit={handleSend}
                        className="max-w-4xl mx-auto flex items-center gap-2 md:gap-3 p-1.5 md:p-2 rounded-[24px] md:rounded-[32px] bg-white/[0.03] border border-white/5 backdrop-blur-3xl focus-within:border-bright-turquoise/30 transition-all duration-500 shadow-2xl"
                    >
                        <div className="flex-1">
                            <Input
                                placeholder={selectedAgentId ? "Напишите..." : "Выберите агента..."}
                                value={inputValue}
                                onChange={(e) => setInputValue(e.target.value)}
                                disabled={!selectedAgentId}
                                className="bg-transparent border-none shadow-none focus:ring-0 text-white placeholder:text-white/20 py-3 md:py-5 px-4 md:px-6 text-base md:text-lg"
                            />
                        </div>

                        <Button
                            variant="gradient"
                            type="submit"
                            className="h-11 md:h-14 px-4 md:px-8 rounded-2xl md:rounded-3xl flex items-center justify-center gap-3 group active:scale-95 transition-transform shrink-0"
                        >
                            <span className="hidden md:inline font-bold tracking-tighter uppercase text-sm">Отправить</span>
                            <Send className="w-5 h-5 group-hover:translate-x-1 group-hover:-translate-y-1 transition-transform" />
                        </Button>
                    </form>

                    <div className="mt-4 md:mt-6 flex justify-center items-center gap-4 md:gap-6">
                        <div className="h-[1px] w-8 md:w-12 bg-white/5" />
                        <span className="text-[8px] md:text-[10px] font-mono text-white/40 tracking-[0.3em] md:tracking-[0.5em] uppercase text-center">
                            MindFlow v1.0 <span className="hidden md:inline">Нейронная связь</span>
                        </span>
                        <div className="h-[1px] w-8 md:w-12 bg-white/5" />
                    </div>
                </div>

            </main>
        </div>
    );
};