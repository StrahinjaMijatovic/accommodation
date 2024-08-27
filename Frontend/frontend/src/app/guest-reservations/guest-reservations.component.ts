import { Component, OnInit } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Router } from '@angular/router';
import { jwtDecode } from 'jwt-decode';
import {  Reservation } from '../models/Accommodation';


@Component({
  selector: 'app-guest-reservations',
  templateUrl: './guest-reservations.component.html',
  styleUrls: ['./guest-reservations.component.css']
})
export class GuestReservationsComponent implements OnInit {
  reservations: Reservation[] = [];

  constructor(private http: HttpClient, private router: Router) {}

  ngOnInit(): void {
    this.loadReservations();
  }

  decodeToken(token: string): any {
    try {
      return jwtDecode(token);
    } catch (error) {
      console.error('Invalid token:', error);
      return null;
    }
  }

  loadReservations(): void {
    const token = localStorage.getItem('token');
    if (token) {
      const decodedToken = this.decodeToken(token);
      if (decodedToken && decodedToken.userID) {
        const headers = new HttpHeaders({
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        });

        this.http.get<Reservation[]>(`http://localhost:8081/guests/${decodedToken.userID}/reservations`, { headers })
          .subscribe({
            next: (reservations: Reservation[]) => {
              this.reservations = reservations;
            },
            error: (error: any) => {
              console.error('Error loading reservations:', error);
              alert('Error occurred while loading reservations. Please try again later.');
            }
          });
      } else {
        console.error('Invalid or missing user ID in token.');
        this.router.navigate(['/login']);
      }
    } else {
      console.error('No token found, redirecting to login.');
      this.router.navigate(['/login']);
    }
  }

  cancelReservation(reservationId?: string): void {
    console.log('Attempting to cancel reservation with ID:', reservationId); // Logovanje reservationId pre slanja zahteva
  
    if (!reservationId || reservationId.trim() === '') {
      console.error('Reservation ID is undefined or empty.');
      alert('Invalid Reservation ID.');
      return;
    }
  
    const token = localStorage.getItem('token');
    console.log('Token:', token); // Logovanje tokena pre slanja zahteva
  
    if (token) {
      const url = `http://localhost:8081/reservations/${reservationId}`;
      console.log('URL being called:', url); // Logovanje URL-a pre slanja zahteva
  
      const headers = new HttpHeaders({
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      });
  
      console.log('Headers being sent:', headers); // Logovanje header-a pre slanja zahteva
  
      this.http.delete(url, { headers })
        .subscribe({
          next: () => {
            console.log('Reservation canceled successfully'); // Logovanje uspeha
            this.reservations = this.reservations.filter(reservation => reservation.id !== reservationId);
            alert('Reservation canceled successfully.');
          },
          error: (error: any) => {
            console.error('Error canceling reservation:', error); // Logovanje greške
            alert('Error occurred while canceling the reservation. Please try again later.');
          }
        });
    } else {
      console.error('No token found, redirecting to login.');
      this.router.navigate(['/login']);
    }
  }
  
  
}

