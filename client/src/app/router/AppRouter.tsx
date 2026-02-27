import { createBrowserRouter, RouterProvider } from "react-router-dom";
import { HomePage } from "@/pages/home/ui/Page";
import { LoginPage } from "@/pages/auth/ui/LoginPage";
import { ChatPage } from "@/pages/chat/ui/Page";
import ErrorPage from "@/pages/error-page/ErrorPage";
import RegisterPage from "@/pages/auth/ui/RegisterPage";
import { ProtectedRoute } from "./ProtectedRoute";
import RelationShipGraph from "@/features/relationship-graph/ui/Graph";

const router = createBrowserRouter([
  {
    path: "/",
    element: <HomePage />,
  },
  {
    path: "/login",
    element: <LoginPage />,
  },
  {
    path: "/register",
    element: <RegisterPage />,
  },
  {
    element: <ProtectedRoute />,
    children: [
      {
        path: "/chat",
        element: <ChatPage />,
      },
      {
        path: "/relationship-graph",
        element: <RelationShipGraph />,
      },
    ],
  },
  {
    path: "*",
    element: <ErrorPage />,
  },
]);

export const AppRouter = () => {
  return <RouterProvider router={router} />;
};
