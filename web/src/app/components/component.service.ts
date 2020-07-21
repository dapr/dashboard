import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { DaprComponentStatus, DaprComponent } from '../types/types';

@Injectable({
  providedIn: 'root',
})
export class ComponentsService {

  constructor(private http: HttpClient) { }

  getComponents(): Observable<DaprComponent[]> {
    return this.http.get<DaprComponent[]>('/api/components');
  }

  getComponentsStatus(): Observable<DaprComponentStatus[]> {
    return this.http.get<DaprComponentStatus[]>('/api/components/status');
  }
}
