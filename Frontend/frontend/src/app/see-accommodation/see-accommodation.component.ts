import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { AccommodationService } from '../services/accommodation.service';
import { Accommodation, Availability, Price, Reservation } from '../models/Accommodation';

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
  showPrice: boolean = false; // Dodano za prikazivanje cena
  showReservationForm: boolean = false; // Prikaz forme za unos datuma rezervacije
  availabilityList: Availability[] = [];
  priceList: Price[] = []; // Lista cena
  role: string = '';

  startDate: string = '';
  endDate: string = '';
  amount: number = 0;
  strategy: 'per_guest' | 'per_unit' = 'per_unit';

  reservationStartDate: string = ''; // Datum početka rezervacije
  reservationEndDate: string = ''; // Datum završetka rezervacije

  constructor(
    private route: ActivatedRoute,
    private accommodationService: AccommodationService
  ) {}

  ngOnInit(): void {
    const id = this.route.snapshot.paramMap.get('id');
    if (id) {
      this.accommodationService.getAccommodationById(id).subscribe(
        data => {
          this.accommodation = data;
          this.getAvailability(id);
          this.getPrice(id); // Pozivamo metodu za dobijanje cena
        },
        error => {
          console.error('Error fetching accommodation details:', error);
        }
      );
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
  
    const reservationData: Reservation = {
      id: 0, // ID generiše backend
      accommodation_id: +this.accommodation.id!, // Pretpostavljamo da je ID broj
      guest_id: 0, // Ovaj ID treba da postavite iz sesije ili autentifikacije
      startDate: startDate, // Konvertovanje stringa u datum
      endDate: endDate, // Konvertovanje stringa u datum
      status: "pending" // Početni status
    };
  
    this.accommodationService.reserveAccommodation(reservationData)
      .subscribe(
        response => {
          console.log('Accommodation reserved successfully', response);
          alert('Reservation successful!');
        },
        error => {
          console.error('Error reserving accommodation:', error);
          alert('Error occurred while reserving accommodation. Please try again later.');
        }
      );
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
