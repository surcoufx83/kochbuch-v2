import { Injectable } from '@angular/core';
import { first, Subject } from 'rxjs';
import { GenericApiReply } from '../types';
import { HttpClient, HttpErrorResponse, HttpResponse } from '@angular/common/http';

@Injectable({
  providedIn: 'root'
})
export class ApiService {

  constructor(
    private http: HttpClient,
  ) { }

  public get(urlfragment: string, etag?: string): Subject<HttpResponse<GenericApiReply> | HttpErrorResponse | null> {
    let reply: Subject<HttpResponse<GenericApiReply> | HttpErrorResponse | null> = new Subject<HttpResponse<GenericApiReply> | HttpErrorResponse | null>();
    this.http.get<GenericApiReply>(`/api/${urlfragment}`, { observe: 'response', headers: etag ? { 'If-None-Match': etag } : undefined })
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
