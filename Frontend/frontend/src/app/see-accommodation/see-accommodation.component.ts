import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { AccommodationService } from '../services/accommodation.service';
import { Accommodation } from '../models/Accommodation';

@Component({
  selector: 'app-see-accommodation',
  templateUrl: './see-accommodation.component.html',
  styleUrls: ['./see-accommodation.component.css']
})
export class SeeAccommodationComponent implements OnInit {
  accommodation: Accommodation | undefined;

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
        },
        error => {
          console.error('Error fetching accommodation details:', error);
        }
      );
    }
  }

  // Primer korišćenja `?.` operatera unutar TypeScript metode
  getAmenities(): string {
    return this.accommodation?.amenities?.join(', ') || 'No amenities listed';
  }

  getImages(): string[] {
    return this.accommodation?.images || [];
  }

  
  
}
