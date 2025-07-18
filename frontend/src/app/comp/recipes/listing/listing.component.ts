import { Component, Input } from '@angular/core';
import { Recipe } from '../../../types';
import { IconLib } from '../../../icons';
import { L10nService } from '../../../svc/l10n.service';
import { L10nLocale } from '../../../svc/locales/types';
import { SharedDataService } from '../../../svc/shared-data.service';

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
    private sharedDataService: SharedDataService,
  ) { }

  public formatDate(date: string | number | Date, formatStr: string): string {
    return this.l10nService.FormatDate(date, formatStr);
  }

  public formatDuration(minutes: number): string {
    return this.l10nService.FormatDuration(minutes);
  }

  public formatVote(n: number): string {
    return this.l10nService.FormatVote(n);
  }

  get Locale(): L10nLocale {
    return this.l10nService.Locale;
  }

  public hover(recipeId: number): void {
    this.sharedDataService.PreloadRecipeToCache(recipeId);
  }

  public replace(content: string, replacements: any[]): string {
    return this.l10nService.Replace(content, replacements);
  }

  public urlencode(content: string): string {
    return this.l10nService.UrlEncode(content);
  }

}
