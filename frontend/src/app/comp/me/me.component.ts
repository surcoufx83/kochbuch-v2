import { Component, OnChanges, OnDestroy, OnInit, signal, SimpleChanges } from '@angular/core';
import { ActivatedRoute, ActivatedRouteSnapshot, Router } from '@angular/router';
import { Subscription } from 'rxjs';
import { IconLib } from '../../icons';
import { L10nService } from '../../svc/l10n.service';
import { L10nLocale } from '../../svc/locales/types';
import { SharedDataService } from '../../svc/shared-data.service';
import { UserSelf } from '../../types';
import { WebSocketService } from '../../svc/web-socket.service';

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
    private l10nService: L10nService,
    private route: ActivatedRoute,
    private router: Router,
    private sharedDataService: SharedDataService,
    private wsService: WebSocketService,
  ) { }

  get Locale(): L10nLocale {
    return this.l10nService.Locale;
  }

  logout(): void {
    this.wsService.Logout();
    setTimeout(() => {
      this.router.navigate(['/']);
    }, 5);
  }

  ngOnDestroy(): void {
    this.subs.forEach((s) => s.unsubscribe());
    this.subs = [];
  }

  ngOnInit(): void {
    console.log('init')
    this.subs.push(this.wsService.isLoggedIn.subscribe((state) => {
      this.LoggedIn.set(state);
      this.User.set(this.wsService.GetUser() ?? false);
    }));
    this.subs.push(this.route.queryParamMap.subscribe((e) => {
      this.PageRef.set(`${e.get('ref')}`)
      this.PageSrc.set(`${e.get('source')}`)
      console.log(e)
    }))
  }

}
