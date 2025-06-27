import { Recipe, RecipeIngredientLocalization, RecipePreparation, RecipePreparationLocalization, RecipeTiming, Unit } from "../types";
import { SharedDataService, UnitsConversion } from "./shared-data.service";

export class IngredientsCalculator {

    private servings: number;
    private steps: CalculatorPreparationStep[];

    private units: { [key: number]: Unit };
    private unitConversion: UnitsConversion[];
    private unitWithOutgoingConversion: { [key: number]: number }; // conversion from unit id -> index of array
    private unitWithIncomingConversion: { [key: number]: number }; // conversion to unit id -> index of array

    constructor(
        private recipe: Recipe,
        sharedDataService: SharedDataService,
    ) {
        this.units = sharedDataService.getUnits();
        this.unitConversion = sharedDataService.UnitConversions;
        this.servings = recipe.servingsCount;
        this.steps = [];
        this.unitWithOutgoingConversion = [];
        this.unitWithIncomingConversion = [];

        this.unitConversion.forEach((uc, i) => {
            this.unitWithOutgoingConversion[uc.fromId] = i;
            this.unitWithIncomingConversion[uc.toId] = i;
        });

        recipe.preparation.forEach((step) => {
            const copystep: CalculatorPreparationStep = {
                id: step.id,
                index: step.index,
                localization: step.localization,
                timing: step.timing,
                ingredients: [],
            };
            step.ingredients.forEach((ing) => {
                const copying: CalculatorIngredient = {
                    id: ing.id,
                    index: ing.index,
                    localization: ing.localization,
                    quantity: ing.quantity,
                    unitId: ing.unitId,
                    unit: this.units[ing.unitId ?? 0] ?? null,
                    displayAsUnit: null,
                    displayAsUnitId: null,
                    displayQuantity: null
                }
                this.convertToBestUnit(copying);
                copystep.ingredients.push(copying);
            });
            this.steps.push(copystep);
        });
        console.log(recipe)
    }

    private convertToBestUnit(ing: CalculatorIngredient) {
        if (!ing.unit || (this.unitWithOutgoingConversion[ing.unit.id] === undefined && this.unitWithIncomingConversion[ing.unit.id] === undefined))
            return;
        console.log(`convertToBestUnit()`, ing)
        console.log(`Unit ${ing.unit.id} (${ing.unit.localization['en'].singular}) => hasOutgoingConversion = ${this.unitWithOutgoingConversion[ing.unit.id] !== undefined} / hasIncomingConversion = ${this.unitWithIncomingConversion[ing.unit.id] !== undefined}`);

        const convout = this.unitWithOutgoingConversion[ing.unit.id];
        const convin = this.unitWithIncomingConversion[ing.unit.id];


    }

}

export type CalculatorPreparationStep = {
    id: number,
    index: number,
    localization: RecipePreparationLocalization,
    timing: RecipeTiming,
    ingredients: CalculatorIngredient[],
}

export type CalculatorIngredient = {
    id: number,
    index: number,
    localization: RecipeIngredientLocalization,
    quantity: number | null,
    unitId: number | null,
    unit: Unit | null,
    displayAsUnit: Unit | null,
    displayAsUnitId: number | null,
    displayQuantity: number | null,
}
