package main

import (
	"fmt"

	"github.com/secsy/goftp"
)

func main() {
	err := setupConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println(conf)

	// Create FTP client
	fmt.Println("Creating FTP client")
	client, err := goftp.DialConfig(goftp.Config{
		User:     conf.Server.User,
		Password: conf.Server.Password,
	}, conf.Server.Address+":"+conf.Server.Port)
	if err != nil {
		fmt.Println("Error connecting to the FTP server:", err)
		return
	}
	fmt.Println("Connected to FTP server")
	// Run pullData
	fmt.Println("Running pullData")
	pullData(client)
	defer client.Close()
	fmt.Println("Done")
}
