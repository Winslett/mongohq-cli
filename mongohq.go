package main

import (
  "fmt"
  "os"
  "github.com/codegangsta/cli"
  "github.com/MongoHQ/mongohq-cli"
  "github.com/MongoHQ/controllers"  // MongoHQ CLI functions
)

func requireArguments(command string, c *cli.Context, argumentsSlice []string, errorMessages []string) {
  err := false

  for _, argument := range argumentsSlice {
    if !c.IsSet(argument) {
      err = true
      fmt.Println("--" + argument + " is required")
    }
  }

  if err {
    fmt.Println("\nMissing arguments, for more information, run: mongohq " + command + " --help\n")
    for _, errorMessage := range errorMessages {
      fmt.Println(errorMessage)
    }
    os.Exit(1)
  }
}

func main() {
  app := cli.NewApp()
  app.Name = "mongohq"
  app.Usage = "Allow MongoHQ interaction from the commandline (enables awesomeness)"
  app.Before = controllers.RequireAuth
  app.Version = mongohq_cli.Version()
  app.Commands = []cli.Command{
    {
      Name:      "databases",
      Usage:     "list databases",
      Action: func(c *cli.Context) {
        controllers.Databases()
      },
    },
    {
      Name:      "databases:create",
      Usage:     "create database on an existing deployment",
      Flags:     []cli.Flag {
        cli.StringFlag { "deployment,dep", "<string>", "Deployment to create database on"},
        cli.StringFlag { "database,db", "<string>", "Name of new database to create"},
      },
      Action: func(c *cli.Context) {
        requireArguments("databases:create", c, []string{"deployment", "database"}, []string{})
        controllers.CreateDatabase(c.String("deployment"), c.String("database"))
      },
    },
    {
      Name:      "databases:info",
      Usage:     "information on database",
      Flags:     []cli.Flag {
        cli.StringFlag { "database,db", "<string>", "Database name for more information"},
      },
      Action: func(c *cli.Context) {
        requireArguments("databases:info", c, []string{"database"}, []string{})
        controllers.Database(c.String("database"))
      },
    },
    {
      Name:      "deployments",
      Usage:     "list deployments",
      Action: func(c *cli.Context) {
        controllers.Deployments()
      },
    },
    {
      Name:      "deployments:create",
      Usage:     "create a new Elastic Deployment",
      Flags:     []cli.Flag {
        cli.StringFlag { "database,db", "<string>", "New database name to be created on your new deployment"},
        cli.StringFlag { "region,r", "<string>", "Region for deployment. For a list of regions, run 'mongohq regions'"},
      },
      Action: func(c *cli.Context) {
        requireArguments("deployments:create", c, []string{"database", "region"}, []string{})
        controllers.CreateDeployment(c.String("database"), c.String("region"))
      },
    },
    {
      Name:      "deployments:info",
      Usage:     "information on deployment",
      Flags:     []cli.Flag {
        cli.StringFlag { "deployment,dep", "<bson_id>", "The id for the deployment for more information"},
      },
      Action: func(c *cli.Context) {
        requireArguments("deployments:info", c, []string{"deployment"}, []string{})
        controllers.Deployment(c.String("deployment"))
      },
    },
    {
      Name:      "deployments:mongostat",
      Usage:     "realtime mongostat",
      Flags:     []cli.Flag {
        cli.StringFlag{"deployment,dep", "<bson_id>", "The id for the deployment for tailing mongostats"},
      },
      Action: func(c *cli.Context) {
        requireArguments("deployments:mongostat", c, []string{"deployment"}, []string{})
        controllers.DeploymentMongoStat(c.String("deployment"))
      },
    },
    {
      Name:      "deployments:logs (pending)",
      Usage:     "tail logs",
      Flags:     []cli.Flag {
        cli.StringFlag{"deployment,dep", "<bson_id>", "The id for the deployment for tailing logs"},
      },
      Action: func(c *cli.Context) {
        requireArguments("deployments:logs", c, []string{"deployment"}, []string{})
        fmt.Println("Pending")
      },
    },
    {
      Name:      "deployments:oplog",
      Usage:     "tail oplog",
      Flags:     []cli.Flag {
        cli.StringFlag{"deployment,dep", "<bson_id>", "The id for the deployment to tail an oplog"},
      },
      Action: func(c *cli.Context) {
        requireArguments("deployments:oplog", c, []string{"deployment"}, []string{})
        controllers.DeploymentOplog(c.String("deployment"))
      },
    },
    {
      Name:      "regions",
      Usage:     "list available regions",
      Action: func(c *cli.Context) {
        controllers.Regions()
      },
    },
    {
      Name:      "users",
      Usage:     "list users on a database",
      Flags:     []cli.Flag {
        cli.StringFlag { "deployment,dep", "<bson_id>", "The deployment id the database is on"},
        cli.StringFlag { "database,db", "<string>", "The specific database to list users"},
      },
      Action: func(c *cli.Context) {
        requireArguments("users", c, []string{"deployment", "database"}, []string{})
        controllers.DatabaseUsers(c.String("deployment"), c.String("database"))
      },
    },
    {
      Name:      "users:create",
      Usage:     "add user to a database",
      Flags:     []cli.Flag {
        cli.StringFlag { "deployment,dep", "<bson_id>", "The deployment id the database is on"},
        cli.StringFlag { "database,db", "<string>", "The database name to create the user on"},
        cli.StringFlag { "username,u", "<string>", "The new user to create"},
      },
      Action: func(c *cli.Context) {
        requireArguments("users:create", c, []string{"deployment", "database", "username"}, []string{})
        controllers.DatabaseCreateUser(c.String("deployment"), c.String("database"), c.String("username"))
      },
    },
    {
      Name:      "users:remove",
      Usage:     "remove user from database",
      Flags:     []cli.Flag {
        cli.StringFlag { "deployment,dep", "<bson_id>", "The deployment id the database is on"},
        cli.StringFlag { "database,db", "<string>", "The database name to remove the user from"},
        cli.StringFlag { "username,u", "<string>", "The user to remove from the deployment"},
      },
      Action: func(c *cli.Context) {
        requireArguments("users:remove", c, []string{"deployment", "database", "username"}, []string{})
        controllers.DatabaseRemoveUser(c.String("deployment"), c.String("database"), c.String("username"))
      },
    },
    {
      Name:      "logout",
      Usage:     "remove stored auth",
      Action:    func(c *cli.Context) {
        controllers.Logout()
      },
    },
  }

  app.Run(os.Args)
}
