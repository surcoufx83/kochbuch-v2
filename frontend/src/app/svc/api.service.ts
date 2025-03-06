import { Injectable } from '@angular/core';
import { BehaviorSubject, first, Subject } from 'rxjs';
import { HttpClient, HttpErrorResponse, HttpResponse, HttpStatusCode } from '@angular/common/http';
import { UserSelf } from '../types';

@Injectable({
  providedIn: 'root'
})
export class ApiService {

  private appParams?: ApiParams;

  private _isLoggedIn = new BehaviorSubject<'unknown' | boolean>('unknown');
  public isLoggedIn = this._isLoggedIn.asObservable();

  private _isInitialized = new BehaviorSubject<boolean>(false);
  public isInitialized = this._isInitialized.asObservable();

  private _user?: UserSelf;

  constructor(
    private http: HttpClient,
  ) {
    this.loadAppParams();
  }

  public get(urlfragment: string, etag?: string): Subject<HttpResponse<unknown> | HttpErrorResponse | null> {
    if (!this._isInitialized.value && urlfragment != 'params') {
      throw {}
    }
    let reply: Subject<HttpResponse<unknown> | HttpErrorResponse | null> = new Subject<HttpResponse<unknown> | HttpErrorResponse | null>();

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

  private loadAppParams(): void {
    const cache = localStorage.getItem('kbParams');
    if (cache !== null)
      this.appParams = JSON.parse(cache) as ApiParams;

    if (this.appParams)
      this._isInitialized.next(true);

    this.get('params').pipe(first()).subscribe((r) => {
      if (r?.status === HttpStatusCode.Ok) {
        this.appParams = (r as HttpResponse<ApiParams>).body || undefined;
        if (this.appParams)
          localStorage.setItem('kbParams', JSON.stringify(this.appParams))
        this._isInitialized.next(true);

        const sub = this.loadUser().pipe(first()).subscribe(() => {
          sub.unsubscribe();
        });
      }
    });
  }

  public loadUser(): Subject<HttpResponse<unknown> | HttpErrorResponse | unknown | null> {
    let reply: Subject<HttpResponse<unknown> | HttpErrorResponse | unknown | null> = new Subject<HttpResponse<unknown> | HttpErrorResponse | unknown | null>();

    this.get('me').pipe(first()).subscribe((r) => {
      reply.next(r);
      reply.complete();
      if (r instanceof HttpResponse && r.status === HttpStatusCode.Ok) {
        this._user = r.body as UserSelf;
        this._isLoggedIn.next(true);
      }
      else {
        this._isLoggedIn.next(false);
      }
    });

    return reply;
  }

  public get LoginUrl(): string | undefined {
    return this.appParams?.loginUrl;
  }

  public logout(): void {
    const sub = this.post('logout', {}).pipe(first()).subscribe((reply) => {
      this.setCookie('session', '', -1);
      localStorage.removeItem('kbParams');
      location.replace('/');
    });
  }

  public oauth2Callback(state: string, code: string): Subject<HttpResponse<unknown> | HttpErrorResponse | unknown | null> {
    let reply: Subject<HttpResponse<unknown> | HttpErrorResponse | unknown | null> = new Subject<HttpResponse<unknown> | HttpErrorResponse | unknown | null>();

    this.post('login', {
      state: state,
      code: code
    }).pipe(first()).subscribe((r) => {
      reply.next(r);
      reply.complete();
    });

    return reply;
  }

  public post(urlfragment: string, payload: any): Subject<HttpResponse<unknown> | HttpErrorResponse | unknown | null> {
    if (!this._isInitialized.value && urlfragment != 'params') {
      throw {}
    }
    let reply: Subject<HttpResponse<unknown> | HttpErrorResponse | unknown | null> = new Subject<HttpResponse<unknown> | HttpErrorResponse | unknown | null>();

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

  public get Session(): string | undefined {
    return this.appParams?.session;
  }

  private setCookie(name: string, value: string, expireDays: number, path: string = '') {
    let d: Date = new Date();
    d.setTime(d.getTime() + expireDays * 24 * 60 * 60 * 1000);
    let expires: string = `expires=${d.toUTCString()}`;
    let cpath: string = path ? `; path=${path}` : '';
    document.cookie = `${name}=${value}; ${expires}${cpath}`;
  }

  public get User(): UserSelf | undefined {
    return this._user;
  }

}

type ApiParams = {
  loginUrl: string,
  session: string,
}

export type PageErrorReport = {
  url: string,
  error: string,
  severity: 'I' | 'W' | 'E',
};