# containerizing_servers
This project served to create a server to test intercontainer communication using container names, with the aid of postman.

After downloading the repo, one can build a container image using the Dockerfile found within the project.

Create a user defined network:
$docker network create my-user-defined-network-name

Launch 2 instance of the server
1. An instance only accessible via the container network just created:
docker run --network my-user-defined-network-name image/name

2. An instance accessible from outside, so we can communicate with it, using postman, and then communicate with the hidden server:
docker run -p 127.0.0.1:9001:9001 --network my-user-defined-network-name image/name

Connecting to the default container network, bridge, wont allow inter contaier communication using names. It's one of the advantages of using a user-defined network.


