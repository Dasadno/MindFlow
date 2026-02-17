// import type { Meta, StoryObj } from '@storybook/react';
// import { Button } from '../Button';

// const meta: Meta<typeof Button> = {
//     title: 'Shared/Button',
//     component: Button,
//     tags: ['autodocs'],
//     argTypes: {
//         variant: {
//             control: { type: 'select' },
//             options: ['primary', 'secondary', 'gradient'],
//         },
//         onClick: { action: 'clicked' },
//     },
// };

// export default meta;
// type Story = StoryObj<typeof Button>;

// export const Primary: Story = {
//     args: {
//         children: 'Primary Button',
//         variant: 'primary',
//     },
// };

// export const Secondary: Story = {
//     args: {
//         children: 'Secondary Button',
//         variant: 'secondary',
//     },
// };

// export const Gradient: Story = {
//     args: {
//         children: 'Gradient Button',
//         variant: 'gradient',
//     },
// };

// export const WithCustomClass: Story = {
//     args: {
//         children: 'Custom Class',
//         variant: 'primary',
//         className: 'uppercase tracking-wider',
//     },
// };


import type { Meta, StoryObj } from '@storybook/react';
import { Button } from '../Button';

const meta: Meta<typeof Button> = {
    title: 'Shared/Button',
    component: Button,
    parameters: {
        layout: 'centered',
        // Ставим темный фон как в приложении
        backgrounds: {
            default: 'mindflow-dark',
            values: [
                { name: 'mindflow-dark', value: '#0B1E3B' },
                { name: 'light', value: '#ffffff' },
            ],
        },
    },
    tags: ['autodocs'],
};

export default meta;
type Story = StoryObj<typeof Button>;

export const Primary: Story = {
    args: {
        children: 'Primary Action',
        variant: 'primary',
    },
};

export const SecondaryGlass: Story = {
    args: {
        children: 'Secondary Action',
        variant: 'secondary',
    },
};

export const GradientFlow: Story = {
    args: {
        children: 'Start Flow',
        variant: 'gradient',
    },
};

export const AccentMint: Story = {
    args: {
        children: 'Confirm Action',
        variant: 'accent',
    },
};