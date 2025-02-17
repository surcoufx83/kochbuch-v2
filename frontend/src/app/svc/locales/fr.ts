import { L10nLocale } from "./types";

export const KB_Fr: L10nLocale = {
    navbar: {
        brand: {
            pageTitle: "Livre de recettes",
            iconLabel: "Icône du livre de recettes de ce site web"
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
    }
}