import { Recipe, RecipeIngredientLocalization, RecipePreparation, RecipePreparationLocalization, RecipeTiming, Unit } from "../types";
import { SharedDataService, UnitsConversion } from "./shared-data.service";

export class IngredientsCalculator {

    private baseServings: number;
    private calcServings: number;
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
        this.baseServings = recipe.servingsCount;
        this.calcServings = recipe.servingsCount;
        this.steps = [];
        this.unitWithOutgoingConversion = [];
        this.unitWithIncomingConversion = [];

        this.unitConversion.forEach((uc, i) => {
            this.unitWithOutgoingConversion[uc.fromId] = i;
            this.unitWithIncomingConversion[uc.toId] = i;
        });

        console.log(recipe)
        recipe.preparation.forEach((step) => {
            const copystep: CalculatorPreparationStep = {
                id: step.id,
                index: step.index,
                localization: step.localization,
                timing: step.timing,
                ingredients: [],
            };
            if (step.ingredients === null)
                step.ingredients = [];

            step.ingredients.forEach((ing) => {
                const copying: CalculatorIngredient = {
                    id: ing.id,
                    index: ing.index,
                    localization: ing.localization,
                    baseQuantity: ing.quantity,
                    calcQuantity: ing.quantity,
                    unitId: ing.unitId,
                    unit: this.units[ing.unitId ?? 0] ?? null,
                    displayAsUnit: null,
                    displayAsUnitId: null,
                    displayQuantity: null,
                    displayStr: ''
                }
                this.convertToBestUnit(copying);
                copystep.ingredients!.push(copying);
            });
            this.steps.push(copystep);
        });
        console.log(recipe)
    }

    private convertToBestUnit(ing: CalculatorIngredient) {
        if (!ing.unit || (this.unitWithOutgoingConversion[ing.unit.id] === undefined && this.unitWithIncomingConversion[ing.unit.id] === undefined) || ing.calcQuantity === null) {
            ing.displayAsUnit = ing.unit;
            ing.displayAsUnitId = ing.unitId;
            ing.displayQuantity = ing.calcQuantity ?? ing.baseQuantity;
            return;
        }

        const convoutI = this.unitWithOutgoingConversion[ing.unit.id];
        const convout = this.unitConversion[convoutI || -1];

        if (!convout || ing.calcQuantity <= convout.fromQuantity) {
            ing.displayAsUnit = ing.unit;
            ing.displayAsUnitId = ing.unitId;
            ing.displayQuantity = ing.calcQuantity ?? ing.baseQuantity;
            return;
        }

        ing.displayQuantity = ing.calcQuantity / convout.fromQuantity * convout.toQuantity;
        ing.displayAsUnitId = convout.toId;
        ing.displayAsUnit = convout.to;

    }

    setServings(newvalue: number): CalculatorPreparationStep[] {
        const outvalue = [...this.steps];
        this.calcServings = newvalue;
        outvalue.forEach((step) => {
            step.ingredients!.forEach((ing) => {
                if (ing.baseQuantity === null)
                    return;
                ing.calcQuantity = ing.baseQuantity / this.baseServings * this.calcServings;
                this.convertToBestUnit(ing);
            })
        });
        return outvalue;
    }

}

export type CalculatorPreparationStep = {
    id: number,
    index: number,
    localization: RecipePreparationLocalization,
    timing: RecipeTiming,
    ingredients: CalculatorIngredient[] | null,
}

export type CalculatorIngredient = {
    id: number,
    index: number,
    localization: RecipeIngredientLocalization,
    baseQuantity: number | null,
    calcQuantity: number | null,
    unitId: number | null,
    unit: Unit | null,
    displayAsUnit: Unit | null,
    displayAsUnitId: number | null,
    displayQuantity: number | null,
    displayStr: string,
}
