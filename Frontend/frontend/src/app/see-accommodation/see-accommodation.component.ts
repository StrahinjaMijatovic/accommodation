import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { AccommodationService } from '../services/accommodation.service';
import { Accommodation, Availability, Price, Reservation } from '../models/Accommodation';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { jwtDecode } from 'jwt-decode';

interface UserResponse {
  firstName: string;
  lastName: string;
  role: string;
}

@Component({
  selector: 'app-see-accommodation',
  templateUrl: './see-accommodation.component.html',
  styleUrls: ['./see-accommodation.component.css']
})
export class SeeAccommodationComponent implements OnInit {
  accommodation: Accommodation | undefined;
  showUpdateForm: boolean = false;
  showAvailabilityForm: boolean = false;
  showAvailability: boolean = false;
  showPrice: boolean = false;
  showReservationForm: boolean = false;
  availabilityList: Availability[] = [];
  priceList: Price[] = [];
  role: string = '';

  startDate: string = '';
  endDate: string = '';
  amount: number = 0;
  strategy: 'per_guest' | 'per_unit' = 'per_unit';
  userIdd: string | undefined = '';

  reservationStartDate: string = '';
  reservationEndDate: string = '';

  constructor(
    private http: HttpClient,
    private route: ActivatedRoute,
    private accommodationService: AccommodationService,
    private router: Router
  ) {
    // Dobijanje ID-ja iz URL-a
    const id = this.route.snapshot.paramMap.get('id');
    if (id) {
      // Poziv metode getAccommodationById sa ID-jem
      this.accommodationService.getAccommodationById(id).subscribe((accommodation: Accommodation) => {
        if (accommodation && accommodation.userID) {
          this.userIdd = accommodation.userID;
          console.log('Accommodation:', accommodation);
          console.log('User ID:', this.userIdd);
        } else {
          console.error('Accommodation or userID is missing.');
        }
      },
      error => {
        console.error('Error fetching accommodation by ID:', error);
      });
    } else {
      console.error('Accommodation ID is missing in the route.');
      this.router.navigate(['/']);
    }
  }


  checkUserStatus() {
    const token = localStorage.getItem('token');
    if (token) {
      this.http.post<UserResponse>('http://localhost:8000/verify-token', { token })
        .subscribe((response: UserResponse) => {
          this.role = response.role;
        }, error => {
          console.error('Provera tokena nije uspela:', error);
          
        });
    }
  }

  ngOnInit(): void {
    this.checkUserStatus();
    const id = this.route.snapshot.paramMap.get('id');
    if (id) {
      this.accommodationService.getAccommodationById(id).subscribe(
        data => {
          this.accommodation = data;
          this.getAvailability(id);
          this.getPrice(id);
          
          // Provera da li postoji userID nakon što su podaci učitani
          console.log('Accommodation object:', this.accommodation);
          console.log('User ID:', this.accommodation?.userID);
          
          if (!this.accommodation?.userID) {
            console.error('User ID is missing in accommodation data.');
          }
        },
        error => {
          console.error('Error fetching accommodation details:', error);
        }
      );
    }
  }

  decodeToken(token: string): any {
    try {
      const decoded = jwtDecode(token);
      console.log('Decoded token:', decoded);
      return decoded;
    } catch (error) {
      console.error('Invalid token:', error);
      return null;
    }
  }

  toggleUpdateForm() {
    this.showUpdateForm = !this.showUpdateForm;
  }

  toggleAvailabilityForm() {
    this.showAvailabilityForm = !this.showAvailabilityForm;
  }

  toggleAvailability() {
    if (!this.showAvailability) {
      this.getAvailability(this.accommodation?.id || '');
    }
    this.showAvailability = !this.showAvailability;
  }

  togglePrice() {
    if (!this.showPrice) {
      this.getPrice(this.accommodation?.id || '');
    }
    this.showPrice = !this.showPrice;
  }

  toggleReservationForm() {
    this.showReservationForm = !this.showReservationForm;
  }

  getAvailability(accommodationId: string) {
    this.accommodationService.getAvailabilityByAccommodationId(accommodationId)
      .subscribe(
        (data: any[]) => {
          console.log('Raw availability data:', data);
          
          const availabilityList = data.map(item => {
            const startDate = item.start_date.split('T')[0]; 
            const endDate = item.end_date.split('T')[0];
            
            console.log('Parsed start date:', startDate);
            console.log('Parsed end date:', endDate);
  
            return {
              id: item.id,
              accommodationId: item.accommodation_id,
              startDate: startDate,
              endDate: endDate
            };
          });
  
          this.availabilityList = availabilityList;
          console.log('Processed availability data:', this.availabilityList);
        },
        error => {
          console.error('Error fetching availability:', error);
        }
      );
  }

