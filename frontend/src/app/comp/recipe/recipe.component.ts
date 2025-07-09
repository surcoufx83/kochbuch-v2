import { Component, OnDestroy, OnInit, signal } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { Subscription } from 'rxjs';
import { SharedDataService } from '../../svc/shared-data.service';
import { L10nService } from '../../svc/l10n.service';
import { Collection, Recipe, Unit, UserSelf } from '../../types';
import { L10nLocale } from '../../svc/locales/types';
import { WebSocketService } from '../../svc/web-socket.service';
import { IconLib } from '../../icons';
import { CalculatorIngredient, CalculatorPreparationStep, IngredientsCalculator } from '../../svc/ingredients-calculator';
import { SwalComponent, SwalPortalTargets } from '@sweetalert2/ngx-sweetalert2';
import { FormControl, FormGroup, Validators } from '@angular/forms';

@Component({
  selector: 'kb-recipe',
  standalone: false,
  templateUrl: './recipe.component.html',
  styleUrl: './recipe.component.scss'
})
export class RecipeComponent implements OnDestroy, OnInit {

  calculatorServings = 0;
  collectionSaving = signal<boolean>(false);
  fileProgress = signal<'none' | 'checking' | 'uploading'>('none');
  icons = IconLib;
  ingredients = signal<CalculatorIngredient[]>([]);
  ingredientsCalc?: IngredientsCalculator;
  langCode = signal<string>('de');
  langCodeVisible = signal<string>('de');
  loadingFailed = signal<boolean>(false);
  localized = signal<boolean>(true);
  recipe?: Recipe;
  recipeFoundInCollections: number[] = [];
  recipeFoundInCollectionsClone: number[] = [];
  recipeSavedToCollection: number[] = [];
  shareLinkServings?: number;
  showCollectionEditor = signal<boolean>(false);
  showCollectionPicker = signal<boolean>(false);
  steps = signal<CalculatorPreparationStep[]>([]);
  uri: string = window.location.href;
  user: UserSelf | null = null;
  userCollections: Collection[] = [];

  collectionEditorForm = new FormGroup({
    title: new FormControl<string>('', {
      validators: [
        Validators.required,
        Validators.minLength(1),
        Validators.maxLength(256),
      ]
    })
  });

  private routeRecipeId?: number;
  private subs: Subscription[] = [];

  constructor(
    private l10nService: L10nService,
    private sharedDataService: SharedDataService,
    private route: ActivatedRoute,
    private router: Router,
    private wsService: WebSocketService,
    public readonly swalTargets: SwalPortalTargets,
  ) {
    this.ngOnInitUser(this.wsService.GetUser());
  }

  calcIngredients(): CalculatorIngredient[] {
    let ingredients: CalculatorIngredient[] = [];
    const steps = this.steps();
    const langcode = this.langCode();

    let ingredientDict: { name: string, ingredient: CalculatorIngredient, values: { [key: number]: { count: number, quantity: number, unit: Unit | null } } }[] = []; // values = unit id -> display quantity
    let ingredientDictMap: { [key: string]: number } = {};

    steps.forEach((step) => {
      step.ingredients!.forEach((ing) => {

        const title = ing.localization[langcode].title ? ing.localization[langcode].title : ing.localization['de'].title

        if (!ingredientDictMap[title]) {
          ingredientDictMap[title] = ingredientDict.length;
          ingredientDict.push({
            name: title,
            ingredient: ing,
            values: {},
          });
        }

        const i = ingredientDictMap[title];
        const unitid = ing.displayAsUnitId ?? ing.unitId ?? -1;
        const quantity = ing.displayQuantity ?? ing.calcQuantity ?? ing.baseQuantity ?? 0;

        if (!ingredientDict[i].values[unitid]) {
          ingredientDict[i].values[unitid] = {
            count: 1,
            quantity: quantity,
            unit: ing.displayAsUnit ?? ing.unit,
          };
        }
        else {
          ingredientDict[i].values[unitid].quantity += quantity;
          ingredientDict[i].values[unitid].count += 1;
        }

      });
    });

    ingredientDict.forEach((ing) => {
      Object.values(ing.values).forEach((val) => {
        ingredients.push({
          baseQuantity: 0,
          id: 0,
          index: 0,
          localization: ing.ingredient.localization,
          calcQuantity: null,
          unitId: null,
          unit: null,
          displayAsUnit: val.unit,
          displayAsUnitId: null,
          displayQuantity: val.quantity,
          displayStr: this.l10nService.FormatIngredient(val.quantity, val.unit),
        });
      });
    })

    return ingredients;
  }

