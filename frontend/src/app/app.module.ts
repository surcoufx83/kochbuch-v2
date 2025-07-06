import { provideHttpClient } from '@angular/common/http';
import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { BrowserModule } from '@angular/platform-browser';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { SweetAlert2Module } from '@sweetalert2/ngx-sweetalert2';
import { PhotoGalleryModule } from '@twogate/ngx-photo-gallery';
import { QRCodeComponent } from 'angularx-qrcode';
import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { EditorComponent } from './comp/editor/editor.component';
import { MeComponent } from './comp/me/me.component';
import { MeNoUserComponent } from './comp/me/no-user/no-user.component';
import { MissingLinkComponent } from './comp/missing-link/missing-link.component';
import { MobileMenuComponent } from './comp/mobile-menu/mobile-menu.component';
import { NavbarComponent } from './comp/navbar/navbar.component';
import { PrimaryComponent } from './comp/navbar/primary/primary.component';
import { SearchComponent } from './comp/navbar/search/search.component';
import { SecondaryComponent } from './comp/navbar/secondary/secondary.component';
import { Oauth2LoginCallbackComponent } from './comp/oauth2-login-callback/oauth2-login-callback.component';
import { RecipeComponent } from './comp/recipe/recipe.component';
import { ListingComponent } from './comp/recipes/listing/listing.component';
import { RecipesComponent } from './comp/recipes/recipes.component';
import { WelcomeHeaderComponent } from './comp/recipes/welcome-header/welcome-header.component';
import { L10nService } from './svc/l10n.service';
import { SharedDataService } from './svc/shared-data.service';
import { WebSocketService } from './svc/web-socket.service';

@NgModule({
  declarations: [
    AppComponent,
    EditorComponent,
    ListingComponent,
    MeComponent,
    MeNoUserComponent,
    MissingLinkComponent,
    MobileMenuComponent,
    NavbarComponent,
    Oauth2LoginCallbackComponent,
    PrimaryComponent,
    RecipeComponent,
    RecipesComponent,
    SearchComponent,
    SecondaryComponent,
    WelcomeHeaderComponent,
  ],
  imports: [
    AppRoutingModule,
    BrowserModule,
    FontAwesomeModule,
    FormsModule,
    PhotoGalleryModule,
    QRCodeComponent,
    ReactiveFormsModule,
    SweetAlert2Module.forRoot(),
  ],
  providers: [
    provideHttpClient(),
    WebSocketService,
    L10nService,
    SharedDataService,
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
