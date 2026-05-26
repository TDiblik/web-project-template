import i18n from "i18next";
import LanguageDetector from "i18next-browser-languagedetector";
import {initReactI18next, type useTranslation} from "react-i18next";
import csTranslation from "../locales/cs/translation.json";
import enTranslation from "../locales/en/translation.json";
import {constants} from "./constants";

export type TranslateFn = ReturnType<typeof useTranslation>["t"];
const resources = {
  en: {
    translation: enTranslation,
  },
  cs: {
    translation: csTranslation,
  },
};

i18n
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    resources,
    load: "languageOnly",
    debug: constants.DEBUG,
    lng: (localStorage.getItem(constants.LOCAL_STORAGE_LOCALIZATION_KEY) ?? "").split("-")[0],
    fallbackLng: constants.DEFAULT_FALLBACK_LANGUAGE,
    detection: {
      lookupLocalStorage: constants.LOCAL_STORAGE_LOCALIZATION_KEY,
    },
  });
