import {create} from "zustand";
import {constants} from "../utils/constants";
import {getRawJWT, type IAuthToken, parseJWT} from "../utils/token";
import {useLoggedUserStore} from "./LoggedUserStore";

interface TokenStoreState {
  tokenRaw: string | null;
  refreshTokenRaw: string | null;
  token: () => IAuthToken | null;
  setTokens: (newToken: string, newRefreshToken: string) => void;
  resetTokens: () => void;
  isAuthenticated: () => boolean;
  isAuthenticatedAndLoaded: () => boolean;
}

export const useAuthTokenStore = create<TokenStoreState>()((set, get) => ({
  tokenRaw: getRawJWT(),
  refreshTokenRaw: localStorage.getItem(constants.LOCAL_STORAGE_REFRESH_TOKEN_KEY),
  token: () => parseJWT(get().tokenRaw),
  setTokens: (newToken, newRefreshToken) => {
    localStorage.setItem(constants.LOCAL_STORAGE_TOKEN_KEY, newToken);
    localStorage.setItem(constants.LOCAL_STORAGE_REFRESH_TOKEN_KEY, newRefreshToken);
    set(() => ({tokenRaw: newToken, refreshTokenRaw: newRefreshToken}));
  },
  resetTokens: () => {
    localStorage.removeItem(constants.LOCAL_STORAGE_TOKEN_KEY);
    localStorage.removeItem(constants.LOCAL_STORAGE_REFRESH_TOKEN_KEY);
    set(() => ({tokenRaw: null, refreshTokenRaw: null}));
  },
  isAuthenticated: () => {
    const tokenParsed = get().token();
    if (!tokenParsed) return false;
    return new Date(tokenParsed.exp * 1000) > new Date();
  },
  isAuthenticatedAndLoaded: () => get().isAuthenticated() && !!useLoggedUserStore.getState().user,
}));
