# go-bbq-monitor
Written in Go this program interfaces with an Inkbird bluetooth thermometer, 
specifically the ibbq-4x, though other Inkbird models may work. Readings are pushed
to AWS IoT.

## Setup
Create an AWS IoT thing with certificates and a policy. Download the thing certificate 
and private key, placing them in the certs directory named `certificate.pem.crt` and
`private.pem.key`. You also need to edit `src/mqtt.go` and change `brokerUri` to your 
correct endpoint.

### Build
In `./src/` run `go install -v ./...`

### Run
`sudo LOGXI=main=INF ./go-bbq-monitor`

## Docker
The program can be run in docker if required:
```
docker build . -t go-bbq-monitor
sudo docker run -d --name bbq-monitor --net host --privileged go-bbq-monitor
```
The program will exit if the thermometer is disconnected or cannot be found.
To start on boot and restart on disconnect add the `--restart=always` option.