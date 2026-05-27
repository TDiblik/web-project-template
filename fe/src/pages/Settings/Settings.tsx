import {type OauthRedirectHandlerRequest, TranslationPossibilities} from "@shared/api-client";
import type React from "react";
import {useState} from "react";
import {FormProvider, useForm} from "react-hook-form";
import {useTranslation} from "react-i18next";
import {FaFacebook, FaGithub, FaGoogle, FaSpotify, FaTimes} from "react-icons/fa";
import {HiOutlineChevronDown} from "react-icons/hi";
import Layout from "../../components/Layout";
import {TextInput} from "../../components/TextInput";
import {usei18nStore} from "../../stores/i18nStore";
import {useLoadingStore} from "../../stores/LoadingStore";
import {useLoggedUser} from "../../stores/LoggedUserStore";
import {useThemeStore} from "../../stores/ThemeStore";
import {oAuthRedirectController} from "../../utils/api";
import {SignUpFirstPageSchema, type SignUpPageFormType, zodResolver} from "../../utils/validations";

export default function SettingsPage() {
  const {i18n, t} = useTranslation();
  const {loggedUser} = useLoggedUser();
  const {theme, changeTheme} = useThemeStore();
  const {changeLanguage} = usei18nStore();
  const {setLoading} = useLoadingStore();
  const form = useForm<SignUpPageFormType>({
    mode: "onChange",
    defaultValues: loggedUser,
    resolver: zodResolver(SignUpFirstPageSchema),
  });

  const onSubmit = async () => {
    setLoading(true);
    try {
      alert("Profile updated!");
    } finally {
      setLoading(false);
    }
  };

  return (
    <Layout>
      <div className="max-w-5xl mx-auto py-8 space-y-8">
        <div className="flex flex-wrap items-center justify-between gap-4">
          <h1 className="text-3xl font-bold">{t("settingsPage.pageTitle")}</h1>
          <div className="flex items-center gap-2">
            <div className="dropdown dropdown-end">
              <button type="button" tabIndex={0} className="btn btn-sm gap-1">
                {t(`layout.changeLanguage.${i18n.language}`)}
                <HiOutlineChevronDown className="w-4 h-4 ml-1" />
              </button>
              <ul className="dropdown-content z-1 menu p-2 shadow bg-base-100 rounded-box w-32">
                {TranslationPossibilities.map((lang) => (
                  <li key={lang}>
                    <button
                      type="button"
                      onClick={() => {
                        changeLanguage(lang);
                        (document.activeElement as HTMLElement)?.blur();
                      }}
                      className={`w-full text-left ${i18n.language === lang ? "font-semibold text-primary" : ""}`}
                    >
                      {t(`layout.changeLanguage.${lang}`)}
                    </button>
                  </li>
                ))}
              </ul>
            </div>
            <div className="w-px h-5 bg-base-300 mx-2"></div>
            <button
              type="button"
              onClick={() => changeTheme(theme === "dark" ? "light" : "dark")}
              className="flex items-center gap-1 px-2.5 py-1 text-sm rounded-md border border-base-300 text-base-content hover:border-primary/50 hover:text-primary transition cursor-pointer"
            >
              <span className="mr-1">{theme === "dark" ? "🌙" : "☀️"}</span>
              {t(`layout.changeTheme.${theme}`)}
            </button>
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div className="space-y-6">
            <div className="p-4 rounded-xl shadow-sm bg-base-100 flex flex-col items-center">
              {loggedUser && (
                <>
                  {loggedUser.avatarUrl ? (
                    <img src={loggedUser.avatarUrl} alt={t("layout.userAvatarAlt")} className="w-28 h-28 rounded-full object-cover mb-3" />
                  ) : (
                    <div className="avatar-placeholder w-32 h-32 rounded-full bg-neutral text-neutral-content flex items-center justify-center mb-3">
                      <span className="text-4xl">{loggedUser.initials}</span>
                    </div>
                  )}
                  <button type="button" className="btn btn-sm btn-outline w-full">
                    {t("settingsPage.changeAvatar")}
                  </button>
                </>
              )}
            </div>

            <div className="p-4 rounded-xl shadow-sm bg-base-100">
              <h2 className="font-semibold mb-3">{t("settingsPage.connectedAccounts")}</h2>
              <div className="flex flex-col gap-2">
                <OAuthButton
                  provider="Google"
                  icon={<FaGoogle />}
                  connected={!!loggedUser?.googleId}
                  textConnect={t("settingsPage.oauth.connect")}
                  textConnected={t("settingsPage.oauth.connected")}
                  onConnect={() =>
                    oAuthRedirectController.apiV1PublicAuthOauthRedirectGoogleGet({
                      redirectBackToAfterOauth: "settings",
                    })
                  }
                />
                <OAuthButton
                  provider="Facebook"
                  icon={<FaFacebook />}
                  connected={!!loggedUser?.facebookId}
                  textConnect={t("settingsPage.oauth.connect")}
                  textConnected={t("settingsPage.oauth.connected")}
                  onConnect={() =>
                    oAuthRedirectController.apiV1PublicAuthOauthRedirectFacebookGet({
                      redirectBackToAfterOauth: "settings",
                    })
                  }
                />
                <OAuthButton
                  provider="Spotify"
                  icon={<FaSpotify />}
                  connected={!!loggedUser?.spotifyId}
                  textConnect={t("settingsPage.oauth.connect")}
                  textConnected={t("settingsPage.oauth.connected")}
                  onConnect={() =>
                    oAuthRedirectController.apiV1PublicAuthOauthRedirectSpotifyGet({
                      redirectBackToAfterOauth: "settings",
                    })
                  }
                />
                <OAuthButton
                  provider="Github"
                  icon={<FaGithub />}
                  connected={!!loggedUser?.githubHandle}
                  textConnect={t("settingsPage.oauth.connect")}
                  textConnected={t("settingsPage.oauth.connected")}
                  onConnect={() =>
                    oAuthRedirectController.apiV1PublicAuthOauthRedirectGithubGet({
                      redirectBackToAfterOauth: "settings",
                    })
                  }
                />
              </div>
            </div>
          </div>

          <div className="md:col-span-2 space-y-6">
            <FormProvider {...form}>
              <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
                <div className="p-6 rounded-xl shadow-sm bg-base-100 space-y-4">
                  <h2 className="font-semibold text-lg">{t("settingsPage.profileInfo")}</h2>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <TextInput label={t("settingsPage.firstName.label")} name="firstName" placeholder={t("settingsPage.firstName.placeholder")} hasBigText />
                    <TextInput label={t("settingsPage.lastName.label")} name="lastName" placeholder={t("settingsPage.lastName.placeholder")} hasBigText />
                  </div>
                  <button type="submit" className="btn btn-primary w-full mt-2">
                    {t("settingsPage.saveChanges")}
                  </button>
                </div>
                <div className="p-6 rounded-xl shadow-sm bg-base-100 space-y-4">
                  <h2 className="font-semibold text-lg">{t("settingsPage.changePassword")}</h2>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <TextInput
                      label={t("settingsPage.password.label")}
                      name="password"
                      placeholder={t("settingsPage.password.placeholder")}
                      inputProps={{type: "password"}}
                      hasBigText
                    />
                    <TextInput
                      label={t("settingsPage.confirmPassword.label")}
                      name="confirmPassword"
                      placeholder={t("settingsPage.confirmPassword.placeholder")}
                      inputProps={{type: "password"}}
                      hasBigText
                    />
                  </div>
                  <button type="submit" className="btn btn-primary w-full mt-2">
                    {t("settingsPage.savePassword")}
                  </button>
                </div>
              </form>
            </FormProvider>
          </div>
        </div>
      </div>
    </Layout>
  );
}

