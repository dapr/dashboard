import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Subject, AnonymousSubject } from 'rxjs/internal/Subject';

@Injectable({
  providedIn: 'root',
})
export class WebsocketService {
  constructor() {
    this.subjects = new Map();
  }

  private subjects: Map<string, Subject<MessageEvent>>;

  public connect(url): Subject<MessageEvent> {
    this.disconnect(url);
    const origin = window.location.origin.replace('https://', 'wss://').replace('http://', 'ws://');
    const subject = this.create(origin + url);
    this.subjects.set(url, subject);
    console.log('Successfully connected: ' + url);
    return subject;
  }

  public disconnect(url): void {
    if (this.subjects.has(url)) {
      const oldSubject = this.subjects.get(url);
      oldSubject.complete();
      this.subjects.delete(url);
      console.log('Successfully disconnected: ' + url);
    }
  }

  private create(url): Subject<MessageEvent> {
    const ws = new WebSocket(url);

    const observable = new Observable<MessageEvent>((obs) => {
      ws.onmessage = obs.next.bind(obs);
      ws.onerror = obs.error.bind(obs);
      ws.onclose = obs.complete.bind(obs);
      return ws.close.bind(ws);
    });
    const observer = {
      next: (data: any) => {
        if (ws.readyState === WebSocket.OPEN) {
          ws.send(JSON.stringify(data));
        }
      },
      error: (err: any) => {},
      complete: () => {}
    };
    return new AnonymousSubject<MessageEvent>(observer, observable);
  }
}
