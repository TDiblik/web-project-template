import type {TranslationPosibilitiesType} from "@shared/api-client";
import type {i18n as i18nType} from "i18next";
import {create} from "zustand";
import {UserController} from "../utils/api";
import {constants} from "../utils/constants";
import {i18n} from "../utils/i18n";
import {useLoggedUserStore} from "./LoggedUserStore";
import {useAuthTokenStore} from "./TokenStore";

interface i18nStoreState {
  i18n: i18nType;
  changeLanguage: (newLanguage: TranslationPosibilitiesType) => void;
}
export const usei18nStore = create<i18nStoreState>()(() => ({
  i18n: i18n,
  changeLanguage: (newLanguage) => {
    i18n.changeLanguage(newLanguage);
    localStorage.setItem(constants.LOCAL_STORAGE_LOCALIZATION_KEY, newLanguage);
    if (useAuthTokenStore.getState().isAuthenticatedAndLoaded()) {
      const userStore = useLoggedUserStore.getState();
      if (userStore.user) {
        userStore.setUser({...userStore.user, preferedLanguage: newLanguage});
      }
      UserController.apiV1PrivateUserMePatch({
        githubComTDiblikProjectTemplateApiHandlersPatchUserMeHandlerRequest: {
          preferedLanguage: newLanguage,
        },
      });
    }
  },
}));
