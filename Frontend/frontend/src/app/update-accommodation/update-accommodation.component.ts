import { Component } from '@angular/core';
import { AccommodationService } from '../services/accommodation.service';
import { ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-update-accommodation',
  templateUrl: './update-accommodation.component.html'
})
export class UpdateAccommodationComponent {
  startDate: string = ''; // Inicijalizacija sa praznim stringom
  endDate: string = '';   // Inicijalizacija sa praznim stringom
  amount: number = 0;     // Inicijalizacija sa 0
  strategy: 'per_guest' | 'per_unit' = 'per_unit'; // Inicijalizacija sa podrazumevanom vrednošću
  accommodationId: string = ''; // Inicijalizacija sa praznim stringom

  constructor(
    private accommodationService: AccommodationService,
    private route: ActivatedRoute
  ) {}

  ngOnInit() {
    this.accommodationId = this.route.snapshot.paramMap.get('id') || '';
  }

  onSubmit() {
    const newPrice = {
      accommodationId: this.accommodationId,
      startDate: this.startDate,
      endDate: this.endDate,
      amount: this.amount,
      strategy: this.strategy
    };

    this.accommodationService.updatePrice(this.accommodationId, newPrice)
      .subscribe(
        response => {
          console.log('Price updated successfully');
        },
        error => {
          console.error('Error updating price:', error);
        }
      );
  }
}
