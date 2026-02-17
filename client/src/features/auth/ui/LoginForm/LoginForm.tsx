import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import Input from '@/shared/ui/Input';
import { Button } from '@/shared/ui/Button';
import { useAuthStore } from '@/features/auth/model/store';

export const LoginForm = () => {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const login = useAuthStore((state) => state.login);
    const navigate = useNavigate();

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        login(email);
        navigate('/chat');
    };

    return (
        <form onSubmit={handleSubmit} className="flex flex-col gap-6 w-full">
            {/* Поле Email */}
            <div className="flex flex-col">
                <label className="text-xl font-bold text-white tracking-wide mb-3 ml-1">
                    Эл. Почта
                </label>
                <Input
                    type="email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    placeholder="your@email.com"
                    required
                />
            </div>

            {/* Поле Password */}
            <div className="flex flex-col">
                <label className="text-xl font-bold text-white tracking-wide mb-3 ml-1">
                    Пароль
                </label>
                <Input
                    type="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    placeholder="********"
                    required
                />
            </div>

            <div className="pt-4">
                <Button
                    variant="gradient"
                    type="submit"
                    className="w-full py-4 text-lg uppercase tracking-widest shadow-[0_0_20px_rgba(38,208,206,0.3)]"
                >
                    Войти
                </Button>
            </div>
        </form>
    );
};