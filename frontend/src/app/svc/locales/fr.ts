import { L10nLocale } from "./types";

export const KB_Fr: L10nLocale = {
    common: {
        language: {
            de: "Allemand",
            en: "Anglais",
            fr: "Français"
        }
    },
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
                    jan: "Terme de recherche (ex. Haggis avec navets et pommes de terre)",
                    feb: "Terme de recherche (ex. Crêpes au citron et sucre)",
                    mar: "Terme de recherche (ex. Agneau rôti au romarin)",
                    apr: "Terme de recherche (ex. Velouté d'asperges)",
                    may: "Terme de recherche (ex. Ragoût de gibier)",
                    jun: "Terme de recherche (ex. Salade estivale fraîche)",
                    jul: "Terme de recherche (ex. Tarte aux fraises)",
                    aug: "Terme de recherche (ex. Légumes grillés avec halloumi)",
                    sep: "Terme de recherche (ex. Soupe de potiron au gingembre)",
                    oct: "Terme de recherche (ex. Gratin de pommes de terre traditionnel)",
                    nov: "Terme de recherche (ex. Oie rôtie avec chou rouge braisé)",
                    dec: "Terme de recherche (ex. Pudding de Noël)"
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
        aiLocalizationSwitch: [
            "Afficher la recette originale",
            "Afficher la traduction"
        ],
        aiLocalizedContent: "Cette recette a été traduite automatiquement.",
        aiLocalizedContentWithSourceLocale: "Cette recette a été traduite automatiquement depuis [0].",
        aiLocalizedContentSourceLocale: {
            de: "l’allemand",
            en: "l’anglais",
            fr: "le français"
        },
        aiSourceLocale: "Titre original : [0]",
        difficulty: {
            0: "Aucune indication",
            1: "Facile",
            2: "Moyen",
            3: "Difficile"
        },
        isLoading: "Chargement des données de recette en cours...",
        loadingFailed: "La recette n'a pas pu être chargée. Vous serez automatiquement redirigé vers la page d'accueil...",
        submittedBy: "Ajouté par [0]"
    },
    recipeGallery: {
        noPicturesUploaded: "Pas encore de photo. Cliquez ici pour en prendre une...",
        uploadBtn: "Téléverser une photo",
        uploadStatus: {
            checking: "Vérification...",
            uploading: "Enregistrement...",
        },
    },
    recipeIngredients: {
        title: "Liste des ingrédients",
        description: [
            "Cette recette est prévue pour 1 portion.",
            "Cette recette est prévue pour [0] portions.",
        ],
        calculator: {
            servings: [
                "Portion",
                "Portions",
            ],
            title: "Calculateur de quantités",
            description: "Vous pouvez recalculer les quantités pour un autre nombre de portions. La conversion s’appliquera également aux étapes de préparation ci-dessous.",
        },
        table: {
            quantityHeader: "Quantité",
            nameHeader: "Description de l’ingrédient",
        },
    },
    recipeOwnerInfo: {
        title: "Note de modification",
        description: "Cette recette a été créée par vous. Vous pouvez la modifier ou en ajuster la visibilité à tout moment. Les options correspondantes se trouvent plus bas, après les étapes de préparation.",
        gothereLink: "Aller directement"
    },
    recipePreparation: {
        title: "Préparation",
        stepFormat: "[0] : [1]",
        stepFallback: "étape",
    },
    recipePreparationTime: {
        title: "Temps de préparation",
        items: {
            cooking: ["Temps de cuisson", "[0] cuisson"],
            preparing: ["Temps de préparation", "[0] préparation"],
            total: "Temps total",
            waiting: ["Temps de repos", "[0] repos"]
        },
        units: {
            days: ['1 jour', '[0] jours'],
            hours: ['1 h', '[0] h'],
            minutes: ['1 min', '[0] min']
        }
    },
}