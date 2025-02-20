import { Component, Input } from '@angular/core';
import { L10nService } from '../../../svc/l10n.service';
import { L10nLocale } from '../../../svc/locales/types';
import { IconLib } from '../../../icons';
import { ApiService } from '../../../svc/api.service';

@Component({
  selector: 'kb-me-no-user',
  standalone: false,
  templateUrl: './no-user.component.html',
  styleUrl: './no-user.component.scss'
})
export class MeNoUserComponent {

  @Input({ required: true }) PageRef: string = '';
  @Input({ required: true }) PageSrc: string = '';

  Icons = IconLib;
  LoginUrl?: string;

  constructor(
    private apiService: ApiService,
    private l10nService: L10nService,
  ) {
    this.LoginUrl = this.apiService.LoginUrl;
  }

  get Locale(): L10nLocale {
    return this.l10nService.Locale;
  }

}
