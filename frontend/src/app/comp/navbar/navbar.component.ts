import { Component, signal } from '@angular/core';
import { L10nService } from '../../svc/l10n.service';
import { L10nLocale } from '../../svc/locales/types';
import { IconLib } from '../../icons';
import { ApiService } from '../../svc/api.service';

@Component({
  selector: 'kb-navbar',
  standalone: false,
  templateUrl: './navbar.component.html',
  styleUrl: './navbar.component.scss'
})
export class NavbarComponent {

  Icons = IconLib;
  LoggedIn = signal<boolean>(false);
  SecondaryNavbar = signal<boolean>(false);

  constructor(
    apiService: ApiService,
    private l10nService: L10nService,
  ) {
    apiService.isLoggedIn.subscribe((state) => this.LoggedIn.set(state));
  }

  get Locale(): L10nLocale {
    return this.l10nService.Locale;
  }

}
