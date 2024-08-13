import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { HttpClient } from '@angular/common/http';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css']
})
export class HomeComponent implements OnInit {
  accommodations: any[] = [];
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
    this.http.get('http://localhost:8000/accommodations')
      .subscribe((data: any) => {
        this.accommodations = data;
      }, error => {
        console.error('Greška prilikom učitavanja smeštaja:', error);
      });
  }

  checkUserStatus() {
    const token = localStorage.getItem('token');
    if (token) {
      this.http.post('http://localhost:8000/verify-token', { token })
        .subscribe((response: any) => {
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
