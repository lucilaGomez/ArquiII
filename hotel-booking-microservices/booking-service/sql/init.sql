-- booking-service/sql/init.sql
CREATE DATABASE IF NOT EXISTS booking_db;
USE booking_db;

-- Tabla de usuarios
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    phone VARCHAR(20),
    date_of_birth DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    
    INDEX idx_email (email),
    INDEX idx_active (is_active)
);

-- Tabla de mapeo entre IDs internos y IDs de Amadeus
CREATE TABLE IF NOT EXISTS hotel_mappings (
    id INT AUTO_INCREMENT PRIMARY KEY,
    internal_hotel_id VARCHAR(100) UNIQUE NOT NULL,
    amadeus_hotel_id VARCHAR(100) NOT NULL,
    hotel_name VARCHAR(255) NOT NULL,
    city VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_internal_hotel_id (internal_hotel_id),
    INDEX idx_amadeus_hotel_id (amadeus_hotel_id),
    INDEX idx_city (city)
);

-- Tabla de reservas
CREATE TABLE IF NOT EXISTS bookings (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    internal_hotel_id VARCHAR(100) NOT NULL,
    amadeus_hotel_id VARCHAR(100),
    amadeus_booking_id VARCHAR(100),
    check_in_date DATE NOT NULL,
    check_out_date DATE NOT NULL,
    guests INT NOT NULL DEFAULT 1,
    room_type VARCHAR(100),
    total_price DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    currency VARCHAR(3) NOT NULL DEFAULT 'ARS',
    status ENUM('pending', 'confirmed', 'cancelled', 'completed') NOT NULL DEFAULT 'pending',
    booking_reference VARCHAR(50) UNIQUE NOT NULL,
    special_requests TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT,
    
    INDEX idx_user_id (user_id),
    INDEX idx_internal_hotel_id (internal_hotel_id),
    INDEX idx_amadeus_hotel_id (amadeus_hotel_id),
    INDEX idx_booking_reference (booking_reference),
    INDEX idx_check_in_date (check_in_date),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
);

-- Insertar mapeos de hoteles de ejemplo  
INSERT INTO hotel_mappings (internal_hotel_id, amadeus_hotel_id, hotel_name, city) VALUES
('68618f6b6113de8e4703ea1d', 'YXPARKPR', 'Hotel Test Córdoba', 'Córdoba'),
('684da3fe381f2aeebaec3d54', 'ADPARADI', 'Hotel Córdoba Plaza', 'Córdoba')
ON DUPLICATE KEY UPDATE 
    amadeus_hotel_id = VALUES(amadeus_hotel_id),
    hotel_name = VALUES(hotel_name),
    city = VALUES(city);