import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';

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
        console.log('Prijava uspešna:', response);
        localStorage.setItem('token', response.token);
        this.router.navigate(['/home']); // Preusmeri na Home stranicu
      }, error => {
        console.error('Greška pri prijavi:', error);
      });
  }
}