  public FormatDate(date: string | number | Date, formatStr: string): string {
    return this.l10nService.FormatDate(date, formatStr);
  }

  FormatDuration(inMinutes: number, longFormat: boolean = false): string {
    return this.l10nService.FormatDuration(inMinutes, longFormat);
  }

  get Locale(): L10nLocale {
    return this.l10nService.Locale;
  }

  public LocaleReplace(content: string, replacements: any[]): string {
    return this.l10nService.Replace(content, replacements);
  }

  ngOnDestroy(): void {
    for (const s of this.subs) {
      s.unsubscribe();
    }
    this.subs = [];
  }

  ngOnInit(): void {

    this.subs.push(this.route.params.subscribe((params) => {
      this.routeRecipeId = params['id'] ? +params['id'] : undefined;
      if (this.routeRecipeId === undefined) {
        this.router.navigate(['/']);
        return;
      }
      this.loadRecipeById(this.routeRecipeId);
    }));

    this.subs.push(this.route.queryParams.subscribe((params) => {
      this.shareLinkServings = params['s'] ? +params['s'] : undefined;
    }));

    this.subs.push(this.sharedDataService.RecipeEvents.subscribe((event) => {
      if (event === false || event.id !== this.routeRecipeId || !this.recipe || event.id !== this.recipe.id || event.etag === this.recipe.modified)
        return;

      this.loadRecipeById(this.routeRecipeId);
    }));

    this.subs.push(this.l10nService.userLocale.subscribe((l) => {
      this.langCode.set(l);
      this.onToggleLocalization(this.localized());
    }));

    this.subs.push(this.wsService.User.subscribe((u) => {
      this.ngOnInitUser(u);
    }))

  }

  ngOnInitUser(user: UserSelf | null): void {
    const colls = user ? Object.values(user.collections).filter(x => x.deleted === null).sort((a, b) => a.title.toLocaleLowerCase().localeCompare(b.title.toLocaleLowerCase())) : [];
    colls.forEach((coll) => {
      if (!coll.items)
        coll.items = [];
    });
    this.userCollections = colls;
    this.user = user;

    if (this.recipe && this.userCollections && !this.showCollectionPicker()) {
      this.recipeFoundInCollections = [];
      this.userCollections
        .filter(c => c.items.filter(i => i.recipeId === this.recipe!.id).length > 0)
        .forEach((c) => this.recipeFoundInCollections.push(c.id));
    }
  }

  loadRecipeById(id: number) {
    this.sharedDataService.getRecipe(id)
      .then((data: { id: number, etag: string, data: Recipe }) => {
        if (!this.recipe || this.recipe.id !== data.data.id || this.recipe.modified !== data.data.modified) {
          this.recipe = data.data;
          this.calculatorServings = this.shareLinkServings ?? this.recipe.servingsCount;
          this.ingredientsCalc = new IngredientsCalculator(this.l10nService, this.recipe, this.sharedDataService);
          this.onSetServingsCount(this.calculatorServings);
          this.loadRecipeCategories(this.recipe);
          this.ngOnInitUser(this.user);
        }
      })
      .catch((err) => {
        console.error(err)
        this.loadingFailed.set(true);
        this.wsService.ReportError({
          url: this.router.url,
          error: `Recipe with id ${id} not found.`,
          severity: 'E'
        });
        setTimeout(() => {
          this.router.navigate(['/']);
        }, 1000);
      });

  }

