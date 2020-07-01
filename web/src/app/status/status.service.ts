import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Injectable({
  providedIn: 'root'
})
export class StatusService {

  constructor(private http: HttpClient) { }

  getControlPlaneStatus() {
    return this.http.get('/api/controlplanestatus');
  }
}
