import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
import { jwtDecode } from 'jwt-decode';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html'
})
export class LoginComponent {
  email: string = '';
  password: string = '';
  errorMessage: string = ''; // Dodaj errorMessage za prikaz grešaka

  constructor(private http: HttpClient, private router: Router) {}

  onLogin() {
    this.errorMessage = ''; // Resetovanje poruke o grešci

    // Validacija email-a
    const emailPattern = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
    if (!emailPattern.test(this.email)) {
      this.errorMessage = 'Invalid email format. Please include "@" in your email.';
      return;
    }

    // Validacija lozinke
    if (this.password.length < 6) {
      this.errorMessage = 'Password must be at least 6 characters long.';
      return;
    }

    // Ako su validacije uspešne, nastavi sa loginom
    const loginData = { email: this.email, password: this.password };
    this.http.post('http://localhost:8000/login', loginData)
      .subscribe((response: any) => {
        console.log('Login successful:', response);

        const token = response.token;
        localStorage.setItem('token', token);

        const decodedToken: any = jwtDecode(token);
        const userID = decodedToken.userID;

        localStorage.setItem('userID', userID);

        this.router.navigate(['/home']);
      }, error => {
        console.error('Login error:', error);
        this.errorMessage = 'Login failed. Please check your email and password.';
      });
  }
}
