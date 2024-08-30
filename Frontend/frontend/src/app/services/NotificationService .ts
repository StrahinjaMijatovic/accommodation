import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class NotificationService {

  private apiUrl = 'http://localhost:8082/notifications';

  constructor(private http: HttpClient) { }

  getNotifications(userID: string): Observable<Notification[]> {
    return this.http.get<Notification[]>(`${this.apiUrl}/${userID}`);
  }
}
