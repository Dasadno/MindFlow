import React, { useState } from 'react';
import { X, Sparkles, Brain, MessageSquare, Terminal } from 'lucide-react';
import { Button } from '../shared/ui';
import Input from '../shared/ui/Input';

interface NewAgentPopUpProps {
    isOpen: boolean;
    onClose: () => void;
    onCreate?: (agentData: any) => void;
}

const PERSONALITY_TYPES = [
    { id: 'analytical', label: 'Analytical', icon: Terminal, desc: 'Logical and data-driven' },
    { id: 'creative', label: 'Creative', icon: Sparkles, desc: 'Imaginative and artistic' },
    { id: 'empathetic', label: 'Empathetic', icon: MessageSquare, desc: 'Understanding and supportive' },
    { id: 'strategic', label: 'Strategic', icon: Brain, desc: 'Planning and goal-oriented' }
];

export const NewAgentPopUp = ({ isOpen, onClose, onCreate }: NewAgentPopUpProps) => {
    const [name, setName] = useState('');
    const [description, setDescription] = useState('');
    const [selectedPersonality, setSelectedPersonality] = useState(PERSONALITY_TYPES[0].id);

    if (!isOpen) return null;

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        console.log('Creating agent:', { name, description, personality: selectedPersonality });

        if (onCreate) {
            onCreate({ name, description, personality: selectedPersonality });
        }

        setName('');
        setDescription('');
        setSelectedPersonality(PERSONALITY_TYPES[0].id);
        onClose();
    };

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
            {/* Backdrop */}
            <div
                className="absolute inset-0 bg-deep-midnight/80 backdrop-blur-sm transition-opacity"
                onClick={onClose}
            />

            {/* Modal */}
            <div className="
                relative w-full max-w-lg 
                bg-deep-midnight 
                border border-white/10 
                rounded-3xl 
                shadow-[0_0_50px_rgba(0,0,0,0.5)]
                overflow-hidden
                animate-in fade-in zoom-in-95 duration-200
            ">
                {/* Header */}
                <div className="flex items-center justify-between p-6 border-b border-white/5 bg-white/[0.02]">
                    <div className="flex items-center gap-3">
                        <div className="w-10 h-10 rounded-xl bg-bright-turquoise/10 flex items-center justify-center text-bright-turquoise">
                            <Sparkles className="w-5 h-5" />
                        </div>
                        <div>
                            <h2 className="text-xl font-bold text-white">Initialising New Agent...</h2>
                            <p className="text-white/40 text-xs">Configure personality matrix</p>
                        </div>
                    </div>
                    <button
                        onClick={onClose}
                        className="p-2 rounded-xl hover:bg-white/5 text-white/40 hover:text-white transition-colors"
                    >
                        <X className="w-5 h-5" />
                    </button>
                </div>

                {/* Form */}
                <form onSubmit={handleSubmit} className="p-6 space-y-6">
                    {/* Name Input */}
                    <div className="space-y-4">
                        <Input
                            label="Sys.Identity"
                            placeholder="Enter agent name..."
                            value={name}
                            onChange={(e) => setName(e.target.value)}
                            className="bg-deep-midnight"
                            autoFocus
                        />
                    </div>

                    {/* Personality Selection */}
                    <div className="space-y-3">
                        <label className="block text-[10px] uppercase tracking-[0.2em] text-bright-turquoise/60 ml-1 font-mono">
                            {'> Sys.Personality_Core'}
                        </label>
                        <div className="grid grid-cols-2 gap-3">
                            {PERSONALITY_TYPES.map((type) => {
                                const Icon = type.icon;
                                const isSelected = selectedPersonality === type.id;
                                return (
                                    <div
                                        key={type.id}
                                        onClick={() => setSelectedPersonality(type.id)}
                                        className={`
                                            cursor-pointer p-4 rounded-xl border transition-all duration-300
                                            ${isSelected
                                                ? 'bg-bright-turquoise/10 border-bright-turquoise/40 shadow-[0_0_20px_rgba(38,208,206,0.1)]'
                                                : 'bg-white/5 border-white/5 hover:bg-white/10 hover:border-white/10'
                                            }
                                        `}
                                    >
                                        <div className={`flex items-center gap-2 mb-2 ${isSelected ? 'text-bright-turquoise' : 'text-white/60'}`}>
                                            <Icon className="w-4 h-4" />
                                            <span className="font-bold text-sm">{type.label}</span>
                                        </div>
                                        <p className="text-xs text-white/40">{type.desc}</p>
                                    </div>
                                );
                            })}
                        </div>
                    </div>

                    {/* Description/Prompt */}
                    <div className="space-y-2 group">
                        <label className="block text-[10px] uppercase tracking-[0.2em] text-bright-turquoise/60 ml-1 font-mono group-focus-within:text-bright-turquoise transition-colors">
                            {'> Sys.Prime_Directive'}
                        </label>
                        <div className="relative">
                            <textarea
                                value={description}
                                onChange={(e) => setDescription(e.target.value)}
                                placeholder="Describe the agent's purpose and behavior..."
                                className="
                                    w-full h-32 px-6 py-4 rounded-2xl font-medium outline-none resize-none
                                    bg-white/5 border border-white/10
                                    text-white placeholder:text-white/20
                                    focus:border-bright-turquoise/40 focus:bg-white/[0.08]
                                    transition-all duration-300
                                "
                            />
                            {/* Decorative corner */}
                            <div className="absolute -right-1 -bottom-1 w-2 h-2 border-r border-b border-white/20 group-focus-within:border-bright-turquoise transition-colors" />
                        </div>
                    </div>

                    {/* Actions */}
                    <div className="pt-4 flex items-center justify-end gap-3 border-t border-white/5">
                        <Button
                            variant="secondary"
                            onClick={onClose}
                            className="w-32"
                        >
                            Cancel
                        </Button>
                        <Button
                            type="submit"
                            variant="primary"
                            className="w-40"
                            disabled={!name.trim()}
                        >
                            Create Agent
                        </Button>
                    </div>
                </form>
            </div>
        </div>
    );
};
