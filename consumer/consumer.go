package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"
	"encoding/json"
	
	"github.com/IBM/sarama"

	"github.com/go-redis/redis/v8"

	// Importar el paquete de MongoDB
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Configuración del consumidor
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// Lista de brokers de Kafka
	brokers := []string{"my-cluster-kafka-bootstrap.kafka.svc:9092"} // Reemplaza con la dirección de tus brokers

	// Crear el cliente de Kafka
	client, err := sarama.NewConsumerGroup(brokers, "my-group", config)
	if err != nil {
		log.Fatalln("Error creando el cliente de Kafka:", err)
	}
	defer client.Close()

	// Canal para manejar las señales de interrupción
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// Iniciar el consumidor de Kafka
	consumer := Consumer{
		ready: make(chan bool),
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := client.Consume(context.Background(), []string{"test"}, &consumer); err != nil {
				log.Fatalln("Error consumiendo mensajes de Kafka:", err)
			}
			// Comprobar si se ha recibido una señal de interrupción
			sig := <-signals
			if sig == os.Interrupt {
				return
			}
		}
	}()

	// Esperar a que el consumidor esté listo
	<-consumer.ready

	// Esperar a que se reciba una señal de interrupción
	<-signals
}

// Consumer representa un consumidor de Kafka
type Consumer struct {
	ready chan bool
}

// Setup implementa la interfaz sarama.ConsumerGroupHandler.Setup
func (c *Consumer) Setup(session sarama.ConsumerGroupSession) error {
	close(c.ready)
	return nil
}

// Cleanup implementa la interfaz sarama.ConsumerGroupHandler.Cleanup
func (c *Consumer) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim implementa la interfaz sarama.ConsumerGroupHandler.ConsumeClaim
func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// Conectar a MongoDB
	client, err := connectMongoDB()
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v\n", err)
	}
	defer client.Disconnect(context.Background())

	// Conexión a Redis
	redisClient := getRedisClient()

	for message := range claim.Messages() {
		// Parsea el mensaje JSON
		var data map[string]string
		if err := json.Unmarshal(message.Value, &data); err != nil {
			log.Printf("Error parsing JSON message: %v\n", err)
			continue
		}

		// Agregar campos de fecha y hora
		data["fecha"] = time.Now().Format("2006-01-02")
		data["hora"] = time.Now().Format("15:04:05")

		// Goroutine para insertar en MongoDB
		go func(data map[string]string) {
			insertMongoDB(client, data)
		}(data)

		// Goroutine para insertar en Redis
		go func(data map[string]string) {
			insertRedis(redisClient, data)
		}(data)

		// Marca el mensaje como procesado
		session.MarkMessage(message, "")
	}

	return nil
}


// connectMongoDB se conecta a una instancia de MongoDB
func connectMongoDB() (*mongo.Client, error) {
    // Define la URI de conexión a MongoDB
    uri := "mongodb://10.101.97.55:27017" 

    // Crea una opción de cliente
    clientOptions := options.Client().ApplyURI(uri)

    // Intenta conectar con MongoDB
    client, err := mongo.Connect(context.Background(), clientOptions)
    if err != nil {
        return nil, err
    }

    // Verifica si la conexión es exitosa
    err = client.Ping(context.Background(), nil)
    if err != nil {
        return nil, err
    }

    fmt.Println("Connected to MongoDB!")
    return client, nil
}

// insertMongoDB inserta un documento en MongoDB
func insertMongoDB(client *mongo.Client, data map[string]string) {
	collection := client.Database("sopes1-p2").Collection("logs")

	_, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		log.Printf("Error inserting document to MongoDB: %v\n", err)
	} else {
		fmt.Println("Mensaje insertado en MongoDB:", data)
	}
}

// getRedisClient obtiene un cliente Redis
func getRedisClient() *redis.Client {
	// Configura el cliente Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "10.109.105.120:6379", // Reemplaza con la dirección de tu servidor Redis
		Password: "",               // Contraseña, si se requiere
		DB:       0,                // Número de base de datos
	})

	return rdb
}


// insertRedis inserta datos en Redis
func insertRedis(client *redis.Client, data map[string]string) {
	// Construye la clave para la hash en Redis
	key := fmt.Sprintf("%s-%s", data["name"], data["album"])

	// Verifica si la clave ya existe en Redis
	exists, err := client.Exists(context.Background(), key).Result()
	if err != nil {
		log.Printf("Error checking key existence in Redis: %v\n", err)
		return
	}

	// Si la clave no existe, establece un nuevo valor
	if exists == 0 {
		err := client.Set(context.Background(), key, 1, 0).Err()
		if err != nil {
			log.Printf("Error setting value in Redis: %v\n", err)
			return
		}
		fmt.Println("New key set in Redis:", key)
	} else {
		// Si la clave ya existe, incrementa el valor
		err := client.IncrBy(context.Background(), key, 1).Err()
		if err != nil {
			log.Printf("Error incrementing value in Redis: %v\n", err)
			return
		}
		fmt.Println("Value incremented in Redis:", key)
	}

	totalKey := "total_votes"
	// Incrementa el total de votos
	err = client.Incr(context.Background(), totalKey).Err()
	if err != nil {
		log.Printf("Error incrementing total votes in Redis: %v\n", err)
		return
	}
	fmt.Println("Total votes incremented in Redis")
}