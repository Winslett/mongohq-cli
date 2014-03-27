package api

import (
  "net/http"
  "io/ioutil"
  "github.com/MongoHQ/mongohq-cli"
  "github.com/gorilla/websocket"
  "errors"
  "encoding/json"
)

var oauth_client_id = "6fb9368538ef061ed73be71cc291e65b"
var oauth_secret    = "028d31d8ca253cc3004b3ae4470c21bb23c3011e2fc8b442ad72f259be7879ce5c66bfda4ff26d5a0ba8d23369ef3355ef4579f6e7a977ba933dc1a37fd2880c"

func rest_url_for(path string) (string) {
   return "https://dblayer-api.herokuapp.com" + path;
}

func socket_url_for(path string, oauthToken string) (string) {
   return "wss://beta-api.mongohq.com/mongo" + path + "?token=Bearer%20" + oauthToken
}

func userAgent() string {
  return "MongoHQ CLI Version " + mongohq_cli.Version()
}

func rest_get(path string, oauthToken string) ([]byte, error) {
  client := &http.Client{}
  request, err := http.NewRequest("GET", rest_url_for(path), nil)
  request.Header.Add("Authorization", "Bearer " + oauthToken)
  request.Header.Add("User-Agent", userAgent())
  response, err := client.Do(request)

  if err != nil {
    return nil, err
  } else {
    responseBody, _ := ioutil.ReadAll(response.Body)
    response.Body.Close()
    return responseBody, err
  }
}

func open_websocket(message SocketMessage, oauthToken string) (*websocket.Conn, error) {
  dialer := websocket.Dialer{}
  header := http.Header{}
  header.Add("User-Agent", userAgent())
  client, _, err := dialer.Dial(socket_url_for("/ws", oauthToken), header)
  if err != nil {
    return client, errors.New("Error initiating connection to websocket.")
  }
  jsonMessage, err := json.Marshal(message)
  if err != nil {
    return client, errors.New("Error marshalling websocket message.")
  }
  err = client.WriteMessage(websocket.TextMessage, jsonMessage)
  if err != nil {
    return client, errors.New("Error subscribing to websocket feed.")
  }
  return client, nil
}
