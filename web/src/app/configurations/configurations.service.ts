import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { DaprConfiguration } from '../types/types';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class ConfigurationsService {

  constructor(private http: HttpClient) { }

  getConfigurations(): Observable<DaprConfiguration[]> {
    return this.http.get<DaprConfiguration[]>('/api/configurations');
  }

  getConfiguration(name: string): Observable<DaprConfiguration> {
    return this.http.get<DaprConfiguration>('/api/configurations/' + name);
  }
}
