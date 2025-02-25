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
    categories: RecipeCategoryitem[],
    created: string,
    difficulty: 0 | 1 | 2 | 3,
    id: number,
    isFork: boolean,
    isPlaceholder: number,
    localization: RecipeLocalization,
    modified: string,
    originalRecipeId: number,
    pictures: RecipePicture[],
    preparation: RecipePreparation[],
    published: string,
    servingsCount: number,
    sharedInternal: boolean,
    sharedPublic: boolean,
    simple: boolean,
    sourceUrl: string,
    statistics: RecipeStatistics,
    timing: RecipeTiming,
    userLocale: string,
}

export type RecipeCategoryitem = {
    categoryitem: number,
    created: string,
    user: User
}

export type RecipeIngredient = {
    id: number,
    index: number,
    localization: RecipeIngredientLocalization,
    quantity: number | null
    unitId: number | null
}

export type RecipeIngredientLocalization = {
    [languageCode: string]: {
        title: string,
    },
}

export type RecipeLocalization = {
    [languageCode: string]: {
        description: string,
        sourceDescription: string,
        title: string,
    },
}

export type RecipePicture = {
    id: number,
    index: number,
    localization: RecipePictureLocalization,
    size: RecipePictureSize,
    uploaded: string,
    user: User,
    filename: string,
}

export type RecipePictureLocalization = {
    [languageCode: string]: {
        description: string,
        name: string,
    },
}

export type RecipePictureSize = {
    height: number,
    width: number,
}

export type RecipePreparation = {
    id: number,
    index: number,
    ingredients: RecipeIngredient[],
    localization: RecipePreparationLocalization,
    timing: RecipeTiming,
}

export type RecipePreparationLocalization = {
    [languageCode: string]: {
        title: string,
        instruct: string,
    },
}

export type RecipeStatistics = {
    cooked: number,
    ratings: RecipeStatisticsItem,
    steps: number,
    views: number,
    votes: RecipeStatisticsItem,
}

export type RecipeStatisticsItem = {
    avg: number,
    count: number,
}

export type RecipeTiming = {
    cooking: number | null,
    preparing: number | null,
    waiting: number | null,
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
