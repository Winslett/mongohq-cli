package api

import (
  "encoding/json"
)

type Database struct {
    Id   string
    Name string
    Status string
    Plan string
    Deployment_id string
}

type DatabaseUser struct {
    Username string `json:"user"`
    PasswordHash string `json:"pwd"`
    ReadOnly bool
}

func GetDatabases(oauthToken string) ([]Database, error) {
  body, err := rest_get(api_url("/databases"), oauthToken)

  if err != nil {
    return nil, err
  }
  var databasesSlice []Database
  err = json.Unmarshal(body, &databasesSlice)
  return databasesSlice, err
}

func GetDatabase(name string, oauthToken string) (Database, error) {
  body, err := rest_get(api_url("/databases/" + name), oauthToken)

  if err != nil {
    return Database{}, err
  }
  var database Database
  err = json.Unmarshal(body, &database)
  return database, err
}

func GetDatabaseUsers(deployment_id, database_name, oauthToken string) ([]DatabaseUser, error) {
  body, err := rest_get(gopher_url("/" + deployment_id + "/" + database_name + "/users"), oauthToken)
  if err != nil {
    return make([]DatabaseUser, 0), err
  }
  var databaseUsersSlice []DatabaseUser
  err = json.Unmarshal(body, &databaseUsersSlice)
  return databaseUsersSlice, err
}

func CreateDatabaseUser(deploymentId, databaseName, username, password, oauthToken string) (OkResponse, error) {
  type UserCreate struct {
    Username string `json:"username"`
    Password string `json:"password"`
    ReadOnly bool `json:"readOnly"`
  }

  userCreate := UserCreate{Username: username, Password: password, ReadOnly: false}
  data, err  := json.Marshal(userCreate)
  if err != nil {
    return OkResponse{}, err
  }

  body, err := rest_post(gopher_url("/" + deploymentId + "/" + databaseName + "/users"), data, oauthToken)

  if err != nil {
    return OkResponse{}, err
  }
  var okResponse OkResponse
  err = json.Unmarshal(body, &okResponse)
  return okResponse, err
}

func RemoveDatabaseUser(deploymentId, databaseName, username, oauthToken string) (OkResponse, error) {
  body, err := rest_delete(gopher_url("/" + deploymentId + "/" + databaseName + "/users/" + username), oauthToken)

  if err != nil {
    return OkResponse{}, err
  }
  var okResponse OkResponse
  err = json.Unmarshal(body, &okResponse)
  return okResponse, err
}
