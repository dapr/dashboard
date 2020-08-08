import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { DaprConfiguration, DaprConfigurationStatus } from '../types/types';
import { Observable } from 'rxjs';
import { ScopesService } from '../scopes/scopes.service';

@Injectable({
  providedIn: 'root'
})
export class ConfigurationsService {

  constructor(
    private http: HttpClient,
    private scopesService: ScopesService,
  ) { }

  getConfigurations(): Observable<DaprConfiguration[]> {
    const scope = this.scopesService.getScope();
    return this.http.get<DaprConfiguration[]>(`/api/configurations${scope}`);
  }

  getConfiguration(name: string): Observable<DaprConfiguration> {
    const scope = this.scopesService.getScope();
    return this.http.get<DaprConfiguration>(`/api/configurations/${scope}/${name}`);
  }

  getConfigurationsStatus(): Observable<DaprConfigurationStatus[]> {
    const scope = this.scopesService.getScope();
    return this.http.get<DaprConfigurationStatus[]>(`/api/configurationsstatus/${scope}`);
  }
}
