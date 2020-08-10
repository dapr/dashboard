import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { DaprComponent } from '../types/types';
import { ScopesService } from '../scopes/scopes.service';

@Injectable({
  providedIn: 'root',
})
export class ComponentsService {

  constructor(
    private http: HttpClient,
    private scopesService: ScopesService
  ) { }

  getComponents(): Observable<DaprComponent[]> {
    const scope = this.scopesService.getScope();
    return this.http.get<DaprComponent[]>(`/api/components/${scope}`);
  }

  getComponent(name: string): Observable<DaprComponent> {
    const scope = this.scopesService.getScope();
    return this.http.get<DaprComponent>(`/api/components/${scope}/${name}`);
  }
}
