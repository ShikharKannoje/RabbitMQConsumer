
# Setup RabbitMQ using docker command below
docker run -it --rm --name rabbitmq -p 15672:15672 -p 5672:5672 rabbitmq:3-management

# setup the PostgresDB & create table, refer .env file
Tar db is in DatabaseTar folder.

# Run Producer
go run main.go


# Run Consumer
go run main.go
