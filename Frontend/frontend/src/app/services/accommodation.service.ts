import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Accommodation, Availability, Price, Reservation } from '../models/Accommodation'; 

@Injectable({
  providedIn: 'root'
})
export class AccommodationService {

  private apiUrl = 'http://localhost:8080/accommodations';
  private apiUrl2 = 'http://localhost:8081/reservations';


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

  // Nova metoda za ažuriranje dostupnosti i cene
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
    return this.http.post<Reservation>(this.apiUrl2, reservation); // Izmenjeni URL sa ispravnim portom
  }
  
}