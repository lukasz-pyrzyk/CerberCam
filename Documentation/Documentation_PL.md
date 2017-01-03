CerberCam
======

## 1. Opis projektu
CerberCam jest to system prywatnego monitoringu opartego o dowolną kamerę podłączoną do komputera (np. wbudowana lub na USB). Dodatkowo, CerberCam potrafi rozpoznawać anomalie w dostarczonym obrazie na podstawie jego wyuczonych wzorców.

## 2. Tryby pracy
System potrafi pracować w dwóch trybach
* ciągłego monitoringu danej przestrzeni i wykrywania anomalii
* ręcznego uczenia się wzorca z podanych obrazów


## 3. Architektura
CerberCam działa w dwóch obszarach
* aplikacji klienckiej
* aplikacji serwerowej (slave)

Klientem jest aplikacja desktopowa napisała w technologii JavaScript/HTML oparta o framework Electron oraz Nodejs. Pozwala to na uruchomienie jej na wielu systemach operacyjnych, np Windows, Linux, MacOS.

Z uwagi na wysokie wymagania wydajnościowe, komunikacja odbywa się za pomocą TCP/IP za pomocą własnego protokołu.

Część serwerowa to dwa mikroserwisy odpowiadające za analizę zdjęcia oraz przetwarzanie informacji zwrotnej do użytkownika.


## 4. Komunikacja i struktury danych

Aby zagwarantować wysoką skalowalność aplikacji, do komunikacji zostały wykorzystane system kolejkowania wiadomości - ``RabbitMQ``. Jest to open-sourcowy projekt pozwalający na tworzenie wydajnych systemów opratych o asynchroniczne wiadomości. 


W aplikacji wyróżniamy dwa rodzaje komunikatów:
* wiadomość (message)
* odpowiedź (response)

```
message Message {
  required string Email = 1;
  required bytes Photo = 2;
}
```

```
message Response {
  required string Email = 1;
  required string Label = 2;
  required float Probability = 3;
}
```

Model tutaj przedstawiony jest oparty o protokuł binarnej serialiacji danych ``Protobuf``, którego twórcą jest ``Google``.
Poprzednio wspomniane struktury danych są odpowiednio interpretowane w danym języku programowania, np w ``Go`` struktura Message przyjmuje formę:
```
type Message struct {
	Email            *string `protobuf:"bytes,1,req" json:"Email,omitempty"`
	Photo            []byte  `protobuf:"bytes,2,req" json:"Photo,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}
```

# 6. Serwer
Z uwagi na preferencje takie jak statyczne typowanie, kompilacjedo kodu natywnego czy wysoką wydajność, aplikacja serwerowa została napisana w języku ``Go``, potocznie zwanym ``Golang``.
Jest to nowoczesny język programowania stworzony i aktywnie używany przez ``Google``. 
Aplikacja działa dwóch trybach, które uruchomione osobno tworzą dwa mikroserwisy.
* Receive - przetwarzanie danych
* SendEmail - przetwarzanie odpowiedzi i alertów

Uruchomienie tychże usług wygląda następująco:
```go
go run *.go -command=receive -config=../../../config.yaml
```
```go
go run *.go -command=sendEmail -config=../../../config.yaml
```

Aby zapewnić swobodną konfigurację, użyty został format ``YAML``, czyli ``Yet another markup language``.
Przykładowy plik konfiguracyjny wygląda następująco:

```yaml
tensorflow:
  modeldir: "model" # working directory for tensorflow
  host: http://hostName.com:8888 
queue:
  host: amqp://login:password@host.com:5672/
  requests: requests
  responses: responses
email:
  host: smtp.host.com
  port: 587
  login: login@host.com
  password: password
```

Taki format jest zgodny ze strukturą
```
type config struct {
	Tensorflow tensorflowConfig
	Queue      queueConfig
	Email      emailConfig
}

type queueConfig struct {
	Host      string
	Requests  string
	Responses string
}

type tensorflowConfig struct {
	ModelDir string
	Host     string
}

type emailConfig struct {
	Host     string
	Port     int
	Login    string
	Password string
}
```
W celu zagwarantowania dostępności, każdy z mikroserwisów pracuje w pętli głównej o czasie 1000ms.

### Receive
Jest to komenda która odpytuje kolekję o nowe wiadomości. 
Kiedy takowa się pojawi, zostaje pobrana i oznaczona jako gotowa do przetworzenia.
Następnym z kroków jest deserializacja wiadomości z tablicy bajtów zserializowanych przez Protobuf do struktury danych. 
Wykorzystuje się do tego nastąpujący fragment kodu:
```
msg := &Message{}
err := proto.Unmarshal(d.Body, msg)
```
Po deserializacji, wiadomość ``msg`` zostaje przesłana do silnika analizy obrazów - ``Tensorflow``.
W momencie udanej analizy, komenda wysyła ponownie serializuje dane i wysyła resultat do następnej kolejki danych w celu odesłania jej użytkownikowi, np. za pomocą maila lub SMSa.


### SendEmail
W porównaniu do poprzedniej komendy, ta jest stosunkowo prosta.
Pobiera one nowe wiadomości do opublikowania do użytkowników i wysyła je poprzez protokuł SMTP.

Cała komenda wygląda następująco:
```
func HandleSendEmail() {
	log.Info("Checking for new data in emails queue...")

	emailManager := NewEmailManager()
	queue := queueManager{}
	msgs := queue.Receive(GlobalConfig.Queue.Responses)
	for d := range msgs {

		msg := &Response{}
		err := proto.Unmarshal(d.Body, msg)
		failOnError(err, "Cannot deserialize response")

		content := fmt.Sprintf("Cerber believes that your picture shows %s (probability %f%%)", *msg.Label, *msg.Probability)
		emailManager.Send(*msg.Email, "Reconginion results", content)
	}
}
```

## 7. Silnik analizy - Tensorflow


// TODO:
* opis technologii (docker, go)
* opis tensorflowa