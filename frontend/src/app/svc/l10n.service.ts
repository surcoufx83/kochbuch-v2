import { Injectable } from '@angular/core';
import { BehaviorSubject, Subject } from 'rxjs';
import { L10nLocale } from './locales/types';
import { KB_De } from './locales/de';
import { KB_En } from './locales/en';
import { KB_Fr } from './locales/fr';
import { enUS as FNS_En, de as FNS_De, fr as FNS_Fr, Locale } from 'date-fns/locale';
import { Duration, format, formatDate, formatDuration, parseISO } from 'date-fns';

@Injectable({
  providedIn: 'root'
})
export class L10nService {

  private availableLocales: { [key: string]: { locale: L10nLocale, datefns: Locale, flag: string, key: string } } = {
    'de': {
      locale: KB_De,
      datefns: FNS_De,
      flag: 'fi-de',
      key: 'de',
    },
    'fr': {
      locale: KB_Fr,
      datefns: FNS_Fr,
      flag: 'fi-fr',
      key: 'fr',
    },
    'en': {
      locale: KB_En,
      datefns: FNS_En,
      flag: 'fi-eu',
      key: 'en',
    },
  };
  private fallbackLocale = 'de';

  private _userLocale$: BehaviorSubject<string> = new BehaviorSubject<string>(this.fallbackLocale);
  public userLocale = this._userLocale$.asObservable();
  private userLocaleStr: string = '';

  private minutesToDurationCache: { [key: number]: { duration: Duration, uivalue: string } } = {};

  constructor() {

    const cache = localStorage.getItem('kbLocale');
    if (cache) {
      const cache2 = JSON.parse(cache) as LocaleStorage;
      if (cache2.locale) {
        this.userLocaleStr = cache2.locale;
      }
    }

    if (this.userLocaleStr === '') {
      for (let i = 0; i < navigator.languages.length; i++) {
        if (Object.keys(this.availableLocales).includes(navigator.languages[i].substring(0, 2))) {
          this.userLocaleStr = navigator.languages[i].substring(0, 2);
          break;
        }
      }
    }

    if (this.userLocaleStr === '') {
      this.userLocaleStr = this.fallbackLocale;
    }

    this._userLocale$.next(this.userLocaleStr);

    this.userLocale.subscribe((langcode) => {
      this.userLocaleStr = langcode;
    });
  }

  public get AvailableLocales(): { [key: string]: { locale: L10nLocale, datefns: Locale, flag: string, key: string } } {
    return this.availableLocales;
  }

  public FormatDate(date: string | number | Date, formatStr: string): string {
    if (typeof (date) === 'string') {
      date = parseISO(date);
    }
    return format(date, formatStr, {
      locale: this.availableLocales[this.LangCode].datefns ?? FNS_De
    });
  }

  public FormatNumber(n: number, o: Intl.NumberFormatOptions): string {
    return n.toLocaleString(this.LangCode, o);
  }

  public FormatDuration(inMinutes: number): string {
    if (this.minutesToDurationCache[inMinutes])
      return this.minutesToDurationCache[inMinutes].uivalue;

    let m = inMinutes;
    let d = 0;
    let h = 0;

    while (m > 1440) {
      d++;
      m -= 1440;
    }
    while (m > 60) {
      h++;
      m -= 60;
    }

    const strvalue = formatDuration({
      minutes: m,
      hours: h,
      days: d,
    }, {
      locale: this.availableLocales[this.LangCode].datefns ?? FNS_De,
      format: ['days', 'hours', 'minutes']
    });

    this.minutesToDurationCache[inMinutes] = {
      duration: {
        minutes: m,
        hours: h,
        days: d,
      },
      uivalue: strvalue
    };

    return strvalue;
  }

  public FormatVote(n: number): string {
    return this.FormatNumber(n, { maximumSignificantDigits: 2 });
  }

  public get LangCode(): string {
    return this.userLocaleStr;
  }

  public get Locale(): L10nLocale {
    return this.availableLocales[this.userLocaleStr].locale;
  }

  public Replace(content: string, replacements: any[]): string {
    for (let i = 0; i < replacements.length; i++) {
      content = content.replace(`[${i}]`, replacements[i]);
    }
    return content;
  }

  SetLocale(code: string | null): void {
    if (code === null) {
      localStorage.removeItem('kbLocale');
      for (let i = 0; i < navigator.languages.length; i++) {
        if (Object.keys(this.availableLocales).includes(navigator.languages[i].substring(0, 2))) {
          code = navigator.languages[i].substring(0, 2);
          break;
        }
      }
      if (code === null)
        code = this.fallbackLocale;
    }
    if (!Object.keys(this.availableLocales).includes(code))
      code = this.fallbackLocale;

    this._userLocale$.next(code);
    localStorage.setItem('kbLocale', JSON.stringify({
      locale: code
    }));
  }

}

type LocaleStorage = {
  locale?: 'de' | 'en' | 'fr' | null,
}