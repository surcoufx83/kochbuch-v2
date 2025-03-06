import { Injectable } from '@angular/core';
import { ApiService } from './api.service';
import { BehaviorSubject } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class WebSocketService {

  private socket?: WebSocket;

  private _events = new BehaviorSubject<WsMessage | null>(null);
  public events = this._events.asObservable();

  constructor(
    private apiService: ApiService,
  ) {
    const sub = setInterval(() => {
      if (!this.apiService.Session)
        return;

      if (this.socket && this.socket.readyState === WebSocket.OPEN)
        return;

      this.connect();
    }, 1000);
  }

  cancel(): void {
    this.socket?.close();
    this.socket = undefined;
  }

  connect(): void {
    this.cancel();

    if (!this.apiService.Session)
      return;

    this.socket = new WebSocket(`/ws`);

    this.socket.onopen = () => {
      console.log('WebSocket connected');
      this.SendMessage({
        type: 'auth',
        content: JSON.stringify({ token: this.apiService.Session })
      })
    };

    this.socket.onmessage = (event) => {
      console.log('Message from server: ', event.data);
      this._events.next(JSON.parse(event.data) as WsMessage);
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

  SendMessage(msg: WsMessage) {
    if (this.socket && this.socket.readyState === WebSocket.OPEN) {
      this.socket.send(JSON.stringify(msg));
    } else {
      console.error('WebSocket is not open');
    }
  }

}

export type WsMessage = {
  type: string,
  content: any,
}
