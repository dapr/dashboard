import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Metadata, Log, Instance, Status } from 'src/app/types/types';
import { Observable } from 'rxjs';
import { ScopesService } from '../scopes/scopes.service';

@Injectable({
  providedIn: 'root',
})
export class InstanceService {

  constructor(
    private http: HttpClient,
    private scopesService: ScopesService
  ) { }

  getInstances(): Observable<Instance[]> {
    const scope = this.scopesService.getScope();
    return this.http.get<Instance[]>(`/api/instances/${scope}`);
  }

  getInstance(id: string): Observable<Instance> {
    const scope = this.scopesService.getScope();
    return this.http.get<Instance>(`/api/instances/${scope}/${id}`);
  }

  deleteInstance(id: string): Observable<Instance[]> {
    const scope = this.scopesService.getScope();
    return this.http.delete<Instance[]>(`/api/instances/${scope}/${id}`);
  }

  getDeploymentConfiguration(id: string): Observable<string> {
    const scope = this.scopesService.getScope();
    return this.http.get(`/api/deploymentconfiguration/${scope}/${id}`, { responseType: 'text' });
  }

  getMetadata(id: string): Observable<Metadata[]> {
    const scope = this.scopesService.getScope();
    return this.http.get<Metadata[]>(`/api/metadata/${scope}/${id}`);
  }

  getControlPlaneStatus(): Observable<Status[]> {
    return this.http.get<Status[]>(`/api/controlplanestatus`);
  }

  getLogs(id: string): Observable<Log[]> {
    const scope = this.scopesService.getScope();
    return this.http.get<Log[]>(`/api/instances/${scope}/${id}/logs`);
  }
}
