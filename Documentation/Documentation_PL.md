CerberCam
======

## 1 - Uruchomienie serwera
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