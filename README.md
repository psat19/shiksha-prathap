# How to run
## Prerequisites
- Docker
- Docker-compose plugin
- Git

## Steps
- Clone this repo: Change directory to a suitable folder and run the command - <code>git clone https://github.com/psat19/shiksha-prathap.git .</code>
- Start the containers: Make sure Docker desktop is running (windows); Make sure Docker Engine is running and docker-compose plugin is installed (Linux); Run the command <code>docker-compose up --build</code> from the root folder of the project
- This will spin up a server at <code>localhost:8080</code>. You can change the port by changing the variable <code>SRV_PORT</code>.
- Please be patient and wait for 3 to 4 minutes as the database spins up. The Go server will be on hold till the DB connection is made.

## Tests
- There is only one basic namesake test case. You can run it by using the command <code>go test ./cmd/</code> from the <code>server</code> folder
