import { Component, OnDestroy, OnInit, signal } from '@angular/core';
import { IconLib } from '../../icons';
import { Router } from '@angular/router';
import { L10nService } from '../../svc/l10n.service';
import { SharedDataService } from '../../svc/shared-data.service';
import { UserSelf } from '../../types';
import { Subscription } from 'rxjs';
import { ApiService } from '../../svc/api.service';

@Component({
  selector: 'kb-mobile-menu',
  standalone: false,
  templateUrl: './mobile-menu.component.html',
  styleUrl: './mobile-menu.component.scss'
})
export class MobileMenuComponent implements OnDestroy, OnInit {

  Icons = IconLib;
  LoggedIn = signal<boolean>(false);
  ShowMenu = signal<boolean>(false);
  User = signal<UserSelf | false>(false);

  private subs: Subscription[] = [];

  constructor(
    private apiService: ApiService,
    private l10nService: L10nService,
    private router: Router,
    private sharedDataService: SharedDataService,
  ) { }

  ngOnDestroy(): void {
    this.subs.forEach((s) => s.unsubscribe());
    this.subs = [];
  }

  ngOnInit(): void {
    this.subs.push(this.apiService.isLoggedIn.subscribe((state) => {
      this.LoggedIn.set(state);
      this.User.set(this.apiService.User ?? false);
    }));
    this.subs.push(this.sharedDataService.ShowMenuBar.subscribe((state) => {
      this.ShowMenu.set(state);
    }));
  }

}
