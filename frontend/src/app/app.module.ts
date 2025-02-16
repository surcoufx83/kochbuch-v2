import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { provideHttpClient } from '@angular/common/http';
import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { ApiService } from './svc/api.service';
import { SharedDataService } from './svc/shared-data.service';
import { NavbarComponent } from './comp/navbar/navbar.component';
import { L10nService } from './svc/l10n.service';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { SearchComponent } from './comp/navbar/search/search.component';
import { PrimaryComponent } from './comp/navbar/primary/primary.component';
import { SecondaryComponent } from './comp/navbar/secondary/secondary.component';
import { RecipesComponent } from './comp/recipes/recipes.component';
import { RecipeComponent } from './comp/recipe/recipe.component';
import { MobileMenuComponent } from './comp/mobile-menu/mobile-menu.component';
import { MeComponent } from './comp/me/me.component';
import { EditorComponent } from './comp/editor/editor.component';

@NgModule({
  declarations: [
    AppComponent,
    NavbarComponent,
    SearchComponent,
    PrimaryComponent,
    SecondaryComponent,
    RecipesComponent,
    RecipeComponent,
    MobileMenuComponent,
    MeComponent,
    EditorComponent,
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    FontAwesomeModule
  ],
  providers: [
    provideHttpClient(),
    ApiService,
    L10nService,
    SharedDataService,
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
