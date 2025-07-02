import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Login from './pages/Login';
import Dashboard from './pages/Dashboard';
import AdminDashboard from './pages/AdminDashboard';
import Detail from './pages/Detail';
import Results from './pages/Results';
import Confirmation from './pages/Confirmation';

function App() {
  return (
    <Router>
      <div className="App">
        <Routes>
          {/* Página principal - Login */}
          <Route path="/" element={<Login />} />
          
          {/* Dashboard para usuarios autenticados */}
          <Route path="/dashboard" element={<Dashboard />} />
          
          {/* Panel de administración (solo admins) */}
          <Route path="/admin" element={<AdminDashboard />} />
          
          {/* Resultados de búsqueda */}
          <Route path="/results" element={<Results />} />
          
          {/* Detalles del hotel */}
          <Route path="/hotel/:id" element={<Detail />} />
          
          {/* Confirmación de reserva */}
          <Route path="/confirmation" element={<Confirmation />} />
          
          {/* Ruta por defecto - redirige a login */}
          <Route path="*" element={<Login />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;