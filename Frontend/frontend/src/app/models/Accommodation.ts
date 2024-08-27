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
  id?: string; // Opcionalno jer će ga backend dodeliti
  accommodation_id: string; // ID smeštaja
  guest_id: string; // ID gosta koji pravi rezervaciju
  start_date: Date; // Datum početka rezervacije (format: YYYY-MM-DD)
  end_date: Date; // Datum kraja rezervacije (format: YYYY-MM-DD)
}

