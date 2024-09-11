import { Component, OnInit } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Router } from '@angular/router';
import { jwtDecode } from 'jwt-decode';

interface Notification {
  id: number;
  host_id: string;
  message: string;
  created_at: Date;
}

@Component({
  selector: 'app-host-notifications',
  templateUrl: './host-notifications.component.html',
  styleUrls: ['./host-notifications.component.css']
})
export class HostNotificationsComponent implements OnInit {
  notifications: Notification[] = [];  
  constructor(private http: HttpClient, private router: Router) {}  

  ngOnInit(): void {
    console.log('HostNotificationsComponent initialized.');
    this.loadNotifications();
  }

  decodeToken(token: string): any {
    try {
      const decoded = jwtDecode(token);
      console.log('Token successfully decoded:', decoded);
      return decoded;
    } catch (error) {
      console.error('Invalid token:', error);
      return null;
    }
  }

  loadNotifications(): void {
    const token = localStorage.getItem('token');
    if (token) {
      console.log('Token found:', token);
      const decodedToken = this.decodeToken(token);
      if (!decodedToken || !decodedToken.userID) {
        console.error('Invalid or missing user ID in token.');
        this.router.navigate(['/login']);
        return;
      }

      console.log('Decoded user ID:', decodedToken.userID);

      const headers = new HttpHeaders({
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      });

      console.log('Sending request to load notifications with headers:', headers);

    
      this.http.get<any>(`http://localhost:8083/notifications/${decodedToken.userID}`, { headers })
        .subscribe({
          next: (response) => {
            console.log('Raw response received:', response);
            this.notifications = response;  // Ažuriramo svojstvo notifications
            console.log('Notifications state updated:', this.notifications);
          },
          error: (error: any) => {
            console.error('Error loading notifications:', error);
            alert('Error occurred while loading notifications. Please try again later.');
          }
        });
    } else {
      console.error('No token found, redirecting to login.');
      this.router.navigate(['/login']);
    }
  }
}
