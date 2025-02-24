export type NameLocalization = {
    [languageCode: string]: { name: string },
}

export type CategoryItem = {
    id: number,
    localization: NameLocalization,
    icon: string,
    modified: string,
}

export type Category = {
    id: number,
    localization: NameLocalization,
    icon: string,
    modified: string,
    items: { [id: number]: CategoryItem },
}

export type Group = {
    displayname: string,
    id: number,
    name: string,
}

export type Recipe = {
    aiGenerated: boolean,
    aiLocalized: boolean,
    categories: number[],
    created: string,
    difficulty: 0 | 1 | 2 | 3,
    id: number,
    isFork: boolean,
    isPlaceholder: number,
    localization: RecipeLocalization,
    modified: string,
    originalRecipeId: number,
    published: string,
    servingsCount: number,
    sharedInternal: boolean,
    sharedPublic: boolean,
    simple: boolean,
    sourceUrl: string,
    userLocale: string,
}

export type RecipeLocalization = {
    [languageCode: string]: {
        description: string,
        sourceDescription: string,
        title: string,
    },
}

export type UserSelf = {
    admin: boolean,
    created: string,
    displayname: string,
    email: string | null,
    enabled: boolean,
    firstname: string,
    groups: Group[],
    id: number,
    lastname: string,
    username: string,
};

export type User = {
    displayname: string,
    id: number,
};
