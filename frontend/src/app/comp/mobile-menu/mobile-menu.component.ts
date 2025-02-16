import { Component, ElementRef, OnDestroy, OnInit, signal, ViewChild } from '@angular/core';
import { IconLib } from '../../icons';
import { NavigationEnd, Router } from '@angular/router';
import { L10nService } from '../../svc/l10n.service';
import { SharedDataService } from '../../svc/shared-data.service';
import { UserSelf } from '../../types';
import { Subscription } from 'rxjs';
import { ApiService } from '../../svc/api.service';
import { L10nLocale } from '../../svc/locales/types';

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
  SearchState = signal<boolean>(false);
  User = signal<UserSelf | false>(false);
  Url = signal<string>('');

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
    this.subs.push(this.router.events.subscribe((e) => {
      if (e instanceof NavigationEnd) {
        console.log(e.urlAfterRedirects);
        if (e.urlAfterRedirects === '/search') {
          this.onStartSearch();
        }
        else {
          this.onCancelSearch();
        }
        this.Url.set(e.urlAfterRedirects);
      }
    }));
    this.subs.push(this.sharedDataService.SearchIsActive.subscribe((state) => {
      console.log('SearchIsActive', state)
      this.SearchState.set(state)
    }));
    this.subs.push(this.sharedDataService.SearchCategories.subscribe((cats) => {
      console.log(cats)
    }));
  }

  onClickSearchButtonIcon($event: MouseEvent): void {
    if (this.SearchState()) {
      this.onCancelSearch();
      $event.stopPropagation();
    }
  }

  private focusInterval?: number;
  private focusCheck: number = 0;
  onStartSearch(): void {
    if (this.SearchState())
      return;
    this.sharedDataService.SetSearchState(true);
    this.focusCheck = 0;
    this.focusInterval = setInterval(() => {
      if (this.searchField || this.focusCheck > 10) {
        clearInterval(this.focusInterval)
        this.searchField?.nativeElement.focus();
      }
      this.focusCheck++;
      console.log(this.searchField)
    }, 100);
  }

  onCancelSearch(): void {
    this.sharedDataService.SetSearchState(false);
  }

}
