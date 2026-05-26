import {ConvertToApiError, type OauthRedirectHandlerRequest} from "@shared/api-client";
import {AnimatePresence, type HTMLMotionProps, motion} from "motion/react";
import {useState} from "react";
import {FormProvider, useForm} from "react-hook-form";
import {useTranslation} from "react-i18next";
import {FaFacebook, FaGithub, FaGoogle, FaSpotify} from "react-icons/fa";
import {useNavigate} from "react-router";
import {HiddenBooleanInput} from "../../components/HiddenBooleanInput";
import {TextInput} from "../../components/TextInput";
import {useLoadingStore} from "../../stores/LoadingStore";
import {useAuthTokenStore} from "../../stores/TokenStore";
import {AuthController, oAuthRedirectController} from "../../utils/api";
import {TranslateApiErrorMessage} from "../../utils/general";
import {routes} from "../../utils/routes";
import {useFormLog} from "../../utils/useFormLog";
import {
  LoginFirstPageSchema,
  type LoginOrSignUpPageFormType,
  type LoginPageFormType,
  SignUpFirstPageSchema,
  type SignUpPageFormType,
  zodResolver,
} from "../../utils/validations";

export default function Login() {
  const {t} = useTranslation();
  const navigate = useNavigate();
  const {setToken} = useAuthTokenStore();
  const {setLoading} = useLoadingStore();
  const [isSignUp, setIsSignUp] = useState(false);
  const [beErrorMessage, setBEErrorMessage] = useState<string | undefined>(undefined);

  const form = useForm<LoginOrSignUpPageFormType>({
    mode: "onChange",
    resolver: zodResolver(!isSignUp ? LoginFirstPageSchema : SignUpFirstPageSchema),
  });
  useFormLog(form);

  const toggleFormType = () => {
    form.resetField("password");
    form.resetField("confirmPassword");
    form.clearErrors();
    setBEErrorMessage(undefined);
    setIsSignUp((prev) => !prev);
  };

  const postLogin = (authToken: string) => {
    setToken(authToken);
    navigate(routes.index);
  };

  const handleLoginErr = async (error: any) => {
    const apiError = await ConvertToApiError(error);
    setBEErrorMessage(TranslateApiErrorMessage(apiError));
  };

  const onSubmit = async (data: LoginOrSignUpPageFormType) => {
    setLoading(true, "loadingStates.loggingIn");
    if (isSignUp) {
      const _data = data as SignUpPageFormType;
      AuthController.apiV1PublicAuthSignupPost({
        githubComTDiblikProjectTemplateApiHandlersSignUpHandlerRequestBody: {
          email: _data.email,
          password: _data.password,
          useUsername: _data.useUsername,
          firstName: _data.firstName,
          lastName: _data.lastName,
          username: _data.username,
        },
      })
        .then((s) => postLogin(s.authToken))
        .catch(handleLoginErr)
        .finally(() => setLoading(false));
    } else {
      const _data = data as LoginPageFormType;
      AuthController.apiV1PublicAuthLoginPost({
        githubComTDiblikProjectTemplateApiHandlersLoginHandlerRequestBody: {
          email: _data.email,
          password: _data.password,
        },
      })
        .then((s) => postLogin(s.authToken))
        .catch(handleLoginErr)
        .finally(() => setLoading(false));
    }
  };

  return (
    <div className="flex min-h-screen items-center justify-center bg-base-200 px-4">
      <motion.div layout transition={{type: "spring", stiffness: 120, damping: 20}} className="card w-full max-w-md shadow-2xl bg-base-100 p-8 rounded-2xl">
        <motion.h2
          key={isSignUp ? "sign-up" : "login"}
          initial={{opacity: 0, y: -15}}
          animate={{opacity: 1, y: 0}}
          exit={{opacity: 0, y: 15}}
          transition={{duration: 0.3}}
          className="text-3xl font-bold text-center mb-3"
        >
          {isSignUp ? t("loginPage.signUpTitle") : t("loginPage.loginTitle")}
        </motion.h2>

        <FormProvider {...form}>
          <AnimatePresence mode="wait">
            <motion.form
              key={isSignUp ? "signup-form" : "login-form"}
              initial={{opacity: 0, y: 15}}
              animate={{opacity: 1, y: 0}}
              exit={{opacity: 0, y: -15}}
              transition={{duration: 0.35}}
              onSubmit={form.handleSubmit(onSubmit)}
            >
              {isSignUp && <NameFields t={t} />}
              <TextInput
                label={t("loginPage.email.label")}
                name="email"
                placeholder={t("loginPage.email.placeholder")}
                inputProps={{type: "email"}}
                hasBigText
              />
              <PasswordFields t={t} isSignUp={isSignUp} />
              {beErrorMessage && (
                <motion.p initial={{opacity: 0}} animate={{opacity: 1}} className="text-red-500 text-sm text-center mt-2">
                  {beErrorMessage}
                </motion.p>
              )}
              <motion.button
                whileHover={{scale: 1.03}}
                whileTap={{scale: 0.97}}
                className="btn btn-primary w-full mt-5 border-0 text-white shadow-md hover:shadow-lg transition-all"
                type="submit"
              >
                {isSignUp ? t("loginPage.signUpButton") : t("loginPage.loginContinue")}
              </motion.button>
            </motion.form>
          </AnimatePresence>
        </FormProvider>

        <div className="divider my-6">{t("loginPage.continueWith")}</div>

        <div className="grid grid-cols-2 gap-3 mb-6">
          <OAuthButton provider="Google" icon={<FaGoogle />} onClick={() => oAuthRedirectController.apiV1PublicAuthOauthRedirectGoogleGet()} />
          <OAuthButton provider="Facebook" icon={<FaFacebook />} onClick={() => oAuthRedirectController.apiV1PublicAuthOauthRedirectFacebookGet()} />
          <OAuthButton provider="Spotify" icon={<FaSpotify />} onClick={() => oAuthRedirectController.apiV1PublicAuthOauthRedirectSpotifyGet()} />
          <OAuthButton provider="GitHub" icon={<FaGithub />} onClick={() => oAuthRedirectController.apiV1PublicAuthOauthRedirectGithubGet()} />
        </div>

        <motion.p className="text-base text-center" initial={{opacity: 0}} animate={{opacity: 1}} transition={{delay: 0.2}}>
          {isSignUp ? t("loginPage.signUpAlreadyHaveAccount") : t("loginPage.loginDontHaveAccount")}{" "}
          <button type="button" onClick={toggleFormType} className="text-primary font-medium underline hover:text-secondary transition-colors cursor-pointer">
            {isSignUp ? t("loginPage.signUpLoginHere") : t("loginPage.loginSignUpHere")}
          </button>
        </motion.p>
      </motion.div>
    </div>
  );
}

