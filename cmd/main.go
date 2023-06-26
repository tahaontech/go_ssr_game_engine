package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"html/template"
	"image/jpeg"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tahaontech/go_ssr_game_engine/internal/game"
	"github.com/tahaontech/go_ssr_game_engine/internal/types"
)

type GameServer struct {
	conn *websocket.Conn
	game *game.GameObj
	sync.Mutex
}

func NewGameServer(conn *websocket.Conn) (*GameServer, error) {
	g, err := game.NewGame("./public/gopher.png")
	if err != nil {
		return nil, err
	}
	return &GameServer{
		conn: conn,
		game: g,
	}, nil
}

func (gs *GameServer) ChangeState() {
	gs.Lock()
	exDir := gs.game.Dir
	var newDir int
	if exDir == 0 {
		newDir = 1
	} else {
		newDir = 0
	}
	gs.game.Dir = newDir
	gs.Unlock()
}

func (gs *GameServer) GetFrame() (string, error) {
	gs.game.Update()
	im := gs.game.GetFrame()
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, im, nil)
	if err != nil {
		return "", err
	}
	str := base64.StdEncoding.EncodeToString(buf.Bytes())
	newStr := "data:image/jpg;base64," + str
	return newStr, nil
}

func (gs *GameServer) RenderLoop() {
	for {
		d, err := gs.GetFrame()
		if err != nil {
			log.Println("write:", err)
			gs.conn.Close()
			break
		}
		msg := types.Message{
			Type: "frame",
			Data: d,
		}
		gs.conn.WriteJSON(msg)
		time.Sleep(time.Microsecond * 10)
	}
}

func (gs *GameServer) EventLoop() {
	for {
		mt, message, err := gs.conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			gs.conn.Close()
			break
		}
		log.Printf("recv: %s", message)

		if string(message) == "change" {
			gs.ChangeState()
		}

		gs.conn.WriteMessage(mt, []byte("state changed"))
	}
}

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

func gameWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	gameServer, err := NewGameServer(conn)
	if err != nil {
		log.Print("game:", err)
		return
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		gameServer.EventLoop()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		gameServer.RenderLoop()
	}()

	wg.Wait()
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/ws")
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/ws", gameWS)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8" />
    <script>
      window.addEventListener("load", function (evt) {
        var output = document.getElementById("output");
        var ws;
        var print = function (message) {
          var d = document.createElement("div");
          d.innerHTML = message;
          output.appendChild(d);
        };
        document.getElementById("open").onclick = function (evt) {
          if (ws) {
            return false;
          }
          ws = new WebSocket("{{.}}");
          ws.onopen = function (evt) {
            print("OPEN");
          };
          ws.onclose = function (evt) {
            print("CLOSE");
            ws = null;
          };
          ws.onmessage = function (evt) {
            console.log(evt);
            if (typeof evt.data == "object") {
              draw(evt.data);
            } else {
                const da = JSON.parse(evt.data);
                if (da.type == "frame") {
                    draw(da.data);
                } else {
                    print("RESPONSE: " + evt.data);
                }
            }
          };
          ws.onerror = function (evt) {
            print("ERROR: " + evt.data);
          };
          return false;
        };
        document.getElementById("send").onclick = function (evt) {
          if (!ws) {
            return false;
          }
          print("SEND: " + "changedir");
          ws.send("change");
          return false;
        };
        document.getElementById("close").onclick = function (evt) {
          if (!ws) {
            return false;
          }
          ws.close();
          return false;
        };
      });
    </script>
  </head>
  <body>
    <table>
      <tr>
        <td valign="top" width="50%">
          <p>
            Click "Start" to create a connection to the server, "Change" to send a
            change direction event to the server and "Close" to close the connection. You can
            send event multiple times.
          </p>

          <p></p>
          <form>
            <button id="open">Start</button>
            <button id="close">Close</button>
            <p>
              <button id="send">Change</button>
            </p>
          </form>

          <canvas
            id="canvas"
            style="width: 400px; height: 400px; border: 2px solid black"
          ></canvas>
        </td>
        <td valign="top" width="50%">
          <div id="output"></div>
        </td>
      </tr>
    </table>
  </body>
  <script>
    var canvas = document.getElementById("canvas");
    var ctx = canvas.getContext("2d");
    function draw(data) {
        var img1 = new Image();
        img1.onload = function() {
            ctx.drawImage(img1, 0, 0);
        }
        img1.src = data;
    }
  </script>
</html>
`))
