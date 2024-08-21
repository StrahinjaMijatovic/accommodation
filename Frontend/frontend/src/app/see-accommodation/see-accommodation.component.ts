import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { AccommodationService } from '../services/accommodation.service';
import { Accommodation, Availability } from '../models/Accommodation';

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
  availabilityList: Availability[] = [];
  role: string = '';

  startDate: string = '';
  endDate: string = '';
  amount: number = 0;
  strategy: 'per_guest' | 'per_unit' = 'per_unit';

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

  getAvailability(accommodationId: string) {
    this.accommodationService.getAvailabilityByAccommodationId(accommodationId)
      .subscribe(
        (data: any[]) => {
          console.log('Raw availability data:', data);
          
          const availabilityList = data.map(item => {
            // Parsiranje datuma iz stringa pre "T"
            const startDate = item.start_date.split('T')[0]; 
            const endDate = item.end_date.split('T')[0];
            
            // Logovanje parsed datuma
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