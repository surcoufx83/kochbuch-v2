import { L10nLocale } from "./types";

export const KB_En: L10nLocale = {
    collections: {
        defaultCollection: {
            description: "",
            title: "My Cookbook",
        },
    },
    common: {
        language: {
            de: "German",
            en: "English",
            fr: "French"
        },
        unknownUser: "Former user",
    },
    errorPages: {
        loginToCreateRecipe: {
            title: "Make it your own cookbook 🧑‍🍳",
            paragraphLine1: "You need to log in to use your own digital cookbook.",
            paragraphLine2: "Usage is voluntary and free, but limited to family and friends. Logging in requires a Nextcloud account provided by us."
        },
        routeNotFound: {
            title1: "Oops!",
            title2: "Page not found 🧑‍🍳🚫",
            paragraphLine1: "It looks like you've stumbled upon an error!",
            paragraphLine2: "But no worries! Here are a few options:",
            optionLink1: "Back to the main menu – You'll surely find some tasty alternatives there.",
            optionLink2: "Try a different ingredient – Maybe a new search will lead you to the desired recipe.",
            optionLink3: "Or something completely different – How about a coffee? ☕ Or maybe a slice of cake? 🍰"
        }
    },
    floatingMenu: {
        searchButton: {
            ariaLabel: "This button opens a search field to search the recipe database. To complete the search, use the Enter key on the keyboard.",
            searchIconAriaLabel: "Icon showing a magnifying glass. If the search is open, a large X is displayed as an icon, indicating that clicking it will close the search field.",
            searchInput: {
                ariaLabel: "Enter your search query into this field. For example, the word \"Shepherd's Pie\" to search for such dishes.",
                placeholder: {
                    jan: "Search term (e.g. Haggis with Neeps and Tatties)",
                    feb: "Search term (e.g. Pancakes with Lemon and Sugar)",
                    mar: "Search term (e.g. Roast Lamb with Rosemary)",
                    apr: "Search term (e.g. Creamy Asparagus Soup)",
                    may: "Search term (e.g. Venison Stew)",
                    jun: "Search term (e.g. Fresh Summer Salad)",
                    jul: "Search term (e.g. Strawberry Shortcake)",
                    aug: "Search term (e.g. Grilled Vegetables with Halloumi)",
                    sep: "Search term (e.g. Pumpkin Soup with Ginger)",
                    oct: "Search term (e.g. Traditional Potato Gratin)",
                    nov: "Search term (e.g. Roast Goose with Braised Red Cabbage)",
                    dec: "Search term (e.g. Christmas Pudding)"
                }
            },
            submitIconAriaLabel: "This button starts the search with the entered search term. After submission, you will be redirected to the results page."
        }
    },
    homeRecipesPage: {
        titleUser: "Hello [0]",
        titleGuest: "Hello Visitor",
        descriptionUser: "Here you will find the latest recipes in the cookbook. You can find the search function at the bottom of the window.",
        descriptionGuest: "Nice to see you checking out the cookbook. Please log in via the user button at the bottom of the window to access all recipes."
    },
    login: {
        loginWithNcButton: "Log in with Nextcloud account"
    },
    me: {
        title: "Your Digital Cookbook",
        description: "This is your personal space — here you’ll find all your recipes, whether you created or saved them.",
        collections: {
            title: "Collections",
            description: "Recipes are organized in collections. We’ve automatically created a default collection for you, but you can create as many as you like — for example, for the best baking recipes, desserts, or cocktails."
        }
    },
    navbar: {
        brand: {
            pageTitle: "Cookbook",
            iconLabel: "Cookbook icon of this website"
        }
    },
    recipe: {
        aiLocalizationSwitch: [
            "Show original recipe",
            "Show translation"
        ],
        aiLocalizedContent: "This recipe has been automatically translated.",
        aiLocalizedContentWithSourceLocale: "This recipe has been automatically translated from [0].",
        aiLocalizedContentSourceLocale: {
            de: "German",
            en: "English",
            fr: "French"
        },
        aiSourceLocale: "Original title: [0]",
        difficulty: {
            0: "No information",
            1: "Easy",
            2: "Medium",
            3: "Difficult"
        },
        hasAdminAccess: "Visible only because you’re logged in as an admin.",
        isLoading: "Recipe data is currently loading...",
        loadingFailed: "The recipe could not be loaded. You will be automatically redirected to the homepage...",
        modifiedAndUser: "[0] by ",
        printBtn: "Print",
        saveBtn: "Save",
        saveBtnAlreadySaved: "Saved",
        submittedBy: "Submitted by [0]",
        shareBtn: "Share",
        share: {
            title: "Recipe for ‘[0]’ on kochbuch.mogul.network",
            message: "*[0]*\n\nYou’ll find this and more recipes by [1] in our digital cookbook. Enjoy cooking!"
        },
    },
    recipeGallery: {
        noPicturesUploaded: "No photo yet. Click here to take one...",
        uploadBtn: "Upload photo",
        uploadStatus: {
            checking: "Checking...",
            uploading: "Saving...",
        },
    },
    recipeIngredients: {
        title: "Ingredients List",
        description: [
            "This recipe is intended for 1 serving.",
            "This recipe is intended for [0] servings.",
        ],
        calculator: {
            servings: [
                "Serving",
                "Servings",
            ],
            title: "Quantity Calculator",
            description: "You can recalculate the quantities for a different number of servings. The conversion also applies to the preparation steps below.",
        },
        table: {
            quantityHeader: "Quantity",
            nameHeader: "Ingredient Description",
        },
    },
    recipeManagement: {
        title: "Manage Recipe",
        adminDescription: "This recipe was created by [0]. As an administrator, you can edit the recipe and change all related settings.",
        gotoOwnerBtn: [
            "Go to [0]'s account",
            "Go to [0]'s account"
        ],
        ownerDescription: "You can edit your own recipes or adjust their visibility at any time. All options are available in this section.",
        created: "Created on [0]",
        modified: "Last modified on [0]",
        delete: {
            title: "Delete",
            description: "Want to get rid of the recipe? Use the delete button to remove it. This action cannot be undone. All data, images, and comments will also be deleted.",
            btn: "Delete Recipe"
        },
        edit: {
            title: "Edit",
            description: "In edit mode, you can revise all content of the recipe.",
            btn: "Edit"
        },
        publish: {
            title: "Publish",
            description: "Recipes are private by default. You can also make your recipe visible to others and share it via the share button.",
            togglePrivate: "Private",
            toggleInternal: "Logged-in Users",
            togglePublic: "Public",
            descriptionPrivate: "Private means private. This recipe belongs only to you and cannot be viewed or accessed by anyone else. This is the default for newly created recipes. You can print it, but sharing via WhatsApp is not possible.",
            descriptionInternal: "Logged-in users – la familia: This level allows other users with a Nextcloud account to view your recipe. It makes it accessible to family and close friends. Sharing via WhatsApp is still not possible.",
            descriptionPublic: "Public: Anyone in the world can view the recipe if they have the link. This level allows free sharing with others and spreading the recipe among your social circle."
        }
    },
    recipeOwnerInfo: {
        title: "Edit Notice",
        description: "This recipe was created by you. You can edit it or adjust its visibility at any time. The corresponding options can be found further below, after the preparation steps.",
        adminDescription: "This recipe was created by [0], but is visible to you because you’re logged in as an admin. Management options for this recipe can be found below the preparation steps.",
        gothereLink: "Jump there"
    },
    recipePreparation: {
        title: "Preparation",
        stepFormat: "[0]. [1]",
        stepFallback: "step",
    },
    recipePreparationTime: {
        title: "Preparation Time",
        longDurationWarning: "Please note that the total preparation time for this recipe is at least [0].",
        recalcWarning: {
            screen: [
                "The preparation time is based on 1 serving and is not automatically adjusted for more servings. Please take that into account when preparing.",
                "The preparation time is based on [0] servings and is not automatically adjusted for more servings. Please take that into account when preparing."
            ],
            print: [
                "Preparation time applies to 1 serving!",
                "Preparation time applies to [0] servings!"
            ]
        },
        items: {
            cooking: ["Cooking/Baking Time", "[0] cook"],
            preparing: ["Preparation Time", "[0] prep"],
            total: "Total Time",
            waiting: ["Resting Time", "[0] rest"]
        },
        units: {
            days: ['1 d', '[0] d', '1 day', '[0] days'],
            hours: ['1 h', '[0] h', '1 hour', '[0] hours'],
            minutes: ['1 min', '[0] min', '1 minute', '[0] minutes']
        },
    },
    saveToCollection: {
        confirmMsg: "Recipe has been saved.",
        gotoCollectionLink: [
            "Go to recipe collection",
            "Go to your collections",
        ],
        pickCollection: {
            description: "Choose a collection to save the recipe to.",
            itemCount: [
                "Contains no recipes yet",
                "Contains one recipe",
                "Contains [0] recipes",
            ],
            titleInputLabel: "Name for the recipe collection:",
            newBtn: "New Collection",
            saveBtn: "Save",
        },
    },
}