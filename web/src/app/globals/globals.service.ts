import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class GlobalsService {

  public changesEnabledRoutes: string[] = ['/overview', '/components', '/configurations', '/controlplane'];

  constructor(private http: HttpClient) { }

  getPlatform(): Observable<string> {
    return this.http.get('/api/platform', { responseType: 'text' });
  }

  getVersion(): Observable<string> {
    return this.http.get('/api/version', { responseType: 'text' });
  }
}
