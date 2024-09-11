import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { AccommodationService } from '../services/accommodation.service';
import { Accommodation } from '../models/Accommodation';
import { HttpHeaders } from '@angular/common/http';
import { jwtDecode } from 'jwt-decode';

@Component({
  selector: 'app-host-accommodations',
  templateUrl: './host-accommodations.component.html',
  styleUrls: ['./host-accommodations.component.css']
})
export class HostAccommodationsComponent implements OnInit {
  accommodations: Accommodation[] = [];

  constructor(private accommodationService: AccommodationService, private router: Router) {}

  ngOnInit(): void {
    this.loadAccommodations();
  }

  decodeToken(token: string): any {
    try {
      return jwtDecode(token);
    } catch (error) {
      console.error('Invalid token:', error);
      return null;
    }
  }

  loadAccommodations(): void {
    console.log('Starting to load accommodations...');

    const token = localStorage.getItem('token');
    if (token) {
        const decodedToken = this.decodeToken(token);
        if (!decodedToken || !decodedToken.userID) {
            console.error('Invalid or missing user ID in token.');
            this.router.navigate(['/login']);
            return;
        }

        console.log('Decoded token:', decodedToken);

        const headers = new HttpHeaders({
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
        });

        console.log('Sending request with headers:', headers);
        console.log('User ID:', decodedToken.userID);

        this.accommodationService.getMyAccommodations(decodedToken.userID).subscribe(
            data => {
                console.log('Accommodations data received:', data);
                this.accommodations = data;
            },
            error => {
                console.error('Error loading accommodations:', error);
                console.error('Error status:', error.status);
                console.error('Error message:', error.message);
                if (error.status === 401) {
                    console.error('Unauthorized access, redirecting to login.');
                    this.router.navigate(['/login']);
                }
            }
        );
    } else {
        console.error('No token found, redirecting to login.');
        this.router.navigate(['/login']);
    }
}

}
