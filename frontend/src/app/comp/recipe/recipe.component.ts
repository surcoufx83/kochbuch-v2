import { Component, OnChanges, OnDestroy, OnInit, signal, SimpleChanges } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { Subscription } from 'rxjs';
import { SharedDataService } from '../../svc/shared-data.service';
import { L10nService } from '../../svc/l10n.service';
import { Recipe } from '../../types';
import { L10nLocale } from '../../svc/locales/types';
import { WebSocketService } from '../../svc/web-socket.service';

@Component({
  selector: 'kb-recipe',
  standalone: false,
  templateUrl: './recipe.component.html',
  styleUrl: './recipe.component.scss'
})
export class RecipeComponent implements OnDestroy, OnInit {

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

  ngOnDestroy(): void {
    for (const s of this.subs) {
      s.unsubscribe();
    }
    this.subs = [];
  }

  ngOnInit(): void {

    this.subs.push(this.route.params.subscribe((params) => {
      console.log(params)
      this.routeRecipeId = params['id'] ? +params['id'] : undefined;
      if (this.routeRecipeId === undefined) {
        this.router.navigate(['/']);
        return;
      }
      this.loadRecipeById(this.routeRecipeId);
    }));

    this.subs.push(this.sharedDataService.RecipeEvents.subscribe((event) => {
      console.log(event)
      if (event === false || event.id !== this.routeRecipeId || !this.recipe || event.id !== this.recipe.id || event.etag !== this.recipe.modified)
        return;
      this.loadRecipeById(this.routeRecipeId);
    }));

  }

  loadRecipeById(id: number) {
    console.log(`RecipeComponent::loadRecipeById(${id})`);

    this.sharedDataService.getRecipe(id)
      .then((data: { id: number, etag: string, data: Recipe }) => {
        this.recipe = data.data;
        console.log(this.recipe);
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

}
