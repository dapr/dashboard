import { Injectable, EventEmitter } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class ScopesService {

  private scope = 'All';
  public scopeChanged: EventEmitter<string> = new EventEmitter();

  constructor(private http: HttpClient) { }

  getScopes(): Observable<string[]> {
    return this.http.get<string[]>('api/scopes');
  }

  getScope(): string {
    return this.scope;
  }

  changeScope(newScope: string): void {
    this.scope = newScope;
    this.scopeChanged.emit(this.scope);
  }
}
