package communication

import (
	"io"
	"log/slog"
	"net"
)

type Options struct {
	Port int16
}

type TcpServer struct {
	opt Options
}

func NewServer(opt Options) *TcpServer {
	return &TcpServer{opt: opt}
}

func (s TcpServer) Start() error {
	listen, err := net.Listen("tcp", ":2525")

	if err != nil {
		slog.Error("Error to start server", err)
		return err
	}

	conn, err := listen.Accept()

	if err != nil {
		slog.Error("Error to receive data", err)
		return err
	}

	defer conn.Close()

	for {
		buf := make([]byte, 1024)

		size, err := conn.Read(buf)

		if err != nil {
			if err == io.EOF {
				slog.Info("End message, size[%d]", size)
				break
			}

			slog.Warn("Error when receiving data", err)
		}
	}

	conn.Write([]byte("+Ok\r\n"))
	return nil
}
