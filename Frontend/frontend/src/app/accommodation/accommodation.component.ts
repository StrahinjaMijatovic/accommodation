import { Component } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Router } from '@angular/router';
import { jwtDecode } from 'jwt-decode';

@Component({
  selector: 'app-accommodation',
  templateUrl: './accommodation.component.html',
  styleUrls: ['./accommodation.component.css']
})
export class AccommodationComponent {
  accommodation: any = {
    name: '',
    location: '',
    guests: 0,
    price: 0,
    amenities: '',
    images: [],
    userID: ''
  };
  imagesString: string = '';

  constructor(private http: HttpClient, private router: Router) {}

  decodeToken(token: string): any {
    try {
      return jwtDecode(token);
    } catch (error) {
      console.error('Invalid token:', error);
      return null;
    }
  }

  onSubmit() {
    this.accommodation.images = this.imagesString.split(',').map(item => item.trim());
  
    const token = localStorage.getItem('token');
    if (token) {
      const decodedToken = this.decodeToken(token);
      if (decodedToken && decodedToken.userID) {
        this.accommodation.userID = decodedToken.userID;
      } else {
        console.error('Invalid or missing user ID in token.');
        this.router.navigate(['/login']);
        return;
      }
    } else {
      console.error('No token found, redirecting to login.');
      this.router.navigate(['/login']);
      return;
    }
  
    const headers = new HttpHeaders({
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    });
  
    this.http.post('http://localhost:8080/accommodations', this.accommodation, { headers })
      .subscribe({
        next: (response) => {
          console.log('Accommodation created:', response);
          this.accommodation = {
            name: '',
            location: '',
            guests: 0,
            price: 0,
            amenities: '',
            images: [],
            userID: ''
          };
          this.imagesString = '';
          this.router.navigate(['/home']);
        },
        error: (error) => {
          console.error('Error creating accommodation:', error);
        }
      });
  }
  
}
