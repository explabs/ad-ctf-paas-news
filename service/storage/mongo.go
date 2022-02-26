package storage

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type NewsStruct struct {
	News []News `json:"news"`
}

func (n *NewsStruct) Load() {
	filename, _ := filepath.Abs("config.yml")
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, &n)
	if err != nil {
		panic(err)
	}
}

type News struct {
	Filename string `yaml:"filename"`
	Text     string `yaml:"text"`
	Round    int    `yaml:"round"`
	Image    string `yaml:"image"`
	Audio    string `yaml:"audio"`
}

func getMongoClient() (*mongo.Client, error) {
	credential := options.Credential{
		Username: "admin",
		Password: os.Getenv("ADMIN_PASS"),
	}
	var host, port = "mongo", 27017
	if os.Getenv("MODE") == "dev" {
		host = "localhost"
	}
	mongoURI := fmt.Sprintf("mongodb://%s:%d", host, port)
	clientOpts := options.Client().ApplyURI(mongoURI).SetAuth(credential)
	return mongo.Connect(context.TODO(), clientOpts)
}
func newsCollection() (*mongo.Collection, error) {
	client, err := getMongoClient()
	if err != nil {
		return nil, err
	}
	coll := client.Database("ad").Collection("news")
	return coll, nil
}

func (n *NewsStruct) UploadNews() error {
	coll, err := newsCollection()
	if err != nil {
		log.Fatal(err)
	}
	many, err := coll.DeleteMany(context.Background(), bson.D{})
	if err != nil {
		log.Println("error:", err)
		return err
	}
	log.Printf("delete %d news", many.DeletedCount)
	for _, news := range n.News {

		_, err = coll.InsertOne(context.Background(), news)
		if err != nil {
			log.Println("error:", err)
			return err
		}
	}
	return nil
}

func FindNews(round int) (*News, error) {
	coll, err := newsCollection()
	if err != nil {
		log.Fatal(err)
	}
	var news News
	findErr := coll.FindOne(context.TODO(), bson.M{"round": round}).Decode(&news)
	if findErr == mongo.ErrNoDocuments {
		return nil, nil
	}
	if findErr != nil {
		log.Fatal(findErr)
	}
	return &news, nil
}
