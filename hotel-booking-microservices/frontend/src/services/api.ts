import axios from 'axios';
import { User, Booking, AuthResponse } from '../types';

const API_BASE = {
  hotel: 'http://localhost:8001/api/v1',
  search: 'http://localhost:8002/api/v1', 
  booking: 'http://localhost:8003/api'     // CORREGIDO: Puerto 8003 y SIN /v1
};

// Configurar axios para incluir token automáticamente
const bookingApi = axios.create({
  baseURL: API_BASE.booking
});

// Interceptor para agregar token a requests
bookingApi.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export const hotelAPI = {
  // Obtener todos los hoteles
  getAllHotels: async () => {
    const response = await axios.get(`${API_BASE.hotel}/hotels`);
    return { data: response.data.data || [] };
  },

  // Obtener hotel por ID
  getHotel: async (id: string) => {
    const response = await axios.get(`${API_BASE.hotel}/hotels/${id}`);
    return { data: response.data.data };
  },

  // Crear hotel
  createHotel: async (hotelData: any) => {
    const token = localStorage.getItem('token');
    const response = await axios.post(`${API_BASE.hotel}/hotels`, hotelData, {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    });
    return { data: response.data.data };
  },

  // Actualizar hotel
  updateHotel: async (id: string, hotelData: any) => {
    const token = localStorage.getItem('token');
    const response = await axios.put(`${API_BASE.hotel}/hotels/${id}`, hotelData, {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    });
    return { data: response.data.data };
  },

  // Eliminar hotel
  deleteHotel: async (id: string) => {
    const token = localStorage.getItem('token');
    await axios.delete(`${API_BASE.hotel}/hotels/${id}`, {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    });
  },

  // Estadísticas de hoteles
  getHotelStats: async () => {
    const token = localStorage.getItem('token');
    const response = await axios.get(`${API_BASE.hotel}/stats`, {
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    });
    return { data: response.data.data };
  }
};

export const searchAPI = {
  // Buscar hoteles
  searchHotels: async (params: {
    city: string;
    checkin?: string;
    checkout?: string;
    guests?: number;
    min_price?: number;
    max_price?: number;
    min_rating?: number;
  }) => {
    const queryParams = new URLSearchParams();
    Object.entries(params).forEach(([key, value]) => {
      if (value !== undefined && value !== '') {
        queryParams.append(key, value.toString());
      }
    });

    const response = await axios.get(
      `${API_BASE.search}/search/hotels?${queryParams.toString()}`
    );
    return response.data;
  }
};

export const bookingAPI = {
  // Registro - CORREGIDO: sin /v1
  register: async (userData: {
    email: string;
    password: string;
    first_name: string;
    last_name: string;
    phone: string;
  }): Promise<User> => {
    const response = await axios.post(`${API_BASE.booking}/auth/register`, userData);
    return response.data.data;
  },

  // Login - CORREGIDO: sin /v1
  login: async (credentials: {
    email: string;
    password: string;
  }): Promise<AuthResponse> => {
    const response = await axios.post(`${API_BASE.booking}/auth/login`, credentials);
    return response.data.data;
  },

  // Verificar disponibilidad
  checkAvailability: async (hotelId: string, params: {
    checkin: string;
    checkout: string;
    guests: number;
  }) => {
    const queryParams = new URLSearchParams({
      checkin: params.checkin,
      checkout: params.checkout,
      guests: params.guests.toString()
    });

    const response = await axios.get(
      `${API_BASE.booking}/availability/${hotelId}?${queryParams.toString()}`
    );
    return response.data.data;
  },

  // Crear reserva
  createBooking: async (bookingData: {
    hotel_id: string;
    check_in_date: string;
    check_out_date: string;
    guests: number;
    room_type: string;
    special_requests: string;
  }): Promise<Booking> => {
    const response = await bookingApi.post('/bookings', bookingData);
    return response.data.data;
  },

  // Obtener reservas del usuario
  getUserBookings: async (): Promise<Booking[]> => {
    const response = await bookingApi.get('/bookings');
    return response.data.data;
  },

  // Obtener perfil
  getProfile: async (): Promise<User> => {
    const response = await bookingApi.get('/profile');
    return response.data.data;
  }
};