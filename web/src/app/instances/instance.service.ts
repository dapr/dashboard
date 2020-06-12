import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Log } from '../pages/logs/log';

@Injectable({
  providedIn: 'root',
})
export class InstanceService {

  constructor(private http: HttpClient) { }

  getInstances() {
    return this.http.get('/api/instances');
  }

  deleteInstance(id: string) {
    return this.http.delete('/api/instances/' + id);
  }

  getLogs(id: string) {
    return this.http.get('/api/instances/' + id + '/logs', {responseType: 'text'});
  }

  getYAML(id: string) {
    return this.http.get('/api/yaml/' + id, {responseType: 'text'});
  }

  getLogsArray(id: string): Log[] {
    const output = [];
    this.http.get('/api/instances/' + id + '/logs', {responseType: 'text'}).subscribe((logData: string) => {
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

  getYAMLArray(id: string): string[] {
    const output = [];
    this.http.get('/api/yaml/' + id, {responseType: 'text'}).subscribe((yamlData: string) => {
      yamlData.split('\n').forEach(line => {
        output.push(line);
      });
    });
    return output;
  }
}
