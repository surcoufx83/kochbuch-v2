import { Injectable } from '@angular/core';
import { Subject } from 'rxjs';
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

  private _userLocale$: Subject<string> = new Subject<string>();
  public userLocale = this._userLocale$.asObservable();
  private userLocaleStr: string = this.fallbackLocale;

  constructor() {
    this.userLocale.subscribe((langcode) => {
      this.userLocaleStr = langcode;
    });
    for (let i = 0; i < navigator.languages.length; i++) {
      if (Object.keys(this.availableLocales).includes(navigator.languages[i].substring(0, 2))) {
        this.userLocaleStr = navigator.languages[i].substring(0, 2);
        break;
      }
    }
    if (this.userLocaleStr === '')
      this.userLocaleStr = this.fallbackLocale;
    this._userLocale$.next(this.userLocaleStr);
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

}
