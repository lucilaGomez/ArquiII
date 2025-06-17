package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connect establece conexión con MongoDB
func Connect(uri string) (*mongo.Client, error) {
	// Configurar opciones de conexión
	opts := options.Client().ApplyURI(uri)
	opts.SetMaxPoolSize(20)
	opts.SetMinPoolSize(5)
	opts.SetMaxConnIdleTime(30 * time.Second)
	opts.SetConnectTimeout(10 * time.Second)
	opts.SetServerSelectionTimeout(5 * time.Second)

	// Crear contexto con timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Conectar
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	// Ping para verificar conexión
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Disconnect cierra la conexión con MongoDB
func Disconnect(client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	return client.Disconnect(ctx)
}

// GetDatabase obtiene la base de datos
func GetDatabase(client *mongo.Client, dbName string) *mongo.Database {
	return client.Database(dbName)
}

// GetCollection obtiene una colección específica
func GetCollection(client *mongo.Client, dbName, collectionName string) *mongo.Collection {
	return client.Database(dbName).Collection(collectionName)
}