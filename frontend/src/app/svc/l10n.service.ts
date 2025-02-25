import { Injectable } from '@angular/core';
import { BehaviorSubject, Subject } from 'rxjs';
import { L10nLocale } from './locales/types';
import { KB_De } from './locales/de';
import { KB_En } from './locales/en';
import { KB_Fr } from './locales/fr';

@Injectable({
  providedIn: 'root'
})
export class L10nService {

  private availableLocales: { [key: string]: L10nLocale } = {
    'de': KB_De,
    'en': KB_En,
    'fr': KB_Fr,
  };
  private fallbackLocale = 'de';

  private _userLocale$: BehaviorSubject<string> = new BehaviorSubject<string>(this.fallbackLocale);
  public userLocale = this._userLocale$.asObservable();
  private userLocaleStr: string = '';

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

  public get LangCode(): string {
    return this.userLocaleStr;
  }

  public get Locale(): L10nLocale {
    return this.availableLocales[this.userLocaleStr];
  }

  public replace(content: string, replacements: any[]): string {
    for (let i = 0; i < replacements.length; i++) {
      content = content.replace(`[${i}]`, replacements[i]);
    }
    return content;
  }

  setLocale(code: 'de' | 'en' | 'fr' | null): void {
    if (code === null) {
      localStorage.removeItem('kbLocale');
      for (let i = 0; i < navigator.languages.length; i++) {
        if (Object.keys(this.availableLocales).includes(navigator.languages[i].substring(0, 2))) {
          code = navigator.languages[i].substring(0, 2) as 'de' | 'en' | 'fr';
          break;
        }
      }
      if (code === null)
        code = this.fallbackLocale as 'de' | 'en' | 'fr';
    }
    this._userLocale$.next(code);
    localStorage.setItem('kbLocale', JSON.stringify({
      locale: code
    }));
  }

}

type LocaleStorage = {
  locale?: 'de' | 'en' | 'fr' | null,
}