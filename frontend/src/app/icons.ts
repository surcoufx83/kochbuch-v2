import { faCloudversify } from '@fortawesome/free-brands-svg-icons';
import { faAnglesUp, faAngleUp, faBars, faBarsStaggered, faCalculator, faCaretLeft, faCaretRight, faCheck, faFireBurner, faHeart, faHome, faHurricane, faLanguage, faMinus, faO, faPen, faPenAlt, faPepperHot, faPlus, faPrint, faQuoteLeft, faQuoteRight, faRotateLeft, faSearch, faShare, faTriangleExclamation, faTurnDown, faUnlockKeyhole, faUser, faX } from '@fortawesome/free-solid-svg-icons';
import { faClock } from '@fortawesome/free-regular-svg-icons';

export const IconLib = {
    accountNone: faUser,
    account: faUser,
    ailocalized: faLanguage,
    attention: faTriangleExclamation,
    calc: faCalculator,
    calcMinus: faMinus,
    calcPlus: faPlus,
    caretLeft: faCaretLeft,
    caretRight: faCaretRight,
    clock: faClock,
    difficulty: {
        0: faO,
        1: faMinus,
        2: faAngleUp,
        3: faAnglesUp,
    },
    edit: faPen,
    goDownToArrow: faTurnDown,
    home: faHome,
    ingredients: faPepperHot,
    like: faHeart,
    login: faUnlockKeyhole,
    menu: faBars,
    menuOpened: faBarsStaggered,
    nextcloud: faCloudversify,
    preparation: faFireBurner,
    print: faPrint,
    quote: faQuoteRight,
    recipeNew: faPlus,
    reset: faRotateLeft,
    search: faSearch,
    searchClose: faX,
    searchStart: faCheck,
    share: faShare,
    spinner: faHurricane
};