import {Configuration, getAuthController, getoAuthController, getoAuthRedirectController, getUserController} from "@shared/api-client";
import {useAuthTokenStore} from "../stores/TokenStore";
import {constants} from "./constants";

let isRefreshing = false;
let refreshSubscribers: ((token: string) => void)[] = [];

const subscribeTokenRefresh = (cb: (token: string) => void) => {
  refreshSubscribers.push(cb);
};

const onRefreshed = (token: string) => {
  refreshSubscribers.forEach((cb) => {
    cb(token);
  });
  refreshSubscribers = [];
};

const config = new Configuration({
  basePath: constants.API_BASE_PATH,
  apiKey: () => useAuthTokenStore.getState().tokenRaw ?? "",
  middleware: [
    {
      post: async (context) => {
        if (context.response.status === 401) {
          const state = useAuthTokenStore.getState();
          if (!state.tokenRaw) return context.response;

          // If a refresh is already in progress, wait for it to complete
          if (isRefreshing) {
            return new Promise((resolve) => {
              subscribeTokenRefresh((newToken) => {
                const newInit = {...context.init};
                if (newInit.headers) {
                  const headers = new Headers(newInit.headers);
                  headers.set(constants.TOKEN_HEADER_KEY, newToken);
                  newInit.headers = headers;
                }
                resolve(context.fetch(context.url, newInit));
              });
            });
          }

          const refreshToken = state.refreshTokenRaw;
          if (!refreshToken) {
            state.resetTokens();
            return context.response;
          }

          isRefreshing = true;
          try {
            const refreshRes = await fetch(`${constants.API_BASE_PATH}/api/v1/public/auth/refresh`, {
              method: "POST",
              headers: {"Content-Type": "application/json"},
              body: JSON.stringify({refresh_token: refreshToken}),
            });

            if (!refreshRes.ok) {
              state.resetTokens();
              onRefreshed("");
              return context.response;
            }

            const data = await refreshRes.json();
            state.setTokens(data.auth_token, data.refresh_token);
            onRefreshed(data.auth_token);

            const newInit = {...context.init};
            if (newInit.headers) {
              const headers = new Headers(newInit.headers);
              headers.set(constants.TOKEN_HEADER_KEY, data.auth_token);
              newInit.headers = headers;
            }
            return context.fetch(context.url, newInit);
          } catch (_err) {
            state.resetTokens();
            onRefreshed("");
            return context.response;
          } finally {
            isRefreshing = false;
          }
        }
        return context.response;
      },
    },
  ],
});

export const AuthController = getAuthController(config);
export const oAuthController = getoAuthController(config);
export const oAuthRedirectController = getoAuthRedirectController(config);
export const UserController = getUserController(config);
