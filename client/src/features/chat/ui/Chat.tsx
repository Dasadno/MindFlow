import { useState, useEffect } from 'react';
import { Menu, X, User, Circle } from 'lucide-react';
import { useChatStore } from '../model/store';
import { ChatSidebar, MessageList, MessageInput } from '../../../widgets';

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

    useEffect(() => {
        fetchAgents();
        initConnection();
        return () => disconnect();
    }, [fetchAgents, initConnection, disconnect]);

    const currentAgent = selectedAgentId ? agents.find(a => a.id === selectedAgentId) : null;

    return (
        <div className="flex h-screen overflow-hidden bg-deep-midnight">
            {/* SIDEBAR (Desktop) */}
            <div className="hidden md:block">
                <ChatSidebar
                    agents={agents}
                    selectedAgentId={selectedAgentId}
                    onSelectAgent={selectAgent}
                />
            </div>

            {/* MOBILE SIDEBAR */}
            {isSidebarOpen && (
                <>
                    <div
                        className="fixed inset-0 bg-black/50 z-40 md:hidden"
                        onClick={() => setIsSidebarOpen(false)}
                    />
                    <div className="fixed top-0 left-0 h-full w-80 z-50 md:hidden transform transition-transform duration-300">
                        <ChatSidebar
                            agents={agents}
                            selectedAgentId={selectedAgentId}
                            onSelectAgent={(id) => {
                                selectAgent(id);
                                setIsSidebarOpen(false);
                            }}
                        />
                    </div>
                </>
            )}

            {/* MAIN CHAT AREA */}
            <div className="flex-1 flex flex-col overflow-hidden">
                {/* HEADER */}
                <header className="bg-dark-ocean border-b border-bright-turquoise/20 px-4 py-3 flex items-center gap-4 sticky top-0 z-30">
                    <button
                        onClick={() => setIsSidebarOpen(!isSidebarOpen)}
                        className="md:hidden p-2 text-text-primary hover:bg-deep-midnight/50 rounded-lg transition-colors"
                    >
                        {isSidebarOpen ? <X className="w-6 h-6" /> : <Menu className="w-6 h-6" />}
                    </button>

                    <div className="w-10 h-10 bg-gradient-primary rounded-full flex items-center justify-center text-white font-bold">
                        {currentAgent ? currentAgent.name[0] : 'G'}
                    </div>

                    <div className="flex-1">
                        <h1 className="text-text-primary font-bold text-lg">
                            {currentAgent ? currentAgent.name : 'Global Stream'}
                        </h1>
                        <div className="flex items-center gap-2">
                            <Circle className={`w-2 h-2 ${currentAgent?.isActive ? 'fill-light-mint text-light-mint' : 'fill-gray-500 text-gray-500'}`} />
                            <span className="text-text-secondary text-sm">
                                {currentAgent
                                    ? `${currentAgent.isActive ? 'online' : 'offline'} • ${currentAgent.personalityType}`
                                    : 'Viewing all messages'}
                            </span>
                        </div>
                    </div>

                    <button className="p-2 text-text-secondary hover:text-text-primary hover:bg-deep-midnight/50 rounded-lg transition-colors">
                        <User className="w-5 h-5" />
                    </button>
                </header>

                {/* MESSAGES */}
                <MessageList messages={messages} />

                {/* INPUT */}
                <MessageInput onSendMessage={sendMessage} disabled={!selectedAgentId} />
            </div>
        </div>
    );
};
