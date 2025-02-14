import { Component } from '@angular/core';
import { L10nService } from '../../svc/l10n.service';
import { L10nLocale } from '../../svc/locales/types';

@Component({
  selector: 'kb-navbar',
  standalone: false,
  templateUrl: './navbar.component.html',
  styleUrl: './navbar.component.scss'
})
export class NavbarComponent {

  constructor(
    private l10nService: L10nService,
  ) { }

  get Locale(): L10nLocale {
    return this.l10nService.Locale;
  }

}
