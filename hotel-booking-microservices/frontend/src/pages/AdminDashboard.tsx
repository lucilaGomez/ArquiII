import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { hotelAPI } from '../services/api';
import { Hotel } from '../types';

const AdminDashboard: React.FC = () => {
  const navigate = useNavigate();
  const [userName, setUserName] = useState('');
  const [hotels, setHotels] = useState<Hotel[]>([]);
  const [loading, setLoading] = useState(true);
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [editingHotel, setEditingHotel] = useState<Hotel | null>(null);
  const [uploading, setUploading] = useState(false);
  
  // ESTRUCTURA CORREGIDA - coincide con lo que espera el backend
  const [hotelForm, setHotelForm] = useState({
    name: '',
    description: '',
    city: '',
    address: '', // REQUERIDO por el backend
    amenities: ['WiFi', 'Desayuno'], // Array con valores por defecto
    rating: 4.0,
    price_range: {
      min_price: 10000,
      max_price: 25000,
      currency: 'ARS' // REQUERIDO por el backend
    },
    contact: {
      phone: '',
      email: '', // REQUERIDO por el backend (debe ser email válido)
      website: ''
    },
    thumbnail: '', // URL de la imagen principal
    images: [] as string[] // URLs de imágenes adicionales
  });

  useEffect(() => {
    // Verificar autenticación y rol de admin
    const token = localStorage.getItem('token');
    const userRole = localStorage.getItem('userRole');
    
    if (!token || userRole !== 'admin') {
      alert('Acceso denegado. Solo administradores pueden acceder a esta página.');
      navigate('/dashboard');
      return;
    }

    const name = localStorage.getItem('userName') || 'Admin';
    setUserName(name);
    
    loadHotels();
  }, [navigate]);

  const loadHotels = async () => {
    try {
      setLoading(true);
      const response = await hotelAPI.getAllHotels();
      setHotels(response.data || []);
    } catch (error) {
      console.error('Error loading hotels:', error);
      alert('Error cargando hoteles');
    } finally {
      setLoading(false);
    }
  };

  // Función para subir imagen principal (thumbnail)
  const handleThumbnailUpload = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;

    // Validar tipo de archivo
    if (!file.type.startsWith('image/')) {
      alert('Por favor selecciona un archivo de imagen válido');
      return;
    }

    // Validar tamaño (máximo 5MB)
    if (file.size > 5 * 1024 * 1024) {
      alert('El archivo es muy grande. Máximo 5MB permitido');
      return;
    }

    setUploading(true);
    const formData = new FormData();
    formData.append('image', file);

    try {
      const response = await fetch('http://localhost:8001/api/v1/hotels/upload-single', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        },
        body: formData
      });

      if (!response.ok) {
        throw new Error('Error subiendo imagen');
      }

      const result = await response.json();
      setHotelForm(prev => ({ ...prev, thumbnail: result.url }));
      alert('Imagen principal subida exitosamente');
    } catch (error) {
      console.error('Error uploading thumbnail:', error);
      alert('Error subiendo imagen principal');
    } finally {
      setUploading(false);
    }
  };

  // Función para subir múltiples imágenes
  const handleMultipleImagesUpload = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const files = Array.from(event.target.files || []);
    if (files.length === 0) return;

    // Validar archivos
    for (const file of files) {
      if (!file.type.startsWith('image/')) {
        alert(`${file.name} no es un archivo de imagen válido`);
        return;
      }
      if (file.size > 5 * 1024 * 1024) {
        alert(`${file.name} es muy grande. Máximo 5MB por archivo`);
        return;
      }
    }

    setUploading(true);
    const formData = new FormData();
    files.forEach(file => {
      formData.append('images', file);
    });

    try {
      const response = await fetch('http://localhost:8001/api/v1/hotels/upload-images', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        },
        body: formData
      });

      if (!response.ok) {
        throw new Error('Error subiendo imágenes');
      }

      const result = await response.json();
      const newImageUrls = result.files.map((file: any) => file.url);
      
      setHotelForm(prev => ({ 
        ...prev, 
        images: [...prev.images, ...newImageUrls] 
      }));
      
      alert(`${files.length} imágenes subidas exitosamente`);
    } catch (error) {
      console.error('Error uploading images:', error);
      alert('Error subiendo imágenes');
    } finally {
      setUploading(false);
    }
  };

  // Función para eliminar imagen de la lista
  const removeImage = (indexToRemove: number) => {
    setHotelForm(prev => ({
      ...prev,
      images: prev.images.filter((_, index) => index !== indexToRemove)
    }));
  };

  const handleCreateHotel = async () => {
    // Validación básica
    if (!hotelForm.name || !hotelForm.city || !hotelForm.address || !hotelForm.contact.email) {
      alert('Por favor completa todos los campos obligatorios: Nombre, Ciudad, Dirección y Email');
      return;
    }

    // Validar email
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(hotelForm.contact.email)) {
      alert('Por favor ingresa un email válido');
      return;
    }

    try {
      await hotelAPI.createHotel(hotelForm);
      alert('Hotel creado exitosamente');
      setShowCreateForm(false);
      resetForm();
      loadHotels();
    } catch (error) {
      console.error('Error creating hotel:', error);
      alert('Error creando hotel. Revisa que todos los campos estén completos.');
    }
  };

  const handleUpdateHotel = async () => {
    if (!editingHotel) return;
    
    // Validación básica
    if (!hotelForm.name || !hotelForm.city || !hotelForm.address || !hotelForm.contact.email) {
      alert('Por favor completa todos los campos obligatorios: Nombre, Ciudad, Dirección y Email');
      return;
    }

    try {
      await hotelAPI.updateHotel(editingHotel.id, hotelForm);
      alert('Hotel actualizado exitosamente');
      setEditingHotel(null);
      resetForm();
      loadHotels();
    } catch (error) {
      console.error('Error updating hotel:', error);
      alert('Error actualizando hotel');
    }
  };

  const handleDeleteHotel = async (hotelId: string) => {
    if (!window.confirm('¿Estás seguro de que quieres eliminar este hotel?')) return;
    
    try {
      await hotelAPI.deleteHotel(hotelId);
      alert('Hotel eliminado exitosamente');
      loadHotels();
    } catch (error) {
      console.error('Error deleting hotel:', error);
      alert('Error eliminando hotel');
    }
  };

  const resetForm = () => {
    setHotelForm({
      name: '',
      description: '',
      city: '',
      address: '',
      amenities: ['WiFi', 'Desayuno'],
      rating: 4.0,
      price_range: {
        min_price: 10000,
        max_price: 25000,
        currency: 'ARS'
      },
      contact: {
        phone: '',
        email: '',
        website: ''
      },
      thumbnail: '',
      images: []
    });
  };

  const startEditing = (hotel: Hotel) => {
    setEditingHotel(hotel);
    setHotelForm({
      name: hotel.name,
      description: hotel.description,
      city: hotel.city,
      address: hotel.address,
      amenities: hotel.amenities || ['WiFi', 'Desayuno'],
      rating: hotel.rating,
      price_range: hotel.price_range || { min_price: 10000, max_price: 25000, currency: 'ARS' },
      contact: hotel.contact || { phone: '', email: '', website: '' },
      thumbnail: hotel.thumbnail || '',
      images: hotel.photos || []
    });
    setShowCreateForm(true);
  };

  const handleLogout = () => {
    localStorage.clear();
    navigate('/');
  };

  const handleGoToDashboard = () => {
    navigate('/dashboard');
  };

  return (
    <div style={{ minHeight: '100vh', backgroundColor: '#f5f5f5' }}>
      
      {/* Header */}
      <header style={{
        backgroundColor: '#ff9800',
        color: 'white',
        padding: '15px 20px',
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center',
        boxShadow: '0 2px 4px rgba(0,0,0,0.1)'
      }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '15px' }}>
          <span style={{ fontSize: '24px' }}>🛠️</span>
          <h1 style={{ margin: 0 }}>Panel de Administración</h1>
        </div>
        
        <div style={{ display: 'flex', alignItems: 'center', gap: '15px' }}>
          <span>👨‍💼 {userName}</span>
          
          <button
            onClick={handleGoToDashboard}
            style={{
              padding: '8px 16px',
              backgroundColor: 'rgba(255,255,255,0.2)',
              color: 'white',
              border: '1px solid rgba(255,255,255,0.3)',
              borderRadius: '6px',
              cursor: 'pointer',
              fontSize: '14px'
            }}
          >
            🏠 Dashboard
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

      <div style={{ maxWidth: '1200px', margin: '0 auto', padding: '40px 20px' }}>
        
        {/* Estadísticas */}
        <div style={{
          display: 'grid',
          gridTemplateColumns: '1fr 1fr 1fr 1fr',
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
            <div style={{ fontSize: '32px', color: '#1976d2', marginBottom: '10px' }}>🏨</div>
            <div style={{ fontSize: '24px', fontWeight: 'bold', color: '#333' }}>{hotels.length}</div>
            <div style={{ color: '#666' }}>Total Hoteles</div>
          </div>
          
          <div style={{
            backgroundColor: 'white',
            padding: '25px',
            borderRadius: '10px',
            textAlign: 'center',
            boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
          }}>
            <div style={{ fontSize: '32px', color: '#2e7d32', marginBottom: '10px' }}>🌟</div>
            <div style={{ fontSize: '24px', fontWeight: 'bold', color: '#333' }}>4.5</div>
            <div style={{ color: '#666' }}>Rating Promedio</div>
          </div>
          
          <div style={{
            backgroundColor: 'white',
            padding: '25px',
            borderRadius: '10px',
            textAlign: 'center',
            boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
          }}>
            <div style={{ fontSize: '32px', color: '#ff9800', marginBottom: '10px' }}>🏙️</div>
            <div style={{ fontSize: '24px', fontWeight: 'bold', color: '#333' }}>6</div>
            <div style={{ color: '#666' }}>Ciudades</div>
          </div>
          
          <div style={{
            backgroundColor: 'white',
            padding: '25px',
            borderRadius: '10px',
            textAlign: 'center',
            boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
          }}>
            <div style={{ fontSize: '32px', color: '#9c27b0', marginBottom: '10px' }}>📊</div>
            <div style={{ fontSize: '24px', fontWeight: 'bold', color: '#333' }}>98%</div>
            <div style={{ color: '#666' }}>Ocupación</div>
          </div>
        </div>

        {/* Acciones principales */}
        <div style={{
          backgroundColor: 'white',
          padding: '25px',
          borderRadius: '10px',
          marginBottom: '30px',
          boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
        }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <h2 style={{ margin: 0, color: '#333' }}>Gestión de Hoteles</h2>
            <button
              onClick={() => {
                setShowCreateForm(true);
                setEditingHotel(null);
                resetForm();
              }}
              style={{
                padding: '10px 20px',
                backgroundColor: '#1976d2',
                color: 'white',
                border: 'none',
                borderRadius: '8px',
                cursor: 'pointer',
                fontSize: '16px',
                fontWeight: 'bold'
              }}
            >
              ➕ Crear Nuevo Hotel
            </button>
          </div>
        </div>

        {/* Formulario de creación/edición */}
        {showCreateForm && (
          <div style={{
            backgroundColor: 'white',
            padding: '30px',
            borderRadius: '10px',
            marginBottom: '30px',
            boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
          }}>
            <h3 style={{ marginTop: 0, color: '#333' }}>
              {editingHotel ? '✏️ Editar Hotel' : '➕ Crear Nuevo Hotel'}
            </h3>
            
            <div style={{
              display: 'grid',
              gridTemplateColumns: '1fr 1fr',
              gap: '20px',
              marginBottom: '20px'
            }}>
              <input
                type="text"
                placeholder="Nombre del hotel *"
                value={hotelForm.name}
                onChange={(e) => setHotelForm(prev => ({ ...prev, name: e.target.value }))}
                style={{
                  padding: '12px',
                  border: '2px solid #e0e0e0',
                  borderRadius: '8px',
                  fontSize: '14px'
                }}
              />
              
              <input
                type="text"
                placeholder="Ciudad *"
                value={hotelForm.city}
                onChange={(e) => setHotelForm(prev => ({ ...prev, city: e.target.value }))}
                style={{
                  padding: '12px',
                  border: '2px solid #e0e0e0',
                  borderRadius: '8px',
                  fontSize: '14px'
                }}
              />
            </div>

            <textarea
              placeholder="Descripción del hotel"
              value={hotelForm.description}
              onChange={(e) => setHotelForm(prev => ({ ...prev, description: e.target.value }))}
              rows={3}
              style={{
                width: '100%',
                padding: '12px',
                border: '2px solid #e0e0e0',
                borderRadius: '8px',
                fontSize: '14px',
                marginBottom: '20px',
                boxSizing: 'border-box',
                resize: 'vertical'
              }}
            />

            <div style={{
              display: 'grid',
              gridTemplateColumns: '2fr 1fr 1fr',
              gap: '20px',
              marginBottom: '20px'
            }}>
              <input
                type="text"
                placeholder="Dirección completa *"
                value={hotelForm.address}
                onChange={(e) => setHotelForm(prev => ({ ...prev, address: e.target.value }))}
                style={{
                  padding: '12px',
                  border: '2px solid #e0e0e0',
                  borderRadius: '8px',
                  fontSize: '14px'
                }}
              />
              
              <input
                type="number"
                placeholder="Precio mínimo"
                value={hotelForm.price_range.min_price}
                onChange={(e) => setHotelForm(prev => ({
                  ...prev,
                  price_range: { ...prev.price_range, min_price: Number(e.target.value) }
                }))}
                style={{
                  padding: '12px',
                  border: '2px solid #e0e0e0',
                  borderRadius: '8px',
                  fontSize: '14px'
                }}
              />
              
              <input
                type="number"
                placeholder="Precio máximo"
                value={hotelForm.price_range.max_price}
                onChange={(e) => setHotelForm(prev => ({
                  ...prev,
                  price_range: { ...prev.price_range, max_price: Number(e.target.value) }
                }))}
                style={{
                  padding: '12px',
                  border: '2px solid #e0e0e0',
                  borderRadius: '8px',
                  fontSize: '14px'
                }}
              />
            </div>

            <div style={{
              display: 'grid',
              gridTemplateColumns: '1fr 1fr 1fr',
              gap: '20px',
              marginBottom: '20px'
            }}>
              <input
                type="tel"
                placeholder="Teléfono"
                value={hotelForm.contact.phone}
                onChange={(e) => setHotelForm(prev => ({
                  ...prev,
                  contact: { ...prev.contact, phone: e.target.value }
                }))}
                style={{
                  padding: '12px',
                  border: '2px solid #e0e0e0',
                  borderRadius: '8px',
                  fontSize: '14px'
                }}
              />
              
              <input
                type="email"
                placeholder="Email *"
                value={hotelForm.contact.email}
                onChange={(e) => setHotelForm(prev => ({
                  ...prev,
                  contact: { ...prev.contact, email: e.target.value }
                }))}
                style={{
                  padding: '12px',
                  border: '2px solid #e0e0e0',
                  borderRadius: '8px',
                  fontSize: '14px'
                }}
              />
              
              <input
                type="url"
                placeholder="Sitio web"
                value={hotelForm.contact.website}
                onChange={(e) => setHotelForm(prev => ({
                  ...prev,
                  contact: { ...prev.contact, website: e.target.value }
                }))}
                style={{
                  padding: '12px',
                  border: '2px solid #e0e0e0',
                  borderRadius: '8px',
                  fontSize: '14px'
                }}
              />
            </div>

            <div style={{ marginBottom: '20px' }}>
              <label style={{ display: 'block', marginBottom: '5px', color: '#666' }}>
                Amenidades (separadas por coma):
              </label>
              <input
                type="text"
                placeholder="WiFi, Desayuno, Gimnasio, Piscina"
                value={hotelForm.amenities.join(', ')}
                onChange={(e) => setHotelForm(prev => ({
                  ...prev,
                  amenities: e.target.value.split(',').map(item => item.trim()).filter(item => item)
                }))}
                style={{
                  width: '100%',
                  padding: '12px',
                  border: '2px solid #e0e0e0',
                  borderRadius: '8px',
                  fontSize: '14px',
                  boxSizing: 'border-box'
                }}
              />
            </div>

            {/* SECCIÓN DE UPLOAD DE IMÁGENES */}
            <div style={{
              backgroundColor: '#f8f9fa',
              padding: '20px',
              borderRadius: '8px',
              marginBottom: '20px'
            }}>
              <h4 style={{ margin: '0 0 15px 0', color: '#333' }}>📸 Imágenes del Hotel</h4>
              
              {/* Upload imagen principal */}
              <div style={{ marginBottom: '20px' }}>
                <label style={{ display: 'block', marginBottom: '8px', color: '#666', fontWeight: 'bold' }}>
                  🖼️ Imagen Principal (Thumbnail):
                </label>
                <input
                  type="file"
                  accept="image/*"
                  onChange={handleThumbnailUpload}
                  disabled={uploading}
                  style={{
                    padding: '10px',
                    border: '2px dashed #ccc',
                    borderRadius: '8px',
                    width: '100%',
                    boxSizing: 'border-box',
                    backgroundColor: 'white'
                  }}
                />
                {hotelForm.thumbnail && (
                  <div style={{ marginTop: '10px' }}>
                    <img 
                      src={`http://localhost:8001${hotelForm.thumbnail}`} 
                      alt="Thumbnail preview" 
                      style={{ 
                        width: '150px', 
                        height: '100px', 
                        objectFit: 'cover', 
                        borderRadius: '8px',
                        border: '2px solid #ddd'
                      }} 
                    />
                    <button
                      onClick={() => setHotelForm(prev => ({ ...prev, thumbnail: '' }))}
                      style={{
                        marginLeft: '10px',
                        padding: '5px 10px',
                        backgroundColor: '#f44336',
                        color: 'white',
                        border: 'none',
                        borderRadius: '4px',
                        cursor: 'pointer',
                        fontSize: '12px'
                      }}
                    >
                      ❌ Quitar
                    </button>
                  </div>
                )}
              </div>

              {/* Upload múltiples imágenes */}
              <div>
                <label style={{ display: 'block', marginBottom: '8px', color: '#666', fontWeight: 'bold' }}>
                  🖼️ Galería de Imágenes (múltiples):
                </label>
                <input
                  type="file"
                  accept="image/*"
                  multiple
                  onChange={handleMultipleImagesUpload}
                  disabled={uploading}
                  style={{
                    padding: '10px',
                    border: '2px dashed #ccc',
                    borderRadius: '8px',
                    width: '100%',
                    boxSizing: 'border-box',
                    backgroundColor: 'white'
                  }}
                />
                
                {/* Preview de imágenes subidas */}
                {hotelForm.images.length > 0 && (
                  <div style={{ marginTop: '15px' }}>
                    <div style={{ marginBottom: '10px', color: '#666', fontSize: '14px' }}>
                      📷 Imágenes subidas ({hotelForm.images.length}):
                    </div>
                    <div style={{ 
                      display: 'grid', 
                      gridTemplateColumns: 'repeat(auto-fill, minmax(120px, 1fr))', 
                      gap: '10px' 
                    }}>
                      {hotelForm.images.map((imageUrl, index) => (
                        <div key={index} style={{ position: 'relative' }}>
                          <img 
                            src={`http://localhost:8001${imageUrl}`} 
                            alt={`Preview ${index + 1}`} 
                            style={{ 
                              width: '100%', 
                              height: '80px', 
                              objectFit: 'cover', 
                              borderRadius: '6px',
                              border: '1px solid #ddd'
                            }} 
                          />
                          <button
                            onClick={() => removeImage(index)}
                            style={{
                              position: 'absolute',
                              top: '2px',
                              right: '2px',
                              backgroundColor: '#f44336',
                              color: 'white',
                              border: 'none',
                              borderRadius: '50%',
                              width: '20px',
                              height: '20px',
                              cursor: 'pointer',
                              fontSize: '10px',
                              display: 'flex',
                              alignItems: 'center',
                              justifyContent: 'center'
                            }}
                          >
                            ×
                          </button>
                        </div>
                      ))}
                    </div>
                  </div>
                )}
                
                <div style={{ marginTop: '10px', fontSize: '12px', color: '#999' }}>
                  💡 Puedes subir múltiples imágenes a la vez. Máximo 5MB por archivo.
                  <br />
                  Formatos aceptados: .jpg, .jpeg, .png, .webp
                </div>
              </div>

              {uploading && (
                <div style={{ 
                  marginTop: '15px', 
                  padding: '10px', 
                  backgroundColor: '#e3f2fd', 
                  borderRadius: '6px',
                  color: '#1976d2',
                  textAlign: 'center'
                }}>
                  ⏳ Subiendo archivos...
                </div>
              )}
            </div>

            <div style={{ display: 'flex', gap: '15px' }}>
              <button
                onClick={editingHotel ? handleUpdateHotel : handleCreateHotel}
                disabled={uploading}
                style={{
                  padding: '12px 24px',
                  backgroundColor: uploading ? '#ccc' : '#2e7d32',
                  color: 'white',
                  border: 'none',
                  borderRadius: '8px',
                  cursor: uploading ? 'not-allowed' : 'pointer',
                  fontSize: '14px',
                  fontWeight: 'bold'
                }}
              >
                {editingHotel ? '💾 Actualizar Hotel' : '✅ Crear Hotel'}
              </button>
              
              <button
                onClick={() => {
                  setShowCreateForm(false);
                  setEditingHotel(null);
                  resetForm();
                }}
                style={{
                  padding: '12px 24px',
                  backgroundColor: '#666',
                  color: 'white',
                  border: 'none',
                  borderRadius: '8px',
                  cursor: 'pointer',
                  fontSize: '14px'
                }}
              >
                ❌ Cancelar
              </button>
            </div>

            <div style={{ marginTop: '15px', fontSize: '12px', color: '#999' }}>
              * Campos obligatorios
            </div>
          </div>
        )}

        {/* Lista de hoteles */}
        <div style={{
          backgroundColor: 'white',
          borderRadius: '10px',
          overflow: 'hidden',
          boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
        }}>
          <div style={{
            padding: '20px',
            backgroundColor: '#f8f9fa',
            borderBottom: '1px solid #e0e0e0'
          }}>
            <h3 style={{ margin: 0, color: '#333' }}>Hoteles Registrados ({hotels.length})</h3>
          </div>

          {loading ? (
            <div style={{ padding: '60px', textAlign: 'center', color: '#666' }}>
              ⏳ Cargando hoteles...
            </div>
          ) : hotels.length === 0 ? (
            <div style={{ padding: '60px', textAlign: 'center', color: '#666' }}>
              🏨 No hay hoteles registrados
            </div>
          ) : (
            <div style={{ maxHeight: '500px', overflowY: 'auto' }}>
              {hotels.map((hotel: Hotel) => (
                <div
                  key={hotel.id}
                  style={{
                    padding: '20px',
                    borderBottom: '1px solid #f0f0f0',
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'center'
                  }}
                >
                  <div style={{ flex: 1, display: 'flex', gap: '15px' }}>
                    {/* Imagen del hotel */}
                    {hotel.thumbnail && (
                      <img 
                        src={`http://localhost:8001${hotel.thumbnail}`} 
                        alt={hotel.name}
                        style={{
                          width: '80px',
                          height: '60px',
                          objectFit: 'cover',
                          borderRadius: '8px',
                          border: '1px solid #ddd'
                        }}
                      />
                    )}
                    
                    <div>
                      <h4 style={{ margin: '0 0 8px 0', color: '#333' }}>
                        🏨 {hotel.name}
                      </h4>
                      <p style={{ margin: '0 0 8px 0', color: '#666' }}>
                        📍 {hotel.city} • ⭐ {hotel.rating}/5 • 
                        💰 ${hotel.price_range?.min_price?.toLocaleString()} - ${hotel.price_range?.max_price?.toLocaleString()} ARS
                      </p>
                      <p style={{ margin: 0, color: '#999', fontSize: '14px' }}>
                        {hotel.description?.substring(0, 100)}...
                      </p>
                      {hotel.photos && hotel.photos.length > 0 && (
                        <p style={{ margin: '5px 0 0 0', color: '#1976d2', fontSize: '12px' }}>
                          📷 {hotel.photos.length} imagen{hotel.photos.length !== 1 ? 'es' : ''}
                        </p>
                      )}
                    </div>
                  </div>
                  
                  <div style={{ display: 'flex', gap: '10px' }}>
                    <button
                      onClick={() => startEditing(hotel)}
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
                      ✏️ Editar
                    </button>
                    
                    <button
                      onClick={() => handleDeleteHotel(hotel.id)}
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
                      🗑️ Eliminar
                    </button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default AdminDashboard;