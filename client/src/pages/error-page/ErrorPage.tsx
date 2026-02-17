import { useNavigate, useRouteError, isRouteErrorResponse } from 'react-router-dom';
import { Footer } from '@/shared/ui/Footer';
import { Button } from '@/shared/ui/Button';
import { useEffect } from 'react';

const ErrorPage = () => {
    const navigate = useNavigate();
    const error = useRouteError();

    // Полезно для отладки: посмотрим в консоли, что за объект ошибки пришел
    useEffect(() => {
        console.error('ErrorPage caught an error:', error);
    }, [error]);

    // --- 1. База описаний ---
    const knownErrors: Record<string, { title: string; desc: React.ReactNode }> = {
        '404': {
            title: 'Поток прерван',
            desc: <>Запрашиваемая страница растворилась в цифровом эфире.<br />Ваше сознание находится вне известных координат.</>
        },
        '500': {
            title: 'Системный сбой',
            desc: <>Критическая ошибка ядра нейросети.<br />Процессы дестабилизированы.</>
        },
        '403': {
            title: 'Доступ запрещен',
            desc: <>Ваши нейронные ключи не подходят к этому сектору.<br />Протокол безопасности активен.</>
        }
        // ... остальные ошибки без изменений
    };

    // --- 2. Улучшенная логика определения ---
    let errorCode = '500'; 
    let technicalDetails = 'Unknown system anomaly.';

    if (isRouteErrorResponse(error)) {
        // Это стандартные ошибки роутинга (404, 401 и т.д.)
        errorCode = error.status.toString();
        technicalDetails = error.statusText || error.data?.message || 'Route Error';
    } else if (error && typeof error === 'object') {
        // Если это не RouteResponse, но у объекта есть статус (бывает при некоторых fetch-ошибках)
        if ('status' in error && typeof error.status === 'number') {
            errorCode = error.status.toString();
        }
        
        // Достаем сообщение
        if ('message' in error && typeof error.message === 'string') {
            technicalDetails = error.message;
        } else if ('data' in error && typeof error.data === 'string') {
            technicalDetails = error.data;
        }
    } else if (typeof error === 'string') {
        technicalDetails = error;
    }

    // Специальная проверка: если страница не найдена, но errorCode почему-то остался 500
    // (Иногда лоадеры выбрасывают ошибки, которые мы хотим трактовать как 404)
    if (technicalDetails.toLowerCase().includes('not found')) {
        errorCode = '404';
    }

    // --- 3. Рендеринг текста ---
    const knownError = knownErrors[errorCode];
    const displayTitle = knownError ? knownError.title : `Ошибка ${errorCode}`;
    const displayDesc = knownError ? knownError.desc : (
        <>Система зафиксировала аномалию {errorCode}.<br />Описание протокола отсутствует.</>
    );

    return (
        <div className="flex flex-col min-h-screen bg-deep-midnight font-sans selection:bg-bright-turquoise selection:text-deep-midnight overflow-hidden">
            <style>{`
                @keyframes float { 0%, 100% { transform: translateY(0px); } 50% { transform: translateY(-20px); } }
                @keyframes pulse-soft { 0%, 100% { opacity: 0.05; } 50% { opacity: 0.15; } }
                .animate-float { animation: float 6s ease-in-out infinite; }
                .animate-pulse-bg { animation: pulse-soft 4s ease-in-out infinite; }
            `}</style>

            <div className="fixed inset-0 overflow-hidden pointer-events-none z-0">
                <div className="absolute top-[-10%] left-[-10%] w-[50%] h-[50%] bg-bright-turquoise/10 blur-[120px] rounded-full animate-pulse" />
                <div className="absolute bottom-[-10%] right-[-10%] w-[50%] h-[50%] bg-light-mint/10 blur-[120px] rounded-full animate-pulse" />
            </div>

            <main className="flex-grow flex flex-col items-center justify-center relative z-10 px-6 py-20">
                <div className="max-w-4xl w-full text-center relative">
                    <h1 className="text-[15rem] md:text-[25rem] font-black leading-none text-white animate-pulse-bg select-none absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 blur-[2px] pointer-events-none tracking-tighter">
                        {errorCode}
                    </h1>

                    <div className="relative z-10">
                        <div className="inline-block px-4 py-1 rounded-full border border-bright-turquoise/30 bg-bright-turquoise/5 text-bright-turquoise font-mono text-sm mb-6 tracking-[0.2em] uppercase">
                            Critical_Status: {errorCode}
                        </div>

                        <h2 className="text-5xl md:text-8xl font-black mb-8 leading-tight">
                            <span className="bg-gradient-to-r from-bright-turquoise via-light-mint to-sky-blue bg-clip-text text-transparent animate-float inline-block">
                                {displayTitle}
                            </span>
                        </h2>

                        <p className="text-lg md:text-2xl text-text-secondary max-w-2xl mx-auto mb-12 leading-relaxed font-light backdrop-blur-sm">
                            {displayDesc}
                        </p>

                        <div className="max-w-lg mx-auto mb-12 bg-black/40 backdrop-blur-xl rounded-2xl border border-white/10 overflow-hidden shadow-2xl text-left font-mono">
                            <div className="bg-white/5 px-4 py-2 border-b border-white/10 flex justify-between items-center">
                                <span className="text-[10px] text-text-secondary uppercase tracking-widest">Error_Log_Stack</span>
                                <div className="flex gap-1">
                                    <div className="w-2 h-2 rounded-full bg-bright-turquoise/40" />
                                    <div className="w-2 h-2 rounded-full bg-light-mint/40" />
                                </div>
                            </div>
                            <div className="p-6 text-xs md:text-sm space-y-3">
                                <div className="flex gap-4">
                                    <span className="text-bright-turquoise opacity-50">STATUS:</span>
                                    <span className="text-white">{errorCode}</span>
                                </div>
                                <div className="flex gap-4">
                                    <span className="text-bright-turquoise opacity-50">SOURCE:</span>
                                    <span className="text-sky-blue">MindFlow_Router_v1.0</span>
                                </div>
                                <div className="flex gap-4">
                                    <span className="text-bright-turquoise opacity-50">DETAILS:</span>
                                    <span className="text-light-mint/80 italic line-clamp-2">{technicalDetails}</span>
                                </div>
                                <div className="pt-2 flex items-center gap-2 text-bright-turquoise">
                                    <span className="animate-pulse">●</span>
                                    <span className="animate-pulse">Awaiting_reconnection...</span>
                                </div>
                            </div>
                        </div>

                        <Button 
                            variant="primary" 
                            onClick={() => navigate('/')}
                            className="w-full md:w-auto"
                        >
                            Вернуться на главную
                        </Button>
                    </div>
                </div>
            </main>

            <div className="relative z-50">
                <Footer />
            </div>
        </div>
    );
};

export default ErrorPage;