## APSARAS

![image](logo.png)

APSARAS (Allocation of PhySicAl devices foR Android teSting) is a distributed testing platform for Android apps.  

### Dependency  

1. Go1.6  

2. [MooseFS](http://www.moosefs.org/) [optional] in each node.  

3. `JAVA Runtime Environment` in each node.  

4. `Android SDK` in each node.  

5. Install Godep:  

	```
	$ go get github.com/tools/godep`  
	```   
	Put godep in envionment variables. In `Apsaras/`, get dependences from the Internet:   

	```
	$ godep restore
	$ godep save ./...
	```   
	You can save all dependnce in `./vendor`, and delete the package in `$GOPATH`.

6. Update dependences.   

	```
	$ get [package]
	$ godep save ./...
	```

### Deployment  

#### Build executable files

In `server/` or `slave/`:    

```
$ go build
```  
Executable files (server, slave) will be generated in relevant files.   

#### Configuration/run

In `server/`, `slave/`diretories, some configuration files should be configured correctly. You can refer **Apsaras Run** for more details in next section.

=========================
========================
## Apsaras docker image  

It is a docker image for [Apsaras](https://github.com/icsnju/apt-core). Get it [apsaras/aps](https://hub.docker.com/r/apsaras/aps/). Run it now:  

```
docker run -it --privileged --net host  --device /dev/fuse apsaras/aps
```

### Dependences
In this image, the following tools are installed:  
1. **Android SDK:** It is in `/opt/sdk`. If you get `command not found` when use `adb`, you should use `source /etc/profile` command firstly.   
2. **JRE:** Java is needed in this project.   
3. **MooseFS:** We use moosefs as our distributed file system. You may use the following command to start the moosefs.  

In file `/etc/mfs`, you will find the configuration file for moosefs, so configure them with proper arguments. Most of configurations are default values. While in `mfschunkserver.cfg`, we set `MASTER_HOST = mfsmaster`, `mfsmaster` defined in `/etc/hosts` ( we use `--net host` here). In `mfshdd.cfg`, we add line `/home/aps/mfs` here. Of course, you can change the configurations to whatever you want.

Start mfsmaster in one mfs master node:  

```
$ mfsmaster start 	//Start mfsmaster
$ mfsmaster stop   //Stop mfsmaster
```

Start mfscgiserv (optional), you can monitor the mfs status in http://127.0.0.1:9425:  

```
$ mfscgiserv start
$ mfscgiserv stop
```

Start mfschunkserver in every mfs slave node:

```
$ mfschunkserver start  
$ mfschunkserver stop
```

Mount mfs file to local file, you can use this local file to share files now. We use `/home/aps/share` as an example here:  

```
$ sudo mfsmount /home/aps/share -H mfsmaster
```   
You can run mfsmaster and mfschunkserver in standalone containers, but you have to mount the correct file, which can be accessed by Apsaras, to the mfsmaster.

4. **MongoDB** MongoDB is not in this image, so you should start a mongodb server in your host or in a container ([tutum/mongodb](https://hub.docker.com/r/tutum/mongodb/)).   

```
docker run -d -p 27017:27017 -p 28017:28017 -e AUTH=no tutum/mongodb   
```

### Apsaras Run
In file `/home/aps/apsaras/`, there are two files, `server` and `slave`. You should run server in one master node and slave in slave nodes that are conneted with many Android devices.  

#### Server   
Firstly, configure your server node.   

In `/home/aps/apsaras/server/config/app.conf`, set the http server configurations:  

```
appname = apsaras
httpport = 8023
runmode = dev
CopyRequestBody = true
AutoRender = false
TemplateLeft = {{<
TemplateRight = >}}
```

You can change the `httpport` to a better one.   

In `/home/aps/apsaras/server/config/master.json`, set the proper configurations for Apsaras master:  

```
{
        "SharePath":"/home/aps/share",
        "Port":"6666",
        "DBUrl":"localhost",
        "DBName":"aptweb-dev"
}
```  

**SharePath:** It is the file to save testing files, which are shared by `MooseFS`. We mount `/home/aps/share` to mfsmaster here as mentioned before.  

**Port:** The port is used by master node to communicate with slave node.  

**DBUrl:** It is the url of Mongodb server.  

**DBName:** Set the DB name.   

It is time to run server:   

```
./server
```     
Now you can get your web client in http://localhost:8023/jobs.

#### Slave  
Same with the server, you should config it.   
{
        "SharePath":"/home/aps/share",
        "ServerIP":"localhost:6666",
        "SDKPath":"/opt/sdk"
}
**SharePath:** It is the file to save testing files, which are shared by `MooseFS`. We mount `/home/aps/share` to mfsmaster here as mentioned before.  

**ServerIP:** The `ip:port` of the master node.   

**SDKPath:** The Android sdk path. We have installed adb in this container with path `/opt/sdk`.  

Connect some Android devices to this slave node, you can run this slave node now:

```
./slave
```
You will see all of your connected Android devices in the web client.   


#### Usage   
You can run server and slave in standalone containers and physical computers. Apsaras is easy to use. With a beautiful web client, you will enjoy the features of Apsaras :).
