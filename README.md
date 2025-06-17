# Sistema de Microservicios de Hoteles

## Descripción del Proyecto

Sistema de microservicios desarrollado para una cadena de hoteles que permite disponibilizar su oferta de forma local y se integra con un proveedor central (Amadeus) para validar las reservas.

**Universidad Católica de Córdoba - Arquitectura de Software II - 2023**

## Arquitectura del Sistema

El sistema está compuesto por 4 microservicios principales:

1. **Frontend** - Interfaz de usuario
2. **Ficha de Hotel** - Gestión de información de hoteles
3. **Búsqueda de Hotel** - Motor de búsqueda con Solr
4. **Usuarios, Reserva y Disponibilidad** - Gestión de usuarios y reservas

### Diagrama de Arquitectura

```
┌─────────────────────────────────────────────────────────────┐
│                        Frontend                             │
└─────────────────┬───────────────┬───────────────────────────┘
                  │               │
        ┌─────────▼─────────┐   ┌─▼─────────────────────────────┐
        │  Ficha de Hotel   │   │ Búsqueda de Hotel             │
        │                   │   │                               │
        │ ┌─────────────┐   │   │ ┌─────────────┐ ┌─────────┐   │
        │ │  MongoDB    │   │   │ │    Solr     │ │RabbitMQ │   │
        │ └─────────────┘   │   │ └─────────────┘ └─────────┘   │
        └───────────────────┘   └───────────────────────────────┘
                  │                               │
                  │               ┌───────────────▼─────────────┐
                  │               │ Usuarios/Reserva/Disponib.  │
                  │               │                             │
                  │               │ ┌─────────┐ ┌─────────────┐ │
                  │               │ │ MySQL   │ │ Memcached   │ │
                  │               │ └─────────┘ └─────────────┘ │
                  │               └─────────────┬───────────────┘
                  │                             │
                  └─────────────────────────────▼─────────────────
                                    ┌─────────────────┐
                                    │    Amadeus      │
                                    │   (Externo)     │
                                    └─────────────────┘
```

## Microservicios

### 1. Microservicio Frontend

**Tecnología:** React/Angular/Vue (según implementación)

**Pantallas:**
- **Inicial:** Búsqueda por ciudad, fecha desde y fecha hasta
- **Resultados:** Listado de hoteles disponibles (nombre, descripción, thumbnail)
- **Detalle:** Información completa del hotel (fotos, amenities, botón reserva)
- **Confirmación:** Éxito o rechazo de la reserva

### 2. Microservicio de Ficha de Hotel

**Tecnología:** Go + MongoDB + RabbitMQ

**Funcionalidades:**
- API RESTful para gestión de hoteles
- Almacenamiento documental en MongoDB
- Validación de creación y modificación de hoteles
- Notificaciones via RabbitMQ para cambios

**Endpoints principales:**
- `GET /hotel/{id}` - Obtener hotel por ID
- `POST /hotel` - Crear nuevo hotel
- `PUT /hotel/{id}` - Actualizar hotel

### 3. Microservicio de Búsqueda de Hotel

**Tecnología:** Go + Apache Solr + RabbitMQ

**Funcionalidades:**
- Motor de búsqueda con Apache Solr
- Consumer de RabbitMQ para sincronización
- Consulta concurrente de disponibilidad
- Filtrado dinámico por disponibilidad

**Características técnicas:**
- Sincronización automática desde ficha de hoteles
- Atributo dinámico "availability" (no persistido)
- Consultas concurrentes al servicio de disponibilidad

### 4. Microservicio de Usuarios, Reserva y Disponibilidad

**Tecnología:** Go + MySQL + Memcached + Amadeus API

**Funcionalidades:**
- Gestión de usuarios y clientes
- Sistema de reservas
- Caché distribuido con TTL de 10 segundos
- Validación externa con Amadeus
- Mapping de IDs internos a IDs de Amadeus

**Características técnicas:**
- Consistencia eventual de 10 segundos
- Validación de reservas con proveedor externo
- Gestión automática de tokens de Amadeus

## Tecnologías Utilizadas

| Componente | Tecnología |
|------------|------------|
| Backend | Go |
| Base de Datos NoSQL | MongoDB |
| Base de Datos SQL | MySQL |
| Motor de Búsqueda | Apache Solr |
| Cache Distribuido | Memcached |
| Message Broker | RabbitMQ |
| API Externa | Amadeus |
| Contenedores | Docker |
| Orquestación | Docker Compose |

## Instalación y Configuración

### Prerrequisitos

- Docker y Docker Compose
- Go 1.19+
- Cuenta en Amadeus Developer Portal

### Configuración de Amadeus

