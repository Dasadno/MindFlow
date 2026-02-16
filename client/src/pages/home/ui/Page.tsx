export const HomePage = () => {
    return (
        <div className="flex flex-col items-center justify-center min-h-screen">
            <h1 className="text-4xl font-bold mb-4 text-text-primary">Welcome to Milk Island AI</h1>
            <p className="text-lg text-text-secondary mb-8">
                A society of autonomous agents.
            </p>
            <button className="bg-gradient-primary text-white px-8 py-4 rounded-lg font-semibold shadow-lg hover:shadow-xl transition-shadow">
                Explore Agents
            </button>
        </div>
    );
};
