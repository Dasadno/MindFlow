// import { useState, useEffect } from 'react';
// import { Menu, X, User, Circle } from 'lucide-react';
// import { useChatStore } from '../model/store';
// import { ChatSidebar, MessageList, MessageInput } from '../../../widgets';

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

//     useEffect(() => {
//         fetchAgents();
//         initConnection();
//         return () => disconnect();
//     }, [fetchAgents, initConnection, disconnect]);

//     const currentAgent = selectedAgentId ? agents.find(a => a.id === selectedAgentId) : null;

//     return (
//         <div className="flex h-screen overflow-hidden bg-deep-midnight">
//             {/* SIDEBAR (Desktop) */}
//             <div className="hidden md:block">
//                 <ChatSidebar
//                     agents={agents}
//                     selectedAgentId={selectedAgentId}
//                     onSelectAgent={selectAgent}
//                 />
//             </div>

//             {/* MOBILE SIDEBAR */}
//             {isSidebarOpen && (
//                 <>
//                     <div
//                         className="fixed inset-0 bg-black/50 z-40 md:hidden"
//                         onClick={() => setIsSidebarOpen(false)}
//                     />
//                     <div className="fixed top-0 left-0 h-full w-80 z-50 md:hidden transform transition-transform duration-300">
//                         <ChatSidebar
//                             agents={agents}
//                             selectedAgentId={selectedAgentId}
//                             onSelectAgent={(id) => {
//                                 selectAgent(id);
//                                 setIsSidebarOpen(false);
//                             }}
//                         />
//                     </div>
//                 </>
//             )}

//             {/* MAIN CHAT AREA */}
//             <div className="flex-1 flex flex-col overflow-hidden">
//                 {/* HEADER */}
//                 <header className="bg-dark-ocean border-b border-bright-turquoise/20 px-4 py-3 flex items-center gap-4 sticky top-0 z-30">
//                     <button
//                         onClick={() => setIsSidebarOpen(!isSidebarOpen)}
//                         className="md:hidden p-2 text-text-primary hover:bg-deep-midnight/50 rounded-lg transition-colors"
//                     >
//                         {isSidebarOpen ? <X className="w-6 h-6" /> : <Menu className="w-6 h-6" />}
//                     </button>

//                     <div className="w-10 h-10 bg-gradient-primary rounded-full flex items-center justify-center text-white font-bold">
//                         {currentAgent ? currentAgent.name[0] : 'G'}
//                     </div>

//                     <div className="flex-1">
//                         <h1 className="text-text-primary font-bold text-lg">
//                             {currentAgent ? currentAgent.name : 'Global Stream'}
//                         </h1>
//                         <div className="flex items-center gap-2">
//                             <Circle className={`w-2 h-2 ${currentAgent?.isActive ? 'fill-light-mint text-light-mint' : 'fill-gray-500 text-gray-500'}`} />
//                             <span className="text-text-secondary text-sm">
//                                 {currentAgent
//                                     ? `${currentAgent.isActive ? 'online' : 'offline'} • ${currentAgent.personalityType}`
//                                     : 'Viewing all messages'}
//                             </span>
//                         </div>
//                     </div>

//                     <button className="p-2 text-text-secondary hover:text-text-primary hover:bg-deep-midnight/50 rounded-lg transition-colors">
//                         <User className="w-5 h-5" />
//                     </button>
//                 </header>

//                 {/* MESSAGES */}
//                 <MessageList messages={messages} />

