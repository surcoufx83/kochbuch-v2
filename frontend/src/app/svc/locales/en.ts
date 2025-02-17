import { L10nLocale } from "./types";

export const KB_En: L10nLocale = {
    navbar: {
        brand: {
            pageTitle: "Cookbook",
            iconLabel: "Cookbook icon of this website"
        }
    },
    floatingMenu: {
        searchButton: {
            ariaLabel: "This button opens a search field to search the recipe database. To complete the search, use the Enter key on the keyboard.",
            searchIconAriaLabel: "Icon showing a magnifying glass. If the search is open, a large X is displayed as an icon, indicating that clicking it will close the search field.",
            searchInput: {
                ariaLabel: "Enter your search query into this field. For example, the word \"Shepherd's Pie\" to search for such dishes.",
                placeholder: {
                    jan: "Enter search term (e.g. Haggis with Neeps and Tatties)",
                    feb: "Enter search term (e.g. Pancakes with Lemon and Sugar)",
                    mar: "Enter search term (e.g. Roast Lamb with Mint Sauce)",
                    apr: "Enter search term (e.g. Spring Asparagus Salad)",
                    may: "Enter search term (e.g. Venison Stew)",
                    jun: "Enter search term (e.g. Fresh Summer Salad)",
                    jul: "Enter search term (e.g. Strawberry Trifle)",
                    aug: "Enter search term (e.g. Grilled Vegetables with Halloumi)",
                    sep: "Enter search term (e.g. Pumpkin Soup with Ginger)",
                    oct: "Enter search term (e.g. Traditional Potato Gratin)",
                    nov: "Enter search term (e.g. Roast Goose with Braised Red Cabbage)",
                    dec: "Enter search term (e.g. Christmas Pudding)"
                }
            },
            submitIconAriaLabel: "This button starts the search with the entered search term. After submission, you will be redirected to the results page."
        }
    }
}