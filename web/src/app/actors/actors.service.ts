import { Injectable, Output } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Injectable({
  providedIn: 'root',
})
export class ActorsService {

  constructor(private http: HttpClient) { }

  getActors(port: string) {
    const output = [];
    this.http.get("/api/actors", { responseType: 'text' }).subscribe((actorData: string) => {
      var actorArray: string = actorData["actors"];
      actorArray.split('\n').forEach(actor => {
          output.push(actor);
      });
    });
    return output;
  }
}
