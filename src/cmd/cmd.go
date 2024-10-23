package cmd

import (
	"biatosh/entity"
	"biatosh/http"
	"biatosh/logging"
	"biatosh/store"
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	log     = logging.New()
	storeDb = store.New(log)
	app     = http.New(storeDb, log)
)

// create new user with command line
// go run main.go create-user --name="Ali" --email="test1" --password="test1" --username="test1" --phone="123"
var userCmd = &cobra.Command{
	Use:   "create-user",
	Short: "User management",
	Long:  `User management`,
	Run: func(cmd *cobra.Command, args []string) {
		name := cmd.Flag("name").Value.String()
		email := cmd.Flag("email").Value.String()
		phone := cmd.Flag("phone").Value.String()
		password := cmd.Flag("password").Value.String()
		username := cmd.Flag("username").Value.String()

		if name == "" || email == "" || password == "" || username == "" {
			log.Fatal("Please provide all required fields")
		}

		_, err := storeDb.CreateUser(context.Background(), &entity.User{
			Name:     name,
			Email:    email,
			Phone:    phone,
			Password: password,
			Username: username,
		})

		if err != nil {
			log.Fatal(err)
		}
	},
}

var rootCmd = &cobra.Command{
	Use:   "biatosh",
	Short: "Biatosh is a very fast static site generator",
	Long:  `Biatosh is a very fast static site generator. It is designed to be simple and easy to use.`,
	Run: func(cmd *cobra.Command, args []string) {
		ip, _ := cmd.Flags().GetString("ip")
		port, _ := cmd.Flags().GetString("port")
		address := fmt.Sprintf("%s:%s", ip, port)

		log.Fatal(app.Listen(address))
	},
}

func init() {
	rootCmd.Flags().String("ip", "127.0.0.1", "IP address to run the app")
	rootCmd.Flags().String("port", "8080", "Port to run the app")

	userCmd.Flags().String("name", "", "Name of the user")
	userCmd.Flags().String("email", "", "Email of the user")
	userCmd.Flags().String("phone", "", "Phone of the user")
	userCmd.Flags().String("password", "", "Password of the user")
	userCmd.Flags().String("username", "", "Username of the user")

	rootCmd.AddCommand(userCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
