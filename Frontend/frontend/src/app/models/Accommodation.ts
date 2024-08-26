export interface Accommodation {
  id?: string;
  name: string;
  location: string;
  guests: number;
  price: number;
  amenities: string; 
  images: string[]; 
  userID?: string;
}


export interface Availability {
  id?: string;
  accommodationId: string;
  startDate: Date; // Format "yyyy-MM-dd"
  endDate: Date;   // Format "yyyy-MM-dd"
}

export interface Price {
  id: string;
  accommodationId: string;
  startDate: Date; // Format "yyyy-MM-dd"
  endDate: Date; // ili Date, ako želite da koristite datumski objekat u TypeScript-u
  amount: number;
  strategy: string;   // "per_guest" ili "per_unit"
}

// src/app/models/Reservation.ts

export interface Reservation {
  id: number;
  accommodation_id: number;
  guest_id: number;
  startDate: Date; // We'll keep these as strings to work with Angular forms easily
  endDate: Date;
  status: string;
}
