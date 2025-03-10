import { Component, ElementRef, OnDestroy, OnInit, signal, ViewChild } from '@angular/core';
import { IconLib } from '../../icons';
import { NavigationEnd, Router } from '@angular/router';
import { L10nService } from '../../svc/l10n.service';
import { SharedDataService } from '../../svc/shared-data.service';
import { UserSelf } from '../../types';
import { Subscription } from 'rxjs';
import { L10nLocale } from '../../svc/locales/types';
import { format } from "date-fns";
import { enGB } from "date-fns/locale";
import { FormControl, FormGroup } from '@angular/forms';
import { WebSocketService } from '../../svc/web-socket.service';

@Component({
  selector: 'kb-mobile-menu',
  standalone: false,
  templateUrl: './mobile-menu.component.html',
  styleUrl: './mobile-menu.component.scss'
})
export class MobileMenuComponent implements OnDestroy, OnInit {

  @ViewChild('search') searchField?: ElementRef;

  Icons = IconLib;
  LoggedIn = signal<boolean>(false);
  SearchPlaceholder = signal<string>('');
  SearchState = signal<boolean>(false);
  User = signal<UserSelf | false>(false);
  IsSearchUrl = signal<boolean>(false);

  SearchForm = new FormGroup({
    phrase: new FormControl<string>('')
  })

  private subs: Subscription[] = [];

  constructor(
    private l10nService: L10nService,
    private router: Router,
    private sharedDataService: SharedDataService,
    private wsService: WebSocketService,
  ) {
    const month = format(Date.now(), 'LLL', { locale: enGB }).toLowerCase();
    const templocale = this.Locale.floatingMenu.searchButton.searchInput.placeholder as { [key: string]: string };
    this.SearchPlaceholder.set(templocale[month]);
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
      this.User.set(this.wsService.GetUser() ?? false);
    }));
    this.subs.push(this.router.events.subscribe((e) => {
      if (e instanceof NavigationEnd) {
        this.IsSearchUrl.set(e.urlAfterRedirects.startsWith('/search'));
        if (e.urlAfterRedirects.startsWith('/search')) {
          this.onStartSearch(this.searchField?.nativeElement);
        }
        else {
          this.onCancelSearch();
        }
      }
    }));
  }

  onCancelSearch(): void {
    this.sharedDataService.SetSearchState(false);
    this.SearchState.set(false);
  }

  onClickSearchButtonIcon($event: MouseEvent): void {
    if (this.SearchState()) {
      this.onCancelSearch();
      $event.stopPropagation();
    }
  }

  onStartSearch(searchField?: HTMLInputElement): void {
    if (this.SearchState())
      return;
    this.sharedDataService.SetSearchState(true);
    this.SearchState.set(true);
    searchField?.focus();
  }

  onSubmitSearch($event: Event): void {
    if (!this.SearchForm.controls.phrase.value)
      return;

    this.router.navigate(['/search'], {
      queryParams: {
        search: this.SearchForm.controls.phrase.value ?? ''
      }
    })
  }

}
