import React, { useState, useEffect } from 'react';
import { useNavigate, useParams, useSearchParams } from 'react-router-dom';
import { hotelAPI, bookingAPI } from '../services/api';
import { Hotel } from '../types';

const Detail: React.FC = () => {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const [searchParams] = useSearchParams();
  
  const [hotel, setHotel] = useState<Hotel | null>(null);
  const [availability, setAvailability] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [bookingLoading, setBookingLoading] = useState(false);
  const [error, setError] = useState<string>('');
  const [showLoginForm, setShowLoginForm] = useState(false);
  const [isLoginMode, setIsLoginMode] = useState(true); // true = login, false = registro
  const [loginData, setLoginData] = useState({ email: '', password: '' });
  const [registerData, setRegisterData] = useState({
    email: '',
    password: '',
    first_name: '',
    last_name: '',
    phone: '',
    date_of_birth: '1990-01-01'
  });
  const [isLoggedIn, setIsLoggedIn] = useState(false);

  // Obtener parámetros de búsqueda
  const city = searchParams.get('city') || '';
  const checkin = searchParams.get('checkin') || '';
  const checkout = searchParams.get('checkout') || '';
  const guests = parseInt(searchParams.get('guests') || '2');

  useEffect(() => {
    // Verificar si el usuario está logueado
    const token = localStorage.getItem('token');
    setIsLoggedIn(!!token);
  }, []);

  useEffect(() => {
    const fetchHotelDetails = async () => {
      if (!id) {
        navigate('/');
        return;
      }

      try {
        setLoading(true);
        
        // Obtener detalles del hotel
        const hotelData = await hotelAPI.getHotel(id);
        setHotel(hotelData);

        // Verificar disponibilidad si tenemos fechas
        if (checkin && checkout) {
          try {
            const availabilityData = await bookingAPI.checkAvailability(id, {
              checkin,
              checkout,
              guests
            });
            setAvailability(availabilityData);
          } catch (availErr) {
            console.error('Error checking availability:', availErr);
            // No mostrar error por disponibilidad, solo continuar
          }
        }
      } catch (err) {
        setError('Error cargando detalles del hotel');
        console.error('Error fetching hotel details:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchHotelDetails();
  }, [id, checkin, checkout, guests, navigate]);

  const handleLogin = async () => {
    try {
      const response = await bookingAPI.login(loginData);
      localStorage.setItem('token', response.token);
      setIsLoggedIn(true);
      setShowLoginForm(false);
      setLoginData({ email: '', password: '' }); // Limpiar formulario
      alert('¡Login exitoso!');
    } catch (err) {
      alert('Error en login. Verifica tus credenciales.');
      console.error('Login error:', err);
    }
  };

  const handleRegister = async () => {
    // Validación básica
    if (!registerData.email || !registerData.password || !registerData.first_name || !registerData.last_name) {
      alert('Por favor completa todos los campos obligatorios');
      return;
    }

    try {
      const registerPayload = {
        email: registerData.email,
        password: registerData.password,
        first_name: registerData.first_name,
        last_name: registerData.last_name,
        phone: registerData.phone,
        date_of_birth: registerData.date_of_birth + 'T00:00:00Z'
      } as any; // Usar 'as any' para evitar errores de tipo

      await bookingAPI.register(registerPayload);
      alert('¡Registro exitoso! Ahora puedes iniciar sesión.');
      setIsLoginMode(true); // Cambiar a modo login
      setRegisterData({
        email: '',
        password: '',
        first_name: '',
        last_name: '',
        phone: '',
        date_of_birth: '1990-01-01'
      }); // Limpiar formulario
    } catch (err) {
      alert('Error en registro. El email puede ya estar registrado.');
      console.error('Register error:', err);
    }
  };

  const handleBooking = async () => {
    if (!isLoggedIn) {
      setShowLoginForm(true);
      return;
    }

    if (!checkin || !checkout) {
      alert('Necesitas fechas de check-in y check-out para hacer una reserva');
      return;
    }

    try {
      setBookingLoading(true);
      
      const bookingData = {
        hotel_id: id!,
        check_in_date: new Date(checkin).toISOString(),
        check_out_date: new Date(checkout).toISOString(),
        guests,
        room_type: 'Standard',
        special_requests: 'Reserva desde frontend'
      };

      const response = await bookingAPI.createBooking(bookingData);
      
      // Navegar a confirmación con datos de la reserva
      const booking = response;
      navigate(`/confirmation?booking_id=${booking.id}&reference=${booking.booking_reference}`);
    } catch (err) {
      alert('Error creando la reserva. Por favor, intenta de nuevo.');
      console.error('Booking error:', err);
    } finally {
      setBookingLoading(false);
    }
  };

  const handleBack = () => {
    if (city) {
      const queryParams = new URLSearchParams({
        city,
        checkin,
        checkout,
        guests: guests.toString()
      });
      navigate(`/results?${queryParams.toString()}`);
    } else {
      navigate('/');
    }
  };

  const resetForms = () => {
    setLoginData({ email: '', password: '' });
    setRegisterData({
      email: '',
      password: '',
      first_name: '',
      last_name: '',
      phone: '',
      date_of_birth: '1990-01-01'
    });
    setIsLoginMode(true);
  };

  if (loading) {
    return (
      <div style={{ 
        display: 'flex', 
        justifyContent: 'center', 
        alignItems: 'center', 
        height: '50vh',
        fontSize: '18px'
      }}>
        🏨 Cargando detalles del hotel...
      </div>
    );
  }

  if (error || !hotel) {
    return (
      <div style={{ padding: '20px', maxWidth: '800px', margin: '0 auto' }}>
        <div style={{ 
          backgroundColor: '#ffebee', 
          color: '#c62828', 
          padding: '20px', 
          borderRadius: '8px',
          textAlign: 'center'
        }}>
          <h3>❌ {error || 'Hotel no encontrado'}</h3>
          <button 
            onClick={handleBack}
            style={{
              backgroundColor: '#1976d2',
              color: 'white',
              padding: '10px 20px',
              border: 'none',
              borderRadius: '5px',
              cursor: 'pointer',
              marginTop: '10px'
            }}
          >
            ← Volver
          </button>
        </div>
      </div>
    );
  }

  return (
    <div style={{ padding: '20px', maxWidth: '1200px', margin: '0 auto' }}>
      {/* Navegación */}
      <button 
        onClick={handleBack}
        style={{
          backgroundColor: '#f5f5f5',
          border: '1px solid #ddd',
          padding: '10px 20px',
          borderRadius: '5px',
          cursor: 'pointer',
          marginBottom: '20px'
        }}
      >
        ← Volver a resultados
      </button>

      {/* Información del hotel */}
      <div style={{ 
        backgroundColor: 'white', 
        borderRadius: '10px', 
        overflow: 'hidden',
        boxShadow: '0 4px 12px rgba(0,0,0,0.1)',
        marginBottom: '20px'
      }}>
        {/* Imagen principal */}
        <div style={{ 
          height: '400px', 
          backgroundColor: '#e3f2fd',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          fontSize: '72px'
        }}>
          {hotel.thumbnail ? (
            <img 
              src={hotel.thumbnail} 
              alt={hotel.name}
              style={{ width: '100%', height: '100%', objectFit: 'cover' }}
            />
          ) : '🏨'}
        </div>

        {/* Información del hotel */}
        <div style={{ padding: '30px' }}>
          <div style={{ display: 'grid', gap: '30px', gridTemplateColumns: '2fr 1fr' }}>
            {/* Columna izquierda - Información principal */}
            <div>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: '20px' }}>
                <h1 style={{ margin: 0, color: '#1976d2' }}>{hotel.name}</h1>
                <div style={{ textAlign: 'right' }}>
                  <div style={{ color: '#ff9800', fontWeight: 'bold', fontSize: '20px' }}>
                    {'⭐'.repeat(Math.floor(hotel.rating))}
                  </div>
                  <div style={{ color: '#666' }}>{hotel.rating}/5</div>
                </div>
              </div>

              <p style={{ color: '#666', fontSize: '16px', lineHeight: '1.6', marginBottom: '20px' }}>
                {hotel.description}
              </p>

              <div style={{ marginBottom: '20px' }}>
                <h3>📍 Ubicación</h3>
                <p style={{ color: '#666' }}>
                  {hotel.address}, {hotel.city}
                </p>
              </div>

              {/* Amenities */}
              <div style={{ marginBottom: '20px' }}>
                <h3>✨ Amenidades</h3>
                <div style={{ display: 'flex', flexWrap: 'wrap', gap: '10px' }}>
                  {hotel.amenities.map((amenity, index) => (
                    <span 
                      key={index}
                      style={{ 
                        backgroundColor: '#e3f2fd', 
                        padding: '8px 16px', 
                        borderRadius: '20px',
                        color: '#1976d2',
                        fontWeight: 'bold'
                      }}
                    >
                      {amenity}
                    </span>
                  ))}
                </div>
              </div>

              {/* Contacto */}
              <div>
                <h3>📞 Contacto</h3>
                <div style={{ color: '#666' }}>
                  {hotel.contact.phone && <div>Teléfono: {hotel.contact.phone}</div>}
                  {hotel.contact.email && <div>Email: {hotel.contact.email}</div>}
                  {hotel.contact.website && (
                    <div>
                      Sitio web: <a href={hotel.contact.website} target="_blank" rel="noopener noreferrer">
                        {hotel.contact.website}
                      </a>
                    </div>
                  )}
                </div>
              </div>
            </div>

            {/* Columna derecha - Reserva */}
            <div style={{ 
              backgroundColor: '#f8f9fa', 
              padding: '25px', 
              borderRadius: '10px',
              height: 'fit-content'
            }}>
              <h3 style={{ marginTop: 0 }}>💰 Precio por noche</h3>
              <div style={{ fontSize: '24px', fontWeight: 'bold', color: '#2e7d32', marginBottom: '15px' }}>
                ${hotel.price_range.min_price.toLocaleString()} - ${hotel.price_range.max_price.toLocaleString()} {hotel.price_range.currency}
              </div>

              {/* Información de búsqueda */}
              {checkin && checkout && (
                <div style={{ marginBottom: '20px', padding: '15px', backgroundColor: 'white', borderRadius: '8px' }}>
                  <h4 style={{ margin: '0 0 10px 0' }}>📅 Tu búsqueda:</h4>
                  <div style={{ fontSize: '14px', color: '#666' }}>
                    <div>Check-in: {checkin}</div>
                    <div>Check-out: {checkout}</div>
                    <div>Huéspedes: {guests}</div>
                  </div>
                </div>
              )}

              {/* Disponibilidad */}
              {availability && (
                <div style={{ 
                  marginBottom: '20px', 
                  padding: '15px', 
                  backgroundColor: availability.available ? '#e8f5e8' : '#ffebee',
                  borderRadius: '8px',
                  color: availability.available ? '#2e7d32' : '#c62828'
                }}>
                  {availability.available ? '✅ Disponible' : '❌ No disponible'}
                  {availability.price && (
                    <div style={{ fontSize: '14px', marginTop: '5px' }}>
                      Precio: ${availability.price} {availability.currency}
                    </div>
                  )}
                </div>
              )}

              {/* Login/Register Form */}
              {showLoginForm && (
                <div style={{ 
                  marginBottom: '20px', 
                  padding: '15px', 
                  backgroundColor: 'white', 
                  borderRadius: '8px',
                  border: '1px solid #ddd'
                }}>
                  {/* Toggle entre Login y Registro */}
                  <div style={{ display: 'flex', marginBottom: '15px' }}>
                    <button
                      onClick={() => setIsLoginMode(true)}
                      style={{
                        flex: 1,
                        padding: '8px',
                        backgroundColor: isLoginMode ? '#1976d2' : '#f5f5f5',
                        color: isLoginMode ? 'white' : '#666',
                        border: '1px solid #ddd',
                        borderRadius: '4px 0 0 4px',
                        cursor: 'pointer'
                      }}
                    >
                      🔐 Iniciar Sesión
                    </button>
                    <button
                      onClick={() => setIsLoginMode(false)}
                      style={{
                        flex: 1,
                        padding: '8px',
                        backgroundColor: !isLoginMode ? '#1976d2' : '#f5f5f5',
                        color: !isLoginMode ? 'white' : '#666',
                        border: '1px solid #ddd',
                        borderRadius: '0 4px 4px 0',
                        cursor: 'pointer'
                      }}
                    >
                      ✏️ Registrarse
                    </button>
                  </div>

                  {/* Formulario de Login */}
                  {isLoginMode ? (
                    <div>
                      <h4 style={{ margin: '0 0 15px 0' }}>🔐 Iniciar Sesión</h4>
                      <input
                        type="email"
                        placeholder="Email"
                        value={loginData.email}
                        onChange={(e) => setLoginData(prev => ({ ...prev, email: e.target.value }))}
                        style={{ 
                          width: '100%', 
                          padding: '8px', 
                          marginBottom: '10px', 
                          borderRadius: '4px',
                          border: '1px solid #ccc',
                          boxSizing: 'border-box'
                        }}
                      />
                      <input
                        type="password"
                        placeholder="Contraseña"
                        value={loginData.password}
                        onChange={(e) => setLoginData(prev => ({ ...prev, password: e.target.value }))}
                        style={{ 
                          width: '100%', 
                          padding: '8px', 
                          marginBottom: '15px', 
                          borderRadius: '4px',
                          border: '1px solid #ccc',
                          boxSizing: 'border-box'
                        }}
                      />
                      <button
                        onClick={handleLogin}
                        style={{
                          width: '100%',
                          backgroundColor: '#1976d2',
                          color: 'white',
                          border: 'none',
                          padding: '10px',
                          borderRadius: '4px',
                          cursor: 'pointer',
                          marginBottom: '10px'
                        }}
                      >
                        🔑 Entrar
                      </button>
                      <div style={{ fontSize: '12px', color: '#666', textAlign: 'center' }}>
                        ¿No tienes cuenta? Haz clic en "Registrarse"
                      </div>
                    </div>
                  ) : (
                    /* Formulario de Registro */
                    <div>
                      <h4 style={{ margin: '0 0 15px 0' }}>✏️ Crear Cuenta</h4>
                      <input
                        type="text"
                        placeholder="Nombre *"
                        value={registerData.first_name}
                        onChange={(e) => setRegisterData(prev => ({ ...prev, first_name: e.target.value }))}
                        style={{ 
                          width: '100%', 
                          padding: '8px', 
                          marginBottom: '8px', 
                          borderRadius: '4px',
                          border: '1px solid #ccc',
                          boxSizing: 'border-box'
                        }}
                      />
                      <input
                        type="text"
                        placeholder="Apellido *"
                        value={registerData.last_name}
                        onChange={(e) => setRegisterData(prev => ({ ...prev, last_name: e.target.value }))}
                        style={{ 
                          width: '100%', 
                          padding: '8px', 
                          marginBottom: '8px', 
                          borderRadius: '4px',
                          border: '1px solid #ccc',
                          boxSizing: 'border-box'
                        }}
                      />
                      <input
                        type="email"
                        placeholder="Email *"
                        value={registerData.email}
                        onChange={(e) => setRegisterData(prev => ({ ...prev, email: e.target.value }))}
                        style={{ 
                          width: '100%', 
                          padding: '8px', 
                          marginBottom: '8px', 
                          borderRadius: '4px',
                          border: '1px solid #ccc',
                          boxSizing: 'border-box'
                        }}
                      />
                      <input
                        type="password"
                        placeholder="Contraseña *"
                        value={registerData.password}
                        onChange={(e) => setRegisterData(prev => ({ ...prev, password: e.target.value }))}
                        style={{ 
                          width: '100%', 
                          padding: '8px', 
                          marginBottom: '8px', 
                          borderRadius: '4px',
                          border: '1px solid #ccc',
                          boxSizing: 'border-box'
                        }}
                      />
                      <input
                        type="tel"
                        placeholder="Teléfono (opcional)"
                        value={registerData.phone}
                        onChange={(e) => setRegisterData(prev => ({ ...prev, phone: e.target.value }))}
                        style={{ 
                          width: '100%', 
                          padding: '8px', 
                          marginBottom: '15px', 
                          borderRadius: '4px',
                          border: '1px solid #ccc',
                          boxSizing: 'border-box'
                        }}
                      />
                      <button
                        onClick={handleRegister}
                        style={{
                          width: '100%',
                          backgroundColor: '#2e7d32',
                          color: 'white',
                          border: 'none',
                          padding: '10px',
                          borderRadius: '4px',
                          cursor: 'pointer',
                          marginBottom: '10px'
                        }}
                      >
                        ✅ Crear Cuenta
                      </button>
                      <div style={{ fontSize: '12px', color: '#666', textAlign: 'center' }}>
                        * Campos obligatorios
                      </div>
                    </div>
                  )}

                  {/* Botón Cancelar */}
                  <button
                    onClick={() => {
                      setShowLoginForm(false);
                      resetForms();
                    }}
                    style={{
                      width: '100%',
                      backgroundColor: '#f5f5f5',
                      border: '1px solid #ddd',
                      padding: '8px',
                      borderRadius: '4px',
                      cursor: 'pointer',
                      marginTop: '10px'
                    }}
                  >
                    ❌ Cancelar
                  </button>
                </div>
              )}

              {/* Botón de reserva */}
              <button
                onClick={handleBooking}
                disabled={bookingLoading || (availability && !availability.available)}
                style={{
                  width: '100%',
                  backgroundColor: (availability && !availability.available) ? '#ccc' : '#ff5722',
                  color: 'white',
                  border: 'none',
                  padding: '15px',
                  borderRadius: '8px',
                  cursor: (availability && !availability.available) ? 'not-allowed' : 'pointer',
                  fontSize: '16px',
                  fontWeight: 'bold'
                }}
              >
                {bookingLoading ? '⏳ Reservando...' : 
                 (availability && !availability.available) ? '❌ No Disponible' : 
                 isLoggedIn ? '🎯 Reservar Ahora' : '🔐 Iniciar Sesión para Reservar'}
              </button>

              {!isLoggedIn && (
                <div style={{ fontSize: '12px', color: '#666', textAlign: 'center', marginTop: '10px' }}>
                  Necesitas estar logueado para hacer una reserva
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Detail;