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

  getInstances() {
    return this.http.get('/api/instances');
  }

  getInstance(id: string): Observable<Instance> {
    return this.http.get<Instance>('/api/instances/' + id);
  }

  deleteInstance(id: string): Observable<Instance[]> {
    return this.http.delete<Instance[]>('/api/instances/' + id);
  }

  getConfiguration(id: string): Observable<string> {
    return this.http.get('/api/configuration/' + id, { responseType: 'text' });
  }

  getMetadata(id: string): Observable<Metadata[]> {
    return this.http.get<Metadata[]>('/api/metadata/' + id);
  }

  getControlPlaneStatus(): Observable<Status[]> {
    return this.http.get<Status[]>('/api/controlplanestatus');
  }

  getLogs(id: string): Observable<Log[]> {
    return this.http.get('/api/instances/' + id + '/logs', { responseType: 'text' }).pipe(
      map(logData => {
        const output = [];
        logData.split('\n').forEach(log => {
          const regEx = RegExp('(?<=level=).*?(?=\s)', '');
          const level: string[] = regEx.exec(log);
          if (level != null && level.length > 0) {
            const currentLog: Log = {
              level: level[0].replace(' m', ''),
              log,
            };
            output.push(currentLog);
          }
        });
        return output;
      })
    );
  }
}
