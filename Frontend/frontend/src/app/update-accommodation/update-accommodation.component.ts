import { Component } from '@angular/core';
import { AccommodationService } from '../services/accommodation.service';
import { ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-update-accommodation',
  templateUrl: './update-accommodation.component.html'
})
export class UpdateAccommodationComponent {
  startDate: Date = new Date(); 
  endDate: Date = new Date();  
  amount: number = 0;
  strategy: 'per_guest' | 'per_unit' = 'per_unit';
  accommodationId: string = '';

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
      startDate: this.formatDate(this.startDate),
      endDate: this.formatDate(this.endDate),
      amount: this.amount,
      strategy: this.strategy
    };

    console.log("Start Date:", newPrice.startDate);
    console.log("End Date:", newPrice.endDate);

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

  formatDate(date: Date): string {
    return date.toISOString().split('T')[0]; 
  }
}