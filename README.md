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
   * 첫 값 제외 평균 측정 가능 ( 첫 값의 경우 kafka 와의 오버헤드로 인해 보다 많은 소요시간이 걸릴 수 있습니다. )
* Pub/Sub 단일 건수 최대/최소 시간 측정 ( 단일 쓰레드 기준 )
   * 첫 값 제외한 단일 pub/sub 중 최대 소요 시간 측정 가능
   * pub/sub 최소 소요 시간 측정 가능
* 메세지 사이즈 조정하여 테스트 가능

### 3) 주의 사항!!!!!
#### 반드시! 테스트 토픽에 진행하길 권고 합니다.
> 주의) 서비스용 토픽에 진행할 경우 consumer group이 다르면 이슈가 없을 수 있으나(장담할 수는 없습니다.) 혹여 consumer group이 같았을 경우 서비스 메시지가 commit되는 이슈가 있을 수 있습니다.

### 4) Download
[Release](https://github.com/GODPARK/KafValidator/releases) 에서 알맞는 버전을 다운로드 하여 실행


### 5) 실행 방법
```bash
# 압축 해제 후
$ tar -xvzf kafvld_VERSION.tar.gz ./

# 압축 해제 후 결과, 이후 config 수정 후 실행
$ ls
kafvld
config_sample.json

# run with config file
$ kafvld -config ${CONFIG_PATH}
```

### 6) 설정 파일
```json
{  
    // kafka broker address (list)
    // ex) kafka 서버 입력 3대일 경우 아래 내용 참조
   "bootstrapServer":[
      "kafka01:9092",
      "kafka02:9092",
      "kafka03.9092"
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
        // 권장: 1
      "acks":"1",
      // default msg : 51byte
      // Sets the bytes to be appended to the send message.
      // ex) if valueByteSize = 1000 --> The size of one message is 1000 + 51 = 1051 byte
      // 기본 메세지는 51 byte 이며, 51 + ? = byte 설정할 때 사용합니다.
      // ex) 만약 valueByteSize 가 1000 일 경우 --> 전체 1000 + 51 = 1051 byte
      "valueByteSize": 1000
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

### 7) 실행 결과 예시
```bash
TEST KEY: 08fc3467-a02d-480d-964c-f7f557e4493d # 테스트 키 ( 테스트 마다 다름 )
 Running... please waiting


################### Reslut ####################
[CONFIG] Kafka Broker: localhost:9092 # 카프카 브로커 리스트
[CONFIG] Target Topic: test-topic # 토픽 명
[CONFIG] Total Msg Count: 10 # 테스트 메시지 발행 갯수
[CONFIG] Msg Size: 1051 byte # 한 건당 메세지 수
[CONFIG] Total Msg Size: 10510 byte # 전송한 총 메세지 크기 ( 한 건당 메세지 수 x  메세지 발생 수)

[PRODUCER] Msg Pub Success Count: 10 # 메세지 pub 성공 수
[PRODUCER] Msg Pub Fail Count: 0  # 메세지 pub 실패 수
[PRODUCER] Send Msg Average ( Enclude First ): 107 ms #(첫 값 포함) 메세지 전송 평균 소요 시간
[PRODUCER] Send Msg Average ( Exclude First ): 9 ms #(첫 값 제외) 메세지 전송 평균 소요 시간
[PRODUCER] Send First Msg Time: 994 ms # 첫 값의 메세지 전송 소요시간
[PRODUCER] Send Msg Max Time ( Exclude First ) : 13 ms # (첫 값 제외) 메세지 전송 건수 중 최대 소요 시간
[PRODUCER] Send Msg Min Time: 9 ms # 메세지 전송 건수 중 최소 소요 시간

[CONSUMER] Msg Sub Success Count: 10 # 메세지 sub 성공 수
[CONSUMER] Msg Sub Fail Count: 0 # 메시지 sub 실패 수
[CONSUMER] Msg Sub Etc Count: 0 # 현재 테스트와 무관한 메시지 수신시 etc 카운트 증가 (보통 consumer group 변경될 경우 커밋되지 않은 메세지 유입 가능)
[CONSUMER] Pub->Sub Msg Time Average ( Enclude First ): 866 ms # (첫 값 포함) 메세지 전송->수신 평균 소요 시간
[CONSUMER] Pub->Sub Msg Time Average ( Exclude First ): 850 ms # (첫 값 제외) 메세지 전송->수신 평균 소요 시간
[CONSUMER] Pub->Sub First Msg Time: 1009 ms # 첫 메세지 전송->수신 소요 시간
[CONSUMER] Pub->Sub Msg Time Max ( Exclude First ) : 925 ms # (첫 값 제외 )메세지 전송->수신 최대 소요 시간
[CONSUMER] Pub->Sub Msg Time Min: 62 ms # 메세지 전송->수신 최소 소요 시간
#################################################
```