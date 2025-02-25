import { Component, OnChanges, OnDestroy, OnInit, signal, SimpleChanges } from '@angular/core';
import { ActivatedRoute, ActivatedRouteSnapshot, Router } from '@angular/router';
import { Subscription } from 'rxjs';
import { IconLib } from '../../icons';
import { ApiService } from '../../svc/api.service';
import { L10nService } from '../../svc/l10n.service';
import { L10nLocale } from '../../svc/locales/types';
import { SharedDataService } from '../../svc/shared-data.service';
import { UserSelf } from '../../types';

@Component({
  selector: 'kb-me',
  standalone: false,
  templateUrl: './me.component.html',
  styleUrl: './me.component.scss'
})
export class MeComponent implements OnDestroy, OnInit {

  Icons = IconLib;
  LoggedIn = signal<boolean | 'unknown'>('unknown');
  PageRef = signal<string>('');
  PageSrc = signal<string>('');
  User = signal<UserSelf | false>(false);

  private subs: Subscription[] = [];

  constructor(
    private apiService: ApiService,
    private l10nService: L10nService,
    private route: ActivatedRoute,
    private router: Router,
    private sharedDataService: SharedDataService,
  ) { }

  get Locale(): L10nLocale {
    return this.l10nService.Locale;
  }

  logout(): void {
    this.apiService.logout();
  }

  ngOnDestroy(): void {
    this.subs.forEach((s) => s.unsubscribe());
    this.subs = [];
  }

  ngOnInit(): void {
    console.log('init')
    this.subs.push(this.apiService.isLoggedIn.subscribe((state) => {
      this.LoggedIn.set(state);
      this.User.set(this.apiService.User ?? false);
    }));
    this.subs.push(this.route.queryParamMap.subscribe((e) => {
      this.PageRef.set(`${e.get('ref')}`)
      this.PageSrc.set(`${e.get('source')}`)
      console.log(e)
    }))
  }

}
