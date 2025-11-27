# chat-application
websocket, chat service (grpc) , postgres etc

github.com/ak-repo/chat-application





docker postgres cmd = docker exec -it chat-postgres psql -U chat-user -d chat


docker compose -f docker/docker-compose.yml down

docker compose -f docker/docker-compose.yml up -d --build
