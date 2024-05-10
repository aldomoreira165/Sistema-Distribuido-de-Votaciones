package main

import (
	"context"
	"fmt"
	"log"
	"net"
	pb "server/proto"
	"time"

	"github.com/IBM/sarama"
	"google.golang.org/grpc"
)

var producer sarama.SyncProducer
var ctx = context.Background()

type server struct {
	pb.UnimplementedGetInfoServer
}

const (
	port = ":3001"
)

type Data struct {
	Name         string
	Album        string
	Year           string
	Rank          string
}

func initProducer() {
    config := sarama.NewConfig()
    config.Producer.RequiredAcks = sarama.WaitForLocal       // Espera a que el líder haya confirmado la recepción del mensaje
    config.Producer.Compression = sarama.CompressionSnappy   // Opcional: utiliza compresión Snappy
    config.Producer.Flush.Frequency = 500 * time.Millisecond // Opcional: establece la frecuencia de los flushes

    // Configura Producer.Return.Successes como true
    config.Producer.Return.Successes = true

    brokers := []string{"my-cluster-kafka-bootstrap.kafka.svc:9092"}  // Reemplaza localhost:9092 con la dirección de tu servidor Kafka

    // Crea un nuevo productor
    p, err := sarama.NewSyncProducer(brokers, config)
    if err != nil {
        log.Fatalln("Failed to start Sarama producer:", err)
    }
    producer = p
}


func (s *server) ReturnInfo(ctx context.Context, in *pb.RequestId) (*pb.ReplyInfo, error) {
	fmt.Println("Recibí de cliente: ", in.GetName())
	data := Data{
		Name:        in.GetName(),
		Album:       in.GetAlbum(),
		Year:          in.GetYear(),
		Rank:         in.GetRank(),
	}
	fmt.Println(data)
	queueData(data)
	return &pb.ReplyInfo{Info: "Hola cliente, recibí el comentario"}, nil
}

func queueData(voto Data) {
    // Serializa el mensaje a JSON
    message := fmt.Sprintf(`{"name": "%s", "album": "%s", "year": "%s", "rank": "%s"}`, voto.Name, voto.Album, voto.Year, voto.Rank)
    // Envia el mensaje al topic "test"
    _, _, err := producer.SendMessage(&sarama.ProducerMessage{
        Topic: "test",
        Value: sarama.StringEncoder(message),
    })
    if err != nil {
        log.Fatalln("Failed to send message:", err)
    }
    fmt.Println("Sent message to Kafka:", message)
}

func main() {
	initProducer()
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Server listening on port", port)
	s := grpc.NewServer()
	pb.RegisterGetInfoServer(s, &server{})

	if err := s.Serve(listen); err != nil {
		log.Fatalln(err)
	}
}