package compose

import (
	"gopkg.in/yaml.v3"
)

func addKafka(m map[string]any) {
	var z map[string]any

	yaml.Unmarshal([]byte(`
    image: bitnami/zookeeper:latest
    restart: on-failure
    ports:
      - "2181:2181"
    environment:
      ZOO_MY_ID: 1
      ZOO_PORT: 2181
      ZOO_SERVERS: server.1=zookeeper:2888:3888
      ALLOW_ANONYMOUS_LOGIN: "yes"`), &z)

	var k map[string]any

	yaml.Unmarshal([]byte(
		`
image: bitnami/kafka:latest
restart: on-failure
ports:
  - "9092:9092"
environment:
  KAFKA_ADVERTISED_LISTENERS: INTERNAL://kafka:29092,EXTERNAL://localhost:9092
  KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: INTERNAL:PLAINTEXT,EXTERNAL:PLAINTEXT
  KAFKA_INTER_BROKER_LISTENER_NAME: INTERNAL
  KAFKA_ZOOKEEPER_CONNECT: "zookeeper:2181"
  KAFKA_BROKER_ID: 1
  KAFKA_LOG4J_LOGGERS: "kafka.controller=INFO,kafka.producer.async.DefaultEventHandler=INFO,state.change.logger=INFO"
  KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
  KAFKA_CFG_ZOOKEEPER_CONNECT: zookeeper:2181
  ALLOW_PLAINTEXT_LISTENER: "yes"
  KAFKA_LISTENERS: "INTERNAL://:29092,EXTERNAL://:9092"
  KAFKA_ZOOKEEPER_SESSION_TIMEOUT: "6000"
  KAFKA_RESTART_ATTEMPTS: "10"
  KAFKA_RESTART_DELAY: "5"
  ZOOKEEPER_AUTOPURGE_PURGE_INTERVAL: "0"
depends_on:
  - zookeeper`), &k)

	var o map[string]any

	yaml.Unmarshal([]byte(`
    image: quay.io/cloudhut/kowl:v1.4.0
    restart: on-failure
    volumes:
    - ./kowl_config:/etc/kowl/
    ports:
    - "8080:8080"
    entrypoint: ./kowl --config.filepath=/etc/kowl/config.yaml
    depends_on:
      - kafka
`), &o)

	m["zookeeper"] = z
	m["kafka"] = k
	m["kowl"] = o

}