1. Registrarse en [Amadeus Developer Portal](https://developers.amadeus.com/)
2. Obtener `CLIENT_ID` y `CLIENT_SECRET`
3. Configurar variables de entorno:

```bash
export AMADEUS_CLIENT_ID=your_client_id
export AMADEUS_CLIENT_SECRET=your_client_secret
```

### Ejecución con Docker Compose

```bash
# Clonar el repositorio
git clone [URL_DEL_REPOSITORIO]
cd hotel-microservices

# Construir y levantar todos los servicios
docker-compose up --build

# En modo desarrollo (con logs)
docker-compose up --build -d && docker-compose logs -f
```

### Ejecución Individual de Servicios

```bash
# Microservicio de Ficha de Hotel
cd hotel-info-service
go run main.go

# Microservicio de Búsqueda
cd hotel-search-service
go run main.go

# Microservicio de Reservas
cd hotel-booking-service
go run main.go
```

## Variables de Entorno

Crear un archivo `.env` en la raíz del proyecto:

```env
# Base de Datos
MONGODB_URI=mongodb://localhost:27017/hotels
MYSQL_URI=user:password@tcp(localhost:3306)/bookings

# Cache y Message Broker
MEMCACHED_URI=localhost:11211
RABBITMQ_URI=amqp://guest:guest@localhost:5672/

# Solr
SOLR_URI=http://localhost:8983/solr/hotels

# Amadeus API
AMADEUS_CLIENT_ID=your_client_id
AMADEUS_CLIENT_SECRET=your_client_secret
AMADEUS_BASE_URL=https://test.api.amadeus.com

# Puertos de Servicios
HOTEL_INFO_PORT=8081
HOTEL_SEARCH_PORT=8082
HOTEL_BOOKING_PORT=8083
FRONTEND_PORT=3000
```

## Testing

### Tests Unitarios

```bash
# Ejecutar tests en todos los microservicios
make test

# Tests individuales por servicio
cd hotel-info-service && go test ./...
cd hotel-search-service && go test ./...
cd hotel-booking-service && go test ./...
```

### Tests de Integración

```bash
# Ejecutar tests de integración
make integration-test
```

## Monitoreo y Logs

### Verificación de Servicios

```bash
# Verificar estado de contenedores
docker-compose ps

# Ver logs de un servicio específico
docker-compose logs hotel-info-service
docker-compose logs hotel-search-service
docker-compose logs hotel-booking-service
```

### Endpoints de Health Check

- Frontend: `http://localhost:3000/health`
- Hotel Info: `http://localhost:8081/health`
- Hotel Search: `http://localhost:8082/health`
- Hotel Booking: `http://localhost:8083/health`

## APIs y Endpoints

### Documentación de API

Una vez levantados los servicios, la documentación de Swagger estará disponible en:

- Hotel Info Service: `http://localhost:8081/swagger`
- Hotel Search Service: `http://localhost:8082/swagger`
- Hotel Booking Service: `http://localhost:8083/swagger`

### Ejemplos de Uso

```bash
# Buscar hoteles en París
curl "http://localhost:8082/search?city=PAR&checkIn=2024-01-22&checkOut=2024-01-24"

# Obtener información de un hotel
curl "http://localhost:8081/hotel/123"

# Verificar disponibilidad
curl "http://localhost:8083/availability?hotelId=123&checkIn=2024-01-22&checkOut=2024-01-24"
```

## Criterios de Evaluación Implementados

- ✅ Frontend con 4 pantallas funcionales
- ✅ Conexión MongoDB en servicio de ficha
- ✅ Conexión RabbitMQ y notificaciones
- ✅ Motor de búsqueda Solr implementado
- ✅ Cálculo concurrente de disponibilidad
- ✅ Conexión MySQL y Memcached
- ✅ Integración con API de Amadeus
- ✅ Dockerización completa
- ✅ Arquitectura MVC en microservicios
- ✅ Manejo de errores HTTP

## Mejoras Futuras (Puntos Adicionales)

- [ ] Load Balancer en microservicios de backend
- [ ] Escalado automático con tests de carga
- [ ] Tests unitarios completos
- [ ] Frontend de administración de infraestructura
- [ ] Cache local para incrementar escalabilidad

## Contribución

1. Fork el proyecto
2. Crear una rama para tu feature (`git checkout -b feature/nueva-funcionalidad`)
3. Commit tus cambios (`git commit -am 'Agrega nueva funcionalidad'`)
4. Push a la rama (`git push origin feature/nueva-funcionalidad`)
5. Crear un Pull Request

## Licencia

Este proyecto está bajo la Licencia MIT - ver el archivo [LICENSE](LICENSE) para detalles.

## Contacto

**Universidad Católica de Córdoba**  
Arquitectura de Software II - 2023

---

*Desarrollado como Trabajo Final Integrador para la materia Arquitectura de Software II*
