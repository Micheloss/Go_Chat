package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

var (
	global_local_ip string = ""
	remote_ip       string = ""
	c               chan string
)

func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}
func listen_msg() {

	for {

		ServerAddr, err := net.ResolveUDPAddr("udp", ":45678")
		CheckError(err)

		/* Now listen at selected port */
		ServerConn, err := net.ListenUDP("udp", ServerAddr)
		CheckError(err)
		defer ServerConn.Close()

		buf := make([]byte, 1024)

		n, _, err := ServerConn.ReadFromUDP(buf)

		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		//fmt.Println("Received ", string(buf[0:n]), " from ", addr)

		conf := string(buf[0:n])

		c <- conf
	}
}

func send_msg() {

	for {
		fmt.Print("Msg> ")
		var i string
		_, _ = fmt.Scanf("%s", &i)
		send_udp(remote_ip, i)
	}
}
func chat() {

	go listen_msg()
	go send_msg()
}

func listen_udp() (string, string) {

	ServerAddr, err := net.ResolveUDPAddr("udp", ":45678")
	CheckError(err)

	/* Now listen at selected port */
	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	CheckError(err)
	defer ServerConn.Close()

	buf := make([]byte, 1024)
	t := time.Now()
	ServerConn.SetDeadline(t.Add(3000 * time.Millisecond))
	n, addr, err := ServerConn.ReadFromUDP(buf)

	if err != nil {
		//fmt.Println("Error: ", err)
		return "", ""
	}
	//fmt.Println("Received ", string(buf[0:n]), " from ", addr)

	conf := string(buf[0:n])

	return conf, addr.IP.String()

}

func send_udp(ip_to string, msg string) bool {

	ServerAddr, err := net.ResolveUDPAddr("udp", ip_to+":45678")
	CheckError(err)

	LocalAddr, err := net.ResolveUDPAddr("udp", global_local_ip+":45678")
	CheckError(err)

	Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
	CheckError(err)

	defer Conn.Close()

	buf := []byte(msg)
	_, err = Conn.Write(buf)
	if err != nil {
		fmt.Println(msg, err)
		return false
	}
	return true
}

func client() {

	fmt.Print("Provide the IP: ")
	var i string
	_, _ = fmt.Scanf("%s", &i)

	send_udp(i, "hello-there")
	conf, addr := listen_udp()
	for {
		send_udp(i, "hello-there")
		conf, addr = listen_udp()
		if conf == "" && addr == "" {
			break
		}
	}
	if conf == "hi-there" {

		remote_ip = addr
		chat()
	}
}

func server() {

	fmt.Println("Pass this IP to your colleague: " + global_local_ip)
	conf := ""
	addr := ""
	for {
		conf, addr = listen_udp()

		if conf != "" && addr != "" {
			break
		}
	}
	if strings.Contains(conf, "hello-there") {
		remote_ip = addr
		fmt.Println("CONFIRMED")
		send_udp(addr, "hi-there")
		chat()
	}

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

	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) < 1 {

		fmt.Println("Provide option")

		fmt.Println("Choose 1 to be the client, choose 2 to be the server\n")

		fmt.Println("-------- ACTIVE (Client) ---------")

		fmt.Println("	Once you have the IP of your fellow partner, introduce it and wait for a confirmation that both of you are connected\n")

		fmt.Println("-------- PASSIVE (Server) ---------")

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
		c = make(chan string)
		for {
			s := <-c
			fmt.Println("Remote> " + s)
		}
	}

}
