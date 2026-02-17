import type { Meta, StoryObj } from '@storybook/react';
import { Settings, User, Bell, Search } from 'lucide-react';
import { IconButton } from '@/shared/ui/IconButton';

const meta: Meta<typeof IconButton> = {
    title: 'Shared/IconButton',
    component: IconButton,
    tags: ['autodocs'],
};

export default meta;
type Story = StoryObj<typeof IconButton>;

export const SettingsIcon: Story = {
    args: {
        icon: Settings,
    },
};

export const UserIcon: Story = {
    args: {
        icon: User,
    },
};

export const BellIcon: Story = {
    args: {
        icon: Bell,
    },
};

export const SearchIcon: Story = {
    args: {
        icon: Search,
    },
};

export const WithClick: Story = {
    args: {
        icon: Settings,
        onClick: () => alert('Button clicked!'),
    },
};