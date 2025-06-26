import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { RecipesComponent } from './comp/recipes/recipes.component';
import { RecipeComponent } from './comp/recipe/recipe.component';
import { SearchComponent } from './comp/search/search.component';
import { MeComponent } from './comp/me/me.component';
import { EditorComponent } from './comp/editor/editor.component';
import { MissingLinkComponent } from './comp/missing-link/missing-link.component';
import { Oauth2LoginCallbackComponent } from './comp/oauth2-login-callback/oauth2-login-callback.component';

const routes: Routes = [
  { path: 'search', component: SearchComponent },
  { path: 'me', component: MeComponent },
  { path: 'new', component: RecipesComponent },
  { path: 'all', component: RecipesComponent },
  { path: 'create', component: EditorComponent },
  { path: 'cat/:id/:name', component: RecipesComponent },
  { path: 'login/oauth2', component: Oauth2LoginCallbackComponent, pathMatch: 'full' },
  { path: 'recipes', redirectTo: '', pathMatch: 'full' },
  { path: 'recipe/:id', component: RecipeComponent, pathMatch: 'full' },
  { path: 'recipe/:id/edit', component: EditorComponent, pathMatch: 'full' },
  { path: '', component: RecipesComponent, pathMatch: 'full' },
  { path: '**', component: MissingLinkComponent }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
