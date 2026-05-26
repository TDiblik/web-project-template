import {AnimatePresence, type HTMLMotionProps, motion} from "motion/react";
import {useEffect, useState} from "react";
import {useTranslation} from "react-i18next";
import {useNavigate, useSearchParams} from "react-router";
import {useLoggedUser} from "../../stores/LoggedUserStore";
import {useAuthTokenStore} from "../../stores/TokenStore";
import {oAuthController} from "../../utils/api";
import {RedirectBackToAfterOauthToRouteMap, routes} from "../../utils/routes";

const delayed = (delay: number) =>
  ({
    initial: {opacity: 0, y: 5},
    animate: {opacity: 1, y: 0},
    transition: {delay},
  }) as HTMLMotionProps<"p"> | HTMLMotionProps<"h1">;

const usedCodes = new Set<string>();

const OAuthRedirect = () => {
  const {t} = useTranslation();
  const navigate = useNavigate();
  const [query] = useSearchParams();
  const oAuthCode = query.get("code");
  const oAuthState = query.get("state");
  const [shouldFail, setShouldFail] = useState(!oAuthCode || !oAuthState);
  const {isAuthenticated} = useAuthTokenStore();
  const {refetchUser} = useLoggedUser();
  const _isAuthenticated = isAuthenticated();

  const {setToken} = useAuthTokenStore();

  useEffect(() => {
    if (!oAuthCode || !oAuthState || usedCodes.has(oAuthCode)) return;
    usedCodes.add(oAuthCode);

    oAuthController
      .apiV1PublicAuthOauthReturnGet({
        state: oAuthState,
        code: oAuthCode,
      })
      .then((s) => {
        setToken(s.authToken);
        refetchUser();
        navigate(RedirectBackToAfterOauthToRouteMap[s.redirectBackToAfterOauth]);
      })
      .catch(() => setShouldFail(true));
  }, [oAuthCode, oAuthState]);

  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-base-100 text-center px-4">
      <AnimatePresence mode="wait">
        {!shouldFail ? (
          <motion.div
            key="loading"
            initial={{opacity: 0, y: 10}}
            animate={{opacity: 1, y: 0}}
            exit={{opacity: 0, y: -10}}
            transition={{duration: 0.6, ease: "easeInOut"}}
            className="flex flex-col items-center"
          >
            <motion.span
              className="loading loading-spinner loading-xl text-primary mb-4"
              animate={{
                scale: [1, 1.15, 1],
                opacity: [1, 0.7, 1],
              }}
              transition={{
                duration: 1.2,
                repeat: Infinity,
                ease: "easeInOut",
              }}
            />
            <motion.p className="text-lg font-medium text-base-content" {...delayed(0.2)}>
              {t("oauthRedirectPage.loginInProgress.title")}
            </motion.p>
            <motion.p className="text-sm text-base-content/70 mt-1" {...delayed(0.4)}>
              {t("oauthRedirectPage.loginInProgress.description")}
            </motion.p>
          </motion.div>
        ) : (
          <motion.div
            key="fail"
            initial={{opacity: 0}}
            animate={{opacity: 1}}
            exit={{opacity: 0}}
            transition={{duration: 0.8, ease: "easeOut"}}
            className="flex flex-col items-center"
          >
            <motion.h1 className="text-2xl font-bold text-error mb-2" {...delayed(0.2)}>
              {t("oauthRedirectPage.loginFailed.title")}
            </motion.h1>
            <motion.p className="text-base text-base-content/80 mb-6" {...delayed(0.4)}>
              {_isAuthenticated ? t("oauthRedirectPage.loginFailed.description-authenticated") : t("oauthRedirectPage.loginFailed.description")}
            </motion.p>
            <motion.button
              whileHover={{scale: 1.05}}
              whileTap={{scale: 0.95}}
              className="btn btn-primary"
              onClick={() => (_isAuthenticated ? navigate(routes.settings) : navigate(routes.login))}
            >
              {_isAuthenticated ? t("oauthRedirectPage.loginFailed.button-authenticated") : t("oauthRedirectPage.loginFailed.button")}
            </motion.button>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
};

export default OAuthRedirect;
