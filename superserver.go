package main

// CTF Super Server
// Connects incoming connections to
// the I/O of a command line program

import (
	"fmt"
	"net"
	"os"
	"os/exec"
)

// ServerArguments holds everything to
// handle a remote connection
type ServerArguments struct {
	Program_path string
	Port         string
	Address      string
}

func handle_connection(conn net.Conn, args *ServerArguments) {

	cmd := exec.Command(args.Program_path)

	cmd.Stdin = conn
	cmd.Stdout = conn
	cmd.Stderr = conn

	// start the program
	println("Executing command")

	err := cmd.Run()
	if err != nil {
		fmt.Errorf("%s\n", err.Error())
	}
	println("Closing connection with", conn.RemoteAddr().String())
}

func listen_and_serve(args *ServerArguments) {
	println("Spinning up on", args.Address)
	listener, err := net.Listen("tcp", args.Address)
	if err != nil {
		panic(err)
	}

	for {
		client, err := listener.Accept()
		if err != nil {
			fmt.Printf("%s\n", err.Error())
		}
		println("Received client from:", client.RemoteAddr().String())
		go handle_connection(client, args)
	}
}

func handle_args() *ServerArguments {
	var args ServerArguments

	// set up program path
	path, set := os.LookupEnv("SUPERSERVER_PROGRAM")

	if !set {
		panic(fmt.Errorf("Missing program path!"))
	}

	args.Program_path = path

	// set up port
	port, set := os.LookupEnv("SUPERSERVER_PORT")
	if set {
		args.Port = port
	} else {
		args.Port = "9001"
	}

	// set up address
	addr, set := os.LookupEnv("SUPERSERVER_ADDR")
	if set {
		args.Address = addr
	} else {
		args.Address = "0.0.0.0"
	}

	args.Address = net.JoinHostPort(args.Address, args.Port)
	return &args
}

// get the arguments
func main() {
	args := handle_args()
	listen_and_serve(args)
}
