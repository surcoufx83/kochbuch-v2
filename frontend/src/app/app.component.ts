import { Component } from '@angular/core';
import { Title } from '@angular/platform-browser';
import { L10nService } from './svc/l10n.service';
import { L10nLocale } from './svc/locales/types';
import { SharedDataService } from './svc/shared-data.service';

@Component({
  selector: 'kb-root',
  templateUrl: './app.component.html',
  standalone: false,
  styleUrl: './app.component.scss'
})
export class AppComponent {

  constructor(
    private l10nService: L10nService,
    sharedDataService: SharedDataService,
    htmlTitleService: Title,
  ) {
    sharedDataService.PageTitle.subscribe((t) => htmlTitleService.setTitle(t));
  }

  Locale(): L10nLocale {
    return this.l10nService.Locale;
  }

}
