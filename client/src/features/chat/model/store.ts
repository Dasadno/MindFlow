import { create } from 'zustand';
import type { AgentSummary, Event, Agent } from '../../../shared/types';
import { chatApi } from '../../../shared/api/api';

interface ChatState {
    agents: AgentSummary[];
    selectedAgentId: string | null;
    agentDetails: Record<string, Agent>; // Cache detailed agent info
    messages: Event[];
    isConnected: boolean;
    isLoading: boolean;
    error: string | null;
    cleanup?: () => void;

    // Actions
    fetchAgents: () => Promise<void>;
    selectAgent: (id: string | null) => void;
    getAgentDetails: (id: string) => Promise<void>;
    addMessage: (message: Event) => void;
    sendMessage: (content: string) => Promise<void>;
    initConnection: () => void;
    disconnect: () => void;
}

export const useChatStore = create<ChatState>((set, get) => ({
    agents: [],
    selectedAgentId: null,
    agentDetails: {},
    messages: [],
    isConnected: false,
    isLoading: false,
    error: null,
    cleanup: undefined,

    fetchAgents: async () => {
        set({ isLoading: true, error: null });
        try {
            const { agents } = await chatApi.getAgents();
            set({ agents, isLoading: false });
        } catch (error: any) {
            set({ error: error.message, isLoading: false });
        }
    },

    selectAgent: (id) => {
        set({ selectedAgentId: id });
        if (id && !get().agentDetails[id]) {
            get().getAgentDetails(id);
        }
    },

    getAgentDetails: async (id) => {
        try {
            const agent = await chatApi.getAgent(id);
            set((state) => ({
                agentDetails: { ...state.agentDetails, [id]: agent },
            }));
        } catch (error) {
            console.error(`Failed to load details for agent ${id}`, error);
        }
    },

    addMessage: (message) => {
        set((state) => ({
            messages: [...state.messages, message],
        }));
    },

    sendMessage: async (content) => {
        const { selectedAgentId } = get();
        if (!selectedAgentId) return;

        try {
            await chatApi.injectMessage(selectedAgentId, content);
            // Optimistic update could be done here, but we rely on SSE for the "true" message
        } catch (error: any) {
            set({ error: error.message });
        }
    },

    initConnection: () => {
        if (get().isConnected) return;

        const cleanup = chatApi.subscribeToEvents((event) => {
            get().addMessage(event);
        });

        set({ isConnected: true, cleanup });
    },

    disconnect: () => {
        const { cleanup } = get();
        if (cleanup) cleanup();
        set({ isConnected: false, cleanup: undefined });
    },
}));
