import { Component, OnInit } from '@angular/core';
import { IconLib } from '../../icons';
import { L10nService } from '../../svc/l10n.service';
import { L10nLocale } from '../../svc/locales/types';
import { ApiService } from '../../svc/api.service';
import { Router } from '@angular/router';

@Component({
  selector: 'kb-missing-link',
  standalone: false,
  templateUrl: './missing-link.component.html',
  styleUrl: './missing-link.component.scss'
})
export class MissingLinkComponent implements OnInit {

  Icons = IconLib;

  constructor(
    private apiService: ApiService,
    private l10nService: L10nService,
    private router: Router,
  ) { }

  get Locale(): L10nLocale {
    return this.l10nService.Locale;
  }

  ngOnInit(): void {
    this.apiService.reportError({
      url: this.router.url,
      error: 'Route not configured',
      severity: 'E'
    });
  }

}
