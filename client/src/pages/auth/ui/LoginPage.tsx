import { LoginForm } from '@/features/auth/ui/LoginForm/LoginForm';
import { useNavigate } from 'react-router-dom';

export const LoginPage = () => {
    const navigate = useNavigate();

    return (
        <div className="relative flex items-center justify-center min-h-screen px-4 overflow-hidden bg-deep-midnight">
            
            <style>{`
                @keyframes scaleIn {
                    from { opacity: 0; transform: scale(0.98); }
                    to { opacity: 1; transform: scale(1); }
                }
                .animate-container { animation: scaleIn 0.5s ease-out forwards; }
            `}</style>

            <div className="fixed inset-0 overflow-hidden pointer-events-none">
                <div className="absolute top-[-10%] left-[-10%] w-[45%] h-[45%] bg-bright-turquoise/10 blur-[130px] rounded-full animate-pulse" />
                <div className="absolute bottom-[-10%] right-[-10%] w-[45%] h-[45%] bg-light-mint/10 blur-[130px] rounded-full animate-pulse" />
            </div>

            <main className="relative z-10 w-full max-w-[440px] animate-container">
                <div className="bg-white/[0.03] border border-white/10 backdrop-blur-3xl p-10 md:p-14 rounded-[3rem] shadow-[0_25px_50px_rgba(0,0,0,0.4)]">
                    
                    <div className="text-center mb-14">
                        <h1 className="text-4xl font-black pb-2 bg-gradient-to-r from-bright-turquoise to-light-mint bg-clip-text text-transparent tracking-tight leading-tight">
                            Вход в систему
                        </h1>
                    </div>

                    <LoginForm />

                    <div className="mt-12 text-center">
                        <button 
                            onClick={() => navigate('/register')}
                            className="text-base font-medium text-white/70 hover:text-bright-turquoise transition-all duration-300 group"
                        >
                            Нет аккаунта? <span className="text-bright-turquoise font-bold underline underline-offset-8 decoration-bright-turquoise/30 group-hover:decoration-bright-turquoise">Зарегистрироваться</span>
                        </button>
                    </div>
                </div>
            </main>
        </div>
    );
};