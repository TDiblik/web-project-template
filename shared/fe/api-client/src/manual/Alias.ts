import {
  RedirectBackToAfterOauth,
  GithubComTDiblikProjectTemplateApiHandlersOauthRedirectHandlerResponse,
  PreferedLanguage,
  PreferedTheme,
  GithubComTDiblikProjectTemplateApiModelsUserModelDB,
} from "../generated";

export type OauthRedirectHandlerResponse = GithubComTDiblikProjectTemplateApiHandlersOauthRedirectHandlerResponse;
export type OauthRedirectHandlerRequest = Promise<OauthRedirectHandlerResponse>;
export type RedirectBackToAfterOauthEnum =
  RedirectBackToAfterOauth;
export type UserModel = GithubComTDiblikProjectTemplateApiModelsUserModelDB | undefined;

export type ThemePosibilitiesType =
  PreferedTheme;
export const ThemePosibilities: ThemePosibilitiesType[] = ["light", "dark"];

export type TranslationPosibilitiesType =
  PreferedLanguage;
export const TranslationPossibilities: TranslationPosibilitiesType[] = ["cs", "en"];
