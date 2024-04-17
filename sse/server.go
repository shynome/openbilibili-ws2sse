package sse

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/shynome/err0"
	"github.com/shynome/err0/try"
	bilibili "github.com/shynome/openapi-bilibili"
	"github.com/shynome/openapi-bilibili/live"
	"github.com/shynome/openapi-bilibili/live/cmd"
)

type Server struct {
	bclient *bilibili.Client
	appid   int64
	games   map[string]bool
	mux     sync.RWMutex
}

var _ http.Handler = (*Server)(nil)

func New(key, secret string, appid int64) *Server {
	return &Server{
		bclient: bilibili.NewClient(key, secret),
		appid:   appid,
		games:   make(map[string]bool),
	}
}

func (srv *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	defer err0.Then(&err, nil, func() {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	})

	flusher, ok := w.(http.Flusher)
	if !ok {
		err = fmt.Errorf("server don't support http.flusher")
		return
	}

	IDCode := r.URL.Query().Get("IDCode")
	if IDCode == "" {
		http.Error(w, "QueryParam IDCode is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	app := try.To1(srv.bclient.Open(ctx, srv.appid, IDCode))
	defer app.Close()

	gid := app.Info().GameInfo.GameId
	srv.addGame(gid)
	defer srv.removeGame(gid)

	room := live.RoomWith(app.Info().WebsocketInfo)
	ch := try.To1(room.Connect(ctx))

	stream := StreamWriter{
		Flusher: flusher,
		Writer:  w,
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.WriteHeader(http.StatusOK)
	init := Msg[string]{
		Type: MsgInit,
		Data: app.Info().WebsocketInfo.AuthBody,
	}
	init.WriteTry(stream)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			try.To1(fmt.Fprintf(stream, ": hack for pass cdn \n"))
			stream.Flush()
		case danmu := <-ch:
			switch danmu.Cmd {
			case cmd.CmdDanmu,
				cmd.CmdGift,
				cmd.CmdGuard,
				cmd.CmdSuperChat,
				cmd.CmdDelSuperChat,
				cmd.CmdLike:
				msg := Msg[cmd.Cmd[json.RawMessage]]{
					Type: MsgDanmu,
					Data: danmu,
				}
				msg.WriteTry(stream)
			default:
				slog.Info("others", "danmu", danmu)
			}
		}
	}

}

func (srv *Server) addGame(gid string) {
	srv.mux.Lock()
	defer srv.mux.Unlock()
	srv.games[gid] = true
}

func (srv *Server) removeGame(gid string) {
	srv.mux.Lock()
	defer srv.mux.Unlock()
	delete(srv.games, gid)
}

type MsgType string

const (
	MsgInit  MsgType = "init"
	MsgDanmu MsgType = "danmu"
)

type Msg[T any] struct {
	Type MsgType `json:"type"`
	Data T       `json:"data"`
}

func (msg Msg[T]) WriteTry(w StreamWriter) {
	defer w.Flush()
	id := time.Now().Unix()
	try.To1(fmt.Fprintf(w, "id:%d\n", id))
	try.To1(io.WriteString(w, "data:"))
	try.To(json.NewEncoder(w).Encode(msg))
	try.To1(io.WriteString(w, "\n"))
	// 结束该 Event
	try.To1(io.WriteString(w, "\n"))
}

type StreamWriter struct {
	http.Flusher
	io.Writer
}
