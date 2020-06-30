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
    this.http.get('/api/environments').subscribe(data => {
      let supportedEnvironments = <Array<any>>data;
      if (supportedEnvironments.includes("kubernetes")) {
        this.kubernetesEnabled = true;
      }
      else if (supportedEnvironments.includes("standalone")) {
        this.standaloneEnabled = true;
      }
    });
  }
}
