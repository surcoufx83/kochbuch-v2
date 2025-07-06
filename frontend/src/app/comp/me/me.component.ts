import { Component, OnDestroy, OnInit, signal } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { Subscription } from 'rxjs';
import { IconLib } from '../../icons';
import { L10nService } from '../../svc/l10n.service';
import { L10nLocale } from '../../svc/locales/types';
import { SharedDataService } from '../../svc/shared-data.service';
import { WebSocketService } from '../../svc/web-socket.service';
import { Collection, UserSelf } from '../../types';

@Component({
  selector: 'kb-me',
  standalone: false,
  templateUrl: './me.component.html',
  styleUrl: './me.component.scss'
})
export class MeComponent implements OnDestroy, OnInit {

  icons = IconLib;
  loggedIn = signal<boolean | 'unknown'>('unknown');
  pageRef = signal<string>('');
  pageSrc = signal<string>('');
  user: UserSelf | false = false;
  userCollections: Collection[] = [];

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

  public LocaleReplace(content: string, replacements: any[]): string {
    return this.l10nService.Replace(content, replacements);
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
    this.subs.push(this.wsService.isLoggedIn.subscribe((state) => {
      this.loggedIn.set(state);
      this.user = this.wsService.GetUser() ?? false;
      console.log(this.user)

      if (this.user) {
        this.userCollections = Object.values(this.user.collections).filter(x => x.deleted === null).sort((a, b) => a.title.toLocaleLowerCase().localeCompare(b.title.toLocaleLowerCase()));
      }
      else {
        this.userCollections = [];
      }

    }));
    this.subs.push(this.route.queryParamMap.subscribe((e) => {
      this.pageRef.set(`${e.get('ref')}`);
      this.pageSrc.set(`${e.get('source')}`);
    }))
  }

}
