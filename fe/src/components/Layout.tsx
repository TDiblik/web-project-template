import {ThemePosibilities, type ThemePosibilitiesType, type TranslationPosibilitiesType, TranslationPossibilities} from "@shared/api-client";
import {AnimatePresence, type HTMLMotionProps, motion} from "motion/react";
import type React from "react";
import {useEffect, useRef, useState} from "react";
import {useTranslation} from "react-i18next";
import {Link, matchPath, useLocation} from "react-router";
import {usei18nStore} from "../stores/i18nStore";
import {useLoggedUser} from "../stores/LoggedUserStore";
import {useThemeStore} from "../stores/ThemeStore";
import {routes} from "../utils/routes";

const Layout: React.FC<React.PropsWithChildren> = ({children}) => {
  const location = useLocation();
  const {t, i18n} = useTranslation();
  const {theme, changeTheme} = useThemeStore();
  const {changeLanguage} = usei18nStore();
  const {loggedUser} = useLoggedUser();

  const menuItems = [
    {name: t("layout.dashboard"), path: routes.index},
    {name: t("layout.settings"), path: routes.settings},
  ];

  const [profileOpen, setProfileOpen] = useState(false);
  const [themeOpen, setThemeOpen] = useState(false);
  const [languageOpen, setLanguageOpen] = useState(false);

  const profileMenuRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (profileMenuRef.current && !profileMenuRef.current.contains(event.target as Node)) {
        setProfileOpen(false);
        setThemeOpen(false);
        setLanguageOpen(false);
      }
    };
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  const toggleProfileMenu = () => {
    setProfileOpen((prev) => !prev);
    if (profileOpen) {
      setThemeOpen(false);
      setLanguageOpen(false);
    }
  };

  const toggleMenu = (menu: "language" | "theme") => {
    if (menu === "language") {
      setLanguageOpen((prev) => !prev);
      setThemeOpen(false);
    } else {
      setThemeOpen((prev) => !prev);
      setLanguageOpen(false);
    }
  };

  const changeThemeAndClose = (newTheme: ThemePosibilitiesType) => {
    changeTheme(newTheme);
    setThemeOpen(false);
    setProfileOpen(false);
  };

  const changeLanguageAndClose = (lang: TranslationPosibilitiesType) => {
    changeLanguage(lang);
    setLanguageOpen(false);
    setProfileOpen(false);
  };

  const toggleMenuClasses = "absolute left-full top-0 menu rounded-box w-48 bg-base-100 p-2 shadow";
  const toggleMenuAnimation = {
    initial: {opacity: 0, y: -8},
    animate: {opacity: 1, y: 0},
    exit: {opacity: 0, y: -8},
    transition: {duration: 0.15, ease: "easeOut"},
  } as HTMLMotionProps<"ul">;

  return (
    <div className="flex h-screen bg-base-200">
      {/* Sidebar */}
      <div className="flex w-64 flex-col border-r border-base-300 bg-base-100">
        <div className="border-b border-base-300 p-6 text-2xl font-bold">{t("projectName")}</div>

        {/* Menu */}
        <ul className="menu flex-1 gap-2 p-4">
          {menuItems.map((item) => (
            <li key={item.name} className={matchPath(item.path, location.pathname) ? "bg-primary text-primary-content rounded-lg" : ""}>
              <Link to={item.path}>{item.name}</Link>
            </li>
          ))}
        </ul>

        {/* Avatar & Settings */}
        <div className="border-t border-base-300 p-4">
          <div ref={profileMenuRef} className={`dropdown dropdown-top dropdown-end w-full cursor-pointer ${profileOpen ? "dropdown-open" : ""}`}>
            <button type="button" onClick={toggleProfileMenu} className="flex items-center w-full">
              {loggedUser && (
                <div className={`btn btn-ghost btn-circle avatar ${!loggedUser.avatarUrl ? "avatar-placeholder" : ""}`}>
                  <div className="w-12 rounded-full bg-neutral text-neutral-content">
                    {loggedUser.avatarUrl ? <img src={loggedUser.avatarUrl} alt={t("layout.userAvatarAlt")} /> : <span>{loggedUser.initials}</span>}
                  </div>
                </div>
              )}
              <span className="ml-2 text-base font-medium normal-case">{loggedUser?.fullName}</span>
            </button>

            <ul className="dropdown-content menu rounded-box z-50 mb-2 w-52 bg-base-100 p-2 shadow">
              <li>
                <Link to={routes.settings} onClick={() => setProfileOpen(false)}>
                  {t("layout.settings")}
                </Link>
              </li>

              <li className="relative">
                <button type="button" className="flex justify-between w-full" onClick={() => toggleMenu("language")}>
                  {t("layout.changeLanguage.label")}
                </button>
                <AnimatePresence>
                  {languageOpen && (
                    <motion.ul {...toggleMenuAnimation} className={toggleMenuClasses}>
                      {TranslationPossibilities.map((lang) => (
                        <li key={lang}>
                          <button
                            type="button"
                            className={i18n.language === lang ? "font-bold text-primary w-full text-left" : "w-full text-left"}
                            onClick={() => changeLanguageAndClose(lang)}
                          >
                            {t(`layout.changeLanguage.${lang}`)}
                          </button>
                        </li>
                      ))}
                    </motion.ul>
                  )}
                </AnimatePresence>
              </li>

              <li className="relative">
                <button type="button" className="flex justify-between w-full" onClick={() => toggleMenu("theme")}>
                  {t("layout.changeTheme.label")}
                </button>
                <AnimatePresence>
                  {themeOpen && (
                    <motion.ul {...toggleMenuAnimation} className={toggleMenuClasses}>
                      {ThemePosibilities.map((s) => (
                        <li key={s}>
                          <button
                            type="button"
                            className={theme === s ? "font-bold text-primary w-full text-left" : "w-full text-left"}
                            onClick={() => changeThemeAndClose(s)}
                          >
                            {t(`layout.changeTheme.${s}`)}
                          </button>
                        </li>
                      ))}
                    </motion.ul>
                  )}
                </AnimatePresence>
              </li>

              <div className="my-1 border-t border-base-300"></div>

              <li>
                <Link to={routes.logout} className="text-error">
                  {t("layout.logout")}
                </Link>
              </li>
            </ul>
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="flex-1 overflow-auto p-6">{children}</div>
    </div>
  );
};

export default Layout;
