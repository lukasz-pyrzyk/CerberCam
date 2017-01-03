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





## Uruchomienie serwera
CerberCam wymaga wystartowania dwóch usług - Receive oraz SendEmail. Są to dwa mikroserwisy.

#### Uruchomienie silnika do analizy zdjęć
```go
go run *.go -command=receive -config=../../../config.yaml
```
#### Uruchomienie wysyłania emaili
```go
go run *.go -command=sendEmail -config=../../../config.yaml
```

// TODO:
* opis projektu
* opis technologii (docker, go)
* opis tensorflowa
* opis klienta