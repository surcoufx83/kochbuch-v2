import { Component } from '@angular/core';
import { Title } from '@angular/platform-browser';
import { SharedDataService } from './svc/shared-data.service';
import { ApiService } from './svc/api.service';
import { first } from 'rxjs';
import { L10nService } from './svc/l10n.service';
import { L10nLocale } from './svc/locales/types';

@Component({
  selector: 'kb-root',
  templateUrl: './app.component.html',
  standalone: false,
  styleUrl: './app.component.scss'
})
export class AppComponent {

  constructor(
    private apiService: ApiService,
    private l10nService: L10nService,
    sharedDataService: SharedDataService,
    htmlTitleService: Title
  ) {
    let catetag: string | undefined = undefined;
    sharedDataService.PageTitle.subscribe((t) => htmlTitleService.setTitle(t));
    apiService.get('').pipe(first()).subscribe((r) => {
      console.log(r)
    })
    apiService.get('categories', catetag).pipe(first()).subscribe((r) => {
      console.log(r)
    })
    apiService.get('units', catetag).pipe(first()).subscribe((r) => {
      console.log(r)
    })
    apiService.get('recipes').pipe(first()).subscribe((r) => {
      console.log(r)
    })

  }

  Locale(): L10nLocale {
    return this.l10nService.Locale;
  }

}
