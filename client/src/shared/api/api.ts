import type { Agent, AgentSummary, Event } from "../types";

const API_BASE_URL = ""; // Use relative path for proxy

export const chatApi = {
  getAgents: async (): Promise<{ agents: AgentSummary[] }> => {
    const response = await fetch(`${API_BASE_URL}/agents`);
    if (!response.ok) {
      throw new Error("Failed to fetch agents");
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
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ type: "message", content }),
    });

    if (!response.ok) {
      throw new Error("Failed to send message");
    }
  },


  // Реализовал переподключение а именно exponential backoff. 
  subscribeToEvents: (onMessage: (event: Event) => void): (() => void) => {
    let eventSource: EventSource | null = null;
    let retryCount = 0;
    const MAX_RETRIES = 10;
    const BASE_DELAY = 1000; // 1 секунда
    let cleanup: (() => void) | null = null;
    let isActive = true;

    const connect = () => {
      if (!isActive) return;

      eventSource = new EventSource(`${API_BASE_URL}/events/stream`);

      eventSource.onmessage = (event) => {
        try {
          const data: Event = JSON.parse(event.data);
          onMessage(data);
          retryCount = 0;
          
        } catch (error) {
          console.error("Failed to parse SSE message:", error);
        }
      };

      eventSource.onerror = (error) => {
        console.error("SSE Error:", error);
        eventSource?.close();

        if (isActive && retryCount < MAX_RETRIES) {
          // Экспоненциальная задержка с джиттером
          const delay = calculateBackoff(retryCount);

          console.log(
            `Reconnecting in ${delay}ms (attempt ${retryCount + 1}/${MAX_RETRIES})`,
          );

          setTimeout(() => {
            retryCount++;
            connect();
          }, delay);
        } else if (retryCount >= MAX_RETRIES) {
          console.error("Max retries reached. Stopping reconnection attempts.");
        }
      };

      cleanup = () => {
        eventSource?.close();
      };
    };

    const calculateBackoff = (attempt: number): number => {
      // Экспоненциальный рост
      const exponentialDelay = BASE_DELAY * Math.pow(2, attempt);

      // Тут бля джиттер чтобы коллизий не было 
      const jitter = exponentialDelay * 0.2 * (Math.random() * 2 - 1);

      return Math.min(exponentialDelay + jitter, 30000);
    };

    connect();

    // Очистка 
    return () => {
      isActive = false;
      if (cleanup) {
        cleanup();
      }
    };
  },
};
