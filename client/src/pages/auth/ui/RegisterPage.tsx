import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import Input from '@/shared/ui/Input';
import { Button } from '@/shared/ui/Button';

const RegisterPage = () => {
    const navigate = useNavigate();
    
    const [name, setName] = useState('');
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        console.log('Registration attempt:', { name, email, password });
    };

    return (
        <div className="relative flex items-center justify-center min-h-screen px-4 overflow-hidden bg-deep-midnight">
            
            {/* Анимация появления как в LoginPage */}
            <style>{`
                @keyframes scaleIn {
                    from { opacity: 0; transform: scale(0.98); }
                    to { opacity: 1; transform: scale(1); }
                }
                .animate-container { animation: scaleIn 0.5s ease-out forwards; }
                ::-webkit-scrollbar { display: none; }
                * { -ms-overflow-style: none; scrollbar-width: none; }
            `}</style>

            {/* Фоновые градиенты */}
            <div className="fixed inset-0 overflow-hidden pointer-events-none">
                <div className="absolute top-[-10%] left-[-10%] w-[45%] h-[45%] bg-bright-turquoise/10 blur-[130px] rounded-full animate-pulse" />
                <div className="absolute bottom-[-10%] right-[-10%] w-[45%] h-[45%] bg-light-mint/10 blur-[130px] rounded-full animate-pulse" />
            </div>

            <main className="relative z-10 w-full max-w-[460px] animate-container overflow-hidden">
                <div className="bg-white/[0.03] border border-white/10 backdrop-blur-3xl p-10 md:p-14 rounded-[3rem] shadow-[0_25px_50px_rgba(0,0,0,0.4)] overflow-hidden">
                    
                    {/* Заголовок с защитой хвостика буквы "я" */}
                    <div className="text-center mb-10">
                        <h1 className="text-4xl font-black pb-2 bg-gradient-to-r from-bright-turquoise to-light-mint bg-clip-text text-transparent tracking-tight leading-tight">
                            Создать аккаунт
                        </h1>
                    </div>

                    <form onSubmit={handleSubmit} className="flex flex-col gap-5 w-full">
                        {/* Поле Имя */}
                        <div className="flex flex-col">
                            <label className="text-xl font-bold text-white tracking-wide mb-3 ml-1">
                                Имя
                            </label>
                            <Input
                                type="text"
                                value={name}
                                onChange={(e) => setName(e.target.value)}
                                placeholder="Иван Иванов"
                                required
                            />
                        </div>

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
                                Зарегистрироваться
                            </Button>
                        </div>
                    </form>

                    <div className="mt-10 text-center">
                        <button 
                            onClick={() => navigate('/login')}
                            className="text-base font-medium text-white/70 hover:text-bright-turquoise transition-all duration-300 group"
                        >
                            Уже есть аккаунт? <span className="text-bright-turquoise font-bold underline underline-offset-8 decoration-bright-turquoise/30 group-hover:decoration-bright-turquoise">Войти</span>
                        </button>
                    </div>
                </div>
            </main>
        </div>
    );
};

export default RegisterPage;