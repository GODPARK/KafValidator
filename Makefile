TARGET=kafvld
CONFIG_PATH=./config_simple_test.json

install: clean build

clean:
	rm -f ./$(TARGET)

run:
	go run main.go -config $(CONFIG_PATH)

build:
	go build -o $(TARGET)