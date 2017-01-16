CerberCam - system monitoringu domowego 
======

| Autorzy          | Jakub Bentkowski, Łukasz Pyrzyk |
|------------------|--------------------------------|
| Uczelnia 		   | Politechnika Opolska           |
| Kierunek studiów | informatyka                    |
| Typ              | niestacjonarne                 |
| Rok studiów      | IV                             |
| Semestr          | VII                            |
| Rok akademicki   | 2016/2017                      |

## 1. Opis projektu
CerberCam jest to system prywatnego monitoringu opartego o dowolną kamerę podłączoną do komputera (np. wbudowana lub na USB). Dodatkowo, CerberCam potrafi rozpoznawać anomalie w dostarczonym obrazie na podstawie jego wyuczonych wzorców.

![Architektura systemu](./system.png)

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

W celu skompilowania plików ``Proto`` do modeli dla wybranego języka programowania należy wykonać polecenie
```bash
protoc --go_out=. *.proto
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
```go
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

Obsługa protokołu SMTP w języku Go jest bardzo prosta, gdyż w standardzie dysponujemy pakietem ``net/smtp``.
W celu hermetyzacji komponentu, zdecydowaliśmy się wyekstrachować komponent zwany EmailManagerem do osobnego typu.
Jest to prosta struktura która czyta dane do serwera z konfiguracji.

```go
type emailManager struct {
	Login    string
	Password string
	Host     string
	Port     int
}

func NewEmailManager() *emailManager {
	manager := new(emailManager)
	manager.Login = GlobalConfig.Email.Login
	manager.Password = GlobalConfig.Email.Password
	manager.Host = GlobalConfig.Email.Host
	manager.Port = GlobalConfig.Email.Port
	return manager
}
```
Dla osoby która nie posiada wiedzy na temat języka GO poprzednio zademonstrowana konstrukcja moze wyglądać dość niespodziewanie. 
Golang nie posiada konceptu konstruktowyów - wedle konwencji metoda do tworzenia instancji powinna być złożona ze słowa ``New`` i nazwy typu, czyli w tym przypadku ``NewEmailManager``.
Oprócz pól, manager Udostępnia jedną, publiczną metodę która pozwala na wysyłanie wiadomości.

```go
func (manager emailManager) Send(to string, subject string, content string) {
	auth := smtp.PlainAuth("", manager.Login, manager.Password, manager.Host)
	recipients := []string{to}

	message := fmt.Sprintf("To: %s \r\n"+
		"Subject: %s !\r\n"+
		"\r\n"+
		"%s \r\n", to, subject, content)

	msg := []byte(message)
	hostString := manager.Host + ":" + strconv.Itoa(manager.Port)
	err := smtp.SendMail(hostString, auth, manager.Login, recipients, msg)
	failOnError(err, "Cannot send email")
}
```

## 7. Silnik analizy - Tensorflow
Tensorflow jest to silnik ``machine learning`` stworzony przez ``Google``. Jego kod jest dostępny publicznie pod adresem github.com/tensorflow/tensorflow. Jest on powszechnie używany przez kilka znanych usług giganta z Mountan View, takich jak Google Photos, analizie mowy, rozpoznawaniu obrazów czy też pisma odręcznego. Udostępniona przez Google biblioteka pozwala na zaawansowane obliczenia numeryczne z wykorzystaniem grafów. Współpracuje ona z klastrami superkomputerów jak i pojedyńczą instancją uruchomioną na dowolnej stacji roboczej czy telefonie z systemem Android. Aby uruchomić Tensorflow lokalnie, potrzebny jest komputer z systemem Linux lub MacOS. Co warte podkreślenia, Tensorflow pozwala na zrównoleglenie operacji poprzez używanie procesorów z kart graficznych NVIDII obsługujących framework CUDA.

## 8. Współpraca Tensorflow z CerberCam
W celu zintegrowania Tensorflow musieliśmy wykonać kilka zaawansowanych zadań, między innymi skompilować cały projekt lokalnie.
Pierwszym z nich było pobranie generatora pozwalającego na uzyskanie wrappera pomiędzy językiem C a Go.
```bash
go get -d github.com/tensorflow/tensorflow/tensorflow/go
```

Aby skompilować projekt, należy użyć narzędzia ``Bazel``, czyli systemu budującego projekty napisanego i wykorzystywanego przez Google. W przypadku Ubuntu 15.10 (Willy) należy wykonać kilkanaście poleceń:

Upewnić się, że installer dla Java 8 jest zainstalowany

```bash
sudo add-apt-repository ppa:webupd8team/java
sudo apt-get update
sudo apt-get install oracle-java8-installer
```

Następnie należy dodać źródła dla narzędzia ``Bazel``
```bash
echo "deb [arch=amd64] http://storage.googleapis.com/bazel-apt stable jdk1.8" | sudo tee /etc/apt/sources.list.d/bazel.list
curl https://bazel.build/bazel-release.pub.gpg | sudo apt-key add -
```

W tym momencie możemy uruchomić process instalacji oraz aktualizacji do najnowsze wersji. 
Łączny czas trwania tych skryptów może zająć do kilkudziesięciu minut.
```bash
sudo apt-get update && sudo apt-get install bazel
sudo apt-get upgrade bazel
```

Gdy Bazel jest już gotowy, mozemy przejść do procesu kompilowania Tensorflow. W tym celu należy wykonać trzy poniższe polecenia:
```bash
cd ${GOPATH}/src/github.com/tensorflow/tensorflow
./configure
bazel build -c opt //tensorflow:libtensorflow.so
```