  loadRecipeCategories(recipe: Recipe): void {
    recipe.categories.forEach((rc) => {
      console.log(rc);
    })
  }

  onCreateCollectionClick(inputEl: HTMLInputElement): void {
    if (!this.recipe || this.showCollectionEditor())
      return;

    this.recipeFoundInCollectionsClone.push(0);
    this.collectionEditorForm.reset({
      title: ''
    });
    this.showCollectionEditor.set(true);

    setTimeout(() => {
      inputEl.focus();
      inputEl.select();
    }, 10);
  }

  onCreateCollectionCancel(): void {
    this.showCollectionEditor.set(false);
    if (this.recipeFoundInCollectionsClone.includes(0)) {
      this.recipeFoundInCollectionsClone.splice(this.recipeFoundInCollectionsClone.indexOf(0), 1);
    }
  }

  onCreateCollectionSubmit(): void {
    if (!this.recipe || !this.showCollectionEditor() || this.collectionSaving())
      return;

    this.collectionSaving.set(true);

    this.sharedDataService.createCollection(
      this.collectionEditorForm.controls.title.value!,
      '',
    ).then((coll) => {
      this.onCreateCollectionCancel();
      this.onToggleRecipeInCollection(coll);
    });


  }

  onPictureUploadChange(event: Event): void {
    const element = event.currentTarget as HTMLInputElement;
    let fileList: FileList | null = element.files;
    if (fileList) {
      this.fileProgress.set('checking');

      const fileReadPromises: Promise<RecipePictureUploadMsgContent | false>[] = [];

      for (let i = 0; i < fileList.length; i++) {
        const file = fileList.item(i);
        if (!file) continue;

        const fileReadPromise = new Promise<RecipePictureUploadMsgContent | false>((resolve) => {
          const reader = new FileReader();

          reader.onload = () => {
            const result = reader.result as string;
            const base64 = result.split(',')[1];

            const fileData: RecipePictureUploadMsgContent = {
              name: file.name,
              type: file.type,
              size: file.size,
              base64: base64,
            };

            resolve(fileData);
          };

          reader.onerror = () => {
            console.error("Error reading file:", file.name);
            resolve(false);
          };

          reader.onabort = () => {
            console.warn("File read aborted:", file.name);
            resolve(false);
          };

          reader.readAsDataURL(file);
        });

        fileReadPromises.push(fileReadPromise);
      }

      // Wait for all file reads to complete
      Promise.all(fileReadPromises).then((results) => {
        const validFiles = results.filter((r): r is RecipePictureUploadMsgContent => r !== false);

        const wspayload: RecipePictureUploadMsg = {
          type: 'recipe_picture_upload',
          content: JSON.stringify({
            recipe: this.recipe!.id,
            files: validFiles,
          }),
        };

        this.fileProgress.set('uploading');

        this.wsService.SendMessageAndWait(wspayload).then((result) => {
          // do nothing.
        }).catch((err) => {
          console.log(err)
        }).finally(() => {
          this.fileProgress.set('none');
        })

      });

    }
  }

  onPrintClick(): void {
    window.print();
  }

  onSaveToCollectionClick(pickCollectionDialog: SwalComponent, confirmDialog: SwalComponent): void {
    if (!this.user || !this.recipe)
      return;

    this.recipeSavedToCollection = [];
    this.showCollectionPicker.set(true);

    this.onSaveToCollectionClick__selectCollectionToSave(pickCollectionDialog).then((colls) => {
      if (!this.recipe)
        return;

      this.sharedDataService.pushRecipeToCollections(colls, this.recipe)
        .then((state) => {
          if (state === true) {
            this.recipeSavedToCollection = colls;
            confirmDialog.fire();
          }
        });

    }).catch(() => {
      // we do nothing here (user cancelled dialog)
    }).finally(() => {
      this.showCollectionPicker.set(false);
    });


  }

