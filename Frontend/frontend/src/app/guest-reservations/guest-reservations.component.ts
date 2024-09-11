import { Component, OnInit } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Router } from '@angular/router';
import { jwtDecode } from 'jwt-decode';
import { Reservation, Rating } from '../models/Accommodation';

@Component({
  selector: 'app-guest-reservations',
  templateUrl: './guest-reservations.component.html',
  styleUrls: ['./guest-reservations.component.css']
})
export class GuestReservationsComponent implements OnInit {
  reservations: Reservation[] = [];
  showHostRatingForm = false;
  showAccommodationRatingForm = false;
  showEditRatingForm = false;
  selectedReservationId?: string;
  hostRating?: number | null = null;
  hostComment = '';
  accommodationRating?: number | null = null;
  accommodationComment = '';

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

  // cancelReservation(reservationId?: string): void {
  //   if (!reservationId || reservationId.trim() === '') {
  //     console.error('Reservation ID is undefined or empty.');
  //     alert('Invalid Reservation ID.');
  //     return;
  //   }

  //   const token = localStorage.getItem('token');

  //   if (token) {
  //     const url = `http://localhost:8081/reservations/${reservationId}`;
  //     const headers = new HttpHeaders({
  //       'Authorization': `Bearer ${token}`,
  //       'Content-Type': 'application/json'
  //     });

  //     this.http.delete(url, { headers })
  //       .subscribe({
  //         next: () => {
  //           this.reservations = this.reservations.filter(reservation => reservation.id !== reservationId);
  //           alert('Reservation canceled successfully.');
  //         },
  //         error: (error: any) => {
  //           console.error('Error canceling reservation:', error);
  //           alert('Error occurred while canceling the reservation. Please try again later.');
  //         }
  //       });
  //   } else {
  //     console.error('No token found, redirecting to login.');
  //     this.router.navigate(['/login']);
  //   }
  // }

  cancelReservation(reservationId?: string): void {
    if (!reservationId || reservationId.trim() === '') {
      console.error('Reservation ID is undefined or empty.');
      alert('Invalid Reservation ID.');
      return;
    }

    const token = localStorage.getItem('token');

    if (token) {
      const url = `http://localhost:8081/reservations/${reservationId}`;
      const headers = new HttpHeaders({
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      });

      this.http.delete(url, { headers })
        .subscribe({
          next: () => {
            this.reservations = this.reservations.filter(reservation => reservation.id !== reservationId);
            alert('Reservation canceled successfully.');
            this.sendNotificationAfterCancelation(reservationId);
          },
          error: (error: any) => {
            console.error('Error canceling reservation:', error);
            alert('Error occurred while canceling the reservation. Please try again later.');
          }
        });
    } else {
      console.error('No token found, redirecting to login.');
      this.router.navigate(['/login']);
    }
  }

  sendNotificationAfterCancelation(reservationId: string): void {
    const token = localStorage.getItem('token');
    if (token) {
      const decodedToken = this.decodeToken(token);
      const message = `Reservation with ID ${reservationId} has been canceled by guest ${decodedToken.userID}.`;

      const notificationData = {
        host_id: decodedToken.userID, 
        message: message
      };

      const headers = new HttpHeaders({
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      });

      this.http.post('http://localhost:8083/notifications', notificationData, { headers })
        .subscribe({
          next: (response) => {
            console.log('Cancellation notification sent successfully', response);
          },
          error: (error) => {
            console.error('Error sending cancellation notification:', error);
          }
        });
    } else {
      console.error('No token found, redirecting to login.');
      this.router.navigate(['/login']);
    }
  }


  toggleHostRatingForm(reservationId: string | undefined): void {
    if (!reservationId) {
      console.error('Reservation ID is undefined.');
      return;
    }
    this.showHostRatingForm = !this.showHostRatingForm;
    this.showAccommodationRatingForm = false;
    this.showEditRatingForm = false;
    this.selectedReservationId = reservationId;
    this.loadHostRating(reservationId); 
  }

  toggleAccommodationRatingForm(reservationId: string | undefined): void {
    if (!reservationId) {
      console.error('Reservation ID is undefined.');
      return;
    }
    this.showAccommodationRatingForm = !this.showAccommodationRatingForm;
    this.showHostRatingForm = false;
    this.showEditRatingForm = false;
    this.selectedReservationId = reservationId;
    this.loadAccommodationRating(reservationId); 
  }

