export interface Accommodation {
  id?: string;
  name: string;
  location: string;
  guests: number;
  price: number;
  amenities: string[]; // Dodaj ovo polje
  images: string[]; // Dodaj ovo polje
}
