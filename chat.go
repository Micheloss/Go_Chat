package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

var (
	global_local_ip string = ""
)

func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func client() {

	fmt.Println("Provide the IP: ")
	var i string
	_, err := fmt.Scanf("%s", &i)

	ServerAddr, err := net.ResolveUDPAddr("udp", i+":45678")
	CheckError(err)

	LocalAddr, err := net.ResolveUDPAddr("udp", global_local_ip+":45678")
	CheckError(err)

	Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
	CheckError(err)

	defer Conn.Close()

	for {
		msg := "hello-there"

		buf := []byte(msg)
		_, err := Conn.Write(buf)
		if err != nil {
			fmt.Println(msg, err)
		}
		time.Sleep(time.Second * 1)
	}
}

func server() {

}

func local_ip() string {
	addrs, _ := net.InterfaceAddrs()
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func main() {

	/*
		Choose 1 to be the client, choose 2 to be the server


		-------- ACTIVE (Client) ---------

			Once you have the IP of your fellow partner, introduce it and wait for a confirmation that both of you are connected

		-------- PASSIVE (Server) ---------

			Give to the client the IP that is showing on the screen and wait for a confirmation that both of you are connected


	*/

	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) < 1 {

		fmt.Println("Provide option")

		fmt.Println("Choose 1 to be the client, choose 2 to be the server\n")

		fmt.Println("-------- CLIENT ---------")

		fmt.Println("	Once you have the IP of your fellow partner, introduce it and wait for a confirmation that both of you are connected\n")

		fmt.Println("-------- SERVER ---------")

		fmt.Println("	Give to the client the IP that is showing on the screen and wait for a confirmation that both of you are connected\n ")

		fmt.Println("Example --> go run ./chat.go 1\n\n")

	} else {
		global_local_ip = local_ip()

		if global_local_ip == "" {

			fmt.Println("Error parsing local IP")
			os.Exit(-1)
		} else {
			if argsWithoutProg[0] == "1" {

				client()
			} else {

				server()

			}
		}
		// var i int
		// _, err := fmt.Scanf("%d", &i)

		// ServerAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:10001")
		// CheckError(err)

		// LocalAddr, err2 := net.ResolveUDPAddr("udp", "127.0.0.1:0")
		// CheckError(err2)

		// Conn, err3 := net.DialUDP("udp", LocalAddr, ServerAddr)
		// CheckError(err3)

		// defer Conn.Close()

		// for {
		// 	msg := strconv.Itoa(i)
		// 	i++
		// 	buf := []byte(msg)
		// 	_, err := Conn.Write(buf)
		// 	if err != nil {
		// 		fmt.Println(msg, err)
		// 	}
		// 	time.Sleep(time.Second * 1)
		// }
	}
}
