import "./index.css";
import "./utils/i18n.ts";
import {QueryClient, QueryClientProvider} from "@tanstack/react-query";
import {StrictMode} from "react";
import {createRoot} from "react-dom/client";
import {BrowserRouter, Route, Routes} from "react-router";
import {ErrorBoundary} from "./components/ErrorBoundary.tsx";
import {I18nProvider} from "./components/I18nProvider.tsx";
import {LoaderProvider} from "./components/LoadingProvider.tsx";
import {LoggedUserProvider} from "./components/LoggedUserProvider.tsx";
import {ThemeProvider} from "./components/ThemeProvider.tsx";
import Home from "./pages/Home/Home.tsx";
import Login from "./pages/Login/Login.tsx";
import OAuthRedirect from "./pages/Login/OauthRedirect.tsx";
import Logout from "./pages/Logout/Logout.tsx";
import Settings from "./pages/Settings/Settings.tsx";
import {routes} from "./utils/routes.ts";

const queryClient = new QueryClient();
const rootElement = document.getElementById("root");
if (rootElement) {
  createRoot(rootElement).render(
    <StrictMode>
      <ErrorBoundary>
        <QueryClientProvider client={queryClient}>
          <BrowserRouter>
            <ThemeProvider>
              <I18nProvider>
                <LoaderProvider>
                  <LoggedUserProvider>
                    <Routes>
                      <Route path={routes.login} element={<Login />} />
                      <Route path={routes.loginOAuthRedirect} element={<OAuthRedirect />} />
                      <Route path={routes.logout} element={<Logout />} />

                      <Route path={routes.index} element={<Home />} />
                      <Route path={routes.settings} element={<Settings />} />
                    </Routes>
                  </LoggedUserProvider>
                </LoaderProvider>
              </I18nProvider>
            </ThemeProvider>
          </BrowserRouter>
        </QueryClientProvider>
      </ErrorBoundary>
    </StrictMode>,
  );
}
