export interface Accommodation {
  id?: string;
  name: string;
  location: string;
  guests: number;
  price: number;
  amenities: string; 
  images: string[]; 
}


export interface Availability {
  id?: string;
  accommodationId: string;
  startDate: Date; // Format "yyyy-MM-dd"
  endDate: Date;   // Format "yyyy-MM-dd"
}
