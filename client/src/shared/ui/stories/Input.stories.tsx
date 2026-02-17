import type { Meta, StoryObj } from '@storybook/react';
import Input from '../Input';

const meta: Meta<typeof Input> = {
    title: 'Shared/Input',
    component: Input,
    tags: ['autodocs'],
    argTypes: {
        type: {
            control: { type: 'select' },
            options: ['text', 'password', 'email', 'number'],
        },
        disabled: { control: 'boolean' },
    },
};

export default meta;
type Story = StoryObj<typeof Input>;

export const Default: Story = {
    args: {
        placeholder: 'Enter text...',
    },
};

export const WithValue: Story = {
    args: {
        value: 'Some value',
        onChange: () => { },
    },
};

export const Password: Story = {
    args: {
        type: 'password',
        placeholder: 'Enter password',
    },
};

export const Disabled: Story = {
    args: {
        disabled: true,
        placeholder: 'Disabled input',
        value: 'Cannot edit this',
    },
};

export const Email: Story = {
    args: {
        type: 'email',
        placeholder: 'user@example.com',
    },
};
