import { L10nLocale } from "./types";

export const KB_De: L10nLocale = {
    navbar: {
        brand: {
            pageTitle: "Kochbuch",
            iconLabel: "Kochbuch-Icon dieser Webseite"
        }
    },
    floatingMenu: {
        searchButton: {
            ariaLabel: "Dieser Button öffnet ein Suchfeld, um in der Rezeptdatenbank zu suchen. Zum Abschließen der Suche wird die Enter-Taste der Tastatur verwendet.",
            searchIconAriaLabel: "Symbol, das eine Lupe anzeigt. Ist die Suche geöffnet, wird ein großes X als Symbol angezeigt, das darauf hinweist, dass ein Klick darauf das Suchfeld schließt.",
            searchInput: {
                ariaLabel: "In dieses Suchfeld wird deine Suchanfrage eingegeben. Zum Beispiel das Wort \"Lasagne\", um nach solchen Gerichten zu suchen.",
                placeholder: {
                    jan: "Suchbegriff eingeben (z.B. Grünkohl mit Pinkel)",
                    feb: "Suchbegriff eingeben (z.B. Krapfen)",
                    mar: "Suchbegriff eingeben (z.B. Lammbraten mit Rosmarin)",
                    apr: "Suchbegriff eingeben (z.B. Spargelcremesuppe)",
                    may: "Suchbegriff eingeben (z.B. Wildragout)",
                    jun: "Suchbegriff eingeben (z.B. Frischer Sommersalat)",
                    jul: "Suchbegriff eingeben (z.B. Obsttorte mit Erdbeeren)",
                    aug: "Suchbegriff eingeben (z.B. Gegrilltes Gemüse mit Feta)",
                    sep: "Suchbegriff eingeben (z.B. Kürbissuppe mit Ingwer)",
                    oct: "Suchbegriff eingeben (z.B. Kartoffelgratin)",
                    nov: "Suchbegriff eingeben (z.B. Gänsebraten mit Rotkohl)",
                    dec: "Suchbegriff eingeben (z.B. Weihnachtsgans)"
                }
            },
            submitIconAriaLabel: "Dieser Button startet die Suche mit dem eingegebenen Suchbegriff. Nach dem Abschicken wirst du zur Ergebnisseite weitergeleitet."
        }
    }
}