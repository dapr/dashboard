import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Injectable({
  providedIn: 'root'
})
export class InstanceService {

  constructor(private http: HttpClient) { }

  getInstances() {
    return this.http.get("/api/instances");
  }

  deleteInstance(id: string) {
    return this.http.delete("/api/instances/" + id)
  }

  getLogs(id: string) {
    return this.http.get("/api/instances/" + id + "/logs", {responseType: 'text'})
  }
}
