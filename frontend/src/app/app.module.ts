import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { provideHttpClient } from '@angular/common/http';
import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { ApiService } from './svc/api.service';
import { SharedDataService } from './svc/shared-data.service';
import { NavbarComponent } from './comp/navbar/navbar.component';

@NgModule({
  declarations: [
    AppComponent,
    NavbarComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule
  ],
  providers: [
    provideHttpClient(),
    ApiService,
    SharedDataService,
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
