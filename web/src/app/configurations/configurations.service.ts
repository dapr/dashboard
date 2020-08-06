import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { DaprConfiguration, DaprConfigurationStatus, Instance } from '../types/types';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

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

  getConfigurationsStatus(): Observable<DaprConfigurationStatus[]> {
    return this.http.get<DaprConfigurationStatus[]>('/api/configurationsstatus');
  }

  getConfigurationApps(name: string): Observable<Instance[]> {
    return this.http.get<Instance[]>('/api/instances').pipe(
      map(instances => {
        return instances.filter(instance => instance.config === name);
      })
    );
  }
}
