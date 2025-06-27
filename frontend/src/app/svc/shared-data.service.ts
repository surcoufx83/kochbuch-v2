import { Injectable } from '@angular/core';
import { BehaviorSubject, Subscription } from 'rxjs';
import { Category, Recipe, RecipePicture, Unit } from '../types';
import { IdbRecipe, IdbService } from './idb.service';
import { WebSocketService, WsMessage } from './web-socket.service';

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
  private _recipePreloadRequests: { [key: number]: number } = {};
  private _recipeUpdated = new BehaviorSubject<RecipeUpdatedEvent | false>(false);
  public RecipeEvents = this._recipeUpdated.asObservable();

  private _units = new BehaviorSubject<{ [key: number]: Unit }>({});
  private _unitsEtag?: string;
  public Units = this._units.asObservable();

  private _searchIsActive = new BehaviorSubject<boolean>(false);
  public SearchIsActive = this._searchIsActive.asObservable();

  private _searchPhrase = new BehaviorSubject<string>('');
  public SearchPhrase = this._searchPhrase.asObservable();

  private _searchCategories = new BehaviorSubject<string[]>([]);
  public SearchCategories = this._searchCategories.asObservable();

  constructor(
    private indexDbService: IdbService,
    private wsService: WebSocketService,
  ) {
    this.initialize();
  }

  private initialize(): void {
    this.loadFromBrowserCache();

    this.wsService.isLoggedIn.subscribe((state) => {
      if (state === 'unknown')
        return;

      if (!this.initialized) {
        this.lastloginstate = state;
        this.initialized = true;
        return;
      }
      if (this.lastloginstate !== state) {
        console.warn('Login state changed!')
        this.clear();
      }
    });

    this.wsService.events.subscribe((event) => {
      if (!event || event.type === 'none')
        return;
      this.wsMessageReceived(event);
    });

  }

  private clear(): void {
    this._recipes.next({});
    this._recipesEtag = '';
    localStorage.removeItem('kbRecipes');
    localStorage.removeItem('kbCategories');
    localStorage.removeItem('kbUnits');
  }

  public getRecipe(id: number, reload: boolean = true): Promise<IdbRecipe> {
    return new Promise((resolve, reject) => {
      this.indexDbService.GetRecipe(id).then((recipe) => {
        setTimeout(() => {
          this.reloadRecipe(id);
        }, 10);
        resolve(recipe);
      }).catch((err) => {
        let sub: Subscription | undefined = this.RecipeEvents.subscribe((rec) => {
          if (!rec || rec.id !== id)
            return;
          resolve({
            id: rec.id,
            etag: rec.etag ?? rec.recipe.modified,
            data: rec.recipe
          });
          sub?.unsubscribe();
          sub = undefined;
        });
        this.reloadRecipe(id);
        setTimeout(() => {
          if (sub) {
            reject('Timeout loading recipe');
            sub.unsubscribe();
          }
        }, 3500);
      });
    });
  }

  private loadFromBrowserCache(): void {
    this.loadCategoriesFromCache();
    this.loadRecipesFromCache();
    this.loadUnitsFromCache();
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

  private loadUnitsFromCache(): void {
    let unitsData: string | null = localStorage.getItem('kbUnits');
    if (unitsData !== null) {
      const units = JSON.parse(unitsData) as UnitsCache;
      this._unitsEtag = units.etag;
      this._units.next(units.units);
    }
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
          sizes.push(`(max-width: ${sizear[i]}px) ${sizear[i]}px`);
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

  public PreloadRecipeToCache(recipeId: number): void {
    if (this._recipePreloadRequests[recipeId]) {
      if ((Date.now() - this._recipePreloadRequests[recipeId]) > 60000)
        delete this._recipePreloadRequests[recipeId];
      else
        return;
    }
    this._recipePreloadRequests[recipeId] = Date.now();
    this.reloadRecipe(recipeId);
  }

  private reloadRecipe(recipeId: number): void {
    this.indexDbService.GetRecipe(recipeId)
      .then((data: { id: number, etag: string, data: Recipe }) => {
        this.wsService.SendMessage({
          type: 'recipe_get',
          content: JSON.stringify({
            id: recipeId,
            etag: data.etag
          })
        });
      })
      .catch(() => {
        this.wsService.SendMessage({
          type: 'recipe_get',
          content: JSON.stringify({
            id: recipeId
          })
        });
      });
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
    Object.values(this._recipes.value).forEach((recipe) => {
      this._recipeUpdated.next({
        id: recipe.id,
        etag: recipe.modified,
        recipe: recipe,
      });
    });
  }

  private saveUnitsToCache(): void {
    const cache: UnitsCache = {
      units: this._units.value,
      etag: this._recipesEtag,
    };
    localStorage.setItem('kbUnits', JSON.stringify(cache));
  }

  public SetTitle(title: string): void {
    this._pageTitle.next(title);
  }

  public SetSearchState(newstate: boolean): void {
    if (newstate !== this._searchIsActive.value)
      this._searchIsActive.next(newstate);
  }

  private wsMessageReceived(msg: WsMessage): void {
    console.log('wsMessageReceived', msg)

    switch (msg.type) {

      case 'categories_etag':
        if (this._categoriesEtag !== <string>msg.content) {
          this.wsService.SendMessage({
            type: 'categories_get_all',
            content: JSON.stringify({})
          });
        }
        break;

      case 'categories':
        const cdata = JSON.parse(msg.content) as CategoriesCache;
        let itemMapping: { [key: number]: number } = {};
        for (const cat of Object.values(cdata.categories)) {
          for (const item of Object.values(cat.items)) {
            itemMapping[item.id] = cat.id;
          }
        }
        this._categoryItemsMapping.next(itemMapping);
        this._categoriesEtag = cdata.etag;
        this._categories.next(cdata.categories);
        this.saveCategoriesToCache();
        break;

      case 'recipe_get':
        const recipedata = JSON.parse(msg.content) as Recipe | ErrorResponse;
        if (Object.hasOwn(recipedata, 'error')) {
          return;
        }
        this.indexDbService.PutRecipe(recipedata as Recipe);
        this._recipeUpdated.next({
          id: (recipedata as Recipe).id,
          etag: (recipedata as Recipe).modified,
          recipe: (recipedata as Recipe),
        });
        break;

      case 'recipes_etag':
        if (this._recipesEtag !== <string>msg.content) {
          this.wsService.SendMessage({
            type: 'recipes_get_all',
            content: JSON.stringify({})
          });
        }
        break;

      case 'recipes':
        const rdata = JSON.parse(msg.content) as RecipesCache;
        this._recipesEtag = rdata.etag;
        this._recipes.next(this.loadRecipes_GeneratePictureSets(rdata.recipes));
        this.saveRecipesToCache();
        break;

      case 'units_etag':
        if (this._unitsEtag !== <string>msg.content) {
          this.wsService.SendMessage({
            type: 'units_get_all',
            content: JSON.stringify({})
          });
        }
        break;

      case 'units':
        const udata = JSON.parse(msg.content) as UnitsCache;
        this._unitsEtag = udata.etag;
        this._units.next(udata.units);
        this.saveUnitsToCache();
        break;

    }
  }

}

type CategoriesCache = {
  categories: { [key: number]: Category };
  categoryItemsMapping: { [key: number]: number };
  etag?: string;
}

type ErrorResponse = {
  error: number;
  message: string;
}

type RecipesCache = {
  recipes: { [key: number]: Recipe };
  etag?: string;
}

type RecipeUpdatedEvent = {
  id: number;
  etag?: string;
  recipe: Recipe;
}

type UnitsCache = {
  units: { [key: number]: Unit };
  etag?: string;
}