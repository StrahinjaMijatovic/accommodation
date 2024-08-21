import { Component } from '@angular/core';
import { Accommodation } from '../models/Accommodation';
import { AccommodationService } from '../services/accommodation.service';

@Component({
  selector: 'app-accommodation',
  templateUrl: './accommodation.component.html',
  styleUrls: ['./accommodation.component.css']
})
export class AccommodationComponent {
  accommodation: Accommodation = {
    name: '',
    location: '',
    guests: 0,
    price: 0,
    amenities: '',
    images: [] // Inicijalizacija niza za slike
  };

  // amenitiesString: string = ''; // String za unos pogodnosti
  imagesString: string = ''; // String za unos slika

  constructor(private accommodationService: AccommodationService) {}

  onSubmit() {
    // Pretvori stringove u nizove
    // this.accommodation.amenities = this.amenitiesString.split(',').map(item => item.trim());
    this.accommodation.images = this.imagesString.split(',').map(item => item.trim());

    this.accommodationService.createAccommodation(this.accommodation)
      .subscribe(response => {
        console.log('Accommodation created:', response);
        // Reset forme nakon slanja
        this.accommodation = {
          name: '',
          location: '',
          guests: 0,
          price: 0,
          amenities: '',
          images: []
        };
        // this.amenitiesString = '';
        this.imagesString = '';
      }, error => {
        console.error('Error creating accommodation:', error);
      });
  }
}