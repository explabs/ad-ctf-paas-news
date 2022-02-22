module github.com/explabs/ad-ctf-paas-exploits

go 1.16

require (
	github.com/explabs/ad-ctf-paas-api v1.0.8-0.20220222195202-059a960a5816
	github.com/rabbitmq/amqp091-go v1.3.0
	github.com/sirupsen/logrus v1.8.1
	go.mongodb.org/mongo-driver v1.8.3
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c
)

//replace github.com/explabs/ad-ctf-paas-api => ../ad-ctf-paas-api
