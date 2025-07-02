#!/bin/bash

echo "🏨 Cargando hoteles en la base de datos..."
echo "=========================================="

# URL del hotel service
HOTEL_SERVICE_URL="http://localhost:8001/api/v1/hotels"

# 1. Hotel Córdoba Centro
echo "1. Creando Hotel Córdoba Centro..."
curl -X POST $HOTEL_SERVICE_URL \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Hotel Córdoba Centro",
    "description": "Hotel céntrico con excelente ubicación en el corazón de Córdoba",
    "city": "Córdoba",
    "country": "Argentina",
    "address": "San Jerónimo 200, Centro",
    "amenities": ["WiFi", "Desayuno", "Aire Acondicionado", "Estacionamiento"],
    "thumbnail": "https://example.com/cordoba-centro.jpg",
    "rating": 4.5,
    "price_range": {
      "min_price": 8000,
      "max_price": 15000,
      "currency": "ARS"
    },
    "contact": {
      "phone": "+54 351 123-4567",
      "email": "info@cordobacentro.com",
      "website": "https://cordobacentro.com"
    }
  }' && echo ""

# 2. Hotel Sierras de Córdoba
echo "2. Creando Hotel Sierras de Córdoba..."
curl -X POST $HOTEL_SERVICE_URL \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Hotel Sierras de Córdoba",
    "description": "Hotel boutique en las sierras con vistas espectaculares",
    "city": "Villa Carlos Paz",
    "country": "Argentina", 
    "address": "Villa Carlos Paz, Córdoba",
    "amenities": ["WiFi", "Spa", "Piscina", "Trekking", "Vista Serrana"],
    "thumbnail": "https://example.com/sierras-cordoba.jpg",
    "rating": 4.6,
    "price_range": {
      "min_price": 12000,
      "max_price": 18000,
      "currency": "ARS"
    },
    "contact": {
      "phone": "+54 351 789-0123",
      "email": "info@sierrascordoba.com",
      "website": "https://sierrascordoba.com"
    }
  }' && echo ""

# 3. Hotel Buenos Aires Plaza
echo "3. Creando Hotel Buenos Aires Plaza..."
curl -X POST $HOTEL_SERVICE_URL \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Hotel Buenos Aires Plaza",
    "description": "Hotel ejecutivo en el microcentro de Buenos Aires",
    "city": "Buenos Aires",
    "country": "Argentina",
    "address": "Av. Corrientes 1500, Microcentro",
    "amenities": ["WiFi", "Gimnasio", "Business Center", "Restaurante"],
    "thumbnail": "https://example.com/ba-plaza.jpg",
    "rating": 4.8,
    "price_range": {
      "min_price": 15000,
      "max_price": 25000,
      "currency": "ARS"
    },
    "contact": {
      "phone": "+54 11 123-4567",
      "email": "info@baplaza.com",
      "website": "https://baplaza.com"
    }
  }' && echo ""

# 4. Hotel Puerto Madero (MÁS VALORADO)
echo "4. Creando Hotel Puerto Madero..."
curl -X POST $HOTEL_SERVICE_URL \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Hotel Puerto Madero",
    "description": "Hotel de lujo exclusivo con vista al río en Puerto Madero",
    "city": "Buenos Aires",
    "country": "Argentina",
    "address": "Alicia Moreau de Justo 1000, Puerto Madero",
    "amenities": ["WiFi", "Spa", "Rooftop", "Vista al Río", "Concierge"],
    "thumbnail": "https://example.com/puerto-madero.jpg",
    "rating": 4.9,
    "price_range": {
      "min_price": 20000,
      "max_price": 35000,
      "currency": "ARS"
    },
    "contact": {
      "phone": "+54 11 456-7890",
      "email": "info@puertomadero.com",
      "website": "https://puertomadero.com"
    }
  }' && echo ""

# 5. Hotel Mendoza Wine
echo "5. Creando Hotel Mendoza Wine..."
curl -X POST $HOTEL_SERVICE_URL \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Hotel Mendoza Wine",
    "description": "Hotel temático de vinos con cata incluida en el corazón de Mendoza",
    "city": "Mendoza",
    "country": "Argentina",
    "address": "San Martín 1200, Centro",
    "amenities": ["WiFi", "Spa", "Cata de Vinos", "Vista Montaña"],
    "thumbnail": "https://example.com/mendoza-wine.jpg",
    "rating": 4.7,
    "price_range": {
      "min_price": 12000,
      "max_price": 20000,
      "currency": "ARS"
    },
    "contact": {
      "phone": "+54 261 123-4567",
      "email": "info@mendozawine.com",
      "website": "https://mendozawine.com"
    }
  }' && echo ""

