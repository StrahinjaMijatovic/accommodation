import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { HttpClient } from '@angular/common/http';

interface Accommodation {
  id: string;
  name: string;
  location: string;
  guests: number;
  price: number;
  amenities: string[];
  images: string[];
}

interface UserResponse {
  firstName: string;
  lastName: string;
  role: string;
}

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css']
})
export class HomeComponent implements OnInit {
  accommodations: Accommodation[] = [];
  isLoggedIn: boolean = false;
  firstName: string = '';
  lastName: string = '';
  role: string = '';

  constructor(private router: Router, private http: HttpClient) {}

  ngOnInit() {
    this.loadAccommodations();
    this.checkUserStatus();
  }

  navigateToCreateAccommodation() {
    this.router.navigate(['/create-accommodation']);
  }

  navigateToLogin() {
    this.router.navigate(['/login']);
  }

  navigateToRegister() {
    this.router.navigate(['/register']);
  }

  logout() {
    localStorage.removeItem('token');
    this.isLoggedIn = false;
    this.firstName = '';
    this.lastName = '';
    this.role = '';
    this.router.navigate(['/home']);
  }

  loadAccommodations() {
    this.http.get<Accommodation[]>('http://localhost:8080/accommodations')
      .subscribe((data: Accommodation[]) => {
        this.accommodations = data;
      }, error => {
        console.error('Greška prilikom učitavanja smeštaja:', error);
      });
  }

  checkUserStatus() {
    const token = localStorage.getItem('token');
    if (token) {
      this.http.post<UserResponse>('http://localhost:8000/verify-token', { token })
        .subscribe((response: UserResponse) => {
          this.isLoggedIn = true;
          this.firstName = response.firstName;
          this.lastName = response.lastName;
          this.role = response.role;
        }, error => {
          console.error('Provera tokena nije uspela:', error);
          this.isLoggedIn = false;
        });
    }
  }
}
