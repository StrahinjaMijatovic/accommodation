import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Accommodation, Availability, Price, Reservation, Rating } from '../models/Accommodation'; 

@Injectable({
  providedIn: 'root'
})
export class AccommodationService {

  private apiUrl = 'http://localhost:8080/accommodations';
  private apiUrl2 = 'http://localhost:8081/reservations';
  private hostRatingUrl = 'http://localhost:8082/hosts';  
  private accommodationRatingUrl = 'http://localhost:8082/accommodations'; 
  private notificationUrl = 'http://localhost:8083/notifications'; 


  constructor(private http: HttpClient) { }

  getAccommodations(): Observable<Accommodation[]> {
    return this.http.get<Accommodation[]>(this.apiUrl);
  }
  
  createAccommodation(accommodation: Accommodation): Observable<Accommodation> {
    return this.http.post<Accommodation>(this.apiUrl, accommodation);
  }

  getAccommodationById(id: string): Observable<Accommodation> {
    return this.http.get<Accommodation>(`${this.apiUrl}/${id}`);
  }

  updatePrice(id: string, price: { startDate: string; endDate: string; amount: number; strategy: string }): Observable<Accommodation> {
    return this.http.put<Accommodation>(`${this.apiUrl}/${id}/price`, price);
  }

  updateAccommodation(id: string, accommodation: Accommodation): Observable<Accommodation> {
    return this.http.put<Accommodation>(`${this.apiUrl}/${id}`, accommodation);
  }

  updateAvailabilityAndPrice(id: string, data: { startDate: string; endDate: string; amount: number; strategy: string }): Observable<Accommodation> {
    return this.http.put<Accommodation>(`${this.apiUrl}/${id}/availability-and-price`, data);
  }

  getAvailabilityByAccommodationId(id: string): Observable<Availability[]> {
    return this.http.get<Availability[]>(`${this.apiUrl}/${id}/availability`);
  }

  getPriceByAccommodationId(id: string): Observable<Price[]> {
    return this.http.get<Price[]>(`${this.apiUrl}/${id}/prices`);
  }
  
  reserveAccommodation(reservation: Reservation): Observable<Reservation> {
    return this.http.post<Reservation>(this.apiUrl2, reservation);
  }
  
 
  rateHost(reservationId: string, rating: Rating): Observable<Rating> {
    return this.http.post<Rating>(`${this.hostRatingUrl}/${reservationId}/rate`, rating);
  }

  
  rateAccommodation(reservationId: string, rating: Rating): Observable<Rating> {
    return this.http.post<Rating>(`${this.accommodationRatingUrl}/${reservationId}/rate`, rating);
  }

  
  deleteHostRating(reservationId: string): Observable<void> {
    return this.http.delete<void>(`${this.hostRatingUrl}/${reservationId}/ratings`);
  }

  
  deleteAccommodationRating(reservationId: string): Observable<void> {
    return this.http.delete<void>(`${this.accommodationRatingUrl}/${reservationId}/ratings`);
  }

sendNotification(hostId: number, message: string): Observable<void> {
  const notificationData = { host_id: hostId, message: message };
  return this.http.post<void>(this.notificationUrl, notificationData);
}

getMyAccommodations(ajdi: string): Observable<Accommodation[]> {
  return this.http.get<Accommodation[]>(`${this.apiUrl}/my-accommodations/${ajdi}`);
}

getNotificationsForUser(userID: string): Observable<Notification[]> {
  const token = localStorage.getItem('token');
  if (token) {
    const headers = new HttpHeaders({
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    });

    return this.http.get<Notification[]>(`${this.notificationUrl}/${userID}`, { headers });
  } else {
    console.error('No token found, redirecting to login.');
    throw new Error('No token found');
  }
}


}
