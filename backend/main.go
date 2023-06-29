package main

import (
	"gitlab.com/koshkaj/macler/backend/api"
)

func main() {
	s := api.InitServer()
	defer s.Cr.Stop()
	s.Logger.Fatal(s.Start(":8080"))
}
