package commands

import (
	"errors"
	"fmt"
	"sync"
	"thiagofo92/scratch-redis/pkg/adapter"
)

type HandlerCommands struct {
	sync.RWMutex
	dataSet map[string]string
	hSet    map[string]map[string]string
	handler map[string]func(args []adapter.RespDataOutPut) adapter.RespDataOutPut
}

func NewHandler() *HandlerCommands {
	handler := make(map[string]func(args []adapter.RespDataOutPut) adapter.RespDataOutPut)
	dataSet := make(map[string]string)
	hSet := make(map[string]map[string]string)

	hc := &HandlerCommands{
		handler: handler,
		dataSet: dataSet,
		hSet:    hSet,
	}

	hc.handler["PING"] = hc.ping
	hc.handler["SET"] = hc.set
	hc.handler["GET"] = hc.get
	hc.handler["HSET"] = hc.hset
	hc.handler["HGET"] = hc.hget

	return hc
}

func (h *HandlerCommands) ResponseCommand(command string, args []adapter.RespDataOutPut) (adapter.RespDataOutPut, error) {
	handler, ok := h.handler[command]

	if !ok {
		strError := fmt.Sprintf("commnad not found [%s]", command)
		return adapter.RespDataOutPut{}, errors.New(strError)
	}

	data := handler(args)

	return data, nil
}

func (h *HandlerCommands) ping(args []adapter.RespDataOutPut) adapter.RespDataOutPut {
	if len(args) == 0 {
		return adapter.RespDataOutPut{Typ: "string", Str: "PONG"}
	}

	return adapter.RespDataOutPut{Typ: "string", Str: args[0].Bulk}
}

func (h *HandlerCommands) set(args []adapter.RespDataOutPut) adapter.RespDataOutPut {
	if len(args) != 2 {
		return adapter.RespDataOutPut{Typ: "error", Str: "ERR wrong number of arguments for 'set' command"}
	}

	key := args[0].Bulk
	value := args[1].Bulk

	h.Lock()
	h.dataSet[key] = value
	h.Unlock()

	return adapter.RespDataOutPut{Typ: "string", Str: "OK"}
}

func (h *HandlerCommands) get(args []adapter.RespDataOutPut) adapter.RespDataOutPut {
	if len(args) != 1 {
		return adapter.RespDataOutPut{Typ: "error", Str: "ERR wrong number of arguments for 'get' command"}
	}

	key := args[0].Bulk

	h.RLock()
	data, ok := h.dataSet[key]
	h.RUnlock()

	if !ok {
		return adapter.RespDataOutPut{Typ: "null"}
	}

	return adapter.RespDataOutPut{Typ: "bulk", Bulk: data}
}

func (h *HandlerCommands) hset(args []adapter.RespDataOutPut) adapter.RespDataOutPut {
	if len(args) != 3 {
		return adapter.RespDataOutPut{Typ: "error", Str: "ERR wrong number of arguments for 'hset' command"}
	}

	hash := args[0].Bulk
	key := args[1].Bulk
	value := args[2].Bulk

	h.Lock()
	if _, ok := h.hSet[hash]; !ok {
		h.hSet[hash] = map[string]string{}
	}
	h.hSet[hash][key] = value
	h.Unlock()

	return adapter.RespDataOutPut{Typ: "string", Str: "OK"}
}

func (h *HandlerCommands) hget(args []adapter.RespDataOutPut) adapter.RespDataOutPut {
	if len(args) != 2 {
		return adapter.RespDataOutPut{Typ: "error", Str: "ERR wrong number of arguments for 'hget' command"}
	}

	hash := args[0].Bulk
	key := args[1].Bulk

	h.RLock()
	value, ok := h.hSet[hash][key]
	h.RUnlock()

	if !ok {
		return adapter.RespDataOutPut{Typ: "null"}
	}

	return adapter.RespDataOutPut{Typ: "bulk", Bulk: value}
}
