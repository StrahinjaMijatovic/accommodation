import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { FormsModule } from '@angular/forms';
import { HttpClientModule } from '@angular/common/http';
import { ReactiveFormsModule } from '@angular/forms';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { LoginComponent } from './login/login.component';
import { RegisterComponent } from './register/register.component';
import { HomeComponent } from './home/home.component';
import { ProfileComponent } from './profile/profile.component';
import { AccommodationComponent } from './accommodation/accommodation.component';
import { SeeAccommodationComponent } from './see-accommodation/see-accommodation.component';
import { UpdateAccommodationComponent } from './update-accommodation/update-accommodation.component';
import { GuestReservationsComponent } from './guest-reservations/guest-reservations.component';
import { HostAccommodationsComponent } from './host-accommodations/host-accommodations.component';
import { HostNotificationsComponent } from './host-notifications/host-notifications.component'; // Importuj komponentu

@NgModule({
  declarations: [
    AppComponent,
    LoginComponent,
    RegisterComponent,
    HomeComponent,
    ProfileComponent,
    AccommodationComponent,
    SeeAccommodationComponent,
    UpdateAccommodationComponent,
    GuestReservationsComponent,
    HostAccommodationsComponent,
    HostNotificationsComponent // Registruj komponentu
  ],
  imports: [
    BrowserModule,
    FormsModule,
    AppRoutingModule,
    HttpClientModule,
    ReactiveFormsModule          
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
