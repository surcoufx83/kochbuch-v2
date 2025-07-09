export type L10nLocale = {
    collections: {
        defaultCollection: {
            title: string,
            description: string,
        },
    },
    common: {
        language: { [key: string]: string },
        unknownUser: string,
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
    },
    me: {
        title: string,
        description: string,
        collections: {
            title: string,
            description: string,
        }
    },
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
        hasAdminAccess: string,
        isLoading: string,
        loadingFailed: string,
        modifiedAndUser: string,
        printBtn: string,
        saveBtn: string,
        saveBtnAlreadySaved: string,
        submittedBy: string,
        shareBtn: string,
        share: {
            title: string,
            message: string,
        },
    },
    recipeGallery: {
        noPicturesUploaded: string,
        uploadBtn: string,
        uploadStatus: {
            checking: string,
            uploading: string,
        },
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
    recipeManagement: {
        title: string,
        adminDescription: string,
        gotoOwnerBtn: string[],
        ownerDescription: string,
        created: string,
        modified: string,
        delete: {
            title: string,
            description: string,
            btn: string,
        },
        edit: {
            title: string,
            description: string,
            btn: string,
        },
        publish: {
            title: string,
            description: string,
            togglePrivate: string,
            toggleInternal: string,
            togglePublic: string,
            descriptionPrivate: string,
            descriptionInternal: string,
            descriptionPublic: string,
        },
    },
    recipeOwnerInfo: {
        title: string,
        description: string,
        adminDescription: string,
        gothereLink: string,
    },
    recipePreparation: {
        title: string,
        stepFormat: string,
        stepFallback: string,
    },
    recipePreparationTime: {
        title: string,
        longDurationWarning: string,
        recalcWarning: {
            screen: string[],
            print: string[],
        },
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
    saveToCollection: {
        confirmMsg: string,
        gotoCollectionLink: string[],
        pickCollection: {
            description: string,
            itemCount: string[],
            titleInputLabel: string,
            newBtn: string,
            saveBtn: string,
        }
    }
}