export interface Personality {
    openness: number;
    conscientiousness: number;
    extraversion: number;
    agreeableness: number;
    neuroticism: number;
    coreValues: string[];
    quirks: string[];
}

export interface Agent {
    id: string;
    name: string;
    personality: Personality;
    currentMood: {
        label: string;
        pad: {
            pleasure: number;
            arousal: number;
            dominance: number;
        };
        activeEmotions: Array<{
            type: string;
            intensity: number;
            trigger: string;
        }>;
    };
    goals: Array<{
        id: string;
        description: string;
        priority: number;
        progress: number;
    }>;
    stats: {
        totalInteractions: number;
        memoriesCount: number;
        relationshipsCount: number;
        daysSinceCreation: number;
    };
    createdAt: string;
}

export interface AgentSummary {
    id: string;
    name: string;
    personalityType: string;
    currentMood: string;
    moodIntensity: number;
    isActive: boolean;
}

export interface Event {
    type: string;
    speaker?: string;
    target?: string;
    content?: string;
    tick?: number;
    topic?: string;
    source?: string;
    affected_agents?: string[];
    payload?: any;
    status?: string;
    created_at?: string;
}

export interface WorldStatus {
    currentTick: number;
    simulationSpeed: number;
    isPaused: boolean;
    activeAgents: number;
    totalEvents: number;
    uptime: string;
}
