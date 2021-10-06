## D7024E
This is a project for the LTU course D7024E, Mobile and distributed computing systems.
## Requirements
- Golang
- Docker
## How to setup the project
1. Start with cloning the repo:
```bash
git clone https://github.com/wilkru-7/D7024E
```
2. Use the following command to build the project:
```bash
docker build --tag kademlia .
```
3. Use the following command to run the build:
```bash
docker-compose up
```

(To edit the number of containers used, open the ``docker-compose.yml`` file and edit ``replicas: X`` on line 38)
## Command Line Interface
To access the CLI from a container, first complete steps 1-3 of the setup process. Then open another terminal and use the following command:
```bash
docker attach d7024e_kademliaNodes_ID
```
The "ID" in the command is to be replaced by the number of the container you want to access.
### put
```bash
put yourtext
```
After running the put command, the hash of "yourtext" will be printed out.
### get
```bash
get yourhash
```
After running the get command, the text corresponding to "yourhash" will be printed out (if the get was successful).
### forget
```bash
forget yourhash
```
After running the forget command, the TTL of the data object corresponding to "yourhash" will no longer be updated and the object will eventually be removed.
### exit
```bash
exit
```
After running the exit command, the node terminates.
## RESTful HTTP Interface
### POST
### GET
## Testing
In order to run the test files, navigate to the ``d7024e`` folder and run the following command:
```bash
go test -cover
```