  toggleEditRatingForm(reservationId: string | undefined): void {
    if (!reservationId) {
      console.error('Reservation ID is undefined.');
      return;
    }
    this.showEditRatingForm = !this.showEditRatingForm;
    this.showHostRatingForm = false;
    this.showAccommodationRatingForm = false;
    this.selectedReservationId = reservationId;
    this.loadHostRating(reservationId);
    this.loadAccommodationRating(reservationId);
  }

  loadHostRating(reservationId: string): void {
    const token = localStorage.getItem('token');
    if (token) {
      const headers = new HttpHeaders({
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      });

      this.http.get<Rating[]>(`http://localhost:8082/hosts/${reservationId}/ratings`, { headers })
        .subscribe({
          next: (ratings: Rating[]) => {
            console.log('Host ratings loaded:', ratings);
            if (ratings.length > 0) {
              this.hostRating = ratings[0].rating;
              this.hostComment = ratings[0].comment;
              console.log('Assigned hostRating:', this.hostRating);
              console.log('Assigned hostComment:', this.hostComment);
            } else {
              console.log('No host ratings found');
            }
          },
          error: (error: any) => {
            console.error('Error loading host rating:', error);
          }
        });
    }
  }

  loadAccommodationRating(reservationId: string): void {
    const token = localStorage.getItem('token');
    if (token) {
      const headers = new HttpHeaders({
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      });

      this.http.get<Rating[]>(`http://localhost:8082/accommodations/${reservationId}/ratings`, { headers })
        .subscribe({
          next: (ratings: Rating[]) => {
            console.log('Accommodation ratings loaded:', ratings);
            if (ratings.length > 0) {
              this.accommodationRating = ratings[0].rating;
              this.accommodationComment = ratings[0].comment;
              console.log('Assigned accommodationRating:', this.accommodationRating);
              console.log('Assigned accommodationComment:', this.accommodationComment);
            } else {
              console.log('No accommodation ratings found');
            }
          },
          error: (error: any) => {
            console.error('Error loading accommodation rating:', error);
          }
        });
    }
  }

  // submitHostRating(reservationId: string): void {
  //   const token = localStorage.getItem('token');
  //   if (token) {
  //     const decodedToken = this.decodeToken(token);
  //     const ratingData: Rating = {
  //       user_id: decodedToken.userID,
  //       targetID: reservationId,
  //       rating: this.hostRating!,
  //       comment: this.hostComment
  //     };

  //     // Log za proveru userID i ratingData
  //     console.log('Decoded userID for host rating:', decodedToken.userID);
  //     console.log('Rating data being sent for host:', ratingData);

  //     const headers = new HttpHeaders({
  //       'Authorization': `Bearer ${token}`,
  //       'Content-Type': 'application/json'
  //     });

  //     this.http.post(`http://localhost:8082/hosts/${reservationId}/rate`, ratingData, { headers })
  //       .subscribe({
  //         next: () => {
  //           alert('Host rating submitted successfully.');
  //           this.showHostRatingForm = false;
  //         },
  //         error: (error: any) => {
  //           console.error('Error submitting host rating:', error);
  //           alert('Error occurred while submitting the rating. Please try again later.');
  //         }
  //       });
  //   } else {
  //     console.error('No token found, redirecting to login.');
  //     this.router.navigate(['/login']);
  //   }
  // }

  sendNotification(message: string): void {
    const token = localStorage.getItem('token');
    if (token) {
      const decodedToken = this.decodeToken(token);
      const notificationData = {
        host_id: decodedToken.userID, 
        message: message
      };

      const headers = new HttpHeaders({
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      });

      this.http.post('http://localhost:8083/notifications', notificationData, { headers })
        .subscribe({
          next: (response) => {
            console.log('Notification sent successfully', response);
          },
          error: (error) => {
            console.error('Error sending notification:', error);
          }
        });
    } else {
      console.error('No token found, redirecting to login.');
      this.router.navigate(['/login']);
    }
  }


  submitHostRating(reservationId: string): void {
    const token = localStorage.getItem('token');
    if (token) {
      const decodedToken = this.decodeToken(token);
      const ratingData: Rating = {
        user_id: decodedToken.userID,
        targetID: reservationId,
        rating: this.hostRating!,
        comment: this.hostComment
      };


      console.log('Decoded userID for host rating:', decodedToken.userID);
      console.log('Rating data being sent for host:', ratingData);

      const headers = new HttpHeaders({
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      });

      this.http.post(`http://localhost:8082/hosts/${reservationId}/rate`, ratingData, { headers })
        .subscribe({
          next: () => {
            alert('Host rating submitted successfully.');
            this.showHostRatingForm = false;
        
            this.sendNotification(`Rate Host for Reservation ID: ${reservationId}`);
          },
          error: (error: any) => {
            console.error('Error submitting host rating:', error);
            alert('Error occurred while submitting the rating. Please try again later.');
          }
        });
    } else {
      console.error('No token found, redirecting to login.');
      this.router.navigate(['/login']);
    }
  }


