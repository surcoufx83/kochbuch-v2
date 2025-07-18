import { L10nLocale } from "./types";

export const KB_Fr: L10nLocale = {
    collections: {
        defaultCollection: {
            description: "",
            title: "Mon livre de cuisine",
        },
    },
    common: {
        language: {
            de: "Allemand",
            en: "Anglais",
            fr: "Français"
        },
        unknownUser: "Ancien utilisateur",
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
    me: {
        title: "Ton livre de cuisine numérique",
        description: "Voici ton espace personnel — tu y trouveras toutes tes recettes, qu’elles aient été créées ou simplement enregistrées.",
        collections: {
            title: "Collections",
            description: "Les recettes sont organisées en collections. Une collection par défaut a été créée automatiquement pour toi, mais tu peux en créer autant que tu veux — par exemple pour les meilleures recettes de pâtisserie, desserts ou cocktails."
        }
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
        hasAdminAccess: "Visible uniquement parce que tu es connecté en tant qu’administrateur.",
        isLoading: "Chargement des données de recette en cours...",
        loadingFailed: "La recette n'a pas pu être chargée. Vous serez automatiquement redirigé vers la page d'accueil...",
        modifiedAndUser: "[0] par ",
        printBtn: "Imprimer",
        saveBtn: "Enregistrer",
        saveBtnAlreadySaved: "Enregistré",
        submittedBy: "Ajouté par [0]",
        shareBtn: "Partager",
        share: {
            title: "Recette de « [0] » sur kochbuch.mogul.network",
            message: "*[0]*\n\nTu trouveras cette recette et bien d’autres de [1] dans notre livre de cuisine numérique. Bonne cuisine !"
        },
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
    recipeManagement: {
        title: "Gérer la recette",
        adminDescription: "Cette recette a été créée par [0]. En tant qu’administrateur, tu peux la modifier et accéder à tous les paramètres associés.",
        gotoOwnerBtn: [
            "Voir le compte de [0]",
            "Voir le compte de [0]"
        ],
        ownerDescription: "Tu peux modifier les recettes que tu as créées ou en ajuster la visibilité à tout moment. Toutes les options se trouvent dans cette section.",
        created: "Créée le [0]",
        modified: "Dernière modification le [0]",
        delete: {
            title: "Supprimer",
            description: "Tu veux te débarrasser de cette recette ? Clique sur le bouton de suppression. Cette action est irréversible. Toutes les informations, images et commentaires seront également supprimés.",
            btn: "Supprimer la recette"
        },
        edit: {
            title: "Modifier",
            description: "En mode édition, tu peux revoir l’intégralité du contenu de la recette.",
            btn: "Modifier"
        },
        publish: {
            title: "Publier",
            description: "Les recettes sont privées par défaut. Tu peux les rendre visibles pour d'autres et les partager via le bouton de partage.",
            togglePrivate: "Privée",
            toggleInternal: "Utilisateurs connectés",
            togglePublic: "Publique",
            descriptionPrivate: "Privée signifie vraiment privée. Cette recette t’appartient et personne d’autre ne peut la consulter. C’est l’option par défaut pour toute nouvelle recette. Tu peux l’imprimer, mais pas la partager via WhatsApp.",
            descriptionInternal: "Utilisateurs connectés – la familia : Ce niveau permet aux utilisateurs possédant un compte Nextcloud d’accéder à ta recette. Parfait pour la famille et les amis proches. Le partage via WhatsApp reste désactivé.",
            descriptionPublic: "Publique : Tout le monde peut voir cette recette s’il a le lien. Ce niveau permet un partage libre avec les autres et de la faire connaître à ton entourage."
        }
    },
    recipeOwnerInfo: {
        title: "Note de modification",
        description: "Cette recette a été créée par vous. Vous pouvez la modifier ou en ajuster la visibilité à tout moment. Les options correspondantes se trouvent plus bas, après les étapes de préparation.",
        adminDescription: "Cette recette a été créée par [0], mais elle est visible pour toi car tu es connecté en tant qu’administrateur. Les options de gestion se trouvent plus bas, après les étapes de préparation.",
        gothereLink: "Aller directement"
    },
    recipePreparation: {
        title: "Préparation",
        stepFormat: "[0] : [1]",
        stepFallback: "étape",
    },
    recipePreparationTime: {
        title: "Temps de préparation",
        longDurationWarning: "Veuillez noter que la durée totale de préparation pour cette recette est d’au moins [0].",
        recalcWarning: {
            screen: [
                "Le temps de préparation est basé sur 1 portion et n’est pas automatiquement ajusté pour un plus grand nombre de portions. Merci d’en tenir compte lors de la préparation.",
                "Le temps de préparation est basé sur [0] portions et n’est pas automatiquement ajusté pour un plus grand nombre de portions. Merci d’en tenir compte lors de la préparation."
            ],
            print: [
                "Temps de préparation pour 1 portion !",
                "Temps de préparation pour [0] portions !"
            ]
        },
        items: {
            cooking: ["Temps de cuisson", "[0] cuisson"],
            preparing: ["Temps de préparation", "[0] préparation"],
            total: "Temps total",
            waiting: ["Temps de repos", "[0] repos"]
        },
        units: {
            days: ['1 j', '[0] j', '1 jour', '[0] jours'],
            hours: ['1 h', '[0] h', '1 heure', '[0] heures'],
            minutes: ['1 min', '[0] min', '1 minute', '[0] minutes']
        },
    },
    saveToCollection: {
        confirmMsg: "Recette enregistrée.",
        gotoCollectionLink: [
            "Voir la collection de recettes",
            "Voir tes collections",
        ],
        pickCollection: {
            description: "Choisis une collection dans laquelle enregistrer la recette.",
            itemCount: [
                "Ne contient encore aucune recette",
                "Contient une recette",
                "Contient [0] recettes",
            ],
            titleInputLabel: "Nom de la collection de recettes :",
            newBtn: "Nouvelle collection",
            saveBtn: "Enregistrer",
        },
    },
}