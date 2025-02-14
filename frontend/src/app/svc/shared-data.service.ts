import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class SharedDataService {

  private _pageTitle = new BehaviorSubject<string>('');
  public PageTitle = this._pageTitle.asObservable();

  constructor() { }

  public SetTitle(title: string) {
    this._pageTitle.next(title);
  }

}
