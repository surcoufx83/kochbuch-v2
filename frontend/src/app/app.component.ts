import { Component, OnDestroy, OnInit, signal } from '@angular/core';
import { Title } from '@angular/platform-browser';
import { L10nService } from './svc/l10n.service';
import { L10nLocale } from './svc/locales/types';
import { SharedDataService } from './svc/shared-data.service';
import { Subscription } from 'rxjs';

@Component({
  selector: 'kb-root',
  templateUrl: './app.component.html',
  standalone: false,
  styleUrl: './app.component.scss'
})
export class AppComponent implements OnInit, OnDestroy {

  ShowMenu = signal<boolean>(false);

  private subs: Subscription[] = [];

  constructor(
    private l10nService: L10nService,
    private sharedDataService: SharedDataService,
    htmlTitleService: Title,
  ) {
    sharedDataService.PageTitle.subscribe((t) => htmlTitleService.setTitle(t));
  }

  Locale(): L10nLocale {
    return this.l10nService.Locale;
  }

  ngOnDestroy(): void {
    this.subs.forEach((s) => s.unsubscribe());
    this.subs = [];
  }

  ngOnInit(): void {
    this.subs.push(this.sharedDataService.ShowMenuBar.subscribe((state) => {
      this.ShowMenu.set(state);
    }));
  }

}
