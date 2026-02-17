import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface User {
    id: string;
    username: string;
}

interface AuthState {
    isAuthenticated: boolean;
    user: User | null;
    login: (username: string) => void;
    logout: () => void;
}

export const useAuthStore = create<AuthState>()(
    persist(
        (set) => ({
            isAuthenticated: false,
            user: null,

            login: (username: string) => {
                // In a real app, this would verify credentials against a backend
                set({
                    isAuthenticated: true,
                    user: { id: '1', username },
                });
            },

            logout: () => {
                set({
                    isAuthenticated: false,
                    user: null,
                });
            },
        }),
        {
            name: 'auth-storage', // name of the item in the storage (must be unique)
        }
    )
);
