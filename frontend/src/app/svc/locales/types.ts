export type L10nLocale = {
    errorPages: {
        loginToCreateRecipe: {
            title: string,
            paragraphLine1: string,
            paragraphLine2: string,
        },
        routeNotFound: {
            title1: string,
            title2: string,
            paragraphLine1: string,
            paragraphLine2: string,
            optionLink1: string,
            optionLink2: string,
            optionLink3: string,
        }
    },
    floatingMenu: {
        searchButton: {
            ariaLabel: string,
            searchIconAriaLabel: string,
            searchInput: {
                ariaLabel: string,
                placeholder: {
                    jan: string,
                    feb: string,
                    mar: string,
                    apr: string,
                    may: string,
                    jun: string,
                    jul: string,
                    aug: string,
                    sep: string,
                    oct: string,
                    nov: string,
                    dec: string,
                },
            }
            submitIconAriaLabel: string,
        }
    },
    homeRecipesPage: {
        titleUser: string,
        titleGuest: string,
        descriptionUser: string,
        descriptionGuest: string,
    },
    login: {
        loginWithNcButton: string,
    }
    navbar: {
        brand: {
            pageTitle: string,
            iconLabel: string,
        }
    },
    recipe: {
        aiLocalizedContent: string,
        submittedBy: string,
    }
}