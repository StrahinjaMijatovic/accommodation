import { Component, OnInit } from '@angular/core';
import { AccommodationService } from '../services/accommodation.service';
import { Reservation } from '../models/Accommodation';

@Component({
  selector: 'app-my-reservations',
  templateUrl: './my-reservations.component.html',
  styleUrls: ['./my-reservations.component.css']
})
export class MyReservationsComponent implements OnInit {

  reservations: Reservation[] = [];

  constructor(private accommodationService: AccommodationService) { }

  ngOnInit(): void {
    this.getReservations();
  }

  getReservations(): void {
    const guestId = 0; // Pretpostavimo da je guest_id već definisan; promenite ovo prema stvarnoj implementaciji
    this.accommodationService.getReservationsByGuestId(guestId).subscribe(
      data => {
        this.reservations = data;
      },
      error => {
        console.error('Error fetching reservations:', error);
      }
    );
  }

}
