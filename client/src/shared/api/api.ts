import type { Agent, AgentSummary, Event } from '../types';

const API_BASE_URL = 'http://localhost:8080';

export const chatApi = {
    getAgents: async (): Promise<{ agents: AgentSummary[] }> => {
        const response = await fetch(`${API_BASE_URL}/agents`);
        if (!response.ok) {
            throw new Error('Failed to fetch agents');
        }
        return response.json();
    },

    getAgent: async (id: string): Promise<Agent> => {
        const response = await fetch(`${API_BASE_URL}/agents/${id}`);
        if (!response.ok) {
            throw new Error(`Failed to fetch agent ${id}`);
        }
        return response.json();
    },

    injectMessage: async (agentId: string, content: string): Promise<void> => {
        const response = await fetch(`${API_BASE_URL}/agents/${agentId}/inject`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ type: 'message', content }),
        });

        if (!response.ok) {
            throw new Error('Failed to send message');
        }
    },

    subscribeToEvents: (onMessage: (event: Event) => void): (() => void) => {
        const eventSource = new EventSource(`${API_BASE_URL}/events/stream`);

        eventSource.onmessage = (event) => {
            try {
                const data: Event = JSON.parse(event.data);
                onMessage(data);
            } catch (error) {
                console.error('Failed to parse SSE message:', error);
            }
        };

        eventSource.onerror = (error) => {
            console.error('SSE Error:', error);
            eventSource.close();
        };

        return () => {
            eventSource.close();
        };
    },
};
