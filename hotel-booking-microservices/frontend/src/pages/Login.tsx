import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { bookingAPI } from '../services/api';

const Login: React.FC = () => {
  const navigate = useNavigate();
  const [isLoginMode, setIsLoginMode] = useState(true);
  const [loading, setLoading] = useState(false);
  
  const [loginData, setLoginData] = useState({
    email: '',
    password: ''
  });
  
  const [registerData, setRegisterData] = useState({
    email: '',
    password: '',
    first_name: '',
    last_name: '',
    phone: '',
    date_of_birth: '1990-01-01'
  });

  // Verificar si ya está logueado
  useEffect(() => {
    const token = localStorage.getItem('token');
    if (token) {
      // Ya está logueado, redirigir según el rol
      const userRole = localStorage.getItem('userRole');
      if (userRole === 'admin') {
        navigate('/admin');
      } else {
        navigate('/dashboard');
      }
    }
  }, [navigate]);

  const handleLogin = async () => {
    setLoading(true);
    try {
      const response = await bookingAPI.login(loginData);
      const token = response.token;
      const user = response.user;
      
      // Guardar en localStorage
      localStorage.setItem('token', token);
      localStorage.setItem('userRole', user.role);
      localStorage.setItem('userName', `${user.first_name} ${user.last_name}`);
      localStorage.setItem('userEmail', user.email);
      
      // Redirigir según rol
      if (user.role === 'admin') {
        navigate('/admin');
      } else {
        navigate('/dashboard');
      }
      
    } catch (err) {
      alert('Error en login. Verifica tus credenciales.');
      console.error('Login error:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleRegister = async () => {
    if (!registerData.email || !registerData.password || !registerData.first_name || !registerData.last_name) {
      alert('Por favor completa todos los campos obligatorios');
      return;
    }

    setLoading(true);
    try {
      const registerPayload = {
        ...registerData,
        date_of_birth: registerData.date_of_birth + 'T00:00:00Z'
      } as any;

      await bookingAPI.register(registerPayload);
      alert('¡Registro exitoso! Ahora puedes iniciar sesión.');
      setIsLoginMode(true);
      setRegisterData({
        email: '',
        password: '',
        first_name: '',
        last_name: '',
        phone: '',
        date_of_birth: '1990-01-01'
      });
    } catch (err) {
      alert('Error en registro. El email puede ya estar registrado.');
      console.error('Register error:', err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ 
      minHeight: '100vh',
      background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      padding: '20px'
    }}>
      <div style={{
        backgroundColor: 'white',
        borderRadius: '20px',
        boxShadow: '0 20px 40px rgba(0,0,0,0.1)',
        overflow: 'hidden',
        width: '100%',
        maxWidth: '1000px',
        display: 'grid',
        gridTemplateColumns: '1fr 1fr'
      }}>
        
        {/* Panel izquierdo - Branding */}
        <div style={{
          background: 'linear-gradient(45deg, #1976d2, #42a5f5)',
          padding: '60px 40px',
          color: 'white',
          display: 'flex',
          flexDirection: 'column',
          justifyContent: 'center',
          alignItems: 'center',
          textAlign: 'center'
        }}>
          <div style={{ fontSize: '72px', marginBottom: '20px' }}>🏨</div>
          <h1 style={{ fontSize: '32px', marginBottom: '20px', margin: 0 }}>
            Hotel Manager
          </h1>
          <p style={{ fontSize: '18px', opacity: 0.9, lineHeight: 1.6 }}>
            Sistema completo de reservas de hoteles con microservicios
          </p>
          <div style={{ marginTop: '40px', fontSize: '14px', opacity: 0.8 }}>
            <div>🔐 Cuentas de prueba:</div>
            <div><strong>Admin:</strong> admin@hotelmanager.com / password</div>
            <div><strong>Usuario:</strong> testfinal@ucc.edu.ar / password</div>
          </div>
        </div>

        {/* Panel derecho - Formulario */}
        <div style={{ padding: '60px 40px' }}>
          
          {/* Toggle Login/Register */}
          <div style={{ display: 'flex', marginBottom: '30px' }}>
            <button
              onClick={() => setIsLoginMode(true)}
              style={{
                flex: 1,
                padding: '12px',
                backgroundColor: isLoginMode ? '#1976d2' : 'transparent',
                color: isLoginMode ? 'white' : '#666',
                border: '2px solid #1976d2',
                borderRadius: '25px 0 0 25px',
                cursor: 'pointer',
                fontWeight: 'bold',
                transition: 'all 0.3s ease'
              }}
            >
              Iniciar Sesión
            </button>
            <button
              onClick={() => setIsLoginMode(false)}
              style={{
                flex: 1,
                padding: '12px',
                backgroundColor: !isLoginMode ? '#1976d2' : 'transparent',
                color: !isLoginMode ? 'white' : '#666',
                border: '2px solid #1976d2',
                borderRadius: '0 25px 25px 0',
                cursor: 'pointer',
                fontWeight: 'bold',
                transition: 'all 0.3s ease'
              }}
            >
              Registrarse
            </button>
          </div>

          {/* Formulario de Login */}
          {isLoginMode ? (
            <div>
              <h2 style={{ marginBottom: '30px', color: '#333' }}>¡Bienvenido de vuelta! 👋</h2>
              
              <input
                type="email"
                placeholder="📧 Email"
                value={loginData.email}
                onChange={(e) => setLoginData(prev => ({ ...prev, email: e.target.value }))}
                style={{
                  width: '100%',
                  padding: '15px 20px',
                  marginBottom: '20px',
                  border: '2px solid #e0e0e0',
                  borderRadius: '10px',
                  fontSize: '16px',
                  boxSizing: 'border-box',
                  transition: 'border-color 0.3s ease',
                  outline: 'none'
                }}
                onFocus={(e) => e.target.style.borderColor = '#1976d2'}
                onBlur={(e) => e.target.style.borderColor = '#e0e0e0'}
              />
              
              <input
                type="password"
                placeholder="🔒 Contraseña"
                value={loginData.password}
                onChange={(e) => setLoginData(prev => ({ ...prev, password: e.target.value }))}
                style={{
                  width: '100%',
                  padding: '15px 20px',
                  marginBottom: '30px',
                  border: '2px solid #e0e0e0',
                  borderRadius: '10px',
                  fontSize: '16px',
                  boxSizing: 'border-box',
                  transition: 'border-color 0.3s ease',
                  outline: 'none'
                }}
                onFocus={(e) => e.target.style.borderColor = '#1976d2'}
                onBlur={(e) => e.target.style.borderColor = '#e0e0e0'}
              />
              
              <button
                onClick={handleLogin}
                disabled={loading}
                style={{
                  width: '100%',
                  padding: '15px',
                  backgroundColor: '#1976d2',
                  color: 'white',
                  border: 'none',
                  borderRadius: '10px',
                  fontSize: '18px',
                  fontWeight: 'bold',
                  cursor: loading ? 'not-allowed' : 'pointer',
                  transition: 'background-color 0.3s ease',
                  opacity: loading ? 0.7 : 1
                }}
                onMouseOver={(e) => {
                  if (!loading) (e.target as HTMLButtonElement).style.backgroundColor = '#1565c0';
                }}
                onMouseOut={(e) => {
                  if (!loading) (e.target as HTMLButtonElement).style.backgroundColor = '#1976d2';
                }}
              >
                {loading ? '⏳ Iniciando sesión...' : '🚀 Iniciar Sesión'}
              </button>
            </div>
          ) : (
            /* Formulario de Registro */
            <div>
              <h2 style={{ marginBottom: '30px', color: '#333' }}>¡Únete a nosotros! ✨</h2>
              
              <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '15px', marginBottom: '15px' }}>
                <input
                  type="text"
                  placeholder="👤 Nombre"
                  value={registerData.first_name}
                  onChange={(e) => setRegisterData(prev => ({ ...prev, first_name: e.target.value }))}
                  style={{
                    padding: '12px 15px',
                    border: '2px solid #e0e0e0',
                    borderRadius: '8px',
                    fontSize: '14px',
                    boxSizing: 'border-box',
                    outline: 'none'
                  }}
                />
                <input
                  type="text"
                  placeholder="👤 Apellido"
                  value={registerData.last_name}
                  onChange={(e) => setRegisterData(prev => ({ ...prev, last_name: e.target.value }))}
                  style={{
                    padding: '12px 15px',
                    border: '2px solid #e0e0e0',
                    borderRadius: '8px',
                    fontSize: '14px',
                    boxSizing: 'border-box',
                    outline: 'none'
                  }}
                />
              </div>
              
              <input
                type="email"
                placeholder="📧 Email"
                value={registerData.email}
                onChange={(e) => setRegisterData(prev => ({ ...prev, email: e.target.value }))}
                style={{
                  width: '100%',
                  padding: '12px 15px',
                  marginBottom: '15px',
                  border: '2px solid #e0e0e0',
                  borderRadius: '8px',
                  fontSize: '14px',
                  boxSizing: 'border-box',
                  outline: 'none'
                }}
              />
              
              <input
                type="password"
                placeholder="🔒 Contraseña"
                value={registerData.password}
                onChange={(e) => setRegisterData(prev => ({ ...prev, password: e.target.value }))}
                style={{
                  width: '100%',
                  padding: '12px 15px',
                  marginBottom: '15px',
                  border: '2px solid #e0e0e0',
                  borderRadius: '8px',
                  fontSize: '14px',
                  boxSizing: 'border-box',
                  outline: 'none'
                }}
              />
              
              <input
                type="tel"
                placeholder="📱 Teléfono (opcional)"
                value={registerData.phone}
                onChange={(e) => setRegisterData(prev => ({ ...prev, phone: e.target.value }))}
                style={{
                  width: '100%',
                  padding: '12px 15px',
                  marginBottom: '25px',
                  border: '2px solid #e0e0e0',
                  borderRadius: '8px',
                  fontSize: '14px',
                  boxSizing: 'border-box',
                  outline: 'none'
                }}
              />
              
              <button
                onClick={handleRegister}
                disabled={loading}
                style={{
                  width: '100%',
                  padding: '15px',
                  backgroundColor: '#2e7d32',
                  color: 'white',
                  border: 'none',
                  borderRadius: '10px',
                  fontSize: '18px',
                  fontWeight: 'bold',
                  cursor: loading ? 'not-allowed' : 'pointer',
                  transition: 'background-color 0.3s ease',
                  opacity: loading ? 0.7 : 1
                }}
                onMouseOver={(e) => {
                  if (!loading) (e.target as HTMLButtonElement).style.backgroundColor = '#1b5e20';
                }}
                onMouseOut={(e) => {
                  if (!loading) (e.target as HTMLButtonElement).style.backgroundColor = '#2e7d32';
                }}
              >
                {loading ? '⏳ Creando cuenta...' : '✅ Crear Cuenta'}
              </button>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default Login;