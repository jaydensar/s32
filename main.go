//go:generate go run github.com/prisma/prisma-client-go generate

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jaydensar/site32-backend/prisma/db"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
)

type PutRequestData struct {
	UserID     int      `json:"UserID"`
	Key        string   `json:"Key"`
	Challenges []string `json:"Challenges"`
	Inventory  []string `json:"Inventory"`
	Points     int      `json:"Points"`
}

type ResponseData struct {
	Challenges []string `json:"Challenges"`
	Inventory  []string `json:"Inventory"`
	Points     int      `json:"Points"`
}

var keys = map[string]string{
	"GET": os.Getenv("GET_KEY"),
	"PUT": os.Getenv("PUT_KEY"),
}

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	client := db.NewClient()
	e := echo.New()

	if err := client.Prisma.Connect(); err != nil {
		log.Println(err)
		return
	}

	ctx := context.Background()

	e.GET("/game", func(c echo.Context) error {
		if c.QueryParam("Key") != keys["GET"] {
			return c.String(401, "Unauthorized (missing or incorrect Key param)")
		}

		if c.QueryParam("UserID") == "" {
			return c.String(403, "Forbidden (missing UserID param)")
		}

		player, err := client.Player.FindUnique(
			db.Player.ID.Equals(c.QueryParam("UserID")),
		).Exec(ctx)

		if err == nil {
			return c.JSON(200, player)
		}

		player, err = client.Player.CreateOne(
			db.Player.ID.Set(c.QueryParam("UserID")),
			db.Player.Challenges.Set([]string{}),
			db.Player.Inventory.Set([]string{}),
		).Exec(ctx)

		if err != nil {
			log.Println(err)
			return c.String(500, "An error occured creating a new player")
		}

		return c.JSON(200, &ResponseData{
			Challenges: player.Challenges,
			Inventory:  player.Inventory,
			Points:     player.Points,
		})
	})

	e.PUT("/game", func(c echo.Context) error {
		data := new(PutRequestData)
		if err := c.Bind(data); err != nil {
			log.Println(err)
			return err
		}

		if data.Key != keys["PUT"] {
			return c.String(401, "Unauthorized (missing or incorrect Key property)")
		}

		_, err := client.Player.FindUnique(
			db.Player.ID.Equals(fmt.Sprint(data.UserID)),
		).Update(
			db.Player.Challenges.Set(data.Challenges),
			db.Player.Inventory.Set(data.Inventory),
			db.Player.Points.Set(data.Points),
		).Exec(ctx)

		if err == nil {
			return c.NoContent(204)
		}

		_, err = client.Player.CreateOne(
			db.Player.ID.Set(fmt.Sprint(data.UserID)),
			db.Player.Challenges.Set(data.Challenges),
			db.Player.Inventory.Set(data.Inventory),
			db.Player.Points.Set(data.Points),
		).Exec(ctx)

		if err == nil {
			return c.NoContent(204)
		}

		return c.String(500, "Failed to update player data")
	})

	defer func() {
		if err := client.Prisma.Disconnect(); err != nil {
			panic(err)
		}
	}()

	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}

	e.Start(":" + port)
}
