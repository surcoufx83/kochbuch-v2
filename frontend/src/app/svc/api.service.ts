import { Injectable } from '@angular/core';
import { first, Subject } from 'rxjs';
import { HttpClient, HttpErrorResponse, HttpResponse } from '@angular/common/http';

@Injectable({
  providedIn: 'root'
})
export class ApiService {

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

}
