import { L10nLocale } from "./types";

export const KB_Fr: L10nLocale = {
    errorPages: {
        loginToCreateRecipe: {
            title: "Faites-en votre propre livre de cuisine 🧑‍🍳",
            paragraphLine1: "Vous devez vous connecter pour utiliser votre propre livre de cuisine numérique.",
            paragraphLine2: "L'utilisation est facultative et gratuite, mais réservée à la famille et aux amis. La connexion nécessite un compte Nextcloud fourni par nous."
        },
        routeNotFound: {
            title1: "Oups !",
            title2: "Page introuvable 🧑‍🍳🚫",
            paragraphLine1: "On dirait que vous êtes tombé sur une erreur !",
            paragraphLine2: "Mais pas d’inquiétude ! Voici quelques options :",
            optionLink1: "Retour au menu principal – Vous y trouverez sûrement de délicieuses alternatives.",
            optionLink2: "Essayer un autre ingrédient – Peut-être qu'une nouvelle recherche vous mènera à la recette souhaitée.",
            optionLink3: "Ou quelque chose de complètement différent – Que diriez-vous d'un café ? ☕ Ou peut-être d'une part de gâteau ? 🍰"
        }
    },
    floatingMenu: {
        searchButton: {
            ariaLabel: "Ce bouton ouvre un champ de recherche pour rechercher dans la base de recettes. Pour valider la recherche, utilisez la touche Entrée du clavier.",
            searchIconAriaLabel: "Icône représentant une loupe. Lorsque la recherche est ouverte, un grand X apparaît comme icône, indiquant qu'un clic dessus fermera le champ de recherche.",
            searchInput: {
                ariaLabel: "Saisissez votre requête de recherche dans ce champ. Par exemple, le mot \"Lasagne\" pour rechercher ce type de plat.",
                placeholder: {
                    jan: "Saisir un terme de recherche (ex. Haggis avec navets et pommes de terre)",
                    feb: "Saisir un terme de recherche (ex. Crêpes au citron et sucre)",
                    mar: "Saisir un terme de recherche (ex. Agneau rôti à la menthe)",
                    apr: "Saisir un terme de recherche (ex. Asperges printanières en salade)",
                    may: "Saisir un terme de recherche (ex. Ragoût de gibier)",
                    jun: "Saisir un terme de recherche (ex. Salade fraîche d'été)",
                    jul: "Saisir un terme de recherche (ex. Fraisier)",
                    aug: "Saisir un terme de recherche (ex. Légumes grillés avec halloumi)",
                    sep: "Saisir un terme de recherche (ex. Soupe de potiron au gingembre)",
                    oct: "Saisir un terme de recherche (ex. Gratin de pommes de terre traditionnel)",
                    nov: "Saisir un terme de recherche (ex. Oie rôtie avec chou rouge braisé)",
                    dec: "Saisir un terme de recherche (ex. Pudding de Noël)"
                }
            },
            submitIconAriaLabel: "Ce bouton lance la recherche avec le terme saisi. Après validation, vous serez redirigé vers la page de résultats."
        }
    },
    homeRecipesPage: {
        titleUser: "Bonjour [0]",
        titleGuest: "Bonjour visiteur",
        descriptionUser: "Vous trouverez ici les dernières recettes du livre de cuisine. La fonction de recherche se trouve en bas de la fenêtre.",
        descriptionGuest: "Ravi de vous voir parcourir le livre de cuisine. Veuillez vous connecter via le bouton utilisateur en bas de la fenêtre pour accéder à toutes les recettes."
    },
    login: {
        loginWithNcButton: "Se connecter avec un compte Nextcloud"
    },
    navbar: {
        brand: {
            pageTitle: "Livre de recettes",
            iconLabel: "Icône du livre de recettes de ce site web"
        }
    },
    recipe: {
        aiLocalizedContent: "Cette recette ou certaines de ses parties ont été traduites automatiquement.",
        submittedBy: "Ajouté par [0]"
    }
}