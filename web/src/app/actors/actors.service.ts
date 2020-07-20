import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Actor } from './Actor';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

@Injectable({
  providedIn: 'root',
})
export class ActorsService {

  constructor(private http: HttpClient) { }

  getActors(id: string): Observable<Actor[]> {
    const output: Actor[] = [];
    return this.http.get(`/api/actors/${id}`, { responseType: 'text' }).pipe(
      map(actorData => {
            let output: Actor[] = [];
            actorData.split('\n').forEach(actor => {
              let actorJSON: Actor;
              try {
                actorJSON = JSON.parse(actor);
                output.push(actorJSON);
              } catch (e) {}
            });
            return output;
      })
    );
  }
}
