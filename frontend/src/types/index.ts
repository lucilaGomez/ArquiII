export interface Hotel {
  id: string;
  name: string;
  description: string;
  city: string;
  address: string;
  photos: string[];
  thumbnail: string;
  amenities: string[];
  rating: number;
  price_range: {
    min_price: number;
    max_price: number;
    currency: string;
  };
  contact: {
    phone: string;
    email: string;
    website: string;
  };
  is_active: boolean;
}

export interface SearchParams {
  city: string;
  checkin: string;
  checkout: string;
  guests: number;
}

export interface SearchResult {
  id: string;
  name: string;
  description: string;
  city: string;
  thumbnail: string;
  amenities: string[];
  rating: number;
  min_price: number;
  max_price: number;
  currency: string;
  available: boolean;
}

export interface User {
  id: number;
  email: string;
  first_name: string;
  last_name: string;
  phone: string;
}

export interface Booking {
  id: number;
  hotel_id: string;
  check_in_date: string;
  check_out_date: string;
  guests: number;
  room_type: string;
  total_price: number;
  currency: string;
  status: string;
  booking_reference: string;
  special_requests: string;
}

export interface AuthResponse {
  user: User;
  token: string;
}
