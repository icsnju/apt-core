go build -o ./master/master ./master/master.go ./master/handler.go ./master/scheduler.go ./master/mgohelper.go
go build -o ./slave/slave ./slave/slave.go 
go build -o ./client/client ./client/client.go
