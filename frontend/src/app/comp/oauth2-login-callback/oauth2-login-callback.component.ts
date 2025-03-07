import { Component, OnDestroy, OnInit, signal } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { first, Subscription } from 'rxjs';
import { IconLib } from '../../icons';
import { L10nService } from '../../svc/l10n.service';
import { L10nLocale } from '../../svc/locales/types';
import { SharedDataService } from '../../svc/shared-data.service';
import { UserSelf } from '../../types';
import { HttpResponse, HttpStatusCode } from '@angular/common/http';
import { WebSocketService } from '../../svc/web-socket.service';

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
    private l10nService: L10nService,
    private route: ActivatedRoute,
    private router: Router,
    private sharedDataService: SharedDataService,
    private wsService: WebSocketService,
  ) { }

  get Locale(): L10nLocale {
    return this.l10nService.Locale;
  }

  ngOnDestroy(): void {
    this.subs.forEach((s) => s.unsubscribe());
    this.subs = [];
  }

  ngOnInit(): void {
    const sub = this.wsService.isConnected.subscribe((state) => {
      if (!state)
        return;
      setTimeout(() => {
        this.init();
        sub.unsubscribe();
      }, 0);
    });
  }

  private init(): void {
    this.subs.push(this.wsService.isLoggedIn.subscribe((state) => {
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

      const sub = this.wsService.Login(state, code).subscribe((state) => {
        if (state === 'wait')
          return;
        if (state === false) {
          this.Failed.set(true);
          this.Busy.set(false);
        } else {
          this.Busy.set(2);
        }
      });
    }));
  }

}