# 6. Hotel Nahuel Huapi
echo "6. Creando Hotel Nahuel Huapi..."
curl -X POST $HOTEL_SERVICE_URL \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Hotel Nahuel Huapi",
    "description": "Hotel premium con vista al lago Nahuel Huapi y actividades de montaña",
    "city": "Bariloche",
    "country": "Argentina",
    "address": "Av. Bustillo Km 2.5, San Carlos de Bariloche",
    "amenities": ["WiFi", "Spa", "Vista al Lago", "Ski", "Kayak"],
    "thumbnail": "https://example.com/nahuel-huapi.jpg",
    "rating": 4.8,
    "price_range": {
      "min_price": 18000,
      "max_price": 30000,
      "currency": "ARS"
    },
    "contact": {
      "phone": "+54 294 123-789",
      "email": "info@nahuelhuapi.com",
      "website": "https://nahuelhuapi.com"
    }
  }' && echo ""

# 7. Hotel Salta Colonial
echo "7. Creando Hotel Salta Colonial..."
curl -X POST $HOTEL_SERVICE_URL \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Hotel Salta Colonial",
    "description": "Hotel histórico con arquitectura colonial en el centro de Salta",
    "city": "Salta",
    "country": "Argentina",
    "address": "Caseros 500, Centro Histórico, Salta",
    "amenities": ["WiFi", "Patio Colonial", "Desayuno Regional", "Turismo"],
    "thumbnail": "https://example.com/salta-colonial.jpg",
    "rating": 4.4,
    "price_range": {
      "min_price": 10000,
      "max_price": 16000,
      "currency": "ARS"
    },
    "contact": {
      "phone": "+54 387 123-456",
      "email": "info@saltacolonial.com",
      "website": "https://saltacolonial.com"
    }
  }' && echo ""

# 8. Hotel Costa Atlántica
echo "8. Creando Hotel Costa Atlántica..."
curl -X POST $HOTEL_SERVICE_URL \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Hotel Costa Atlántica",
    "description": "Hotel frente al mar con playa privada en Mar del Plata",
    "city": "Mar del Plata",
    "country": "Argentina",
    "address": "Boulevard Marítimo 2500, Mar del Plata",
    "amenities": ["WiFi", "Playa Privada", "Pileta", "Actividades Acuáticas"],
    "thumbnail": "https://example.com/costa-atlantica.jpg",
    "rating": 4.3,
    "price_range": {
      "min_price": 9000,
      "max_price": 14000,
      "currency": "ARS"
    },
    "contact": {
      "phone": "+54 223 456-789",
      "email": "info@costaatlantica.com",
      "website": "https://costaatlantica.com"
    }
  }' && echo ""

# 9. Hotel Rosario Centro (MÁS ECONÓMICO)
echo "9. Creando Hotel Rosario Centro..."
curl -X POST $HOTEL_SERVICE_URL \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Hotel Rosario Centro",
    "description": "Hotel económico y cómodo en el centro de Rosario",
    "city": "Rosario",
    "country": "Argentina", 
    "address": "San Martín 1200, Centro, Rosario",
    "amenities": ["WiFi", "Business Center", "Gimnasio", "Cerca del Río"],
    "thumbnail": "https://example.com/rosario-centro.jpg",
    "rating": 4.2,
    "price_range": {
      "min_price": 7500,
      "max_price": 12000,
      "currency": "ARS"
    },
    "contact": {
      "phone": "+54 341 789-012",
      "email": "info@rosariocentro.com",
      "website": "https://rosariocentro.com"
    }
  }' && echo ""

echo ""
echo "✅ ¡Proceso completado!"
echo "=========================================="
echo "📊 Verificando hoteles creados..."

# Verificar hoteles creados
curl -s $HOTEL_SERVICE_URL | jq '.data | length' 2>/dev/null || echo "Instalá jq para ver el conteo: brew install jq"

echo ""
echo "🎉 ¡Todos los hoteles han sido creados exitosamente!"
echo "   Puedes verificarlos en: http://localhost:3000"
