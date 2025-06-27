import { Component, OnDestroy, OnInit, signal } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { Subscription } from 'rxjs';
import { SharedDataService } from '../../svc/shared-data.service';
import { L10nService } from '../../svc/l10n.service';
import { Recipe } from '../../types';
import { L10nLocale } from '../../svc/locales/types';
import { WebSocketService } from '../../svc/web-socket.service';
import { IconLib } from '../../icons';

@Component({
  selector: 'kb-recipe',
  standalone: false,
  templateUrl: './recipe.component.html',
  styleUrl: './recipe.component.scss'
})
export class RecipeComponent implements OnDestroy, OnInit {

  localized = signal<boolean>(true);
  icons = IconLib;
  langCode = signal<string>('de');
  langCodeVisible = signal<string>('de');
  loadingFailed = signal<boolean>(false);
  recipe?: Recipe;
  private routeRecipeId?: number;
  private subs: Subscription[] = [];

  constructor(
    private l10nService: L10nService,
    private sharedDataService: SharedDataService,
    private route: ActivatedRoute,
    private router: Router,
    private wsService: WebSocketService,
  ) { }

  get Locale(): L10nLocale {
    return this.l10nService.Locale;
  }

  public LocaleReplace(content: string, replacements: any[]): string {
    return this.l10nService.Replace(content, replacements);
  }

  ngOnDestroy(): void {
    for (const s of this.subs) {
      s.unsubscribe();
    }
    this.subs = [];
  }

  ngOnInit(): void {

    this.subs.push(this.route.params.subscribe((params) => {
      this.routeRecipeId = params['id'] ? +params['id'] : undefined;
      if (this.routeRecipeId === undefined) {
        this.router.navigate(['/']);
        return;
      }
      this.loadRecipeById(this.routeRecipeId);
    }));

    this.subs.push(this.sharedDataService.RecipeEvents.subscribe((event) => {
      if (event === false || event.id !== this.routeRecipeId || !this.recipe || event.id !== this.recipe.id || event.etag !== this.recipe.modified)
        return;

      this.loadRecipeById(this.routeRecipeId);
    }));

    this.subs.push(this.l10nService.userLocale.subscribe((l) => {
      this.langCode.set(l);
      this.onToggleLocalization(this.localized());
    }));

  }

  loadRecipeById(id: number) {
    this.sharedDataService.getRecipe(id)
      .then((data: { id: number, etag: string, data: Recipe }) => {
        if (!this.recipe || this.recipe.id !== data.data.id || this.recipe.modified !== data.data.modified) {
          this.recipe = data.data;
        }
      })
      .catch((err) => {
        console.error(err)
        this.loadingFailed.set(true);
        this.wsService.ReportError({
          url: this.router.url,
          error: `Recipe with id ${id} not found.`,
          severity: 'E'
        });
        setTimeout(() => {
          this.router.navigate(['/']);
        }, 1000);
      });

  }

  onToggleLocalization(to: boolean): void {
    this.localized.set(to);
    this.langCodeVisible.set(to === true ? this.langCode() : (this.recipe?.userLocale ?? this.langCode()));
  }

}
