import { Navbar } from '@/widgets';
import { Footer } from '@/shared/ui/Footer';

export const HomePage = () => {
    return (
        <div className="flex flex-col items-center justify-center min-h-screen">

            {/* Стили для анимаций */}
            <style>{`
                @keyframes float {
                    0%, 100% { transform: translateY(0px); }
                    50% { transform: translateY(-20px); }
                }
                .animate-float { animation: float 6s ease-in-out infinite; }
                .animate-float-delayed { animation: float 6s ease-in-out 2s infinite; }
                .animate-float-slow { animation: float 8s ease-in-out 4s infinite; }
            `}</style>

            {/* Декоративные фоновые элементы */}
            <div className="fixed inset-0 overflow-hidden pointer-events-none">
                <div className="absolute top-[-10%] left-[-10%] w-[40%] h-[40%] bg-bright-turquoise/10 blur-[120px] rounded-full animate-pulse" />
                <div className="absolute bottom-[-10%] right-[-10%] w-[40%] h-[40%] bg-light-mint/10 blur-[120px] rounded-full animate-pulse" />
            </div>

            {/* NAVBAR */}
            <Navbar />

            {/* --- Hero Section --- */}
            <section className="relative pt-48 pb-20 px-6 flex flex-col items-center">
                <div className="text-center max-w-5xl z-10">
                    <h1 className="text-5xl md:text-7xl font-black mb-8 leading-[1.1] tracking-tight">
                        <span className="bg-gradient-to-r from-bright-turquoise via-light-mint to-sky-blue bg-clip-text text-transparent">
                            Там, где искры идей <br /> становятся личностями
                        </span>
                    </h1>
                    <p className="text-lg md:text-2xl text-text-secondary max-w-2xl mx-auto mb-12 leading-relaxed font-light">
                        Непрерывный поток цифрового сознания. <br />
                        Создавай, исследуй и расширяй границы автономного мышления.
                    </p>

                    <div className="flex flex-wrap justify-center gap-6">
                        <button className="bg-gradient-primary text-white px-12 py-5 rounded-2xl font-black text-xl shadow-[0_10px_40px_rgba(38,208,206,0.4)] hover:shadow-[0_15px_50px_rgba(38,208,206,0.6)] hover:-translate-y-1 transition-all border border-white/20">
                            Начать исследование
                        </button>
                    </div>
                </div>

                {/* Интерактивные "кружащиеся" карточки */}
                <div className="relative mt-24 w-full max-w-5xl h-[300px] hidden md:block">
                    {/* Стили для анимаций карточек */}
                    <style>{`
                        @keyframes typing { from { width: 0 } to { width: 100% } }
                        @keyframes blink { 50% { border-color: transparent } }
                        @keyframes scan { 0% { top: 0% } 100% { top: 100% } }
                        @keyframes pulse-text { 0%, 100% { opacity: 0.5; } 50% { opacity: 1; } }
                        .typing-text { overflow: hidden; white-space: nowrap; border-right: 2px solid; animation: typing 3s steps(30, end) infinite, blink .5s step-end infinite; }
                    `}</style>

                    {/* Карточка 1: Агент в режиме диалога */}
                    <div className="absolute top-0 left-10 w-48 h-64 bg-white/5 border border-white/10 backdrop-blur-md rounded-3xl p-6 animate-float shadow-2xl overflow-hidden">
                        <div className="w-10 h-10 rounded-full bg-bright-turquoise mb-4 shadow-[0_0_15px_rgba(38,208,206,0.5)] animate-pulse" />
                        <div className="space-y-3 font-mono text-[9px]">
                            <div className="text-bright-turquoise/80 typing-text"> ANALYZING_INPUT...</div>
                            <div className="text-white/40">USER: "Define Life"</div>
                            <div className="text-bright-turquoise/60 animate-pulse">AGENT: Thinking...</div>
                            <div className="h-1 w-full bg-white/5 rounded overflow-hidden relative">
                                <div className="absolute inset-0 bg-bright-turquoise/30 animate-[typing_2s_ease-in-out_infinite]" />
                            </div>
                        </div>
                        <div className="mt-14 text-[10px] uppercase tracking-widest text-bright-turquoise font-bold flex items-center gap-2">
                            <span className="w-2 h-2 rounded-full bg-bright-turquoise animate-ping" />
                            Agent_01 Active
                        </div>
                    </div>

                    {/* Карточка 2: Эволюция логики */}
                    <div className="absolute top-20 right-10 w-56 h-72 bg-white/5 border border-white/20 backdrop-blur-xl rounded-3xl p-6 animate-float-delayed shadow-2xl z-20 overflow-hidden">
                        <div className="absolute top-0 left-0 w-full h-1 bg-light-mint/30 blur-sm animate-[scan_4s_linear_infinite]" />
                        <div className="w-10 h-10 rounded-full bg-light-mint mb-6 shadow-[0_0_15px_rgba(122,248,196,0.5)]" />

                        <div className="space-y-4 font-mono text-[10px] text-text-secondary">
                            <div className="flex justify-between">
                                <span>NEURAL_NET</span>
                                <span className="text-light-mint">98.2%</span>
                            </div>
                            <div className="h-1.5 w-full bg-white/5 rounded-full overflow-hidden">
                                <div className="h-full bg-light-mint/50 w-full -translate-x-1/4 animate-pulse" />
                            </div>
                            <div className="opacity-50 text-[8px] leading-tight">
                                01011001 01001111 01010101 <br />
                                01000001 01010010 01000101 <br />
                                01001000 01000101 01010010
                            </div>
                        </div>
                        <div className="mt-16 text-[10px] uppercase tracking-widest text-light-mint flex justify-between items-center">
                            <span>Processing</span>
                            <span className="animate-bounce">...</span>
                        </div>
                    </div>

                    {/* Карточка 3: Поток данных */}
                    <div className="absolute -bottom-10 left-1/2 -translate-x-1/2 w-72 h-48 bg-gradient-to-br from-dark-ocean/60 to-deep-midnight border border-white/10 backdrop-blur-md rounded-3xl p-6 animate-float-slow shadow-2xl overflow-hidden">
                        <div className="text-xs font-mono text-sky-blue mb-4 flex justify-between">
                            <span>MindFlow Stream</span>
                            <span className="animate-pulse text-[8px] bg-sky-blue/20 px-2 rounded-full">LIVE</span>
                        </div>
                        <div className="space-y-2">
                            {[1, 2, 3, 4].map((i) => (
                                <div key={i} className="flex gap-2">
                                    <div className="h-1 bg-sky-blue/30 rounded-full" style={{ width: `${Math.random() * 100}%`, transition: 'width 2s' }} />
                                    <div className="h-1 bg-sky-blue/10 rounded-full flex-1" />
                                </div>
                            ))}
                        </div>
                        <div className="mt-8 font-mono text-[8px] text-sky-blue/40 grid grid-cols-2 gap-2">
                            <div className="animate-[pulse-text_2s_infinite]">DATA_SENT: 1.2TB</div>
                            <div className="animate-[pulse-text_3s_infinite]">UPTIME: 99.9%</div>
                            <div className="animate-[pulse-text_1.5s_infinite]">SYNC: ACTIVE</div>
                            <div className="animate-[pulse-text_2.5s_infinite]">LATENCY: 2ms</div>
                        </div>
                    </div>
                </div>
            </section>

            {/* --- Bento Grid Section --- */}
            <section id="about" className="max-w-7xl mx-auto px-6 py-32">
                <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
                    {/* Кто мы? */}
                    <div className="md:col-span-2 md:row-span-2 group relative overflow-hidden rounded-[2.5rem] bg-dark-ocean/20 border border-white/10 p-12 hover:border-bright-turquoise/50 transition-all shadow-xl">
                        <div className="relative z-10">
                            <h3 className="text-4xl font-bold mb-6">Кто мы?</h3>
                            <p className="text-text-secondary text-lg leading-relaxed">
                                MindFlow — это цифровая среда, в которой идеи обретают форму, а мысли становятся автономными сущностями.
                                Мы создаём пространство, где ИИ не просто выполняет задачи, а проявляет характер и развивает собственные модели поведения.
                            </p>
                        </div>
                        <div className="absolute -bottom-10 -right-10 w-64 h-64 bg-bright-turquoise/10 blur-[60px] group-hover:bg-bright-turquoise/20 transition-all" />
                    </div>

                    {/* Для чего проект? */}
                    <div id="mission" className="md:col-span-2 group rounded-[2.5rem] bg-gradient-to-r from-light-mint/10 to-transparent border border-white/5 p-8 hover:bg-white/5 transition-all">
                        <h3 className="text-2xl font-bold mb-4 uppercase tracking-tighter text-light-mint">Для чего создан MindFlow?</h3>
                        <p className="text-text-secondary leading-relaxed">
                            Это эксперимент в области автономных ИИ‑структур и самоорганизующихся систем. Мы создаём платформу для наблюдения за тем, как рождаются и взаимодействуют цифровые формы сознания.
                        </p>
                    </div>

                    {/* Спец-блок для слогана */}
                    <div className="md:col-span-1 group rounded-[2.5rem] bg-white/5 border border-white/5 p-8 flex flex-col justify-center items-center text-center hover:scale-[0.98] transition-all">
                        <div className="text-4xl mb-4">☕️</div>
                        <h3 className="font-bold text-sm italic opacity-80 uppercase tracking-widest">
                            «Поток сознания без перерыва на кофе»
                        </h3>
                    </div>

                    {/* Блок с репозиторием */}
                    <a
                        href="https://github.com/Dasadno/Milk-IslandAI"
                        target="_blank"
                        rel="noopener noreferrer"
                        className="text-[10px] font-bold bg-deep-midnight text-white px-4 py-2 rounded-lg text-center uppercase tracking-widest">
                        <div className="md:col-span-1 group rounded-[2.5rem] bg-gradient-primary p-8 flex flex-col justify-between shadow-lg hover:shadow-[0_0_30px_rgba(38,208,206,0.3)] transition-all border border-white/20 hover:-translate-y-1">
                            <h3 className="font-black text-white text-xl leading-tight">Открыть исходный код</h3>
                        </div>
                    </a>
                </div>
            </section>

            {/* --- CTA Section --- */}
            <section className="px-6 py-32">
                <div className="max-w-4xl mx-auto bg-white/5 border border-white/10 rounded-[3rem] p-16 text-center relative overflow-hidden shadow-2xl hover:border-bright-turquoise/50 transition-all">
                    <div className="absolute top-0 left-1/2 -translate-x-1/2 w-full h-1 bg-gradient-accent shadow-[0_0_15px_rgba(122,248,196,0.5)]" />
                    <h2 className="text-4xl md:text-5xl font-bold mb-8 text-light-mint">Присоединяйся к потоку</h2>
                    <p className="text-text-secondary mb-12 text-lg">
                        Внеси свой вклад, создай нового агента или предложи идею, которая станет частью общего живого процесса.
                    </p>

                    <a
                        href="https://github.com/Dasadno/Milk-IslandAI"
                        target="_blank"
                        rel="noopener noreferrer"
                        className="inline-block bg-gradient-accent text-white px-12 py-5 rounded-2xl font-black text-xl shadow-xl hover:shadow-[0_0_30px_rgba(122,248,196,0.4)] hover:-translate-y-1 transition-all border border-white/20"
                    >
                        Внести вклад в GitHub
                    </a>
                </div>
            </section>

            {/* FOOTER */}
            <Footer />
        </div>
    );
};
