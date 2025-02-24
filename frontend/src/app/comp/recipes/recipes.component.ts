import { Component, OnDestroy, OnInit, signal } from '@angular/core';
import { SharedDataService } from '../../svc/shared-data.service';
import { Router } from '@angular/router';
import { L10nService } from '../../svc/l10n.service';
import { ApiService } from '../../svc/api.service';
import { Subscription } from 'rxjs';
import { IconLib } from '../../icons';
import { L10nLocale } from '../../svc/locales/types';
import { Recipe, UserSelf } from '../../types';

@Component({
  selector: 'kb-recipes',
  standalone: false,
  templateUrl: './recipes.component.html',
  styleUrl: './recipes.component.scss'
})
export class RecipesComponent implements OnDestroy, OnInit {

  Icons = IconLib;
  LoggedIn = signal<boolean>(false);
  Recipes = signal<Recipe[]>([]);
  User = signal<UserSelf | false>(false);

  private subs: Subscription[] = [];

  constructor(
    private apiService: ApiService,
    private l10nService: L10nService,
    private router: Router,
    private sharedDataService: SharedDataService,
  ) { }

  get Locale(): L10nLocale {
    return this.l10nService.Locale;
  }

  ngOnDestroy(): void {
    this.subs.forEach((s) => s.unsubscribe());
    this.subs = [];
  }

  ngOnInit(): void {
    this.subs.push(this.apiService.isLoggedIn.subscribe((state) => {
      this.LoggedIn.set(state);
      this.User.set(this.apiService.User ?? false);
    }));
    this.subs.push(this.sharedDataService.Recipes.subscribe((items) => {
      console.log(items)
      this.Recipes.set(Object.values(items).sort((a, b) => a.published > b.published ? 1 : -1));
    }));
  }

}
