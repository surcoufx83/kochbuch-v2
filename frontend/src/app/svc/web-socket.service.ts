import { Injectable } from '@angular/core';
import { BehaviorSubject, Subject } from 'rxjs';
import { UserSelf } from '../types';
import { addDays } from 'date-fns';
import { HttpStatusCode } from '@angular/common/http';
import { v4 as uuidv4 } from 'uuid';

@Injectable({
  providedIn: 'root'
})
export class WebSocketService {

  private appParams?: WsHelloMessageContent;
  private socket?: WebSocket;

  private _events = new BehaviorSubject<WsMessage | null>(null);
  public events = this._events.asObservable();

  private msgQueue: WsMessage[] = [];

  private _isLoggedIn = new BehaviorSubject<'unknown' | boolean>('unknown');
  public isLoggedIn = this._isLoggedIn.asObservable();

  private _isConnected = new BehaviorSubject<boolean>(false);
  public isConnected = this._isConnected.asObservable();

  private _user = new BehaviorSubject<UserSelf | null>(null);
  public User = this._user.asObservable();

  constructor() {
    this.loadSession();
    this.connect();
    this.reconnect();
  }

  private cancel(): void {
    this._isConnected.next(false);
    this.socket?.close();
    this.socket = undefined;
  }

  private connect(): void {
    this.cancel();

    this.socket = new WebSocket(`/ws`);

    this.socket.onopen = () => {
      this.SendMessage({
        type: 'auth',
        content: JSON.stringify({ token: this.appParams?.connection.session ?? '' })
      });
    };

    this.socket.onmessage = (event) => {
      const message = JSON.parse(event.data) as WsMessage;
      if (message.type === 'hello') {
        this.appParams = JSON.parse(message.content) as WsHelloMessageContent;
        this.saveSession();
        this._user.next(this.appParams.user && this.appParams.loggedIn ? this.appParams.user : null);
        if (this.appParams.loggedIn && this.appParams.connection.session) {
          document.cookie = `session=${this.appParams.connection.session}; exires=${addDays(new Date(), 365).toUTCString()}; path=/`;
        }
        this._isLoggedIn.next(this.appParams.loggedIn);
        this._isConnected.next(true);
        this.ResendFromQueue();
      }
      else
        this._events.next(message);
    };

    this.socket.onerror = (error) => {
      console.error('WebSocket error: ', error);
    };

    this.socket.onclose = (event) => {
      if (event.wasClean) {
        console.log(`Closed cleanly with code ${event.code}`);
      } else {
        console.error(`Closed with error code ${event.code}`);
      }
    };

  }

  public GetLoginUrl(): string | undefined {
    return this.appParams?.connection.loginUrl;
  }

  public GetUser(): UserSelf | null {
    return this._user.value;
  }

  private loadSession(): void {
    const data = localStorage.getItem('kbSession');
    if (data != null) {
      this.appParams = JSON.parse(data) as WsHelloMessageContent;
    }
  }

  public Login(state: string, code: string): Subject<boolean | 'wait'> {
    let reply = new Subject<boolean | 'wait'>();
    const sub = this.events.subscribe((e) => {
      if (!e || e.type !== 'oauth2_response') {
        return;
      }
      reply.next(e.content === '202/Accepted');
      reply.complete();
      sub.unsubscribe();
    });
    this.SendMessage({
      type: 'oauth2_callback',
      content: JSON.stringify({
        state: state,
        code: code
      })
    });
    return reply;
  }

  public Logout(): void {
    this.SendMessage({
      type: 'bye',
      content: ""
    });
  }

  private reconnect(): void {
    const sub = setInterval(() => {
      if (this.socket && this.socket.readyState === WebSocket.OPEN)
        return;

      this.connect();
    }, 1000);
  }

  ReportError(data: {
    url: string,
    error: string,
    severity: 'I' | 'E' | 'W',
  }): void {
    this.SendMessage({
      type: 'error_report',
      content: data
    });
  }

  private saveSession(): void {
    if (!this.appParams)
      return;
    localStorage.setItem('kbSession', JSON.stringify(this.appParams));
  }

  SendMessage(msg: WsMessage): boolean {
    if (this.socket && this.socket.readyState === WebSocket.OPEN) {
      if (typeof msg.content !== 'string') {
        msg.content = JSON.stringify(msg.content);
      }
      this.socket.send(JSON.stringify(msg));
      return true;
    } else {
      // console.warn('WebSocket is not open => continue after reconnect');
      this.msgQueue.push(msg);
      return false;
    }
  }

  SendMessageAndWait(msg: WsMessage, waittime: number = 30000): Promise<[number, any]> {
    return new Promise<[number, any]>((resolve, reject) => {
      if (!this.socket || this.socket.readyState !== WebSocket.OPEN) {
        reject(HttpStatusCode.ServiceUnavailable);
        return;
      }

      if (typeof msg.content !== 'string') {
        msg.content = JSON.stringify(msg.content);
      }

      msg.state = uuidv4();

      const sub = this.events.subscribe((ev) => {
        if (!ev || ev.state !== msg.state)
          return;

        const content = JSON.parse(ev.content) as WsCommonResponse;

        if (content.error === HttpStatusCode.Accepted) {
          resolve([HttpStatusCode.Accepted, content]);
          sub.unsubscribe();
        }
        else {
          resolve([content.error ?? HttpStatusCode.Conflict, content]);
          sub.unsubscribe();
        }
      });

      this.socket.send(JSON.stringify(msg));

      setTimeout(() => {
        if (sub) {
          reject(HttpStatusCode.RequestTimeout);
          sub.unsubscribe();
        }
      }, waittime);

    });

  }

  ResendFromQueue(): void {
    if (!this.socket || this.socket.readyState !== WebSocket.OPEN)
      return;
    const datacopy = [...this.msgQueue];
    while (datacopy.length > 0) {
      const msg = datacopy.splice(0, 1);
      if (msg.length === 1) {
        if (!this.SendMessage(msg[0]))
          return;
      }
    }
  }

}

export type WsMessage = {
  type: string,
  content: any,
  state?: string,
}

export type WsCommonResponse = {
  error?: number,
  message?: string,
}

export type WsHelloMessageContent = {
  connection: {
    loginUrl: string,
    session: string,
  },
  loggedIn: boolean,
  user?: UserSelf | null,
}
