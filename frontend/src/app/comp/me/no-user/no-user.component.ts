import { Component, Input } from '@angular/core';
import { L10nService } from '../../../svc/l10n.service';
import { L10nLocale } from '../../../svc/locales/types';

@Component({
  selector: 'kb-me-no-user',
  standalone: false,
  templateUrl: './no-user.component.html',
  styleUrl: './no-user.component.scss'
})
export class MeNoUserComponent {

  @Input({ required: true }) PageRef: string = '';
  @Input({ required: true }) PageSrc: string = '';

  constructor(
    private l10nService: L10nService,
  ) { }

  get Locale(): L10nLocale {
    return this.l10nService.Locale;
  }

}
