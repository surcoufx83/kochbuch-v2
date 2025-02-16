import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { RecipesComponent } from './comp/recipes/recipes.component';
import { RecipeComponent } from './comp/recipe/recipe.component';
import { SearchComponent } from './comp/search/search.component';
import { MeComponent } from './comp/me/me.component';
import { EditorComponent } from './comp/editor/editor.component';

const routes: Routes = [
  { path: 'search', component: SearchComponent },
  { path: 'me', component: MeComponent },
  { path: 'new', component: RecipesComponent },
  { path: 'all', component: RecipesComponent },
  { path: 'create', component: EditorComponent },
  { path: 'cat/:id/:name', component: RecipesComponent },
  { path: 'r', component: RecipesComponent, pathMatch: 'full' },
  { path: 'r/:id', component: RecipeComponent, pathMatch: 'full' },
  { path: 'r/:id/edit', component: EditorComponent, pathMatch: 'full' },
  { path: '**', component: RecipesComponent }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