const NameFields = ({t}: {t: any}) => {
  const [useUsername, setUseUsername] = useState(false);
  const animation = {
    initial: {opacity: 0, scale: 0.95},
    animate: {opacity: 1, scale: 1},
    exit: {opacity: 0, scale: 0.95},
    transition: {duration: 0.25},
  } as HTMLMotionProps<"div">;

  return (
    <div className="flex flex-col gap-4">
      <div className="flex justify-center">
        <button
          type="button"
          onClick={() => setUseUsername(false)}
          className={`btn btn-sm rounded-l-full ${!useUsername ? "btn-primary text-white" : "btn-outline"}`}
        >
          {t("loginPage.useName")}
        </button>
        <button
          type="button"
          onClick={() => setUseUsername(true)}
          className={`btn btn-sm rounded-r-full ${useUsername ? "btn-primary text-white" : "btn-outline"}`}
        >
          {t("loginPage.useUsername")}
        </button>
      </div>
      {!useUsername ? (
        <motion.div key="name-fields" {...animation} className="flex gap-4">
          <TextInput label={t("loginPage.firstName.label")} name="firstName" placeholder={t("loginPage.firstName.placeholder")} hasBigText />
          <TextInput label={t("loginPage.lastName.label")} name="lastName" placeholder={t("loginPage.lastName.placeholder")} hasBigText />
        </motion.div>
      ) : (
        <motion.div key="username-field" {...animation}>
          <TextInput label={t("loginPage.username.label")} name="username" placeholder={t("loginPage.username.placeholder")} hasBigText />
        </motion.div>
      )}
      <HiddenBooleanInput name="useUsername" value={useUsername} />
    </div>
  );
};

const PasswordFields = ({t, isSignUp}: {t: any; isSignUp: boolean}) => (
  <div className="flex gap-4">
    <TextInput
      label={t("loginPage.password.label")}
      name="password"
      placeholder={t("loginPage.password.placeholder")}
      inputProps={{type: "password"}}
      hasBigText
    />
    {isSignUp && (
      <TextInput
        label={t("loginPage.confirmPassword.label")}
        name="confirmPassword"
        placeholder={t("loginPage.confirmPassword.placeholder")}
        inputProps={{type: "password"}}
        hasBigText
      />
    )}
  </div>
);

const OAuthButton = ({provider, icon, onClick}: {provider: string; icon: React.ReactNode; onClick: () => OauthRedirectHandlerRequest}) => {
  const {loading, setLoading} = useLoadingStore();
  return (
    <motion.button
      whileHover={{scale: 1.05}}
      whileTap={{scale: 0.97}}
      className="btn btn-outline w-full flex items-center gap-2 transition-all hover:border-primary"
      onClick={() => {
        setLoading(true);
        onClick()
          .then((s) => {
            setLoading(false);
            window.location.href = s.redirectUrl!;
          })
          .finally(() => {
            if (loading) setLoading(false);
          });
      }}
    >
      {icon} {provider}
    </motion.button>
  );
};
