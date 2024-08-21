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
  showUpdateForm: boolean = false; // Kontroliše prikaz forme za ažuriranje
  showAvailabilityForm: boolean = false; // Kontroliše prikaz forme za dostupnost
  showAvailability: boolean = false; // Kontroliše prikaz dostupnosti
  availabilityList: Availability[] = []; // Lista dostupnih termina
  role: string = '';

  // Polja za ažuriranje cene i dostupnosti
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
          this.getAvailability(id); // Učitaj dostupnost kada se učita smeštaj
        },
        error => {
          console.error('Error fetching accommodation details:', error);
        }
      );
    }
  }
  
  // Funkcija za prikazivanje/sakrivanje forme za ažuriranje smeštaja
  toggleUpdateForm() {
    this.showUpdateForm = !this.showUpdateForm;
  }

  toggleAvailabilityForm() {
    this.showAvailabilityForm = !this.showAvailabilityForm;
  }

  // Funkcija za prikazivanje/sakrivanje dostupnosti
  toggleAvailability() {
    this.showAvailability = !this.showAvailability;
  }

  // Funkcija za dobijanje dostupnih termina
  getAvailability(accommodationId: string) {
    this.accommodationService.getAvailabilityByAccommodationId(accommodationId)
      .subscribe(
        (data: Availability[]) => {
          this.availabilityList = data;
        },
        error => {
          console.error('Error fetching availability:', error);
        }
      );
  }

  // Funkcija koja se poziva prilikom slanja forme za ažuriranje smeštaja
  onUpdateSubmit() {
    if (this.accommodation) {
      this.accommodationService.updateAccommodation(this.accommodation.id!, this.accommodation)
        .subscribe(
          response => {
            console.log('Accommodation updated successfully', response);
            this.showUpdateForm = false; // Sakrij formu nakon uspešnog ažuriranja
          },
          error => {
            console.error('Error updating accommodation:', error);
          }
        );
    }
  }

  // Funkcija za slanje forme za ažuriranje dostupnosti i cene
  // onSubmit() {
  //   if (this.accommodation) {
  //     const formattedStartDate = this.formatDate(this.startDate);
  //     const formattedEndDate = this.formatDate(this.endDate);
  
  //     const requestData = {
  //       startDate: formattedStartDate,
  //       endDate: formattedEndDate,
  //       amount: this.amount,
  //       strategy: this.strategy
  //     };
  
  //     this.accommodationService.updateAvailabilityAndPrice(this.accommodation.id!, requestData)
  //       .subscribe(
  //         response => {
  //           console.log('Availability and price updated successfully', response);
  //           this.showUpdateForm = false; // Sakrij formu nakon uspešnog ažuriranja
  //         },
  //         error => {
  //           console.error('Error updating availability and price:', error);
  //         }
  //       );
  //   }
  // }

  onSubmit() {
    if (this.accommodation) {
      // Direktno koristite vrednosti iz input polja
      const requestData = {
        startDate: this.startDate, // Već je u formatu YYYY-MM-DD
        endDate: this.endDate,     // Već je u formatu YYYY-MM-DD
        amount: this.amount,
        strategy: this.strategy
      };
  
      this.accommodationService.updateAvailabilityAndPrice(this.accommodation.id!, requestData)
        .subscribe(
          response => {
            console.log('Availability and price updated successfully', response);
            this.showUpdateForm = false; // Sakrij formu nakon uspešnog ažuriranja
          },
          error => {
            console.error('Error updating availability and price:', error);
          }
        );
    }
  }
  

  // formatDate(date: string): string {
  //   const parsedDate = new Date(date);
  //   const year = parsedDate.getFullYear();
  //   const month = ('0' + (parsedDate.getMonth() + 1)).slice(-2);
  //   const day = ('0' + parsedDate.getDate()).slice(-2);
  //   return `${year}-${month}-${day}`;
  // }
  // formatDate(date: string): string {
  //   const [day, month, year] = date.split('.');
  //   return `${year}-${month}-${day}`;
  // }
//   formatDate(date: string): string {
//     const [year, month, day] = date.split('-');
//     return `${year}-${month}-${day}`;
// }

  

  // Metoda za prikazivanje pogodnosti
  getAmenities(): string {
    return this.accommodation?.amenities || 'No amenities listed';
  }

  // Metoda za prikazivanje slika
  getImages(): string[] {
    return this.accommodation?.images || [];
  }
}
