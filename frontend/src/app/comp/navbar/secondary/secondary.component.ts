import { Component } from '@angular/core';
import { IconLib } from '../../../icons';
import { L10nService } from '../../../svc/l10n.service';
import { L10nLocale } from '../../../svc/locales/types';

@Component({
  selector: 'kb-navbar-secondary',
  standalone: false,
  templateUrl: './secondary.component.html',
  styleUrl: './secondary.component.scss'
})
export class SecondaryComponent {

  Icons = IconLib;

  constructor(
    private l10nService: L10nService,
  ) { }

  get Locale(): L10nLocale {
    return this.l10nService.Locale;
  }

}