  onSaveToCollectionClick__selectCollectionToSave(pickCollectionDialog: SwalComponent): Promise<number[]> {
    return new Promise<number[]>((resolve, reject) => {

      if (!this.user) {
        reject();
        return;
      }

      if (!this.userCollections) {
        this.ngOnInitUser(this.user);
        if (!this.userCollections) {
          this.userCollections = [];
        }
      }

      this.recipeFoundInCollectionsClone = [...this.recipeFoundInCollections];

      if (this.userCollections.length === 0) {
        // user has no collection yet, create the default onw
        this.sharedDataService.createCollection(
          this.Locale.collections.defaultCollection.title,
          this.Locale.collections.defaultCollection.description,
        ).then((coll) => {
          // server has created a new collection so we're ready to show the dialog
          this.onSaveToCollectionClick__selectCollectionToSave__fire(pickCollectionDialog, resolve, reject)
        }).catch(() => {
          reject();
        });
        return;

      }
      else {
        this.onSaveToCollectionClick__selectCollectionToSave__fire(pickCollectionDialog, resolve, reject);
      }

    });
  }

  onSaveToCollectionClick__selectCollectionToSave__fire(pickCollectionDialog: SwalComponent, resolve: (value: number[] | PromiseLike<number[]>) => void, reject: (reason?: any) => void): void {
    // shows the collection select dialog
    pickCollectionDialog.fire().then((result) => {
      if (result.isConfirmed) {
        // dialog is closed with ok button
        resolve(this.recipeFoundInCollectionsClone);
      }
      else {
        // dialog is closed with cancel or backdrop -> do nothing
        reject();
      }
    }).catch((reason) => {
      // dialog as error
      reject();
    }).finally(() => {
      // cleanup dialog modifications
      this.onCreateCollectionCancel();
    });
  }

  onShareClick(): void {
    if (!this.recipe || !this.recipe.sharedPublic)
      return;
    // https://stackoverflow.com/questions/62350936/how-to-add-a-preview-image-to-navigator-share
    // this suggests to remove text and title
    // should be checked for live page
    navigator.share({
      url: `https://${window.location.host}${window.location.pathname}?s=${this.calculatorServings}&n=${this.urlencode(this.recipe.localization[this.recipe.localized ? this.langCodeVisible() : this.recipe.userLocale].title)}`,
      title: this.LocaleReplace(this.Locale.recipe.share.title, [this.recipe.localization[this.recipe.localized ? this.langCodeVisible() : this.recipe.userLocale].title]),
      text: this.LocaleReplace(this.Locale.recipe.share.message, [this.recipe.localization[this.recipe.localized ? this.langCodeVisible() : this.recipe.userLocale].title, this.recipe.user?.displayname]),
    });
  }

  onSetServingsCount(value: number): void {
    if (value < 1)
      value = 1;
    if (value > 100)
      value = 100;
    this.calculatorServings = value;
    this.steps.set(this.ingredientsCalc?.setServings(value) ?? []);
    this.ingredients.set(this.calcIngredients());
  }

  onToggleLocalization(to: boolean): void {
    this.localized.set(to);
    this.langCodeVisible.set(to === true ? this.langCode() : (this.recipe?.userLocale ?? this.langCode()));
    this.onSetServingsCount(this.calculatorServings);
  }

  onToggleRecipeInCollection(coll: Collection): void {
    if (this.recipeFoundInCollectionsClone.includes(coll.id)) {
      this.recipeFoundInCollectionsClone.splice(this.recipeFoundInCollectionsClone.indexOf(coll.id), 1);
    }
    else {
      this.recipeFoundInCollectionsClone.push(coll.id);
    }
  }

  public urlencode(content: string): string {
    return this.l10nService.UrlEncode(content);
  }

}

export type RecipePictureUploadMsg = {
  type: 'recipe_picture_upload',
  content: string,
}

export type RecipePictureUploadMsgContent = {
  name: string,
  type: string,
  size: number,
  base64: string,
}