package main

import (
"flag"
"fmt"
"log"
"os"

"github.com/golang-migrate/migrate/v4"
_ "github.com/golang-migrate/migrate/v4/database/postgres"
_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
var direction string
flag.StringVar(&direction, "direction", "up", "migration direction (up or down)")
flag.Parse()

dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
os.Getenv("DB_USER"),
os.Getenv("DB_PASSWORD"),
os.Getenv("DB_HOST"),
os.Getenv("DB_PORT"),
os.Getenv("DB_NAME"),
)

m, err := migrate.New(
"file://migrations",
dbURL,
)
if err != nil {
log.Fatal(err)
}

switch direction {
case "up":
if err := m.Up(); err != nil && err != migrate.ErrNoChange {
log.Fatal(err)
}
case "down":
if err := m.Down(); err != nil && err != migrate.ErrNoChange {
log.Fatal(err)
}
default:
log.Fatal("invalid direction")
}
}
