import { Injectable } from "@angular/core";
import { Observable } from "rxjs";
import { ScopesService } from '../scopes/scopes.service';
import { WebsocketService } from "../websocket/websocket.service";
import { Log } from 'src/app/types/types';
import { map } from 'rxjs/operators';

@Injectable({
  providedIn: 'root',
})
export class LogStreamService {

  constructor(
    private wsService: WebsocketService,
    private scopesService: ScopesService
  ) { }

  startStream(id: string, container: string): Observable<Log> {
    const scope = this.scopesService.getScope();
    const url = `/api/instances/${scope}/${id}/logstreams/${container}`;
    return this.wsService.connect(url).pipe(map(
      (response: MessageEvent): Log => {
        return JSON.parse(response.data);
      }
    ));
  }

  endStream(id: string, container: string): void {
    const scope = this.scopesService.getScope();
    const url = `/api/instances/${scope}/${id}/logstreams/${container}`;
    this.wsService.disconnect(url);
  }
}