import { L10nLocale } from "./types";

export const KB_En: L10nLocale = {
    common: {
        language: {
            de: "German",
            en: "English",
            fr: "French"
        }
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
        isLoading: "Recipe data is currently loading...",
        loadingFailed: "The recipe could not be loaded. You will be automatically redirected to the homepage...",
        submittedBy: "Submitted by [0]"
    },
    recipePreparationTime: {
        title: "Preparation Time",
        items: {
            cooking: ["Cooking/Baking Time", "Cooking/Baking"],
            preparing: ["Preparation Time", "Preparing"],
            total: "Total Time",
            waiting: ["Resting Time", "Resting"]
        }
    }
}