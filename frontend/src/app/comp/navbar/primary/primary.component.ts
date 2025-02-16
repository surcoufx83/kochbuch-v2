import { Component } from '@angular/core';
import { IconLib } from '../../../icons';
import { L10nService } from '../../../svc/l10n.service';
import { L10nLocale } from '../../../svc/locales/types';

@Component({
  selector: 'kb-navbar-primary',
  standalone: false,
  templateUrl: './primary.component.html',
  styleUrl: './primary.component.scss'
})
export class PrimaryComponent {

  Icons = IconLib;

  constructor(
    private l10nService: L10nService,
  ) { }

  get Locale(): L10nLocale {
    return this.l10nService.Locale;
  }

}
