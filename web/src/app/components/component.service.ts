import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Injectable({
  providedIn: 'root',
})
export class ComponentsService {

  constructor(private http: HttpClient) { }

  getComponents() {
    return this.http.get('/api/components');
  }

  getComponentsStatus() {
    return this.http.get('/api/components/status');
  }
}
