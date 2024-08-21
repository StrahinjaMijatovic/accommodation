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
  id?: string;  // Opcionalno polje za UUID dostupnosti
  accommodationId: string; // ID smeštaja kojem dostupnost pripada
  startDate: string; // Datum početka dostupnosti (kao string, u formatu "yyyy-mm-dd")
  endDate: string; // Datum završetka dostupnosti (kao string, u formatu "yyyy-mm-dd")
}
