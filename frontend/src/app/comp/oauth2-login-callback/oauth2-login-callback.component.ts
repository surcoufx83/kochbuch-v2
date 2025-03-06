import { Component, OnDestroy, OnInit, signal } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { first, Subscription } from 'rxjs';
import { IconLib } from '../../icons';
import { ApiService } from '../../svc/api.service';
import { L10nService } from '../../svc/l10n.service';
import { L10nLocale } from '../../svc/locales/types';
import { SharedDataService } from '../../svc/shared-data.service';
import { UserSelf } from '../../types';
import { HttpResponse, HttpStatusCode } from '@angular/common/http';

@Component({
  selector: 'kb-oauth2-login-callback',
  standalone: false,
  templateUrl: './oauth2-login-callback.component.html',
  styleUrl: './oauth2-login-callback.component.scss'
})
export class Oauth2LoginCallbackComponent implements OnDestroy, OnInit {

  Icons = IconLib;
  Busy = signal<false | 1 | 2>(false);
  Failed = signal<boolean>(false);
  LoggedIn = signal<boolean>(false);
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

  ngOnDestroy(): void {
    this.subs.forEach((s) => s.unsubscribe());
    this.subs = [];
  }

  ngOnInit(): void {
    const sub = this.apiService.isInitialized.subscribe((state) => {
      if (!state)
        return;
      setTimeout(() => {
        this.init();
        sub.unsubscribe();
      }, 0);
    });
  }

  private init(): void {
    this.subs.push(this.apiService.isLoggedIn.subscribe((state) => {
      if (state === true)
        this.router.navigate(['/']);
    }));

    this.subs.push(this.route.queryParamMap.subscribe((e) => {
      if (this.Busy() !== false)
        return;

      this.Busy.set(1);
      this.Failed.set(false);

      const state = e.get('state');
      const code = e.get('code');

      if (!state || !code) {
        this.router.navigate(['/']);
        return;
      }

      const sub = this.apiService.oauth2Callback(state, code).pipe(first()).subscribe((response) => {
        sub.unsubscribe();
        console.log('=====', response)
        if (response instanceof HttpResponse) {
          this.Busy.set(2);
          setTimeout(() => {
            this.postLoginQeryProfile();
          }, 500);
        }
        else {
          this.Failed.set(true);
          this.Busy.set(false);
        }
      });

    }));
  }

  private profileQueryCount = 0;
  postLoginQeryProfile(): void {
    if (this.profileQueryCount > 10) {
      this.Failed.set(true);
      this.Busy.set(false);
      return;
    }

    this.profileQueryCount++;

    const sub = this.apiService.loadUser().pipe(first()).subscribe((reply) => {
      sub.unsubscribe();
      if (!(reply instanceof HttpResponse && reply.status === HttpStatusCode.Ok)) {
        setTimeout(() => {
          this.postLoginQeryProfile();
        }, 1000);
      }
    });

  }

}
