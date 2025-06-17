import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Home from './pages/Home';
import Results from './pages/Results';
import Detail from './pages/Detail';
import Confirmation from './pages/Confirmation'; // Agregar esta línea

function App() {
  return (
    <Router>
      <div style={{ minHeight: '100vh', backgroundColor: '#f0f2f5' }}>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/results" element={<Results />} />
          <Route path="/hotel/:id" element={<Detail />} />
          <Route path="/confirmation" element={<Confirmation />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;