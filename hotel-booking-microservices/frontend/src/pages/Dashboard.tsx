import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { bookingAPI } from '../services/api';

interface Booking {
  id: number;
  internal_hotel_id: string;
  hotel_name?: string;
  check_in_date: string;
  check_out_date: string;
  guests: number;
  room_type: string;
  total_price: number;
  currency: string;
  status: string;
  booking_reference: string;
  special_requests?: string;
  created_at: string;
}

const Dashboard: React.FC = () => {
  const navigate = useNavigate();
  const [searchData, setSearchData] = useState({
    city: '',
    checkin: '',
    checkout: '',
    guests: 2
  });

  const [userName, setUserName] = useState('');
  const [userRole, setUserRole] = useState('');
  const [showBookings, setShowBookings] = useState(false);
  const [bookings, setBookings] = useState<Booking[]>([]);
  const [loadingBookings, setLoadingBookings] = useState(false);

  useEffect(() => {
    // Verificar autenticación
    const token = localStorage.getItem('token');
    if (!token) {
      navigate('/');
      return;
    }

    // Obtener datos del usuario
    const name = localStorage.getItem('userName') || 'Usuario';
    const role = localStorage.getItem('userRole') || 'user';
    setUserName(name);
    setUserRole(role);
  }, [navigate]);

  const loadUserBookings = async () => {
    try {
      setLoadingBookings(true);
      // Hacer petición al booking-service para obtener reservas del usuario
      const response = await fetch('http://localhost:8003/api/bookings/my-bookings', {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      });
      
      if (response.ok) {
        const data = await response.json();
        setBookings(data.data || []); // Cambié de data.bookings a data.data
      } else {
        console.error('Error fetching bookings:', response.status);
        setBookings([]);
      }
    } catch (error) {
      console.error('Error loading bookings:', error);
      setBookings([]);
    } finally {
      setLoadingBookings(false);
    }
  };

  const handleShowBookings = () => {
    setShowBookings(true);
    loadUserBookings();
  };

  const getStatusColor = (status: string) => {
    switch (status.toLowerCase()) {
      case 'confirmed':
      case 'confirmada':
        return '#2e7d32';
      case 'pending':
      case 'pendiente':
        return '#ff9800';
      case 'cancelled':
      case 'cancelada':
        return '#d32f2f';
      default:
        return '#666';
    }
  };

  const getStatusText = (status: string) => {
    switch (status.toLowerCase()) {
      case 'confirmed':
        return '✅ Confirmada';
      case 'pending':
        return '⏳ Pendiente';
      case 'cancelled':
        return '❌ Cancelada';
      default:
        return status;
    }
  };

  const handleSearch = () => {
    if (!searchData.city || !searchData.checkin || !searchData.checkout) {
      alert('Por favor completa todos los campos obligatorios');
      return;
    }

    if (new Date(searchData.checkin) >= new Date(searchData.checkout)) {
      alert('La fecha de check-out debe ser posterior al check-in');
      return;
    }

    // Navegar a resultados con parámetros
    const queryParams = new URLSearchParams({
      city: searchData.city,
      checkin: searchData.checkin,
      checkout: searchData.checkout,
      guests: searchData.guests.toString()
    });
    
    navigate(`/results?${queryParams.toString()}`);
  };

  const handleLogout = () => {
    localStorage.clear();
    navigate('/');
  };

  const handleAdminPanel = () => {
    navigate('/admin');
  };

  if (showBookings) {
    return (
      <div style={{ minHeight: '100vh', backgroundColor: '#f5f5f5' }}>
        
        {/* Header */}
        <header style={{
          backgroundColor: 'white',
          boxShadow: '0 2px 4px rgba(0,0,0,0.1)',
          padding: '15px 20px',
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center'
        }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '15px' }}>
            <span style={{ fontSize: '24px' }}>🏨</span>
            <h1 style={{ margin: 0, color: '#1976d2' }}>Mis Reservas</h1>
          </div>
          
          <div style={{ display: 'flex', alignItems: 'center', gap: '15px' }}>
            <span style={{ color: '#666' }}>
              👋 Hola, <strong>{userName}</strong>
            </span>
            
            <button
              onClick={() => setShowBookings(false)}
              style={{
                padding: '8px 16px',
                backgroundColor: '#1976d2',
                color: 'white',
                border: 'none',
                borderRadius: '6px',
                cursor: 'pointer',
                fontSize: '14px'
              }}
            >
              🏠 Volver al Dashboard
            </button>
            
            <button
              onClick={handleLogout}
              style={{
                padding: '8px 16px',
                backgroundColor: '#f44336',
                color: 'white',
                border: 'none',
                borderRadius: '6px',
                cursor: 'pointer',
                fontSize: '14px'
              }}
            >
              🚪 Salir
            </button>
          </div>
        </header>

        {/* Contenido de reservas */}
        <div style={{ maxWidth: '1200px', margin: '0 auto', padding: '40px 20px' }}>
          
          {/* Estadísticas de reservas */}
          <div style={{
            display: 'grid',
            gridTemplateColumns: '1fr 1fr 1fr',
            gap: '20px',
            marginBottom: '40px'
          }}>
            <div style={{
              backgroundColor: 'white',
              padding: '25px',
              borderRadius: '10px',
              textAlign: 'center',
              boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
            }}>
              <div style={{ fontSize: '32px', color: '#1976d2', marginBottom: '10px' }}>📋</div>
              <div style={{ fontSize: '24px', fontWeight: 'bold', color: '#333' }}>{bookings.length}</div>
              <div style={{ color: '#666' }}>Total Reservas</div>
            </div>
            
            <div style={{
              backgroundColor: 'white',
              padding: '25px',
              borderRadius: '10px',
              textAlign: 'center',
              boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
            }}>
              <div style={{ fontSize: '32px', color: '#2e7d32', marginBottom: '10px' }}>✅</div>
              <div style={{ fontSize: '24px', fontWeight: 'bold', color: '#333' }}>
                {bookings.filter(b => b.status.toLowerCase() === 'confirmed').length}
              </div>
              <div style={{ color: '#666' }}>Confirmadas</div>
            </div>
            
            <div style={{
              backgroundColor: 'white',
              padding: '25px',
              borderRadius: '10px',
              textAlign: 'center',
              boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
            }}>
              <div style={{ fontSize: '32px', color: '#ff9800', marginBottom: '10px' }}>⏳</div>
              <div style={{ fontSize: '24px', fontWeight: 'bold', color: '#333' }}>
                {bookings.filter(b => b.status.toLowerCase() === 'pending').length}
              </div>
              <div style={{ color: '#666' }}>Pendientes</div>
            </div>
          </div>

          {/* Lista de reservas */}
          <div style={{
            backgroundColor: 'white',
            borderRadius: '10px',
            boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
            overflow: 'hidden'
          }}>
            <div style={{
              padding: '20px',
              backgroundColor: '#f8f9fa',
              borderBottom: '1px solid #e0e0e0'
            }}>
              <h3 style={{ margin: 0, color: '#333' }}>📅 Historial de Reservas</h3>
              <button
                onClick={loadUserBookings}
                style={{
                  marginTop: '10px',
                  padding: '8px 16px',
                  backgroundColor: '#1976d2',
                  color: 'white',
                  border: 'none',
                  borderRadius: '6px',
                  cursor: 'pointer',
                  fontSize: '14px'
                }}
              >
                🔄 Actualizar
              </button>
            </div>

            {loadingBookings ? (
              <div style={{ padding: '60px', textAlign: 'center', color: '#666' }}>
                ⏳ Cargando reservas...
              </div>
            ) : bookings.length === 0 ? (
              <div style={{ padding: '60px', textAlign: 'center', color: '#666' }}>
                <div style={{ fontSize: '48px', marginBottom: '20px' }}>🏨</div>
                <h3>No tienes reservas aún</h3>
                <p>Cuando hagas tu primera reserva, aparecerá aquí.</p>
                <button
                  onClick={() => setShowBookings(false)}
                  style={{
                    marginTop: '20px',
                    padding: '12px 24px',
                    backgroundColor: '#1976d2',
                    color: 'white',
                    border: 'none',
                    borderRadius: '8px',
                    cursor: 'pointer',
                    fontSize: '16px'
                  }}
                >
                  🔍 Buscar Hoteles
                </button>
              </div>
            ) : (
              <div style={{ maxHeight: '600px', overflowY: 'auto' }}>
                {bookings.map((booking) => (
                  <div
                    key={booking.id}
                    style={{
                      padding: '20px',
                      borderBottom: '1px solid #f0f0f0',
                      display: 'flex',
                      justifyContent: 'space-between',
                      alignItems: 'center'
                    }}
                  >
                    <div style={{ flex: 1 }}>
                      <div style={{ display: 'flex', alignItems: 'center', gap: '15px', marginBottom: '10px' }}>
                        <h4 style={{ margin: 0, color: '#333' }}>
                          🏨 {booking.hotel_name || `Hotel ID: ${booking.internal_hotel_id}`}
                        </h4>
                        <span 
                          style={{ 
                            backgroundColor: getStatusColor(booking.status),
                            color: 'white',
                            padding: '4px 12px',
                            borderRadius: '12px',
                            fontSize: '12px',
                            fontWeight: 'bold'
                          }}
                        >
                          {getStatusText(booking.status)}
                        </span>
                      </div>
                      
                      <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr 1fr', gap: '20px', marginBottom: '10px' }}>
                        <div>
                          <div style={{ color: '#666', fontSize: '12px' }}>📅 Check-in</div>
                          <div style={{ fontWeight: 'bold' }}>
                            {new Date(booking.check_in_date).toLocaleDateString()}
                          </div>
                        </div>
                        <div>
                          <div style={{ color: '#666', fontSize: '12px' }}>📅 Check-out</div>
                          <div style={{ fontWeight: 'bold' }}>
                            {new Date(booking.check_out_date).toLocaleDateString()}
                          </div>
                        </div>
                        <div>
                          <div style={{ color: '#666', fontSize: '12px' }}>👥 Huéspedes</div>
                          <div style={{ fontWeight: 'bold' }}>{booking.guests}</div>
                        </div>
                      </div>

                      <div style={{ display: 'flex', gap: '20px', fontSize: '14px', color: '#666' }}>
                        <span><strong>🛏️ Habitación:</strong> {booking.room_type}</span>
                        <span><strong>🎫 Referencia:</strong> {booking.booking_reference}</span>
                        <span><strong>📅 Reservado:</strong> {new Date(booking.created_at).toLocaleDateString()}</span>
                      </div>

                      {booking.special_requests && (
                        <div style={{ marginTop: '10px', fontSize: '14px', color: '#666' }}>
                          <strong>📝 Solicitudes especiales:</strong> {booking.special_requests}
                        </div>
                      )}
                    </div>
                    
                    <div style={{ textAlign: 'right', marginLeft: '20px' }}>
                      <div style={{ fontSize: '20px', fontWeight: 'bold', color: '#2e7d32' }}>
                        ${booking.total_price?.toLocaleString() || 'N/A'} {booking.currency}
                      </div>
                      <div style={{ fontSize: '12px', color: '#666', marginTop: '5px' }}>
                        Total de la reserva
                      </div>
                      
                      {booking.status.toLowerCase() === 'confirmed' && (
                        <div style={{ 
                          marginTop: '10px',
                          padding: '8px 12px',
                          backgroundColor: '#e8f5e8',
                          borderRadius: '6px',
                          fontSize: '12px',
                          color: '#2e7d32'
                        }}>
                          ✅ Reserva confirmada
                        </div>
                      )}
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
      </div>
    );
  }

  return (
    <div style={{ minHeight: '100vh', backgroundColor: '#f5f5f5' }}>
      
      {/* Header */}
      <header style={{
        backgroundColor: 'white',
        boxShadow: '0 2px 4px rgba(0,0,0,0.1)',
        padding: '15px 20px',
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center'
      }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '15px' }}>
          <span style={{ fontSize: '24px' }}>🏨</span>
          <h1 style={{ margin: 0, color: '#1976d2' }}>Hotel Manager</h1>
        </div>
        
        <div style={{ display: 'flex', alignItems: 'center', gap: '15px' }}>
          <span style={{ color: '#666' }}>
            👋 Hola, <strong>{userName}</strong>
            {userRole === 'admin' && <span style={{ 
              backgroundColor: '#ff9800', 
              color: 'white', 
              padding: '2px 8px', 
              borderRadius: '12px', 
              fontSize: '12px',
              marginLeft: '8px'
            }}>ADMIN</span>}
          </span>
          
          {/* Botón Mis Reservas */}
          <button
            onClick={handleShowBookings}
            style={{
              padding: '8px 16px',
              backgroundColor: '#2e7d32',
              color: 'white',
              border: 'none',
              borderRadius: '6px',
              cursor: 'pointer',
              fontSize: '14px'
            }}
          >
            📋 Mis Reservas
          </button>
          
          {userRole === 'admin' && (
            <button
              onClick={handleAdminPanel}
              style={{
                padding: '8px 16px',
                backgroundColor: '#ff9800',
                color: 'white',
                border: 'none',
                borderRadius: '6px',
                cursor: 'pointer',
                fontSize: '14px'
              }}
            >
              🛠️ Panel Admin
            </button>
          )}
          
          <button
            onClick={handleLogout}
            style={{
              padding: '8px 16px',
              backgroundColor: '#f44336',
              color: 'white',
              border: 'none',
              borderRadius: '6px',
              cursor: 'pointer',
              fontSize: '14px'
            }}
          >
            🚪 Salir
          </button>
        </div>
      </header>

      {/* Contenido principal */}
      <div style={{
        maxWidth: '1200px',
        margin: '0 auto',
        padding: '60px 20px',
        textAlign: 'center'
      }}>
        
        {/* Título y descripción */}
        <div style={{ marginBottom: '60px' }}>
          <h1 style={{ 
            fontSize: '48px', 
            margin: '0 0 20px 0', 
            color: '#333',
            fontWeight: '300'
          }}>
            Encuentra tu hotel perfecto
          </h1>
          <p style={{ 
            fontSize: '20px', 
            color: '#666', 
            maxWidth: '600px',
            margin: '0 auto'
          }}>
            Busca entre miles de hoteles en todo el mundo y encuentra la mejor opción para tu viaje
          </p>
        </div>

        {/* Formulario de búsqueda */}
        <div style={{
          backgroundColor: 'white',
          borderRadius: '20px',
          padding: '40px',
          boxShadow: '0 10px 30px rgba(0,0,0,0.1)',
          maxWidth: '800px',
          margin: '0 auto'
        }}>
          
          <div style={{
            display: 'grid',
            gridTemplateColumns: '2fr 1fr 1fr 1fr',
            gap: '20px',
            marginBottom: '30px'
          }}>
            
            {/* Ciudad */}
            <div style={{ textAlign: 'left' }}>
              <label style={{ 
                display: 'block', 
                marginBottom: '8px', 
                fontWeight: 'bold',
                color: '#333'
              }}>
                🏙️ Ciudad de destino
              </label>
              <input
                type="text"
                placeholder="ej. Barcelona, Madrid, Córdoba..."
                value={searchData.city}
                onChange={(e) => setSearchData(prev => ({ ...prev, city: e.target.value }))}
                style={{
                  width: '100%',
                  padding: '15px',
                  border: '2px solid #e0e0e0',
                  borderRadius: '10px',
                  fontSize: '16px',
                  boxSizing: 'border-box',
                  outline: 'none'
                }}
                onFocus={(e) => e.target.style.borderColor = '#1976d2'}
                onBlur={(e) => e.target.style.borderColor = '#e0e0e0'}
              />
            </div>

            {/* Check-in */}
            <div style={{ textAlign: 'left' }}>
              <label style={{ 
                display: 'block', 
                marginBottom: '8px', 
                fontWeight: 'bold',
                color: '#333'
              }}>
                📅 Check-in
              </label>
              <input
                type="date"
                value={searchData.checkin}
                onChange={(e) => setSearchData(prev => ({ ...prev, checkin: e.target.value }))}
                min={new Date().toISOString().split('T')[0]}
                style={{
                  width: '100%',
                  padding: '15px',
                  border: '2px solid #e0e0e0',
                  borderRadius: '10px',
                  fontSize: '16px',
                  boxSizing: 'border-box',
                  outline: 'none'
                }}
                onFocus={(e) => e.target.style.borderColor = '#1976d2'}
                onBlur={(e) => e.target.style.borderColor = '#e0e0e0'}
              />
            </div>

            {/* Check-out */}
            <div style={{ textAlign: 'left' }}>
              <label style={{ 
                display: 'block', 
                marginBottom: '8px', 
                fontWeight: 'bold',
                color: '#333'
              }}>
                📅 Check-out
              </label>
              <input
                type="date"
                value={searchData.checkout}
                onChange={(e) => setSearchData(prev => ({ ...prev, checkout: e.target.value }))}
                min={searchData.checkin || new Date().toISOString().split('T')[0]}
                style={{
                  width: '100%',
                  padding: '15px',
                  border: '2px solid #e0e0e0',
                  borderRadius: '10px',
                  fontSize: '16px',
                  boxSizing: 'border-box',
                  outline: 'none'
                }}
                onFocus={(e) => e.target.style.borderColor = '#1976d2'}
                onBlur={(e) => e.target.style.borderColor = '#e0e0e0'}
              />
            </div>

            {/* Huéspedes */}
            <div style={{ textAlign: 'left' }}>
              <label style={{ 
                display: 'block', 
                marginBottom: '8px', 
                fontWeight: 'bold',
                color: '#333'
              }}>
                👥 Huéspedes
              </label>
              <select
                value={searchData.guests}
                onChange={(e) => setSearchData(prev => ({ ...prev, guests: parseInt(e.target.value) }))}
                style={{
                  width: '100%',
                  padding: '15px',
                  border: '2px solid #e0e0e0',
                  borderRadius: '10px',
                  fontSize: '16px',
                  boxSizing: 'border-box',
                  outline: 'none',
                  backgroundColor: 'white'
                }}
                onFocus={(e) => e.target.style.borderColor = '#1976d2'}
                onBlur={(e) => e.target.style.borderColor = '#e0e0e0'}
              >
                <option value={1}>1 huésped</option>
                <option value={2}>2 huéspedes</option>
                <option value={3}>3 huéspedes</option>
                <option value={4}>4 huéspedes</option>
                <option value={5}>5+ huéspedes</option>
              </select>
            </div>
          </div>

          {/* Botón de búsqueda */}
          <button
            onClick={handleSearch}
            style={{
              width: '100%',
              padding: '18px',
              backgroundColor: '#1976d2',
              color: 'white',
              border: 'none',
              borderRadius: '15px',
              fontSize: '20px',
              fontWeight: 'bold',
              cursor: 'pointer',
              transition: 'background-color 0.3s ease'
            }}
            onMouseOver={(e) => (e.target as HTMLButtonElement).style.backgroundColor = '#1565c0'}
            onMouseOut={(e) => (e.target as HTMLButtonElement).style.backgroundColor = '#1976d2'}
          >
            🔍 Buscar Hoteles
          </button>
        </div>

        {/* Información adicional */}
        <div style={{
          marginTop: '60px',
          display: 'grid',
          gridTemplateColumns: '1fr 1fr 1fr',
          gap: '40px',
          maxWidth: '900px',
          margin: '60px auto 0'
        }}>
          <div style={{ textAlign: 'center' }}>
            <div style={{ fontSize: '48px', marginBottom: '15px' }}>🌟</div>
            <h3 style={{ color: '#333', marginBottom: '10px' }}>Hoteles de calidad</h3>
            <p style={{ color: '#666', fontSize: '14px' }}>Más de 1000 hoteles verificados en todo el mundo</p>
          </div>
          
          <div style={{ textAlign: 'center' }}>
            <div style={{ fontSize: '48px', marginBottom: '15px' }}>💳</div>
            <h3 style={{ color: '#333', marginBottom: '10px' }}>Reserva segura</h3>
            <p style={{ color: '#666', fontSize: '14px' }}>Proceso de pago seguro y confirmación inmediata</p>
          </div>
          
          <div style={{ textAlign: 'center' }}>
            <div style={{ fontSize: '48px', marginBottom: '15px' }}>📞</div>
            <h3 style={{ color: '#333', marginBottom: '10px' }}>Soporte 24/7</h3>
            <p style={{ color: '#666', fontSize: '14px' }}>Atención al cliente disponible las 24 horas</p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;