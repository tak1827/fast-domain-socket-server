package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tak1827/fast-domain-socket-server/data"
	srv "github.com/tak1827/fast-domain-socket-server/server"
)

func main() {
	if len(os.Args) > 1 {
		client()
	} else {
		server()
	}
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func server() {
	var (
		// ctx, cancel = context.WithCancel(context.Background())
		sigCh = make(chan os.Signal, 1)
		s     = srv.NewServer(srv.SockFilePath)
	)

	go func() {
		s.ListenAndServeUNIX()
	}()

	signal.Notify(sigCh, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)

	for {
		select {
		case err := <-s.ErrCh:
			log.Println(err.Error())
			continue
		case <-sigCh:
		}
		break
	}

	err := s.Shutdown()
	handleErr(err)
}

func client() {
	for i := 0; i < 5; i++ {
		conn, err := net.Dial("unix", srv.SockFilePath)
		handleErr(err)

		tx := data.Message{
			Type:    "hoge",
			Payload: "fuga",
		}

		d, _ := tx.Marshal()

		d = append(d, srv.EOFByte)

		_, err = conn.Write(d)
		handleErr(err)

		dst := make([]byte, 1024)

		dst, err = srv.ReadConn(conn, dst)
		handleErr(err)

		err = tx.Unmarshal(dst)
		handleErr(err)

		fmt.Printf("Recieved from client: %v\n", tx)

		conn.Close()

		time.Sleep(1 * time.Second)
	}
}
