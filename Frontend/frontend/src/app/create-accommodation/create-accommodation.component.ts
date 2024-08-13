import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';

@Component({
  selector: 'app-create-accommodation',
  templateUrl: './create-accommodation.component.html',
  styleUrls: ['./create-accommodation.component.css']
})
export class CreateAccommodationComponent {
  accommodation = {
    name: '',
    location: '',
    amenities: '' as string,  // Promeniti ovo na string
    minGuests: 0,
    maxGuests: 0,
    images: '' as string   // Promeniti ovo na string
  };

  constructor(private http: HttpClient, private router: Router) {}

  onCreateAccommodation() {
    // Pretvori stringove u nizove
    const amenitiesArray = this.accommodation.amenities.split(',').map(item => item.trim());
    const imagesArray = this.accommodation.images.split(',').map(item => item.trim());

    // Stvori objekat koji šalješ
    const accommodationData = {
      ...this.accommodation,
      amenities: amenitiesArray,
      images: imagesArray
    };

    this.http.post('http://localhost:8001/accommodations', accommodationData)
      .subscribe(response => {
        console.log('Accommodation created:', response);
        this.router.navigate(['/home']); // Preusmeravanje nazad na Home stranicu
      }, error => {
        console.error('Error creating accommodation:', error);
      });
  }
}
