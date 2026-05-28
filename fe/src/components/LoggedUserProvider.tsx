import {useEffect} from "react";
import {useNavigate} from "react-router";
import {useLoadingStore} from "../stores/LoadingStore";
import {useFetchLoggedUser, useLoggedUserStore} from "../stores/LoggedUserStore";
import {useAuthTokenStore} from "../stores/TokenStore";
import {routes} from "../utils/routes";

export const LoggedUserProvider: React.FC<React.PropsWithChildren> = ({children}) => {
  const {isAuthenticated} = useAuthTokenStore();

  const _isAuthenticated = isAuthenticated();
  if (!_isAuthenticated) {
    return <>{children}</>;
  }

  return <LoggedUserProviderInternal>{children}</LoggedUserProviderInternal>;
};

const LoggedUserProviderInternal: React.FC<React.PropsWithChildren> = ({children}) => {
  const navigate = useNavigate();
  const {tokenRaw, resetTokens} = useAuthTokenStore();
  const {loading: loadingIndicator, setLoading: setLoadingIndicator} = useLoadingStore();
  const {isLoading: loadingRequest, isError: isFetchingError, isFetching} = useFetchLoggedUser();

  const user = useLoggedUserStore((s) => s.user);
  const isError = isFetchingError && !loadingRequest && !isFetching && !user;

  useEffect(() => {
    if (!user && loadingRequest) {
      setLoadingIndicator(true);
    }
    if (isError || (user && loadingIndicator && !isError)) {
      setLoadingIndicator(false);
    }
    if (isError) {
      resetTokens();
      navigate(routes.loginOAuthRedirect);
    }
  }, [user, loadingRequest, isError, tokenRaw]);

  return <>{children}</>;
};
