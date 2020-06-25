package haws

import (
	"encoding/json"
	"strconv"

	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

type Connection struct {
	conn                         *websocket.Conn
	onOpenCallback               func()
	subscribedToggleableEntities map[string][]func(ToggleableEntity)
	id                           int
	accessToken                  string
}

func handleMessage(c *Connection, message []byte) {
	var f interface{}
	json.Unmarshal(message, &f)

	m := f.(map[string]interface{})

	// TODO: I need to stop with this large handling nonsense.
	// 	     Instead, I should take advantage of HA's id system.
	//       I can set up a map of ids and callbacks that are set up by whomever makes the request
	if m["type"] != nil {
		switch m["type"] {
		case "auth_required":
			// Send an auth message to HA
			authMessage := `{"type": "auth", "access_token": "` + c.accessToken + `"}`

			if err := c.conn.WriteMessage(websocket.TextMessage, []byte(authMessage)); err != nil {
				log.Println("write:", err)
				return
			}

			authMessage = `{"id": 1, "type": "subscribe_events", "event_type": "state_changed"}`

			log.Println("Sent subscribe message")
			if err := c.conn.WriteMessage(websocket.TextMessage, []byte(authMessage)); err != nil {
				log.Println("write:", err)
				return
			}

		case "auth_ok":
			log.Println("Auth is okay!")
			if c.onOpenCallback != nil {
				c.onOpenCallback()
			}

		case "event":
			if event := m["event"].(map[string]interface{}); event["event_type"] == "state_changed" {
				data := event["data"].(map[string]interface{})
				// entity_id here? It should be the same as new state right? Does HA let me know when it changes too?

				newState := data["new_state"].(map[string]interface{})

				entity := ToggleableEntity{
					EntityID:     newState["entity_id"].(string),
					FriendlyName: newState["attributes"].(map[string]interface{})["friendly_name"].(string),
					State:        (newState["state"] == "on"),
				}

				notifySubscribedToggleableEntities(c, entity)
			}
		case "result":
			if m["result"] != nil {
				if result, ok := m["result"].([]interface{}); ok {
					for _, v := range result {
						v := v.(map[string]interface{})
						entityID := v["entity_id"].(string)

						if c.subscribedToggleableEntities[entityID] != nil {
							entity := ToggleableEntity{
								EntityID:     entityID,
								FriendlyName: v["attributes"].(map[string]interface{})["friendly_name"].(string),
								State:        (v["state"] == "on"),
							}
							notifySubscribedToggleableEntities(c, entity)
						}
					}
				}
			}

			// default:
			// 	log.Print("Unrecognized type:", authType)
		}
	}
}

func notifySubscribedToggleableEntities(c *Connection, entity ToggleableEntity) {
	for _, callbackfn := range c.subscribedToggleableEntities[entity.EntityID] {
		callbackfn(entity)
	}
}

type serviceMessage struct {
	Domain      string            `json:"domain"`
	Service     string            `json:"service"`
	ServiceData map[string]string `json:"service_data"`

	// TODO: It'd be good to find out how to only add these to the json and not actually have them present in the struct
	ID          int    `json:"id"`
	ServiceType string `json:"type"`
}

func callHAService(c *Connection, message serviceMessage) {
	c.id++
	message.ID = c.id
	message.ServiceType = "call_service"
	messageStr, _ := json.Marshal(message)

	if err := c.conn.WriteMessage(websocket.TextMessage, []byte(messageStr)); err != nil {
		log.Println("write:", err)
		return
	}
}

func ToggleDevice(c *Connection, entity ToggleableEntity) {
	serviceData := map[string]string{"entity_id": entity.EntityID}

	message := serviceMessage{
		Domain:      entity.Domain,
		Service:     "toggle",
		ServiceData: serviceData,
	}
	callHAService(c, message)
}

type ToggleableEntity struct {
	EntityID     string
	Domain       string
	FriendlyName string
	State        bool
}

func SubscribeToggleableEntity(c *Connection, entityID string, callback func(ToggleableEntity)) {
	if c.subscribedToggleableEntities[entityID] == nil {
		var test []func(ToggleableEntity)
		c.subscribedToggleableEntities[entityID] = test
	}
	c.subscribedToggleableEntities[entityID] = append(c.subscribedToggleableEntities[entityID], callback)
}

func RefreshAllEntities(c *Connection) {
	c.id++
	message := `{"id": ` + strconv.Itoa(c.id) + `, "type": "get_states"}`
	if err := c.conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		log.Println("write:", err)
		return
	}
}

func OpenConnection(serverAddress string, accessToken string, onOpen func(), onErr func()) *Connection {

	u := url.URL{Scheme: "ws", Host: serverAddress, Path: "/api/websocket"}
	log.Printf("connecting to %s", u.String())

	c, resp, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Println("resp:", resp)
		log.Fatal("dial:", err)
	}

	connection := &Connection{
		conn:                         c,
		onOpenCallback:               onOpen,
		subscribedToggleableEntities: make(map[string][]func(ToggleableEntity)),
		id:                           1,
		accessToken:                  accessToken,
	}

	// Handle incoming messages
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				if onErr != nil {
					onErr()
				}
				return
			}
			handleMessage(connection, message)
		}
	}()

	return connection
}

// Close closes the underlying network connection without sending or waiting
// for a close message
func Close(c *Connection) {
	if c == nil || c.conn == nil {
		return // TODO: Return error
	}

	c.conn.Close()
}
