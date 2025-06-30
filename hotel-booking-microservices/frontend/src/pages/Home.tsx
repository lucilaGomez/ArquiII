import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { SearchParams } from '../types';

const Home: React.FC = () => {
  const navigate = useNavigate();
  const [searchParams, setSearchParams] = useState<SearchParams>({
    city: '',
    checkin: '',
    checkout: '',
    guests: 2
  });
  const [error, setError] = useState<string>('');

  const handleSearch = () => {
    // Validaciones básicas
    if (!searchParams.city.trim()) {
      setError('La ciudad es obligatoria');
      return;
    }

    if (!searchParams.checkin || !searchParams.checkout) {
      setError('Las fechas son obligatorias');
      return;
    }

    if (new Date(searchParams.checkout) <= new Date(searchParams.checkin)) {
      setError('La fecha de check-out debe ser posterior al check-in');
      return;
    }

    setError('');

    // Navegar a resultados
    const queryParams = new URLSearchParams({
      city: searchParams.city,
      checkin: searchParams.checkin,
      checkout: searchParams.checkout,
      guests: searchParams.guests.toString()
    });

    navigate(`/results?${queryParams.toString()}`);
  };

  const handleInputChange = (field: keyof SearchParams) => (
    event: React.ChangeEvent<HTMLInputElement>
  ) => {
    setSearchParams(prev => ({
      ...prev,
      [field]: field === 'guests' ? parseInt(event.target.value) || 1 : event.target.value
    }));
    if (error) setError('');
  };

  const today = new Date().toISOString().split('T')[0];

  return (
    <div style={{ padding: '20px', maxWidth: '800px', margin: '0 auto' }}>
      <h1 style={{ textAlign: 'center', color: '#1976d2' }}>🏨 Hotel Booking</h1>
      <h2 style={{ textAlign: 'center', color: '#666' }}>
        Encuentra el hotel perfecto para tu próximo viaje
      </h2>

      <div style={{ 
        backgroundColor: '#f5f5f5', 
        padding: '30px', 
        borderRadius: '10px', 
        marginTop: '30px' 
      }}>
        <h3 style={{ textAlign: 'center', marginBottom: '30px' }}>
          ✈️ Buscar Hoteles
        </h3>

        {error && (
          <div style={{ 
            backgroundColor: '#ffebee', 
            color: '#c62828', 
            padding: '10px', 
            borderRadius: '5px', 
            marginBottom: '20px' 
          }}>
            {error}
          </div>
        )}

        <div style={{ display: 'grid', gap: '20px', gridTemplateColumns: 'repeat(auto-fit, minmax(250px, 1fr))' }}>
          <div>
            <label style={{ display: 'block', marginBottom: '5px', fontWeight: 'bold' }}>
              🏙️ Ciudad de destino:
            </label>
            <input
              type="text"
              value={searchParams.city}
              onChange={handleInputChange('city')}
              placeholder="Ej: Córdoba, Buenos Aires, Mendoza..."
              style={{ 
                width: '100%', 
                padding: '10px', 
                borderRadius: '5px', 
                border: '1px solid #ccc',
                fontSize: '16px'
              }}
              required
            />
          </div>

          <div>
            <label style={{ display: 'block', marginBottom: '5px', fontWeight: 'bold' }}>
              👥 Número de huéspedes:
            </label>
            <input
              type="number"
              value={searchParams.guests}
              onChange={handleInputChange('guests')}
              min="1"
              max="10"
              style={{ 
                width: '100%', 
                padding: '10px', 
                borderRadius: '5px', 
                border: '1px solid #ccc',
                fontSize: '16px'
              }}
              required
            />
          </div>

          <div>
            <label style={{ display: 'block', marginBottom: '5px', fontWeight: 'bold' }}>
              📅 Fecha de check-in:
            </label>
            <input
              type="date"
              value={searchParams.checkin}
              onChange={handleInputChange('checkin')}
              min={today}
              style={{ 
                width: '100%', 
                padding: '10px', 
                borderRadius: '5px', 
                border: '1px solid #ccc',
                fontSize: '16px'
              }}
              required
            />
          </div>

          <div>
            <label style={{ display: 'block', marginBottom: '5px', fontWeight: 'bold' }}>
              📅 Fecha de check-out:
            </label>
            <input
              type="date"
              value={searchParams.checkout}
              onChange={handleInputChange('checkout')}
              min={searchParams.checkin || today}
              style={{ 
                width: '100%', 
                padding: '10px', 
                borderRadius: '5px', 
                border: '1px solid #ccc',
                fontSize: '16px'
              }}
              required
            />
          </div>
        </div>

        <div style={{ textAlign: 'center', marginTop: '30px' }}>
          <button
            onClick={handleSearch}
            style={{
              backgroundColor: '#1976d2',
              color: 'white',
              padding: '15px 40px',
              fontSize: '18px',
              border: 'none',
              borderRadius: '8px',
              cursor: 'pointer',
              fontWeight: 'bold'
            }}
          >
            🔍 Buscar Hoteles
          </button>
        </div>
      </div>

      <div style={{ 
        display: 'grid', 
        gap: '20px', 
        gridTemplateColumns: 'repeat(auto-fit, minmax(250px, 1fr))',
        marginTop: '40px'
      }}>
        <div style={{ 
          backgroundColor: 'white', 
          padding: '20px', 
          borderRadius: '10px', 
          textAlign: 'center',
          boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
        }}>
          <h3>🌟 Hoteles de Calidad</h3>
          <p>Trabajamos con los mejores hoteles de Argentina para garantizar tu comodidad.</p>
        </div>

        <div style={{ 
          backgroundColor: 'white', 
          padding: '20px', 
          borderRadius: '10px', 
          textAlign: 'center',
          boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
        }}>
          <h3>💰 Mejores Precios</h3>
          <p>Encuentra las mejores ofertas y precios competitivos para tu estadía.</p>
        </div>

        <div style={{ 
          backgroundColor: 'white', 
          padding: '20px', 
          borderRadius: '10px', 
          textAlign: 'center',
          boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
        }}>
          <h3>🔒 Reserva Segura</h3>
          <p>Tus datos y pagos están protegidos con la mejor tecnología de seguridad.</p>
        </div>
      </div>
    </div>
  );
};

export default Home;
