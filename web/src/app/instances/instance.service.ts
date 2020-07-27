import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Metadata, Log, Instance, Status } from 'src/app/types/types';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

@Injectable({
  providedIn: 'root',
})
export class InstanceService {

  constructor(private http: HttpClient) { }

  getInstances(): Observable<Instance[]> {
    return this.http.get<Instance[]>('/api/instances');
  }

  getInstance(id: string): Observable<Instance> {
    return this.http.get<Instance>('/api/instances/' + id);
  }

  deleteInstance(id: string): Observable<Instance[]> {
    return this.http.delete<Instance[]>('/api/instances/' + id);
  }

  getDeploymentConfiguration(id: string): Observable<string> {
    return this.http.get('/api/deploymentconfiguration/' + id, { responseType: 'text' });
  }

  getMetadata(id: string): Observable<Metadata[]> {
    return this.http.get<Metadata[]>('/api/metadata/' + id);
  }

  getControlPlaneStatus(): Observable<Status[]> {
    return this.http.get<Status[]>('/api/controlplanestatus');
  }

  getLogs(id: string): Observable<Log[]> {
    return this.http.get<Log[]>('/api/instances/' + id + '/logs');
  }
}
