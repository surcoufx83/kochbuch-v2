export type NameLocalization = {
    [languageCode: string]: { name: string };
};

export type CategoryItem = {
    id: number;
    localization: NameLocalization;
    icon: string;
    modified: string;
};

export type Category = {
    id: number;
    localization: NameLocalization;
    icon: string;
    modified: string;
    items: { [id: number]: CategoryItem };
};

export type UserSelf = {};

export type User = {};
