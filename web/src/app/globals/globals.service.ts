import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Injectable({
  providedIn: 'root'
})
export class GlobalsService {

  public kubernetesEnabled = false;
  public standaloneEnabled = false;

  constructor(private http: HttpClient) { }

  getSupportedEnvironments() {
    return this.http.get('/api/environments');
  }
}