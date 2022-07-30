# KafValidator
> * go 언어로 개발
> * kafka lib 로는 [Kafka Go Client](https://docs.confluent.io/kafka-clients/go/current/overview.html) 을 사용하여 개발
> * [Kafka go Client API Docs](https://docs.confluent.io/platform/current/clients/confluent-kafka-go/index.html#Consumer)
> * [Kakfa Docs](https://kafka.apache.org/documentation/)

### 1) 목적
Kafka 구성 시 정상 동작 여부를 확인하기 위한 기본적인 Test Application 입니다.

### 2) 주요 기능
* Pub/Sub 기능 검증
* Pub/Sub 평균 소요 시간 측정 (단일 쓰레드 기준)

### 3) 주의 사항!!!!!
#### 반드시! 테스트 토픽에 진행하길 권고 합니다.
> 주의) 서비스용 토픽에 진행할 경우 consumer group이 다르면 이슈가 없을 수 있으나(장담할 수는 없습니다.) 혹여 consumer group이 같았을 경우 서비스 메시지가 commit되는 이슈가 있을 수 있습니다.

### 4) Download
Release 에서 알맞는 버전을 다운로드 하여 실행


### 5) 실행 방법
```bash
# run with config file
$ kafvld -config ${CONFIG_PATH}
```

### 6) 설정 파일
```json
{  
    // kafka broker address (list)
   "bootstrapServer":[
      "localhost:9092"
   ],
   // kafka test topic
   "topic":"test-topic",
   // producer
   "producer":{
        // producer client id
      "clientId":"test-producer",
        // producer acks
        // 0 : broker does not send any response or acks
        // -1 or all: all broker committed ( sync replica )
        // 1 : one broker committed
      "acks":"1"
   },
   // consumer
   "consumer":{
        // consumer group id
      "groupId":"test-consumer-group"
   },
   // test info
   "simple":{
        // send msg count ex) 10, 100 ..
      "msgCount":100
   }
}
```