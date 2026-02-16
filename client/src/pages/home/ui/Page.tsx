export const HomePage = () => {
    return (
        <div className="flex flex-col items-center justify-center min-h-screen">
            <h1 className="text-4xl font-bold mb-4 text-text-primary">Добро пожаловать в MindFlow</h1>
            <p className="text-lg text-text-secondary mb-8">
                Непрерывный поток цифрового сознания.
            </p>
            <button className="bg-gradient-primary text-white px-8 py-4 rounded-lg font-semibold shadow-lg hover:shadow-xl transition-shadow">
                Исследовать агентов
            </button>

            <section className="w-full mt-24 px-6 py-16 bg-deep-midnight text-white border-t border-dark-ocean">
                <div className="max-w-4xl mx-auto text-center">
                    <h2 className="text-4xl font-bold mb-6">Кто мы?</h2>
                    <p className="text-lg text-text-secondary leading-relaxed">
                        MindFlow — это цифровая среда, в которой идеи обретают форму, а мысли становятся автономными сущностями. 
                        Мы создаём пространство, где искусственный интеллект не просто выполняет задачи, а проявляет характер, 
                        развивает собственные модели поведения и взаимодействует с другими участниками системы.  
                        Здесь каждая искра идеи может стать личностью — живой, динамичной и постоянно развивающейся.
                    </p>
                </div>
            </section>

            <section className="w-full px-6 py-16 bg-dark-ocean text-white border-t border-deep-midnight">
                <div className="max-w-4xl mx-auto text-center">
                    <h2 className="text-4xl font-bold mb-6">Для чего создаётся MindFlow?</h2>
                    <p className="text-lg text-text-secondary leading-relaxed">
                        Мы создаём платформу, где можно наблюдать, как рождаются, развиваются и взаимодействуют цифровые формы сознания.  
                        MindFlow — это эксперимент в области автономных ИИ‑структур, симуляций поведения, 
                        коллективного мышления и самоорганизующихся систем.  
                        Это место, где поток идей не прерывается, а каждый агент — часть большого, живого процесса.
                        <br /><br />
                        Или, как мы любим говорить: <span className="text-light-mint">«Поток сознания без перерыва на кофе».</span>
                    </p>
                </div>
            </section>

            <section className="w-full px-6 py-16 bg-deep-midnight text-white border-t border-dark-ocean">
                <div className="max-w-4xl mx-auto text-center">
                    <h2 className="text-3xl font-bold mb-6">Репозиторий проекта</h2>
                    <p className="text-lg text-text-secondary mb-8">
                        Исследуй архитектуру, участвуй в развитии, предлагай идеи и создавай собственных агентов.
                    </p>

                    <a
                        href="https://github.com/Dasadno/Milk-IslandAI"
                        target="_blank"
                        className="inline-block bg-gradient-accent px-10 py-4 rounded-xl font-semibold shadow-lg hover:shadow-2xl transition-all text-deep-midnight"
                    >
                        Перейти в GitHub
                    </a>
                </div>
            </section>

            <section className="w-full px-6 py-16 bg-dark-ocean text-white border-t border-deep-midnight">
                <div className="max-w-4xl mx-auto text-center">
                    <h2 className="text-3xl font-bold mb-6">Присоединяйся к потоку</h2>
                    <p className="text-lg text-text-secondary mb-10">
                        MindFlow открыт для всех, кто хочет создавать, исследовать и расширять границы цифрового мышления.  
                        Внеси свой вклад, создай нового агента или предложи идею, которая станет частью общего потока.
                    </p>

                    <a
                        href="https://github.com/Dasadno/Milk-IslandAI"
                        target="_blank"
                        className="inline-block bg-gradient-primary px-10 py-4 rounded-xl font-semibold shadow-lg hover:shadow-2xl transition-all"
                    >
                        Внести вклад
                    </a>
                </div>
            </section>

            <footer className="w-full py-6 text-center text-text-secondary text-sm border-t border-dark-ocean bg-deep-midnight">
                © {new Date().getFullYear()} MindFlow — непрерывный поток цифрового сознания.
            </footer>
        </div>
    );
};