Z uwagi na rozmiar oraz ilość plików zawartych w solucji, na komputerze z procesorem Intel i3 2.6GHz i 12GB pamięci ram, process ten trwał 51 minut. 

Wynikiem operacji jest plik ``libtensorflow.so`` czyli biblioteka, którą jesteśmy zainteresowani.
Aby CerberCam mógł jej użyć, musimy sprawić, by była ona widoczna dla Linkera. Najprostszym rozwiązaniem będzie dodanie jej do folderu ``/usr/local/lib``. Operację kopiowania można wykonać prostym poleceniem
```bash
cp ${GOPATH}/src/github.com/tensorflow/tensorflow/bazel-bin/tensorflow/libtensorflow.so /usr/local/lib
```

Ostatnim krokiem będzie samo wygenerowanie nagłówków dla Go, czyli wywołanie polecenia
```bash
go generate github.com/tensorflow/tensorflow/tensorflow/go/op
```

Jeśli wszystko przebiegło pomyśle, Tensorflow jest gotowy do współpracy z CerberCam. Możemy to zweryfikować uruchamiając poniższy test
```bash
go test -v github.com/tensorflow/tensorflow/tensorflow/go
```

## 8. Wdrożenie
Proces wdrażania aplikacji serwerowej przebiega w trzech krokach
* Uruchom nową instancję kolejek dla wiadomości, jeśli takowe nie istnieją. 
* Uruchom CerberCam z komendą ``receive``
* Uruchom CerberCam z komenda ``sendEmail``.

Co ważne, wdrażanie nowej wersji nie powoduje niedostępności systemu. CerberCam może zostać wyłączony, a przychodzące wiadomości będą nadal gromadzone na kolejkach.

Całość deploymentu odbywa się za pomocą usługi ``Docker``, ułatwiającej korzystanie z linuxowych kontenerów, będących odseparowanymi od siebie środowiskami, zwirtualizowanymi na poziomie systemu operacyjnego.
Aby zbudować własny obraz Dockera, który będzie można wykorzystywać wielokrotnie na różnych maszynach, należy stworzyć plik konfiguracyjny  ``Dockerfile``. Dla CerberCam wygląda on następująco:

```yaml
# bazuj na obrazie Linuxa zawierającym Tensorflow
FROM bentou/tensorflowgo:latest 

# dodaj informacje o autorach
MAINTAINER Lukasz Pyrzyk <lukasz.pyrzyk@gmail.com>, Jakub Bentkowski <bentkowski.jakub@gmail.com>

# skopiuj wszystkie pliki
COPY ./Src/Server /go/src/Cerber

# zainstaluj zależności
RUN go get github.com/op/go-logging & go get github.com/streadway/amqp & go get github.com/golang/protobuf/proto & go get gopkg.in/mgo.v2 & go get gopkg.in/yaml.v2 & go get -d github.com/tensorflow/tensorflow/tensorflow/go

# skopiuj potrzebne pliki dla Tensorflow
COPY ./Patches $GOPATH/src/github.com/tensorflow/tensorflow/tensorflow/go

# dodaj potrzebne repozytoria i pakiety
RUN add-apt-repository ppa:maarten-fonville/protobuf
RUN apt-get update && apt-get install -y --allow-unauthenticated protobuf-compiler python-protobuf

# wygeneruj nagłówki dla Golanga
RUN go get github.com/tensorflow/tensorflow/tensorflow/go/op & go generate -v -x github.com/tensorflow/tensorflow/tensorflow/go/op

# zainstaluj aplikacje w kontenerze
RUN go install Cerber

# wystaw entrypoint dla stdin
ENTRYPOINT ["/go/bin/Cerber"]
```

Po stworzeniu obrazu i opublikowaniu go w globalnym repozytorium, CerberCama można zainstalować za pomocą poniższego polecenia
```bash
docker run lukaszpyrzyk/cerbercam
```

W celu większej automatyzacji procesu wdrażania aplikacji oraz kolejek, stworzyliśmy plik ``docker-compose``, który pozwala uruchomić kilka obrazów Dockera na raz za pomocą jednego polecenia.
```yaml
rabbitmq:
  image: rabbitmq:3.6.5-management
  hostname: rabitmqhost
  ports:
    - 8080:15672
    - 5672:5672
TensorFlow:
  image: gcr.io/tensorflow/tensorflow
  ports:
    - "8888:8888"
  restart: always
CerberReceive:
  image: lukaszpyrzyk/cerbercam
  command: "-command=receive"
CerberSend:
  image: lukaszpyrzyk/cerbercam
  command: "-command=sendEmail"
```

Dzięki temu, uruchomienie całego systemu jest możliwe za pomocą jednego polecenia
```bash
docker-compose up
```

## 9. Chmura
Platformą hostingową dla projektu jest ``Microsoft Azure``. Serwer jest uruchomiony na wirtualnej maszynie (VM) z systemem Ubuntu Server 16.04 o mocy 3.5GB pamięci RAM, dwóch rdzeni i dysku SSD. Koszt takiej maszyny to 50,82 euro miesięcznie. Z uwagi na potrzebę wysokiej wydajności, zdecydowaliśmy się wybrać centrum danych w Europie Zachodniej, co pozwoli ograniczyć czas tranferu danych. 