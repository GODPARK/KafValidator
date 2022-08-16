TARGET=kafvld
CONFIG_PATH=./config_simple_test.json
SERVER_CONFIG_PATH=./config_server_test.json

install: clean build

clean:
	rm -f ./$(TARGET)

run:
	go run main.go -config $(CONFIG_PATH)

run_server:
	go run main.go -config ${SERVER_CONFIG_PATH}

build:
	go build -o $(TARGET)

build_ubuntu:
	-rm -f ./build/ubuntu/kafvld_ubuntu_x86_64.tar.gz
	-docker build -f build/ubuntu/Dockerfile -t kafvld-ubuntu .
	-docker run -d --name kafvld_ubuntu kafvld-ubuntu
	-docker cp kafvld_ubuntu:/root/kafvld_ubuntu_x86_64.tar.gz ./build/ubuntu/
	-docker stop kafvld_ubuntu
	-docker rm kafvld_ubuntu

build_centos:
	-rm -f ./build/ubuntu/kafvld_centos_x86_64.tar.gz
	-docker build -f build/centos/Dockerfile -t kafvld-centos .
	-docker run -d --name kafvld_centos kafvld-centos
	-docker cp kafvld_centos:/root/kafvld_centos_x86_64.tar.gz ./build/centos/
	-docker stop kafvld_centos
	-docker rm kafvld_centos

build_mac:
	-rm -f kafvld_mac_x86_64.tar.gz 
	go build -o kafvld
	tar -zcvf kafvld_mac_x86_64.tar.gz kafvld config_sample.json