interface OAuthButtonProps {
  provider: string;
  icon: React.ReactNode;
  connected: boolean;
  textConnect: string;
  textConnected: string;
  onConnect: () => OauthRedirectHandlerRequest;
  // onDisconnect: () => Promise<any>;
}
const OAuthButton: React.FC<OAuthButtonProps> = ({provider, icon, connected, textConnect, textConnected, onConnect}) => {
  const [hovered, setHovered] = useState(false);

  return (
    <button
      type="button"
      onClick={() =>
        !connected &&
        onConnect().then((s) => {
          if (s.redirectUrl) window.location.href = s.redirectUrl;
        })
      }
      className={`group relative flex items-center justify-between w-full px-4 py-2 rounded-md border transition ${
        connected ? "bg-green-100 text-green-700 border-green-200 cursor-not-allowed" : "border-gray-300 hover:bg-gray-50 cursor-pointer"
      }`}
      onMouseEnter={() => setHovered(true)}
      onMouseLeave={() => setHovered(false)}
    >
      <span className="flex items-center gap-2">
        {icon}
        <span className="font-medium">{provider}</span>
      </span>

      <span className={`text-sm opacity-80 ${connected && hovered && "mr-3"}`}>{connected ? textConnected : textConnect}</span>

      {connected && (
        <FaTimes
          size={16}
          className="absolute right-2 top-1/2 -translate-y-1/2 opacity-0 group-hover:opacity-100 text-red-500 hover:text-red-600 transition cursor-pointer"
          onClick={() => console.log("todo")}
        />
      )}
    </button>
  );
};
