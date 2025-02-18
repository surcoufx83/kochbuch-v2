import { Injectable } from '@angular/core';
import { BehaviorSubject, first, Subject } from 'rxjs';
import { HttpClient, HttpErrorResponse, HttpResponse } from '@angular/common/http';
import { UserSelf } from '../types';

@Injectable({
  providedIn: 'root'
})
export class ApiService {

  private _isLoggedIn = new BehaviorSubject<boolean>(false);
  public isLoggedIn = this._isLoggedIn.asObservable();

  private _user?: UserSelf;

  constructor(
    private http: HttpClient,
  ) { }

  public get(urlfragment: string, etag?: string): Subject<HttpResponse<unknown> | HttpErrorResponse | null> {
    let reply: Subject<HttpResponse<unknown> | HttpErrorResponse | null> = new Subject<HttpResponse<unknown> | HttpErrorResponse | null>();
    console.log(`GET /api/${urlfragment}`, etag);
    this.http.get<unknown>(`/api/${urlfragment}`, { observe: 'response', headers: etag ? { 'If-None-Match': etag } : undefined })
      .pipe(first()).subscribe({
        next: (res) => {
          reply.next(res)
        },
        error: (err) => {
          reply.next(err)
        }
      });
    return reply;
  }

  public post(urlfragment: string, payload: any): Subject<HttpResponse<unknown> | HttpErrorResponse | unknown | null> {
    let reply: Subject<HttpResponse<unknown> | HttpErrorResponse | unknown | null> = new Subject<HttpResponse<unknown> | HttpErrorResponse | unknown | null>();
    console.log(`POST /api/${urlfragment}`, payload);
    this.http.post<unknown>(`/api/${urlfragment}`, payload, { observe: 'response' })
      .pipe(first()).subscribe({
        next: (res) => {
          reply.next(res)
        },
        error: (err) => {
          reply.next(err)
        }
      });
    return reply;
  }

  public reportError(report: PageErrorReport): void {
    this.post('errorreport', report).pipe(first()).subscribe(() => { });
  }

  public get User(): UserSelf | undefined {
    return this._user;
  }

}

export type PageErrorReport = {
  url: string,
  error: string,
  severity: 'I' | 'W' | 'E',
};