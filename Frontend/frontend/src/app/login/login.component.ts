import { Component } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html'
})
export class LoginComponent {
  email: string = '';
  password: string = '';

  constructor(private http: HttpClient) {}

  onLogin() {
    const loginData = { email: this.email, password: this.password };
    this.http.post('http://localhost:8000/login', loginData)
      .subscribe(response => {
        console.log('Login successful:', response);
        // Ovde možeš dodati redirekciju ili čuvanje tokena
      }, error => {
        console.error('Login error:', error);
      });
  }
}
