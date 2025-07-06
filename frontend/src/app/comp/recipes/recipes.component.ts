import { Component, OnDestroy, OnInit, signal } from '@angular/core';
import { SharedDataService } from '../../svc/shared-data.service';
import { Router } from '@angular/router';
import { L10nService } from '../../svc/l10n.service';
import { Subscription } from 'rxjs';
import { IconLib } from '../../icons';
import { L10nLocale } from '../../svc/locales/types';
import { Recipe, UserSelf } from '../../types';
import { WebSocketService } from '../../svc/web-socket.service';

@Component({
  selector: 'kb-recipes',
  standalone: false,
  templateUrl: './recipes.component.html',
  styleUrl: './recipes.component.scss'
})
export class RecipesComponent implements OnDestroy, OnInit {

  Icons = IconLib;
  LangCode = signal<string>('de');
  LoggedIn = signal<boolean>(false);
  Recipes = signal<Recipe[]>([]);
  User = signal<UserSelf | false>(false);

  private subs: Subscription[] = [];

  constructor(
    private l10nService: L10nService,
    private router: Router,
    private sharedDataService: SharedDataService,
    private wsService: WebSocketService,
  ) {
    this.LangCode.set(this.l10nService.LangCode);
  }

  get Locale(): L10nLocale {
    return this.l10nService.Locale;
  }

  ngOnDestroy(): void {
    this.subs.forEach((s) => s.unsubscribe());
    this.subs = [];
  }

  ngOnInit(): void {
    this.subs.push(this.wsService.isLoggedIn.subscribe((state) => {
      if (state === 'unknown')
        return;
      this.LoggedIn.set(state);
    }));

    this.subs.push(this.wsService.User.subscribe((u) => this.User.set(u ?? false)));

    this.subs.push(this.l10nService.userLocale.subscribe((l) => {
      if (l !== this.LangCode()) {
        this.LangCode.set(l);
      }
    }));

    this.subs.push(this.sharedDataService.Recipes.subscribe((items) => {
      console.log(items)
      if (Object.keys(items).length == 0)
        return;
      this.Recipes.set(
        Object.values(items)
          .filter((a) => a.pictures != null && a.pictures.length > 0)
          .sort((a, b) => (a.published ?? a.edited ?? a.modified) > (b.published ?? b.edited ?? b.modified) ? -1 : 1)
      );
    }));
  }

}
