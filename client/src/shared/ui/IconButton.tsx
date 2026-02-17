import type { LucideIcon } from 'lucide-react';

interface IconButtonProps {
    icon: LucideIcon;
    onClick?: () => void;
}

export const IconButton = ({ icon: Icon, onClick }: IconButtonProps) => {
    return (
    <button
        onClick={onClick}
        className="
        p-2
        text-text-secondary
        hover:text-text-primary
        hover:bg-deep-midnight/50
        rounded-lg
        transition-colors"
    >
        <Icon className="w-5 h-5" />
    </button>
    );
};