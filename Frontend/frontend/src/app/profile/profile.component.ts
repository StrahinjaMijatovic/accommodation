import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';

@Component({
  selector: 'app-profile',
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.css']
})
export class ProfileComponent implements OnInit {
  profile: any = {
    firstName: '',
    lastName: '',
    email: '',
    age: null,
    location: ''
  };
  
  passwords = {
    currentPassword: '',
    newPassword: '',
    confirmNewPassword: ''
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
          const email = response.email;
          if (email) {
            this.http.get(`http://localhost:8001/profile?email=${email}`)
              .subscribe((profile: any) => {
                this.profile = profile;
              }, error => {
                console.error('Error loading profile:', error);
              });
          } else {
            console.error('Email is missing from token response');
            this.router.navigate(['/login']);
          }
        }, error => {
          console.error('Token verification failed:', error);
          this.router.navigate(['/login']);
        });
    } else {
      this.router.navigate(['/login']);
    }
  }

  onSaveProfile() {
    this.http.put('http://localhost:8001/profile', this.profile, {
      headers: { 'Content-Type': 'application/json' }
    })
    .subscribe(() => {
      alert('Profile successfully updated');
    }, error => {
      console.error('Error updating profile:', error);
    });
  }

  onChangePassword() {
    if (this.passwords.newPassword !== this.passwords.confirmNewPassword) {
      alert('New password and confirmation do not match');
      return;
    }
  
    const token = localStorage.getItem('token');
    if (token) {
      this.http.post('http://localhost:8000/verify-token', { token })
        .subscribe((response: any) => {
          const email = response.email;
          if (email) {
            this.http.post('http://localhost:8001/change-password', {
              email: email,
              currentPassword: this.passwords.currentPassword,
              newPassword: this.passwords.newPassword
            }, {
              headers: { 'Content-Type': 'application/json' }
            })
            .subscribe(() => {
              alert('Password successfully changed');
              this.passwords = {
                currentPassword: '',
                newPassword: '',
                confirmNewPassword: ''
              };
            }, error => {
              console.error('Error changing password:', error);
              alert('Error changing password: ' + error.error);
            });
          } else {
            console.error('Email is missing from token response');
            this.router.navigate(['/login']);
          }
        }, error => {
          console.error('Token verification failed:', error);
          this.router.navigate(['/login']);
        });
    } else {
      this.router.navigate(['/login']);
    }
  }
  onDeleteProfile() {
    if (confirm('Are you sure you want to delete your profile? This action cannot be undone.')) {
      const token = localStorage.getItem('token');
      if (token) {
        this.http.post('http://localhost:8000/verify-token', { token })
          .subscribe((response: any) => {
            const userID = response.userID;
            if (userID) {
              this.http.delete(`http://localhost:8001/profile/${userID}`)
                .subscribe(() => {
                  alert('Profile successfully deleted');
                  localStorage.removeItem('token');
                  this.router.navigate(['/register']);
                }, error => {
                  console.error('Error deleting profile:', error);
                  alert('Error deleting profile: ' + error.error);
                });
            } else {
              console.error('userID is missing from token response');
              this.router.navigate(['/login']);
            }
          }, error => {
            console.error('Token verification failed:', error);
            this.router.navigate(['/login']);
          });
      } else {
        this.router.navigate(['/login']);
      }
    }
  }
  
}

