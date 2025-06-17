import React, { useState, useEffect } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { searchAPI } from '../services/api';
import { SearchResult } from '../types';

const Results: React.FC = () => {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const [hotels, setHotels] = useState<SearchResult[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string>('');

  // Obtener parámetros de búsqueda de la URL
  const city = searchParams.get('city') || '';
  const checkin = searchParams.get('checkin') || '';
  const checkout = searchParams.get('checkout') || '';
  const guests = parseInt(searchParams.get('guests') || '2');

  useEffect(() => {
    const searchHotels = async () => {
      if (!city) {
        navigate('/');
        return;
      }

      try {
        setLoading(true);
        const response = await searchAPI.searchHotels({
          city,
          checkin,
          checkout,
          guests
        });
        
        // Manejar tanto la respuesta del search-service básico como avanzado
        if (response.data?.hotels) {
          setHotels(response.data.hotels);
        } else if (Array.isArray(response.data)) {
          setHotels(response.data);
        } else {
          setHotels([]);
        }
      } catch (err) {
        setError('Error buscando hoteles. Por favor, intenta de nuevo.');
        console.error('Error searching hotels:', err);
      } finally {
        setLoading(false);
      }
    };

    searchHotels();
  }, [city, checkin, checkout, guests, navigate]);

  const handleSelectHotel = (hotelId: string) => {
    // Navegar a detalle del hotel pasando los parámetros de búsqueda
    const queryParams = new URLSearchParams({
      city,
      checkin,
      checkout,
      guests: guests.toString()
    });
    
    navigate(`/hotel/${hotelId}?${queryParams.toString()}`);
  };

  const handleNewSearch = () => {
    navigate('/');
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
        🔍 Buscando hoteles en {city}...
      </div>
    );
  }

  if (error) {
    return (
      <div style={{ padding: '20px', maxWidth: '800px', margin: '0 auto' }}>
        <div style={{ 
          backgroundColor: '#ffebee', 
          color: '#c62828', 
          padding: '20px', 
          borderRadius: '8px',
          textAlign: 'center'
        }}>
          <h3>❌ {error}</h3>
          <button 
            onClick={handleNewSearch}
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
            🔄 Nueva Búsqueda
          </button>
        </div>
      </div>
    );
  }

  return (
    <div style={{ padding: '20px', maxWidth: '1200px', margin: '0 auto' }}>
      {/* Header con información de búsqueda */}
      <div style={{ 
        backgroundColor: 'white', 
        padding: '20px', 
        borderRadius: '10px', 
        marginBottom: '20px',
        boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
      }}>
        <h1>🔍 Resultados de Búsqueda</h1>
        <div style={{ display: 'flex', gap: '20px', flexWrap: 'wrap', marginTop: '10px' }}>
          <span><strong>🏙️ Ciudad:</strong> {city}</span>
          <span><strong>📅 Check-in:</strong> {checkin}</span>
          <span><strong>📅 Check-out:</strong> {checkout}</span>
          <span><strong>👥 Huéspedes:</strong> {guests}</span>
        </div>
        <button 
          onClick={handleNewSearch}
          style={{
            backgroundColor: '#f5f5f5',
            border: '1px solid #ddd',
            padding: '8px 16px',
            borderRadius: '5px',
            cursor: 'pointer',
            marginTop: '15px'
          }}
        >
          🔄 Nueva Búsqueda
        </button>
      </div>

      {/* Resultados */}
      {hotels.length === 0 ? (
        <div style={{ 
          textAlign: 'center', 
          padding: '40px', 
          backgroundColor: 'white',
          borderRadius: '10px'
        }}>
          <h3>😔 No se encontraron hoteles</h3>
          <p>No hay hoteles disponibles en {city} para las fechas seleccionadas.</p>
          <button 
            onClick={handleNewSearch}
            style={{
              backgroundColor: '#1976d2',
              color: 'white',
              padding: '12px 24px',
              border: 'none',
              borderRadius: '5px',
              cursor: 'pointer',
              marginTop: '15px'
            }}
          >
            🔄 Probar otra búsqueda
          </button>
        </div>
      ) : (
        <>
          <h2>🏨 {hotels.length} hotel(es) encontrado(s) en {city}</h2>
          
          <div style={{ 
            display: 'grid', 
            gap: '20px', 
            gridTemplateColumns: 'repeat(auto-fill, minmax(400px, 1fr))',
            marginTop: '20px'
          }}>
            {hotels.map((hotel) => (
              <div 
                key={hotel.id}
                style={{ 
                  backgroundColor: 'white', 
                  borderRadius: '10px', 
                  overflow: 'hidden',
                  boxShadow: '0 4px 12px rgba(0,0,0,0.1)',
                  cursor: 'pointer',
                  transition: 'transform 0.2s'
                }}
                onClick={() => handleSelectHotel(hotel.id)}
                onMouseEnter={(e) => {
                  e.currentTarget.style.transform = 'translateY(-2px)';
                }}
                onMouseLeave={(e) => {
                  e.currentTarget.style.transform = 'translateY(0)';
                }}
              >
                {/* Imagen del hotel */}
                <div style={{ 
                  height: '200px', 
                  backgroundColor: '#e3f2fd',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  fontSize: '48px'
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
                <div style={{ padding: '20px' }}>
                  <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
                    <h3 style={{ margin: '0 0 10px 0', color: '#1976d2' }}>
                      {hotel.name}
                    </h3>
                    <div style={{ textAlign: 'right' }}>
                      <div style={{ color: '#ff9800', fontWeight: 'bold' }}>
                        {'⭐'.repeat(Math.floor(hotel.rating || 4))}
                      </div>
                      <div style={{ fontSize: '14px', color: '#666' }}>
                        {hotel.rating || 4}/5
                      </div>
                    </div>
                  </div>

                  <p style={{ 
                    color: '#666', 
                    margin: '10px 0',
                    display: '-webkit-box',
                    WebkitLineClamp: 2,
                    WebkitBoxOrient: 'vertical',
                    overflow: 'hidden'
                  }}>
                    {hotel.description}
                  </p>

                  {/* Amenities */}
                  {hotel.amenities && hotel.amenities.length > 0 && (
                    <div style={{ margin: '15px 0' }}>
                      <div style={{ fontSize: '14px', fontWeight: 'bold', marginBottom: '5px' }}>
                        ✨ Amenities:
                      </div>
                      <div style={{ display: 'flex', flexWrap: 'wrap', gap: '5px' }}>
                        {hotel.amenities.slice(0, 3).map((amenity, index) => (
                          <span 
                            key={index}
                            style={{ 
                              backgroundColor: '#e3f2fd', 
                              padding: '4px 8px', 
                              borderRadius: '12px',
                              fontSize: '12px',
                              color: '#1976d2'
                            }}
                          >
                            {amenity}
                          </span>
                        ))}
                        {hotel.amenities.length > 3 && (
                          <span style={{ fontSize: '12px', color: '#666' }}>
                            +{hotel.amenities.length - 3} más
                          </span>
                        )}
                      </div>
                    </div>
                  )}

                  {/* Precio y disponibilidad */}
                  <div style={{ 
                    display: 'flex', 
                    justifyContent: 'space-between', 
                    alignItems: 'center',
                    marginTop: '15px',
                    paddingTop: '15px',
                    borderTop: '1px solid #eee'
                  }}>
                    <div>
                      <div style={{ fontSize: '18px', fontWeight: 'bold', color: '#2e7d32' }}>
                        ${hotel.min_price?.toLocaleString() || 'N/A'} - ${hotel.max_price?.toLocaleString() || 'N/A'}
                      </div>
                      <div style={{ fontSize: '12px', color: '#666' }}>
                        {hotel.currency || 'ARS'} por noche
                      </div>
                    </div>
                    <div style={{ textAlign: 'right' }}>
                      <div style={{ 
                        color: hotel.available ? '#2e7d32' : '#d32f2f',
                        fontWeight: 'bold',
                        fontSize: '14px'
                      }}>
                        {hotel.available ? '✅ Disponible' : '❌ No disponible'}
                      </div>
                      <button
                        style={{
                          backgroundColor: '#1976d2',
                          color: 'white',
                          border: 'none',
                          padding: '8px 16px',
                          borderRadius: '5px',
                          cursor: 'pointer',
                          marginTop: '5px',
                          fontSize: '14px'
                        }}
                        onClick={(e) => {
                          e.stopPropagation();
                          handleSelectHotel(hotel.id);
                        }}
                      >
                        Ver Detalles
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </>
      )}
    </div>
  );
};

export default Results;
