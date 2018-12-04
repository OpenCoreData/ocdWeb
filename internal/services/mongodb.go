package services

import (
	"os"

	"gopkg.in/mgo.v2"
)

func GetMongoCon() (*mgo.Session, error) {
	host := os.Getenv("MONGO_HOST")

	return mgo.Dial(host)
}
