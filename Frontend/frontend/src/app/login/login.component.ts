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

  constructor(private http: HttpClient, private router: Router) {}

  onLogin() {
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
      });
  }
}
