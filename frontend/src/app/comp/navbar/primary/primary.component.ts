import { Component, EventEmitter, Input, OnDestroy, OnInit, Output, signal } from '@angular/core';
import { IconLib } from '../../../icons';
import { L10nService } from '../../../svc/l10n.service';
import { L10nLocale } from '../../../svc/locales/types';
import { UserSelf } from '../../../types';
import { SharedDataService } from '../../../svc/shared-data.service';
import { Subscription } from 'rxjs';

@Component({
  selector: 'kb-navbar-primary',
  standalone: false,
  templateUrl: './primary.component.html',
  styleUrl: './primary.component.scss'
})
export class PrimaryComponent implements OnInit, OnDestroy {

  Icons = IconLib;
  ActiveLocale = signal<string>('');
  ShownLocales: { flag: string, key: string }[];
  ShowLanguageSelector = signal<boolean>(false);

  private subs: Subscription[] = [];

  constructor(
    private l10nService: L10nService,
    private sharedDataService: SharedDataService,
  ) {
    this.ShownLocales = Object.values(this.l10nService.AvailableLocales);
  }

  get Locale(): L10nLocale {
    return this.l10nService.Locale;
  }

  ngOnDestroy(): void {
    this.subs.forEach((s) => s.unsubscribe());
    this.subs = [];
  }

  ngOnInit(): void {
    this.subs.push(this.l10nService.userLocale.subscribe((l) => this.ActiveLocale.set(l)));
  }

  setLocale(code: string): void {
    this.l10nService.SetLocale(code);
  }

}
