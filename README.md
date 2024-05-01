# Dns LookUp Application

# Dependencies

## Ubuntu

1) 
``` 
    $ apt update && apt install docker.io make
    $ sudo snap install --classic go
```
2) 
```
    $ sudo apt-get update

    $ sudo apt-get install apt-transport-https ca-certificates curl gnupg lsb-release git
```
3) 
```
    $ curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
```
4) 
```
    $ echo "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

    $ sudo apt-get update
```
5) 
```
$ sudo apt-get install docker-ce docker-ce-cli containerd.io
```

## windows 

1) download docker desktop from https://docs.docker.com/desktop/install/windows-install/

2) set up docker compose 

3) download go from https://golang.org/dl/

4) download make from https://sourceforge.net/projects/gnuwin32/files/make/3.81/make-3.81.exe/download

# Setup

1) Clone the Repository with 
``` 
$ git clone git@github.com:Dev-Heaven69/db.git
```
2) cd into the repository
```
$ cd db
```
3) run the make command to build the app
##### UBUNTU 
```
$ make run
```
##### Windows
```
$ make run
```

# Endpoints

1) send csv file as form-data and get back personal data from db

```
$ curl --location 'localhost:5000/getPersonalEmail' \
--form 'csv=@<CSV FILE>' \
--form 'responseType=<JSON OR CSV>' \
--form 'discordUsername=<GITHUB USERNAME>' \
--form 'email=<EMAIL>'
```

2) send csv file as form-data and get back professional data from db

```
$ curl --location 'localhost:5000/getProfessionalEmails' \
--form 'csv=@<CSV FILE>' \
--form 'responseType=<JSON OR CSV>' \
--form 'discordUsername=<GITHUB USERNAME>' \
--form 'email=<EMAIL>'
```

3) send csv file as form-data and get back both data from db

```
$ curl --location 'localhost:5000/getBothEmails' \
--form 'csv=@<CSV FILE>' \
--form 'responseType=<JSON OR CSV>' \
--form 'discordUsername=<GITHUB USERNAME>' \
--form 'email=<EMAIL>'
```

4) send csv file as form-data and get back both data from db

```
$ curl --location 'localhost:5000/scandb' \
--form 'csv=@<CSV FILE>' \
--form 'responseType=<JSON OR CSV>' \
--form 'discordUsername=<GITHUB USERNAME>' \
--form 'email=<EMAIL>'
```

5) change webhook url

```
$ curl --location 'localhost:5000/changewebhook' \
--header 'Content-Type: application/json' \
--data '{
    "url":<URL>
}'
``` 
# Docker Image

1) run this docker command to run the app on the go 
```
$ sudo docker run -p 5000:5000 seew0/db:1.0  
```