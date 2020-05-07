# Album.com

Album.com - Sample CRUD application using GO, shows notification in case of delete operation.  
Album.com provides api for basic CRUD operation [fetch, insert, delete  all albums, images contained within a album under a specified directory].

## Prerequisites 
[Install GO](https://golang.org/doc/install) - Refer the link to install go lang
[Install Docker](https://docs.docker.com/) - Download and install Docker with nsqd

## Installation
Extract the project to go src

```bash
$gopath -> $src -> $album.com
go get
```
## Run Docker mode 
cd $application_folder ->  
cd $albumpath: docker build -t album .
cd $albumpath/messageclient: docker build -t msgcli .

docker run -it --rm --name nsqd -p 4150:4150 -p 4151:4151 nsqio/nsq /nsqd
docker run -it --rm --name albumcli msgcli
docker run -it --rm --name album -p 9000:9000 album




Sample message data.
-------------------
If u want to attach use this  (screen shot alos i have attached)
2020/05/07 05:18:29 Connect to the NSQD[172.17.0.1:4150], topic[GALLERY]
2020/05/07 05:18:29 INF    1 [GALLERY/ch] (172.17.0.1:4150) connecting to nsqd
2020/05/07 05:21:32 Seq: [1], Message: Album[2020] created successfully
2020/05/07 05:21:35 Seq: [2], Message: Image[test-image.png] uploaded successfully to the album[2020]
2020/05/07 05:21:43 Seq: [3], Message: Image[test-image.png] delted successfully from the album[2020]
2020/05/07 05:21:53 Seq: [4], Message: Album[2020] deleted successfully.

**NOTE**
Nsqd ip address: "172.17.0.1:4150"(hard coded for docker).(refer client.go) to Communicate between containers. Try to ping nsqd running container and used it instead of 172.17.0.1 in messageQueue server, client


## Storage
```bash
For storage of albums isn't auto generated a folder named gallery is already present 
```   

## Concept

Simple CRUD operation which will create, fetch and delete 

```bash
type Album  struct {
	AlbumName    string 	`json:"albumName"`
	Images     []string		`json:"imageName"`
}
```

## API's
```bash
POST "/album" - Adds a new album directory
DELETE "/album" - Deletes an album directory

POST "/image" - Adds a new image under album directory
DELETE "/image" - Deletes an image under album directory

GET  "/image" - Get an image under album directory
GET  "/images" - Get list of images under album directory
```

## Swagger configuration 
```bash
go get -u github.com/go-swagger/go-swagger/cmd/swagger under root directory
Refer swagger.yaml to try out the api end points

swagger serve -F=swagger swagger.yaml - opens swagger UI under your browser as host specified as swagger.yaml

```