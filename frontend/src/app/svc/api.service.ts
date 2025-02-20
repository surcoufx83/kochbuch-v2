import { Injectable } from '@angular/core';
import { BehaviorSubject, first, Subject } from 'rxjs';
import { HttpClient, HttpErrorResponse, HttpResponse, HttpStatusCode } from '@angular/common/http';
import { UserSelf } from '../types';

@Injectable({
  providedIn: 'root'
})
export class ApiService {

  private appParams?: ApiParams;

  private _isLoggedIn = new BehaviorSubject<boolean>(false);
  public isLoggedIn = this._isLoggedIn.asObservable();

  private _user?: UserSelf;

  constructor(
    private http: HttpClient,
  ) {
    this.loadAppParams();
  }

  public get(urlfragment: string, etag?: string): Subject<HttpResponse<unknown> | HttpErrorResponse | null> {

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
    this.get('params').pipe(first()).subscribe((r) => {
      if (r?.status === HttpStatusCode.Ok) {
        this.appParams = (r as HttpResponse<ApiParams>).body || undefined;
        if (this.appParams)
          localStorage.setItem('kbParams', JSON.stringify(this.appParams))
      }
    });
  }

  public get LoginUrl(): string | undefined {
    return this.appParams?.loginUrl;
  }

  public loadUser(): Subject<HttpResponse<unknown> | HttpErrorResponse | unknown | null> {

    let reply: Subject<HttpResponse<unknown> | HttpErrorResponse | unknown | null> = new Subject<HttpResponse<unknown> | HttpErrorResponse | unknown | null>();

    this.get('me').pipe(first()).subscribe((r) => {
      reply.next(r);
      reply.complete();
    });

    return reply;
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

  public get User(): UserSelf | undefined {
    return this._user;
  }

}

type ApiParams = {
  loginUrl: string,
}

export type PageErrorReport = {
  url: string,
  error: string,
  severity: 'I' | 'W' | 'E',
};