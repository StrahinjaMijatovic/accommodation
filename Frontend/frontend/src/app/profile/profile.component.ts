import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';

@Component({
  selector: 'app-profile',
  templateUrl: './profile.component.html'
})
export class ProfileComponent implements OnInit {
  profile: any = {
    firstName: '',
    lastName: '',
    email: '',
    gender: '',
    age: null,
    location: ''
  };

  constructor(private http: HttpClient, private router: Router) {}

  ngOnInit() {
    this.loadProfile();
  }

  loadProfile() {
    const token = localStorage.getItem('token');
    if (token) {
      this.http.post('http://localhost:8000/verify-token', { token })
        .subscribe((response: any) => {
          this.http.get(`http://localhost:8000/profile?email=${response.email}`)
            .subscribe((profile: any) => {
              this.profile = profile;
            }, error => {
              console.error('Greška prilikom učitavanja profila:', error);
            });
        }, error => {
          console.error('Provera tokena nije uspela:', error);
          this.router.navigate(['/login']);
        });
    } else {
      this.router.navigate(['/login']);
    }
  }

  onSave() {
    this.http.put('http://localhost:8000/profile', this.profile)
      .subscribe(() => {
        alert('Profil uspešno ažuriran');
      }, error => {
        console.error('Greška prilikom ažuriranja profila:', error);
      });
  }
}
