import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { RecipesComponent } from './comp/recipes/recipes.component';
import { RecipeComponent } from './comp/recipe/recipe.component';

const routes: Routes = [
  { path: 'new', component: RecipesComponent, pathMatch: 'full' },
  { path: 'all', component: RecipesComponent, pathMatch: 'full' },
  { path: 'cat/:id/:name', component: RecipesComponent },
  { path: 'r', component: RecipesComponent, pathMatch: 'full' },
  { path: 'r/:id', component: RecipeComponent, pathMatch: 'full' },
  { path: '**', component: RecipesComponent }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
