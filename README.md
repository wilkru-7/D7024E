## D7024E
This is a project for the LTU course D7024E, Mobile and distributed computing systems.
### Requirements
- Golang
- Docker
### How to setup the project
1. Start with cloning the repo
```bash
git clone https://github.com/wilkru-7/D7024E
```
2. Use the following command to build the project
```bash
docker build --tag your_name .
```
3. Use the following command to run the project
```bash
docker-compose up
```
### Testing
In order to run the test files, navigate to the ``d7024e`` folder and run the following command
```bash
go test-cover
```
