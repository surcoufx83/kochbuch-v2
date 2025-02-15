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
    let unitsetag: string | undefined = undefined;
    let recipesetag: string | undefined = undefined;
    sharedDataService.PageTitle.subscribe((t) => htmlTitleService.setTitle(t));

    setInterval(() => {
      apiService.get('categories', catetag).pipe(first()).subscribe((r) => {
        if (r?.status === 304) {
          console.log('304');
          return;
        }

        catetag = r?.headers.get('etag') ?? undefined;
        console.log(r)
      })
    }, 5000);

    setInterval(() => {
      apiService.get('units', unitsetag).pipe(first()).subscribe((r) => {
        if (r?.status === 304) {
          console.log('304');
          return;
        }

        unitsetag = r?.headers.get('etag') ?? undefined;
        console.log(r)
      })
    }, 7500);

    setInterval(() => {
      apiService.get('recipes', recipesetag).pipe(first()).subscribe((r) => {
        if (r?.status === 304) {
          console.log('304');
          return;
        }

        recipesetag = r?.headers.get('etag') ?? undefined;
        console.log(r)
      })
    }, 10500);

  }

  Locale(): L10nLocale {
    return this.l10nService.Locale;
  }

}