  getPrice(accommodationId: string) {
    this.accommodationService.getPriceByAccommodationId(accommodationId)
      .subscribe(
        (data: Price[]) => {
          console.log('Raw price data:', data);
          this.priceList = data.map(item => {
            return {
              id: item.id,
              accommodationId: item.accommodationId,
              startDate: item.startDate,
              endDate: item.endDate,
              amount: item.amount,
              strategy: item.strategy
            };
          });
          console.log('Processed price data:', this.priceList);
        },
        error => {
          console.error('Error fetching price:', error);
        }
      );
  }

//   reserveAccommodation() {
//     if (!this.accommodation) {
//         console.error('Accommodation details are missing');
//         alert('Accommodation details are missing.');
//         return;
//     }

//     if (!this.reservationStartDate || !this.reservationEndDate) {
//         alert('Please select both start and end dates for your reservation.');
//         return;
//     }

//     const startDate = new Date(this.reservationStartDate);
//     const endDate = new Date(this.reservationEndDate);

//     // Provera da li su izabrani datumi unutar dostupnosti
//     const isWithinAvailability = this.availabilityList.some(availability => {
//         const availableStart = new Date(availability.startDate);
//         const availableEnd = new Date(availability.endDate);
//         return startDate >= availableStart && endDate <= availableEnd;
//     });

//     if (!isWithinAvailability) {
//         alert('Selected dates are not within the available dates for this accommodation.');
//         return;
//     }

//     const token = localStorage.getItem('token');
//     if (token) {
//         const decodedToken = this.decodeToken(token);
//         if (decodedToken && decodedToken.userID) {
//             if (this.accommodation && this.accommodation.id) {
//                 const reservationData: Reservation = {
//                     accommodation_id: this.accommodation.id,
//                     guest_id: decodedToken.userID,
//                     start_date: startDate,  // Date objekat
//                     end_date: endDate       // Date objekat
//                 };

//                 console.log('Reservation Data:', reservationData);

//                 const headers = new HttpHeaders({
//                     'Authorization': `Bearer ${token}`,
//                     'Content-Type': 'application/json'
//                 });

//                 // Provera dostupnosti i kreiranje rezervacije na backendu
//                 this.http.post('http://localhost:8081/reservations', reservationData, { headers })
//                     .subscribe({
//                         next: (response) => {
//                             console.log('Accommodation reserved successfully', response);
//                             alert('Reservation successful!');
//                         },
//                         error: (error) => {
//                             if (error.status === 409) {  // HTTP status 409 - Conflict
//                                 alert('The selected dates have already been reserved. Please choose different dates.');
//                             } else {
//                                 console.error('Error reserving accommodation:', error);
//                                 alert('Error occurred while reserving accommodation. Please try again later.');
//                             }
//                         }
//                     });
//             } else {
//                 console.error('Invalid or missing accommodation ID.');
//                 this.router.navigate(['/']);
//             }
//         } else {
//             console.error('Invalid or missing user ID in token.');
//             this.router.navigate(['/login']);
//         }
//     } else {
//         console.error('No token found, redirecting to login.');
//         this.router.navigate(['/login']);
//     }
// }

reserveAccommodation() {
  if (!this.accommodation) {
    console.error('Accommodation details are missing');
    alert('Accommodation details are missing.');
    return;
  }

  if (!this.reservationStartDate || !this.reservationEndDate) {
    alert('Please select both start and end dates for your reservation.');
    return;
  }

  const startDate = new Date(this.reservationStartDate);
  const endDate = new Date(this.reservationEndDate);

  const isWithinAvailability = this.availabilityList.some(availability => {
    const availableStart = new Date(availability.startDate);
    const availableEnd = new Date(availability.endDate);
    return startDate >= availableStart && endDate <= availableEnd;
  });

  if (!isWithinAvailability) {
    alert('Selected dates are not within the available dates for this accommodation.');
    return;
  }

  const token = localStorage.getItem('token');
  if (token) {
    const decodedToken = this.decodeToken(token);
    if (decodedToken && decodedToken.userID) {
      if (this.accommodation && this.accommodation.id) {
        const reservationData: Reservation = {
          accommodation_id: this.accommodation.id,
          guest_id: decodedToken.userID,
          start_date: startDate,
          end_date: endDate
        };

        const headers = new HttpHeaders({
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        });

        this.http.post('http://localhost:8081/reservations', reservationData, { headers })
          .subscribe({
            next: (response) => {
              console.log('Accommodation reserved successfully', response);
              alert('Reservation successful!');

              // Check if userIdd is defined before sending notification
              if (this.userIdd) {
                this.sendNotification(this.userIdd, `Reservation made for accommodation ID: ${this.accommodation?.id}`);
              } else {
                console.error('User ID is missing, notification not sent.');
              }
            },
            error: (error) => {
              if (error.status === 409) {
                alert('The selected dates have already been reserved. Please choose different dates.');
              } else {
                console.error('Error reserving accommodation:', error);
                alert('Error occurred while reserving accommodation. Please try again later.');
              }
            }
          });
      } else {
        console.error('Invalid or missing accommodation ID.');
        this.router.navigate(['/']);
      }
    } else {
      console.error('Invalid or missing user ID in token.');
      this.router.navigate(['/login']);
    }
  } else {
    console.error('No token found, redirecting to login.');
    this.router.navigate(['/login']);
  }
}




sendNotification(hostID: string, message: string) {
  const notificationData = {
      host_id: hostID,
      message: message
  };

  this.http.post('http://localhost:8083/notifications', notificationData)
      .subscribe({
          next: (response) => {
              console.log('Notification sent successfully', response);
          },
          error: (error) => {
              console.error('Error sending notification:', error);
          }
      });
}





  onUpdateSubmit() {
    if (this.accommodation) {
      this.accommodationService.updateAccommodation(this.accommodation.id!, this.accommodation)
        .subscribe(
          response => {
            console.log('Accommodation updated successfully', response);
            this.showUpdateForm = false;
          },
          error => {
            console.error('Error updating accommodation:', error);
          }
        );
    }
  }

  onSubmit() {
    if (this.accommodation) {
      const requestData = {
        startDate: this.startDate,
        endDate: this.endDate,
        amount: this.amount,
        strategy: this.strategy
      };

      this.accommodationService.updateAvailabilityAndPrice(this.accommodation.id!, requestData)
        .subscribe(
          response => {
            console.log('Availability and price updated successfully', response);
            this.showUpdateForm = false;
          },
          error => {
            console.error('Error updating availability and price:', error);
          }
        );
    }
  }

  getAmenities(): string {
    return this.accommodation?.amenities || 'No amenities listed';
  }

  getImages(): string[] {
    return this.accommodation?.images || [];
  }
}
