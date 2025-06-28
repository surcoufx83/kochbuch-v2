import { Component, OnDestroy, OnInit, signal } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { Subscription } from 'rxjs';
import { SharedDataService } from '../../svc/shared-data.service';
import { L10nService } from '../../svc/l10n.service';
import { Recipe, Unit, UserSelf } from '../../types';
import { L10nLocale } from '../../svc/locales/types';
import { WebSocketService } from '../../svc/web-socket.service';
import { IconLib } from '../../icons';
import { CalculatorIngredient, CalculatorPreparationStep, IngredientsCalculator } from '../../svc/ingredients-calculator';

@Component({
  selector: 'kb-recipe',
  standalone: false,
  templateUrl: './recipe.component.html',
  styleUrl: './recipe.component.scss'
})
export class RecipeComponent implements OnDestroy, OnInit {

  calculatorServings = 0;
  icons = IconLib;
  ingredients = signal<CalculatorIngredient[]>([]);
  ingredientsCalc?: IngredientsCalculator;
  langCode = signal<string>('de');
  langCodeVisible = signal<string>('de');
  loadingFailed = signal<boolean>(false);
  localized = signal<boolean>(true);
  recipe?: Recipe;
  steps = signal<CalculatorPreparationStep[]>([]);
  user = signal<UserSelf | null>(null);

  private routeRecipeId?: number;
  private subs: Subscription[] = [];

  constructor(
    private l10nService: L10nService,
    private sharedDataService: SharedDataService,
    private route: ActivatedRoute,
    private router: Router,
    private wsService: WebSocketService,
  ) {
    this.user.set(this.wsService.GetUser());
  }

  FormatDuration(inMinutes: number): string {
    return this.l10nService.FormatDuration(inMinutes);
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

    this.subs.push(this.sharedDataService.RecipeEvents.subscribe((event) => {
      if (event === false || event.id !== this.routeRecipeId || !this.recipe || event.id !== this.recipe.id || event.etag !== this.recipe.modified)
        return;

      this.loadRecipeById(this.routeRecipeId);
    }));

    this.subs.push(this.l10nService.userLocale.subscribe((l) => {
      this.langCode.set(l);
      this.onToggleLocalization(this.localized());
    }));

    this.subs.push(this.wsService.User.subscribe((u) => {
      this.user.set(u);
    }))

  }

  loadRecipeById(id: number) {
    this.sharedDataService.getRecipe(id)
      .then((data: { id: number, etag: string, data: Recipe }) => {
        if (!this.recipe || this.recipe.id !== data.data.id || this.recipe.modified !== data.data.modified) {
          this.recipe = data.data;
          this.calculatorServings = this.recipe.servingsCount;
          this.ingredientsCalc = new IngredientsCalculator(this.l10nService, this.recipe, this.sharedDataService);
          this.onSetServingsCount(this.calculatorServings);
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

  onSetServingsCount(value: number): void {
    if (value < 1)
      value = 1;
    if (value > 100)
      value = 100;
    this.calculatorServings = value;
    this.steps.set(this.ingredientsCalc?.setServings(value) ?? []);
    this.ingredients.set(this.calcIngredients());
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

  onToggleLocalization(to: boolean): void {
    this.localized.set(to);
    this.langCodeVisible.set(to === true ? this.langCode() : (this.recipe?.userLocale ?? this.langCode()));
    this.onSetServingsCount(this.calculatorServings);
  }

}
