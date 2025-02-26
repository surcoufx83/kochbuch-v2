import { Component, Input, signal } from '@angular/core';
import { UserSelf } from '../../../types';
import { L10nLocale } from '../../../svc/locales/types';
import { IconLib } from '../../../icons';
import { L10nService } from '../../../svc/l10n.service';

@Component({
  selector: 'kb-recipes-welcome-header',
  standalone: false,
  templateUrl: './welcome-header.component.html',
  styleUrl: './welcome-header.component.scss'
})
export class WelcomeHeaderComponent {

  Icons = IconLib;

  @Input({ required: true }) LoggedIn?: boolean;
  @Input({ required: true }) User?: UserSelf | false;

  constructor(
    private l10nService: L10nService,
  ) { }

  get Locale(): L10nLocale {
    return this.l10nService.Locale;
  }

  public replace(content: string, replacements: any[]): string {
    return this.l10nService.Replace(content, replacements);
  }

}
