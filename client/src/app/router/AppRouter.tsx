import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import { HomePage } from '@/pages/home/ui/Page';
import { LoginPage } from '@/pages/auth/ui/LoginPage';

const router = createBrowserRouter([
    {
        path: '/',
        element: <HomePage />,
    },
    {
        path: '/login',
        element: <LoginPage />,
    },
]);

export const AppRouter = () => {
    return <RouterProvider router={router} />;
};
