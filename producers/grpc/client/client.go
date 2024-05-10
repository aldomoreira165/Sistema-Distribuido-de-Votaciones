package main

import (
	pb "client/proto" // nombre_proyecto/carpeta
	"context"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var ctx = context.Background()

type Data struct {
	Name        string
	Album    	string
	Year 		  string
	Rank         string
}

func insertData(c *fiber.Ctx) error {
	fmt.Println("Insertando data")
	var data map[string]string
	e := c.BodyParser(&data)
	if e != nil {
		return e
	}

	voto := Data{
		Name:         data["name"],
		Album:        data["album"],
		Year:           data["year"],
		Rank:          data["rank"],
	}

	go sendServer(voto)
	
	//retornar respuesta si se encoló o no
	return nil
}

func sendServer(voto Data) {
	fmt.Println("Enviando al server")
	conn, err := grpc.Dial("localhost:3001", grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock())
	if err != nil {
		log.Fatalln(err)
		fmt.Println("No se pudo conectar al server: ", err)
	}

	cl := pb.NewGetInfoClient(conn)
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(conn)

	ret, err := cl.ReturnInfo(ctx, &pb.RequestId{
		Name:     voto.Name,
		Album:    voto.Album,
		Year:       voto.Year,
		Rank:      voto.Rank,
	})
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Respuesta del server " + ret.GetInfo())
}

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"res": "todo bien",
		})
	})
	app.Post("/grpc", insertData)

	err := app.Listen(":3000")
	if err != nil {
		return
	}
}