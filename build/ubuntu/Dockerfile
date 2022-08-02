FROM ubuntu:focal

ENV WORK_DIR=/root/work/source/ \
    GO_BIN=/root/work/go/bin \
    TARGET_NAME=kafvld  \
    TAR_NAME=kafvld_ubuntu_x86_64.tar.gz \
    GO_VERSION=go1.16.5.linux-amd64.tar.gz

RUN apt-get update -y
RUN apt-get -y upgrade

RUN apt-get install wget git -y

RUN mkdir -p /root/work

RUN cd /root/work && wget https://golang.org/dl/${GO_VERSION}
RUN cd /root/work && tar -xvf ${GO_VERSION}

RUN apt-get install -y make gcc

RUN mkdir -p ${WORK_DIR}
COPY ./ ${WORK_DIR}
RUN cd ${WORK_DIR} && ${GO_BIN}/go install 
RUN cd ${WORK_DIR} && ${GO_BIN}/go build -o ${TARGET_NAME}

RUN cd ${WORK_DIR} && tar -zcvf ${TAR_NAME} ${TARGET_NAME} config_sample.json
RUN mv ${WORK_DIR}/${TAR_NAME} /root/
RUN touch /root/sleep.txt

ENTRYPOINT ["tail", "-f", "/dev/null"]