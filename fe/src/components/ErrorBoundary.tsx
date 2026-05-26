import {Component, type ErrorInfo, type ReactNode} from "react";

interface Props {
  children?: ReactNode;
}

interface State {
  hasError: boolean;
  error?: Error;
}

export class ErrorBoundary extends Component<Props, State> {
  public state: State = {
    hasError: false,
  };

  public static getDerivedStateFromError(error: Error): State {
    return {hasError: true, error};
  }

  public componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    console.error("Uncaught error:", error, errorInfo);
  }

  public render() {
    if (this.state.hasError) {
      return (
        <div className="flex min-h-screen w-full flex-col items-center justify-center bg-gray-50 p-4 text-center dark:bg-gray-900">
          <h1 className="mb-4 text-3xl font-bold text-red-600 dark:text-red-400">Oops, something went wrong.</h1>
          <p className="mb-8 text-gray-600 dark:text-gray-400">We apologize for the inconvenience. Please try refreshing the page.</p>
          <button onClick={() => window.location.reload()} className="rounded bg-blue-600 px-4 py-2 font-semibold text-white transition hover:bg-blue-700">
            Refresh Page
          </button>
          {this.state.error && (
            <pre className="mt-8 max-w-2xl overflow-auto rounded bg-gray-200 p-4 text-left text-xs text-gray-800 dark:bg-gray-800 dark:text-gray-300">
              {this.state.error.toString()}
            </pre>
          )}
        </div>
      );
    }

    return this.props.children;
  }
}