//                 {/* INPUT */}
//                 <MessageInput onSendMessage={sendMessage} disabled={!selectedAgentId} />
//             </div>
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

    return (
        <div className="flex h-screen overflow-hidden bg-deep-midnight relative selection:bg-bright-turquoise/30 selection:text-white font-sans">
            
            {/* --- Advanced Animations & Global Styles --- */}
            <style>{`
                @keyframes float {
                    0%, 100% { transform: translateY(0px) opacity: 0.3; }
                    50% { transform: translateY(-20px) opacity: 0.6; }
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

            {/* Фоновый градиентный ландшафт */}
            <div className="fixed inset-0 overflow-hidden pointer-events-none">
                <div className="absolute top-[-10%] left-[-5%] w-[50%] h-[50%] bg-bright-turquoise/10 blur-[140px] rounded-full animate-pulse" />
                <div className="absolute bottom-[-5%] right-[-5%] w-[40%] h-[40%] bg-soft-teal/10 blur-[120px] rounded-full" style={{ animation: 'float 10s infinite' }} />
                <div className="absolute top-[30%] right-[10%] w-[20%] h-[20%] bg-sky-blue/5 blur-[100px] rounded-full" />
            </div>

            {/* --- SIDEBAR (Desktop) --- */}
            <aside className="hidden md:flex flex-col w-72 bg-white/[0.02] border-r border-white/5 backdrop-blur-3xl relative z-20">
                <div className="p-6 mb-2">
                    <div className="flex items-center gap-3 group transition-transform duration-300 hover:scale-[1.02]">
                        <img 
                            src="/cover2.png" 
                            alt="Logo" 
                            className="w-9 h-9 rounded-xl object-cover shadow-[0_0_20px_rgba(38,208,206,0.3)]" 
                        />
                        <div className="text-xl font-bold bg-gradient-accent bg-clip-text text-transparent tracking-tighter uppercase text-white">
                            MindFlow
                        </div>
                    </div>
                </div>
                <div className="flex-1 overflow-y-auto no-scrollbar px-3">
                    <ChatSidebar
                        agents={agents}
                        selectedAgentId={selectedAgentId}
                        onSelectAgent={selectAgent}
                    />
                </div>
            </aside>

            {/* --- MAIN CHAT AREA --- */}
            <main className="flex-1 flex flex-col overflow-hidden relative z-10">
                
                {/* HEADER */}
                <header className="h-20 flex items-center justify-between px-8 border-b border-white/5 glass-panel shrink-0">
                    <div className="flex items-center gap-5">
                        <button
                            onClick={() => setIsSidebarOpen(true)}
                            className="md:hidden p-2 text-white/70 hover:text-bright-turquoise transition-colors"
                        >
                            <Menu className="w-6 h-6" />
                        </button>

                        <div className="flex items-center gap-4 animate-slide-up">
                            <div className="relative group">
                                <div className="w-11 h-11 bg-gradient-primary rounded-2xl flex items-center justify-center text-white font-bold text-lg shadow-lg border border-white/10 group-hover:shadow-bright-turquoise/20 transition-all">
                                    {currentAgent ? currentAgent.name[0] : 'M'}
                                </div>
                                {currentAgent?.isActive && (
                                    <div className="absolute -bottom-1 -right-1 w-3.5 h-3.5 bg-light-mint rounded-full border-2 border-deep-midnight animate-pulse" />
                                )}
                            </div>

                            <div>
                                <h1 className="text-white font-semibold tracking-tight text-lg leading-none mb-1">
                                    {currentAgent ? currentAgent.name : 'System Stream'}
                                </h1>
                                <span className="text-[10px] uppercase tracking-[0.2em] text-bright-turquoise/60 font-medium">
                                    {currentAgent?.isActive ? 'Connection Active' : 'Waiting for link'}
                                </span>
                            </div>
                        </div>
                    </div>

                    <button className="w-10 h-10 flex items-center justify-center rounded-xl text-white/30 hover:text-white hover:bg-white/5 transition-all">
                        <MoreVertical className="w-5 h-5" />
                    </button>
                </header>

                {/* MESSAGES AREA */}
                <div className="flex-1 overflow-y-auto no-scrollbar scroll-smooth px-6 py-10 flex flex-col gap-6">
                    <MessageList messages={messages} />
                    
                    {messages.length === 0 && (
                        <div className="flex-1 flex flex-col items-center justify-center animate-pulse">
                            <div className="w-16 h-16 rounded-3xl bg-white/[0.03] flex items-center justify-center mb-4 border border-white/5">
                                <img src="/cover2.png" alt="Logo" className="w-8 h-8 opacity-20 grayscale" />
                            </div>
                            <p className="text-white/20 font-mono text-[10px] tracking-[0.3em] uppercase">
                                System ready for input
                            </p>
                        </div>
                    )}
                </div>

                {/* INPUT AREA */}
                <div className="p-6 md:p-10 bg-gradient-to-t from-deep-midnight via-deep-midnight/90 to-transparent">
                    <form 
                        onSubmit={handleSend} 
                        className="max-w-4xl mx-auto flex items-end gap-3 p-2 rounded-[24px] bg-white/[0.03] border border-white/5 backdrop-blur-2xl focus-within:border-bright-turquoise/30 transition-all duration-500 shadow-2xl"
                    >
                        <div className="flex-1">
                            <Input
                                placeholder={selectedAgentId ? "Type your message..." : "Select an agent to begin..."}
                                value={inputValue}
                                onChange={(e) => setInputValue(e.target.value)}
                                disabled={!selectedAgentId}
                                className="bg-transparent border-none shadow-none focus:ring-0 text-white placeholder:text-white/20 py-4 px-4"
                            />
                        </div>

                        <Button
                            variant="gradient"
                            type="submit"
                            // disabled={!selectedAgentId || !inputValue.trim()}
                            className="h-12 w-12 md:w-auto md:px-6 rounded-[16px] flex items-center justify-center gap-2 group overflow-hidden relative"
                        >
                            <span className="hidden md:inline font-bold tracking-tighter uppercase text-sm">Send</span>
                            <Send className="w-4 h-4 group-hover:translate-x-1 group-hover:-translate-y-1 transition-transform" />
                        </Button>
                    </form>
                    
                    <div className="mt-4 flex justify-center items-center gap-4">
                        <div className="h-[1px] w-8 bg-white/5" />
                            <span className="text-[9px] font-mono text-white/10 tracking-[0.4em] uppercase">
                                MindFlow Neural Link v1.0
                            </span>
                        <div className="h-[1px] w-8 bg-white/5" />
                    </div>
                </div>

            </main>

            {/* MOBILE SIDEBAR OVERLAY */}
            {isSidebarOpen && (
                <div className="fixed inset-0 z-[100] md:hidden animate-scale-in">
                    <div className="absolute inset-0 bg-deep-midnight/60 backdrop-blur-md" onClick={() => setIsSidebarOpen(false)} />
                    <nav className="absolute top-0 left-0 bottom-0 w-72 bg-deep-midnight border-r border-white/10 p-6 flex flex-col shadow-2xl">
                        <div className="flex items-center justify-between mb-8">
                            <div className="text-lg font-bold bg-gradient-accent bg-clip-text text-transparent uppercase tracking-tighter">
                                MindFlow
                            </div>
                            <button onClick={() => setIsSidebarOpen(false)} className="text-white/40"><X /></button>
                        </div>
                        <div className="flex-1 overflow-y-auto no-scrollbar">
                            <ChatSidebar
                                agents={agents}
                                selectedAgentId={selectedAgentId}
                                onSelectAgent={(id) => { selectAgent(id); setIsSidebarOpen(false); }}
                            />
                        </div>
                    </nav>
                </div>
            )}
        </div>
    );
};