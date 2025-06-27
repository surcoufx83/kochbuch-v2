import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';
import { Recipe } from '../types';

@Injectable({
  providedIn: 'root'
})
export class IdbService {

  private dbname = 'kbdb';
  private schemaversion = 1;
  private recipesStore = 'recipes';

  private _isOpen = new BehaviorSubject<boolean>(false);
  public isOpen = this._isOpen.asObservable();

  constructor() {
    setTimeout(() => {
      this.open();
    }, 0);
  }

  public GetRecipe(id: number): Promise<IdbRecipe> {
    return new Promise((resolve, reject) => {
      const dbRequest = indexedDB.open(this.dbname);

      dbRequest.onsuccess = (event: any) => {
        const db = event.target.result;
        const transaction = db.transaction([this.recipesStore], 'readonly');
        const objectStore = transaction.objectStore(this.recipesStore);

        const request = objectStore.get(id);

        request.onsuccess = () => {
          const result = request.result;
          if (result) {
            result.data = JSON.parse(result.data) as Recipe;
            resolve(result);
          } else {
            reject('Entity not found');
          }
        };

        request.onerror = (event: any) => reject(event);
      };

      dbRequest.onerror = (event: any) => reject(event);
    });
  }

  public PutRecipe(recipe: Recipe): Promise<void> {
    return new Promise((resolve, reject) => {
      const dbRequest = indexedDB.open(this.dbname);

      dbRequest.onsuccess = (event: any) => {
        const db = event.target.result;
        const transaction = db.transaction([this.recipesStore], 'readwrite');
        const objectStore = transaction.objectStore(this.recipesStore);

        const entity = {
          id: recipe.id,
          etag: recipe.modified,
          data: JSON.stringify(recipe)
        };

        const request = objectStore.put(entity);

        request.onsuccess = () => resolve();
        request.onerror = (event: any) => reject(event);
      };

      dbRequest.onerror = (event: any) => reject(event);
    });
  }

  private open(): void {
    const request = indexedDB.open(this.dbname, this.schemaversion);

    request.onerror = (error) => {
      console.error('IndexedDB error: ', error);
    }

    request.onupgradeneeded = (event) => {
      console.log('IndexedDB upgrade needed: ', event);
      this.upgradeSchema(event);
    }
  }

  private upgradeSchema(event: IDBVersionChangeEvent): void {
    if (event.oldVersion === 0) {
      this.upgradeSchema_v1(event);
    }
  }

  private upgradeSchema_v1(event: IDBVersionChangeEvent): void {
    if (!event.target) {
      console.warn('Event target is null', event);
      return;
    }
    const db = (event.target as IDBOpenDBRequest).result;
    if (!db.objectStoreNames.contains(this.recipesStore)) {
      const objectStore = db.createObjectStore(this.recipesStore, { keyPath: 'id' });
      objectStore.createIndex('etag', 'etag', { unique: false });
    }
  }

}

export type IdbRecipe = {
  id: number;
  etag: string;
  data: Recipe;
}