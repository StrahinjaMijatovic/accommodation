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
  startDate: Date;
  endDate: Date;   
}

export interface Price {
  id: string;
  accommodationId: string;
  startDate: Date;
  endDate: Date; 
  amount: number;
  strategy: string;  
}

// src/app/models/Reservation.ts
export interface Reservation {
  id?: string; 
  accommodation_id: string; 
  guest_id: string; 
  start_date: Date; 
  end_date: Date; 
}

export interface Rating {
  id?: string;           
  user_id: string;       
  targetID: string;     
  rating: number;        
  comment: string;      
  ratedAt?: Date;      
}


