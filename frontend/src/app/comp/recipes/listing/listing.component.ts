import { Component, Input } from '@angular/core';
import { Recipe } from '../../../types';
import { IconLib } from '../../../icons';
import { L10nService } from '../../../svc/l10n.service';
import { L10nLocale } from '../../../svc/locales/types';

@Component({
  selector: 'kb-recipes-listing',
  standalone: false,
  templateUrl: './listing.component.html',
  styleUrl: './listing.component.scss'
})
export class ListingComponent {

  Icons = IconLib;

  @Input({ required: true }) Recipes: Recipe[] = [];
  @Input({ required: true }) LangCode: string = 'de';

  constructor(
    private l10nService: L10nService,
  ) { }

  get Locale(): L10nLocale {
    return this.l10nService.Locale;
  }

  public replace(content: string, replacements: any[]): string {
    return this.l10nService.replace(content, replacements);
  }

}
