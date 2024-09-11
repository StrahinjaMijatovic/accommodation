import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { LoginComponent } from './login/login.component';
import { RegisterComponent } from './register/register.component';
import { HomeComponent } from './home/home.component';
import { AccommodationComponent } from './accommodation/accommodation.component';
import { ProfileComponent } from './profile/profile.component';
import { SeeAccommodationComponent } from './see-accommodation/see-accommodation.component'; 
import { UpdateAccommodationComponent } from './update-accommodation/update-accommodation.component';
import { GuestReservationsComponent } from './guest-reservations/guest-reservations.component';
import { HostAccommodationsComponent } from './host-accommodations/host-accommodations.component';
import { HostNotificationsComponent } from './host-notifications/host-notifications.component';


const routes: Routes = [
  { path: 'host-notifications', component: HostNotificationsComponent },
  { path: 'host-accommodations', component: HostAccommodationsComponent },
  { path: 'my-reservations', component: GuestReservationsComponent },
  { path: 'proba', component: UpdateAccommodationComponent },
  { path: 'accommodation/:id', component: SeeAccommodationComponent },
  { path: 'create-accommodation', component: AccommodationComponent },
  { path: '', redirectTo: '/home', pathMatch: 'full' },
  { path: 'login', component: LoginComponent },
  { path: 'register', component: RegisterComponent },
  { path: 'home', component: HomeComponent },
  { path: 'profile', component: ProfileComponent }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
