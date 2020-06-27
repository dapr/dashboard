import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Log } from '../pages/detail/logs/log';

@Injectable({
  providedIn: 'root',
})
export class InstanceService {

  constructor(private http: HttpClient) { }

  getInstances() {
    return this.http.get('/api/instances');
  }

  getInstance(id: string) {
    return this.http.get('/api/instances/' + id);
  }

  deleteInstance(id: string) {
    return this.http.delete('/api/instances/' + id);
  }

  getLogs(id: string) {
    return this.http.get('/api/instances/' + id + '/logs', { responseType: 'text' });
  }

  getConfiguration(id: string) {
    return this.http.get('/api/configuration/' + id, { responseType: 'text' });
  }

  getLogsArray(id: string): Log[] {
    const output = [];
    this.http.get('/api/instances/' + id + '/logs', { responseType: 'text' }).subscribe((logData: string) => {
      logData.split('\n').forEach(log => {
        const regEx = RegExp('(?<=level=).*?(?=\s)', '');
        const level: string[] = regEx.exec(log);
        if (level != null && level.length > 0) {
          const currentLog: Log = {
            level: level[0].replace(' m', ''),
            log: log,
          };
          output.push(currentLog);
        }
      });
    });
    return output;
  }
}
