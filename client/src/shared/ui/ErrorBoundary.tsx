import { Component, type ReactNode } from "react";

// Пропс чтобы принимал дочерний компонент
interface ErrorBoundaryProps {
  children: ReactNode;
}

// Пропс для состояния ошибки
interface ErrorBoundaryState {
  error: Error | null;
}

class ErrorBoundary extends Component<ErrorBoundaryProps, ErrorBoundaryState> {
  state: ErrorBoundaryState = {
    error: null,
  };
  static getDerivedStateFromError(error: unknown): ErrorBoundaryState {
    if (error instanceof Error) {
      return { error };
    }

    return { error: new Error(String(error)) };
  }

  render(): ReactNode {
    const { error } = this.state;

    if (error) {
      return (
        <div>
          <p>Seems like an error occured!</p>
          <p>{error.message}</p>
        </div>
      );
    }
    return this.props.children;
  }
}

export default ErrorBoundary;
