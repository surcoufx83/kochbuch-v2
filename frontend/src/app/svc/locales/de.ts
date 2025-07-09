import { L10nLocale } from "./types";

export const KB_De: L10nLocale = {
    collections: {
        defaultCollection: {
            description: "",
            title: "Mein Kochbuch",
        },
    },
    common: {
        language: {
            de: "Deutsch",
            en: "Englisch",
            fr: "Französisch"
        },
        unknownUser: "Ehemaliger Benutzer",
    },
    errorPages: {
        loginToCreateRecipe: {
            title: "Mach es zu deinem Kochbuch 🧑‍🍳",
            paragraphLine1: "Du musst dich anmelden, um dein eigenes digitales Kochbuch zu nutzen.",
            paragraphLine2: "Die Nutzung ist freiwillig und kostenfrei, aber auf Familie und Freunde beschränkt. Die Anmeldung erfordert ein von uns bereitgestelltes Nextcloud-Konto."
        },
        routeNotFound: {
            title1: "Hoppla!",
            title2: "Seite nicht gefunden 🧑‍🍳🚫",
            paragraphLine1: "Es sieht so aus, als wärst du auf einen Fehler gestoßen!",
            paragraphLine2: "Aber keine Sorge! Hier sind ein paar Möglichkeiten:",
            optionLink1: "Zurück zum Hauptmenü – Dort findest du garantiert schmackhafte Alternativen.",
            optionLink2: "Eine andere Zutat ausprobieren – Vielleicht führt eine neue Suche zum gewünschten Rezept.",
            optionLink3: "Oder etwas ganz anderes – Wie wäre es mit einem Kaffee? ☕ Oder doch lieber ein Stück Kuchen? 🍰"
        }
    },
    floatingMenu: {
        searchButton: {
            ariaLabel: "Dieser Button öffnet ein Suchfeld, um in der Rezeptdatenbank zu suchen. Zum Abschließen der Suche wird die Enter-Taste der Tastatur verwendet.",
            searchIconAriaLabel: "Symbol, das eine Lupe anzeigt. Ist die Suche geöffnet, wird ein großes X als Symbol angezeigt, das darauf hinweist, dass ein Klick darauf das Suchfeld schließt.",
            searchInput: {
                ariaLabel: "In dieses Suchfeld wird deine Suchanfrage eingegeben. Zum Beispiel das Wort \"Lasagne\", um nach solchen Gerichten zu suchen.",
                placeholder: {
                    jan: "Suchbegriff (z.B. Grünkohl mit Pinkel)",
                    feb: "Suchbegriff (z.B. Krapfen)",
                    mar: "Suchbegriff (z.B. Lammbraten mit Rosmarin)",
                    apr: "Suchbegriff (z.B. Spargelcremesuppe)",
                    may: "Suchbegriff (z.B. Wildragout)",
                    jun: "Suchbegriff (z.B. Frischer Sommersalat)",
                    jul: "Suchbegriff (z.B. Obsttorte mit Erdbeeren)",
                    aug: "Suchbegriff (z.B. Gegrilltes Gemüse mit Feta)",
                    sep: "Suchbegriff (z.B. Kürbissuppe mit Ingwer)",
                    oct: "Suchbegriff (z.B. Kartoffelgratin)",
                    nov: "Suchbegriff (z.B. Gänsebraten mit Rotkohl)",
                    dec: "Suchbegriff (z.B. Weihnachtsgans)"
                }
            },
            submitIconAriaLabel: "Dieser Button startet die Suche mit dem eingegebenen Suchbegriff. Nach dem Abschicken wirst du zur Ergebnisseite weitergeleitet."
        }
    },
    homeRecipesPage: {
        titleUser: "Hallo [0]",
        titleGuest: "Hallo Besucher",
        descriptionUser: "Hier findest du die neuesten Rezepte im Kochbuch. Die Suchfunktion befindet sich am unteren Rand des Fensters.",
        descriptionGuest: "Schön, dass du ins Kochbuch schaust. Bitte melde dich über den Benutzer-Button unten im Fenster an, um Zugriff auf alle Rezepte zu erhalten."
    },
    login: {
        loginWithNcButton: "Anmelden mit Nextcloud-Konto"
    },
    me: {
        title: "Dein digitales Kochbuch",
        description: "Das ist dein persönlicher Bereich, hier findest du alle deine Rezepte - egal ob von dir erstellt oder gespeichert.",
        collections: {
            title: "Sammlungen",
            description: "Rezepte werden in Sammlungen organisiert. Wir haben für dich automatisch eine Standard-Sammlung angelegt, du kannst aber eine beliebige Anzahl weiterer anlegen, zum Beispiel für die besten Backrezepte, Deserts, Cocktails."
        },
    },
    navbar: {
        brand: {
            pageTitle: "Kochbuch",
            iconLabel: "Kochbuch-Icon dieser Webseite"
        }
    },
    recipe: {
        aiLocalizationSwitch: [
            "Originalrezept anzeigen",
            "Übersetzung anzeigen"
        ],
        aiLocalizedContent: "Dieses Rezept wurde automatisch übersetzt.",
        aiLocalizedContentWithSourceLocale: "Dieses Rezept wurde automatisch aus [0] übersetzt.",
        aiLocalizedContentSourceLocale: {
            de: "dem Deutschen",
            en: "dem Englischen",
            fr: "dem Französischen"
        },
        aiSourceLocale: "Originaltitel: [0]",
        difficulty: {
            0: "Keine Angabe",
            1: "Leicht",
            2: "Mittelschwer",
            3: "Schwierig"
        },
        hasAdminAccess: "Nur sichtbar, weil du als Admin angemeldet bist.",
        isLoading: "Rezeptdaten werden gerade geladen ...",
        loadingFailed: "Das Rezept konnte nicht geladen werden. Du wirst automatisch auf die Startseite umgeleitet...",
        modifiedAndUser: "[0] von ",
        printBtn: "Drucken",
        saveBtn: "Speichern",
        saveBtnAlreadySaved: "Gespeichert",
        submittedBy: "Eingetragen von [0]",
        shareBtn: "Teilen",
        share: {
            title: "Rezept für ‘[0]’ auf kochbuch.mogul.network",
            message: "*[0]*\n\nDieses und weitere Rezepte von [1] findest du in unserem digitalen Kochbuch. Happy Nachkochen!",
        },
    },
    recipeGallery: {
        noPicturesUploaded: "Noch kein Foto vorhanden. Klicke hier um eines aufzunehmen...",
        uploadBtn: "Foto hochladen",
        uploadStatus: {
            checking: "Überprüfen...",
            uploading: "Wird gespeichert...",
        },
    },
    recipeIngredients: {
        title: "Zutatenliste",
        description: [
            "Das Rezept ist für 1 Portion ausgelegt.",
            "Das Rezept ist für [0] Portionen ausgelegt.",
        ],
        calculator: {
            servings: [
                "Portion",
                "Portionen",
            ],
            title: "Mengenrechner",
            description: "Du kannst dir die Mengen auf eine andere Anzahl von Portionen umrechnen lassen. Die Umrechnung wird auch auf die Zubereitungsschritte weiter unten angewendet.",
        },
        table: {
            quantityHeader: "Mengenangabe",
            nameHeader: "Zutatenbeschreibung",
        },
    },
    recipeOwnerInfo: {
        title: "Bearbeitungshinweis",
        description: "Dieses Rezept wurde von dir selbst erstellt. Du kannst es jederzeit bearbeiten oder die Sichtbarkeit anpassen. Die entsprechenden Optionen findest du weiter unten nach den Zubereitungsschritten.",
        adminDescription: "Dieses Rezept wurde von [0] erstellt, ist aber für dich sichtbar, da du als Admin angemeldet bist. Optionen zum Verwalten des Rezepts findest du weiter unten nach den Zubereitungsschritten.",
        gothereLink: "Dahin springen"
    },
    recipePreparation: {
        title: "Zubereitung",
        stepFormat: "[0]. [1]",
        stepFallback: "Schritt",
    },
    recipePreparationTime: {
        title: "Zubereitungsdauer",
        longDurationWarning: "Bitte beachte, dass für dieses Rezept eine Gesamt-Zubereitungsdauer von mindestens [0] angegeben ist.",
        recalcWarning: {
            screen: [
                "Die Zubereitungszeit basiert auf der Angabe von 1 Portion und wird nicht automatisch für mehr Portionen angepasst. Bitte berücksichtige das bei deiner Vorbereitung.",
                "Die Zubereitungszeit basiert auf der Angabe von [0] Portionen und wird nicht automatisch für mehr Portionen angepasst. Bitte berücksichtige das bei deiner Vorbereitung.",
            ],
            print: [
                "Die Zubereitungszeit gilt für 1 Portion!",
                "Die Zubereitungszeit gilt für [0] Portionen!",
            ],
        },
        items: {
            cooking: ["Koch-/Backzeit", "[0] Kochen/Backen"],
            preparing: ["Vorbereitungszeit", "[0] Zubereitung"],
            total: "Gesamtzeit",
            waiting: ["Ruhezeit", "[0] Ruhe"],
        },
        units: {
            days: ['1 Tag', '[0] Tage', '1 Tag', '[0] Tage'],
            hours: ['1 Std.', '[0] Std.', '1 Stunde', '[0] Stunden'],
            minutes: ['1 Min.', '[0] Min.', '1 Minute', '[0] Minuten'],
        },
    },
    saveToCollection: {
        confirmMsg: "Rezept wurde gespeichert.",
        gotoCollectionLink: [
            "Zur Rezeptesammlung",
            "Zu deinen Sammlungen",
        ],
        pickCollection: {
            description: "Wähle eine Sammlung in die das Rezept gespeichert wird.",
            itemCount: [
                "Enthält noch keine Rezepte",
                "Enthält ein Rezept",
                "Enthält [0] Rezepte",
            ],
            titleInputLabel: "Name für die Rezeptesammlung:",
            newBtn: "Neue Sammlung",
            saveBtn: "Speichern",
        },
    },
}