  // submitAccommodationRating(reservationId: string): void {
  //   const token = localStorage.getItem('token');
  //   if (token) {
  //     const decodedToken = this.decodeToken(token);
  //     const ratingData: Rating = {
  //       user_id: decodedToken.userID,
  //       targetID: reservationId,
  //       rating: this.accommodationRating!,
  //       comment: this.accommodationComment
  //     };

  //     // Log za proveru userID i ratingData
  //     console.log('Decoded userID for accommodation rating:', decodedToken.userID);
  //     console.log('Rating data being sent for accommodation:', ratingData);

  //     const headers = new HttpHeaders({
  //       'Authorization': `Bearer ${token}`,
  //       'Content-Type': 'application/json'
  //     });

  //     this.http.post(`http://localhost:8082/accommodations/${reservationId}/rate`, ratingData, { headers })
  //       .subscribe({
  //         next: () => {
  //           alert('Accommodation rating submitted successfully.');
  //           this.showAccommodationRatingForm = false;
  //         },
  //         error: (error: any) => {
  //           console.error('Error submitting accommodation rating:', error);
  //           alert('Error occurred while submitting the rating. Please try again later.');
  //         }
  //       });
  //   } else {
  //     console.error('No token found, redirecting to login.');
  //     this.router.navigate(['/login']);
  //   }
  // }

  submitAccommodationRating(reservationId: string): void {
    const token = localStorage.getItem('token');
    if (token) {
      const decodedToken = this.decodeToken(token);
      const ratingData: Rating = {
        user_id: decodedToken.userID,
        targetID: reservationId,
        rating: this.accommodationRating!,
        comment: this.accommodationComment
      };


      console.log('Decoded userID for accommodation rating:', decodedToken.userID);
      console.log('Rating data being sent for accommodation:', ratingData);

      const headers = new HttpHeaders({
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      });

      this.http.post(`http://localhost:8082/accommodations/${reservationId}/rate`, ratingData, { headers })
        .subscribe({
          next: () => {
            alert('Accommodation rating submitted successfully.');
            this.showAccommodationRatingForm = false;
        
            this.sendNotification(`Rate Accommodation for Reservation ID: ${reservationId}`);
          },
          error: (error: any) => {
            console.error('Error submitting accommodation rating:', error);
            alert('Error occurred while submitting the rating. Please try again later.');
          }
        });
    } else {
      console.error('No token found, redirecting to login.');
      this.router.navigate(['/login']);
    }
  }


  deleteHostRating(reservationId: string): void {
    const token = localStorage.getItem('token');
    if (token) {
      const headers = new HttpHeaders({
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      });

      this.http.delete(`http://localhost:8082/hosts/${reservationId}/ratings`, { headers })
        .subscribe({
          next: () => {
            alert('Host rating deleted successfully.');
            this.hostRating = null;
            this.hostComment = '';
          },
          error: (error: any) => {
            console.error('Error deleting host rating:', error);
            alert('Error occurred while deleting the rating. Please try again later.');
          }
        });
    } else {
      console.error('No token found, redirecting to login.');
      this.router.navigate(['/login']);
    }
  }

  deleteAccommodationRating(reservationId: string): void {
    const token = localStorage.getItem('token');
    if (token) {
      const headers = new HttpHeaders({
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      });

      this.http.delete(`http://localhost:8082/accommodations/${reservationId}/ratings`, { headers })
        .subscribe({
          next: () => {
            alert('Accommodation rating deleted successfully.');
            this.accommodationRating = null;
            this.accommodationComment = '';
          },
          error: (error: any) => {
            console.error('Error deleting accommodation rating:', error);
            alert('Error occurred while deleting the rating. Please try again later.');
          }
        });
    } else {
      console.error('No token found, redirecting to login.');
      this.router.navigate(['/login']);
    }
  }

  submitEditRating(reservationId: string): void {
    this.submitHostRating(reservationId);
    this.submitAccommodationRating(reservationId);
  }
}
