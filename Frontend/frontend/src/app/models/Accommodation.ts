export interface Accommodation {
  id?: string;
  name: string;
  location: string;
  guests: number;
  price: number;
  amenities: string[]; 
  images: string[]; 
}