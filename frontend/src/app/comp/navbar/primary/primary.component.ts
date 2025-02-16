import { Component, EventEmitter, Input, Output, signal } from '@angular/core';
import { IconLib } from '../../../icons';
import { L10nService } from '../../../svc/l10n.service';
import { L10nLocale } from '../../../svc/locales/types';
import { ApiService } from '../../../svc/api.service';
import { UserSelf } from '../../../types';

@Component({
  selector: 'kb-navbar-primary',
  standalone: false,
  templateUrl: './primary.component.html',
  styleUrl: './primary.component.scss'
})
export class PrimaryComponent {

  @Input({ required: true }) isSecondaryNavbarVisible = signal<boolean>(false);
  @Output() onToggleSecondaryNavbar = new EventEmitter();

  Icons = IconLib;
  LoggedIn = signal<boolean>(false);
  User = signal<UserSelf | false>(false);

  constructor(
    apiService: ApiService,
    private l10nService: L10nService,
  ) {
    apiService.isLoggedIn.subscribe((state) => {
      this.LoggedIn.set(state);
      this.User.set(apiService.User ?? false);
    });
  }

  get Locale(): L10nLocale {
    return this.l10nService.Locale;
  }

}
