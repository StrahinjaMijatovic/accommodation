import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class ProfileService {
  private baseUrl = 'http://localhost:8001/profile'; // Prilagoditi URL-u backend-a

  constructor(private http: HttpClient) {}

  checkActiveReservations(userId: string): Observable<boolean> {
    return this.http.get<boolean>(`${this.baseUrl}/check-reservations?user_id=${userId}`);
  }

  deleteProfile(email: string): Observable<any> {
    return this.http.delete(`${this.baseUrl}/delete?email=${email}`);
  }
}
