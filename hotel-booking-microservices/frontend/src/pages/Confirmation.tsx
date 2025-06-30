import React from 'react';
import { useSearchParams, useNavigate } from 'react-router-dom';

const Confirmation: React.FC = () => {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  
  const bookingId = searchParams.get('booking_id');
  const reference = searchParams.get('reference');

  return (
    <div style={{ 
      padding: '40px 20px', 
      maxWidth: '800px', 
      margin: '0 auto',
      textAlign: 'center'
    }}>
      {/* Header de éxito */}
      <div style={{ 
        backgroundColor: '#e8f5e8', 
        padding: '40px', 
        borderRadius: '15px',
        marginBottom: '30px',
        border: '2px solid #4caf50'
      }}>
        <div style={{ fontSize: '60px', marginBottom: '20px' }}>🎉</div>
        <h1 style={{ 
          color: '#2e7d32', 
          margin: '0 0 15px 0',
          fontSize: '28px'
        }}>
          ¡Reserva Confirmada!
        </h1>
        <p style={{ 
          fontSize: '18px', 
          color: '#388e3c',
          margin: 0
        }}>
          Tu reserva se ha procesado exitosamente
        </p>
      </div>

      {/* Detalles de la reserva */}
      <div style={{ 
        backgroundColor: 'white', 
        padding: '30px', 
        borderRadius: '10px',
        boxShadow: '0 4px 12px rgba(0,0,0,0.1)',
        marginBottom: '30px',
        textAlign: 'left'
      }}>
        <h3 style={{ 
          margin: '0 0 20px 0', 
          color: '#1976d2',
          textAlign: 'center'
        }}>
          📋 Detalles de tu Reserva
        </h3>
        
        <div style={{ 
          display: 'grid', 
          gap: '15px',
          gridTemplateColumns: '1fr 1fr',
          marginBottom: '20px'
        }}>
          <div>
            <strong>🆔 ID de Reserva:</strong>
            <div style={{ 
              color: '#666', 
              fontFamily: 'monospace',
              backgroundColor: '#f5f5f5',
              padding: '8px',
              borderRadius: '4px',
              marginTop: '5px'
            }}>
              {bookingId || 'No disponible'}
            </div>
          </div>
          
          <div>
            <strong>📞 Referencia:</strong>
            <div style={{ 
              color: '#666',
              fontFamily: 'monospace',
              backgroundColor: '#f5f5f5',
              padding: '8px',
              borderRadius: '4px',
              marginTop: '5px'
            }}>
              {reference || 'No disponible'}
            </div>
          </div>
        </div>

        <div style={{ 
          backgroundColor: '#e3f2fd', 
          padding: '15px', 
          borderRadius: '8px',
          marginTop: '20px'
        }}>
          <p style={{ margin: 0, color: '#1565c0' }}>
            <strong>📧 Próximos pasos:</strong><br/>
            Recibirás un email de confirmación en los próximos minutos con todos los detalles de tu reserva.
          </p>
        </div>
      </div>

      {/* Acciones */}
      <div style={{ 
        display: 'flex', 
        gap: '15px', 
        justifyContent: 'center',
        flexWrap: 'wrap'
      }}>
        <button
          onClick={() => navigate('/')}
          style={{
            backgroundColor: '#1976d2',
            color: 'white',
            border: 'none',
            padding: '15px 30px',
            borderRadius: '8px',
            cursor: 'pointer',
            fontSize: '16px',
            fontWeight: 'bold'
          }}
        >
          🏠 Volver al Inicio
        </button>
        
        <button
          onClick={() => navigate('/results')}
          style={{
            backgroundColor: '#ff9800',
            color: 'white',
            border: 'none',
            padding: '15px 30px',
            borderRadius: '8px',
            cursor: 'pointer',
            fontSize: '16px',
            fontWeight: 'bold'
          }}
        >
          🔍 Buscar Más Hoteles
        </button>
      </div>

      {/* Información adicional */}
      <div style={{ 
        marginTop: '40px', 
        padding: '20px', 
        backgroundColor: '#f8f9fa',
        borderRadius: '8px',
        color: '#666'
      }}>
        <h4 style={{ margin: '0 0 10px 0' }}>ℹ️ Información Importante</h4>
        <ul style={{ textAlign: 'left', paddingLeft: '20px' }}>
          <li>Guarda este número de referencia para futuras consultas</li>
          <li>Revisa tu email para más detalles de la reserva</li>
          <li>Si necesitas modificar tu reserva, contacta al hotel directamente</li>
          <li>Presenta tu ID al momento del check-in</li>
        </ul>
      </div>
    </div>
  );
};

export default Confirmation;
