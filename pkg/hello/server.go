package hello

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net"
)

type conn struct {
	id   int
	conn net.Conn
}

func Server(server net.Listener) {
	quit := make(chan struct{})
	conns := make(chan conn)

	go connManager(server, conns)

	for {
		select {
		case c := <-conns:
			go echo(c, quit)
		case <-quit:
			return
		}
	}
}

func connManager(server net.Listener, conns chan conn) {
	i := 0
	for {
		i++
		cn, err := server.Accept()
		if err != nil {
			log.Printf("Error accepting connection %d: %s", i, err)
			continue
		}
		log.Printf("Connection %d: accepted from %s", i, cn.RemoteAddr())
		conns <- conn{i, cn}
	}
}

func echo(c conn, quit chan struct{}) {
	buf := bufio.NewReader(c.conn)
	firstLine := true
	for {
		line, err := buf.ReadBytes('\n')
		if firstLine && bytes.HasPrefix(line, []byte{'Z', 'Z', 'Z'}) {
			quit <- struct{}{}
		}
		firstLine = false

		if len(line) > 0 {
			c.conn.Write(line)
			log.Printf("Connection %d: Wrote %s", c.id, line)
		}

		if err != nil {
			if err == io.EOF {
				log.Printf("Connection %d: Reached end of input", c.id)
			} else {
				log.Printf("Error reading connection %d: %s", c.id, err)
			}
			err = c.conn.Close()
			if err != nil {
				log.Printf("Error closing connection %d: %s", c.id, err)
			}
			break
		}
	}
}
