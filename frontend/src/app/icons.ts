import { faCloudversify } from '@fortawesome/free-brands-svg-icons';
import { faAnglesUp, faAngleUp, faBars, faBarsStaggered, faCaretLeft, faCaretRight, faCheck, faHeart, faHome, faHurricane, faLanguage, faMinus, faO, faPlus, faSearch, faTriangleExclamation, faUnlockKeyhole, faUser, faX } from '@fortawesome/free-solid-svg-icons';
import { faClock } from '@fortawesome/free-regular-svg-icons';

export const IconLib = {
    accountNone: faUser,
    account: faUser,
    ailocalized: faLanguage,
    attention: faTriangleExclamation,
    caretLeft: faCaretLeft,
    caretRight: faCaretRight,
    clock: faClock,
    difficulty: {
        0: faO,
        1: faMinus,
        2: faAngleUp,
        3: faAnglesUp,
    },
    home: faHome,
    like: faHeart,
    login: faUnlockKeyhole,
    menu: faBars,
    menuOpened: faBarsStaggered,
    nextcloud: faCloudversify,
    recipeNew: faPlus,
    search: faSearch,
    searchClose: faX,
    searchStart: faCheck,
    spinner: faHurricane
};