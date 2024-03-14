package main

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"thiagofo92/scratch-redis/pkg/adapter"
	"thiagofo92/scratch-redis/pkg/commands"
)

func main() {
	r, w := io.Pipe()
	input := "*3\r\n$4\r\nping\r\n$5\r\nadmin\r\n$5\r\nahmed"

	reader := bufio.NewReader(strings.NewReader(input))

	resp := adapter.NewRespInput(reader)

	data, err := resp.Read()

	if err != nil {
		slog.Error("error to read input", err)
		return
	}

	command := strings.ToUpper(data.Value[0].Bulk)
	args := data.Value[1:]

	go func() {
		defer w.Close()

		handler := commands.NewHandler()
		writerResp := adapter.NewWriter(w)
		fmt.Printf("Command %s\n", command)
		resp, err := handler.ResponseCommand(command, args)

		if err != nil {
			slog.Error("Error to writer", err)
		}

		writerResp.Write(resp)
	}()

	_, err = io.Copy(os.Stdout, r)

	if err != nil {
		slog.Error("Error to read data", err)
	}

	// server := communication.NewServer(communication.Options{})

	// err := server.Start()

	// if err != nil {
	// 	slog.Error("Server error", err)
	// }

}
