export type L10nLocale = {
    common: {
        language: { [key: string]: string },
    },
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
        aiLocalizationSwitch: string[],
        aiLocalizedContent: string,
        aiLocalizedContentWithSourceLocale: string,
        aiLocalizedContentSourceLocale: { [key: string]: string },
        aiSourceLocale: string,
        difficulty: {
            0: string,
            1: string,
            2: string,
            3: string,
        },
        isLoading: string,
        loadingFailed: string,
        submittedBy: string,
    },
    recipeIngredients: {
        title: string,
        description: string[],
        calculator: {
            title: string,
            servings: string[],
            description: string,
        },
        table: {
            quantityHeader: string,
            nameHeader: string,
        },
    },
    recipeOwnerInfo: {
        title: string,
        description: string,
        gothereLink: string,
    },
    recipePreparation: {
        title: string,
        stepFormat: string,
        stepFallback: string,
    },
    recipePreparationTime: {
        title: string,
        items: {
            cooking: string[],
            preparing: string[],
            total: string,
            waiting: string[],
        },
        units: {
            days: string[],
            hours: string[],
            minutes: string[],
        }
    },
}