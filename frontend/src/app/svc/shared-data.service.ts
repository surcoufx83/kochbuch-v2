import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';
import { ApiService } from './api.service';
import { Category, Recipe, RecipePicture } from '../types';
import { HttpStatusCode } from '@angular/common/http';

@Injectable({
  providedIn: 'root'
})
export class SharedDataService {

  private initialized = false;
  private lastloginstate: 'unknown' | boolean = false;

  private _pageTitle = new BehaviorSubject<string>('');
  public PageTitle = this._pageTitle.asObservable();

  private _categories = new BehaviorSubject<{ [key: number]: Category }>({});
  private _categoriesEtag?: string;
  private _categoryItemsMapping = new BehaviorSubject<{ [key: number]: number }>({}); // item_id -> category_id
  public Categories = this._categories.asObservable();

  private _recipes = new BehaviorSubject<{ [key: number]: Recipe }>({});
  private _recipesEtag?: string;
  public Recipes = this._recipes.asObservable();

  private _searchIsActive = new BehaviorSubject<boolean>(false);
  public SearchIsActive = this._searchIsActive.asObservable();

  private _searchPhrase = new BehaviorSubject<string>('');
  public SearchPhrase = this._searchPhrase.asObservable();

  private _searchCategories = new BehaviorSubject<string[]>([]);
  public SearchCategories = this._searchCategories.asObservable();

  constructor(
    private apiService: ApiService,
  ) {
    this.apiService.isLoggedIn.subscribe((state) => {
      if (state === 'unknown')
        return;
      if (!this.initialized) {
        this.lastloginstate = state;
        this.initialized = true;
        return;
      }
      if (this.lastloginstate !== state) {
        console.warn('Login state changed!')
        const localecache = localStorage.getItem('kbLocale');
        localStorage.clear();
        if (localecache)
          localStorage.setItem('kbLocale', localecache);
        this.reloadEntitiesFromServer();
      }
    });
    this.loadFromBrowserCache();
    this.reloadEntitiesFromServer();
  }

  private loadFromBrowserCache(): void {
    this.loadCategoriesFromCache();
    this.loadRecipesFromCache();
  }

  private loadCategoriesFromCache(): void {
    let categoriesData: string | null = localStorage.getItem('kbCategories');
    if (categoriesData !== null) {
      const categories = JSON.parse(categoriesData) as CategoriesCache;
      this._categories.next(categories.categories);
      this._categoryItemsMapping.next(categories.categoryItemsMapping);
      this._categoriesEtag = categories.etag;
    }
  }

  private loadRecipesFromCache(): void {
    let recipesData: string | null = localStorage.getItem('kbRecipes');
    if (recipesData !== null) {
      const recipes = JSON.parse(recipesData) as RecipesCache;
      this._recipes.next(
        this.loadRecipes_GeneratePictureSets(recipes.recipes)
      );
      this._recipesEtag = recipes.etag;
    }
  }

  private loadCategoriesFromServer(): void {
    this.apiService.get('categories', this._categoriesEtag).subscribe((res) => {
      if (res?.status === HttpStatusCode.Ok) {
        const categories = (res as CategoriesResponse).body.categories;
        let itemMapping: { [key: number]: number } = {};
        for (const cat of Object.values(categories)) {
          for (const item of Object.values(cat.items)) {
            itemMapping[item.id] = cat.id;
          }
        }
        this._categories.next(categories);
        this._categoryItemsMapping.next(itemMapping);
        this._categoriesEtag = res.headers.get('etag') ?? undefined;
        this.saveCategoriesToCache();
      }
    });
  }

  private loadRecipesFromServer(): void {
    this.apiService.get('recipes', this._recipesEtag).subscribe((res) => {
      if (res?.status === HttpStatusCode.Ok) {
        const recipes =
          this.loadRecipes_GeneratePictureSets((res as RecipesResponse).body.recipes);
        this._recipes.next(Object.values(recipes));
        this._recipesEtag = res.headers.get('etag') ?? undefined;
        this.saveRecipesToCache();
      }
    });
  }

  private loadRecipes_GeneratePictureSets(recipes: { [key: number]: Recipe }): { [key: number]: Recipe } {
    for (const key of Object.keys(recipes)) {
      const recipe = recipes[+key];
      if (!recipe.pictures || !Array.isArray(recipe.pictures)) {
        continue;
      }
      for (let i = 0; i < recipe.pictures.length; i++) {
        recipe.pictures[i] = this.loadRecipes_GeneratePictureSet(recipe.id, recipe.pictures[i]);
      }
    }
    return recipes;
  }

  private loadRecipes_GeneratePictureSet(recipeid: number, picture: RecipePicture): RecipePicture {
    picture.htmlSrc = `/api/media/uploads/${recipeid}/${picture.id}/${picture.filename}`;
    if (picture.size.thbSizes.length > 0) {
      let srcset: string[] = [];
      let sizes: string[] = [];
      const sizear = picture.size.thbSizes.sort((a, b) => a - b);
      for (let i = 0; i < sizear.length; i++) {
        srcset.push(`/api/media/uploads/${recipeid}/${picture.id}/thb/${sizear[i]}/${picture.filename} ${sizear[i]}w`);
        if (sizear[i + 1]) {
          const maxw = sizear[i] + Math.floor((sizear[i + 1] - sizear[i]) * .66);
          sizes.push(`(max-width: ${maxw}px) ${sizear[i]}px`);
        }
        else {
          sizes.push(`${sizear[i]}px`);
        }
      }
      picture.htmlSrcSet = srcset.join(", ")
      picture.htmlSizes = sizes.join(", ")
    }
    return picture
  }

  private reloadEntitiesFromServer(): void {
    this.loadCategoriesFromServer();
    this.loadRecipesFromServer();
  }

  private saveCategoriesToCache(): void {
    const cache: CategoriesCache = {
      categories: this._categories.value,
      categoryItemsMapping: this._categoryItemsMapping.value,
      etag: this._categoriesEtag,
    };
    localStorage.setItem('kbCategories', JSON.stringify(cache));
  }

  private saveRecipesToCache(): void {
    const cache: RecipesCache = {
      recipes: this._recipes.value,
      etag: this._recipesEtag,
    };
    localStorage.setItem('kbRecipes', JSON.stringify(cache));
  }

  public SetTitle(title: string): void {
    this._pageTitle.next(title);
  }

  public SetSearchState(newstate: boolean): void {
    if (newstate !== this._searchIsActive.value)
      this._searchIsActive.next(newstate);
  }

}

type CategoriesResponse = {
  body: {
    categories: { [key: number]: Category };
  }
}

type CategoriesCache = {
  categories: { [key: number]: Category };
  categoryItemsMapping: { [key: number]: number };
  etag?: string;
}

type RecipesResponse = {
  body: {
    recipes: { [key: number]: Recipe };
  }
}

type RecipesCache = {
  recipes: { [key: number]: Recipe };
  etag?: string;
}