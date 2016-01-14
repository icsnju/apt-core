## APSARAS
-------

### Introduction  
APSARAS (Allocation of PhySicAl devices foR Android teSting) is a distributed testing platform for Android apps.  

### Dependency  
`github.com/bitly/go-simplejson`.  

[MooseFS](http://www.moosefs.org/) [optional] in each node.  

`JAVA Runtime Environment` in each node.  

`Android SDK` in each node.

### Deployment  

#### Build executable files
   
```
$ ./make.sh
```  
Executable files (master, slave and client) will be generated in relevant files. You can edit `make.sh` to output executable files in anywhere you want.  

#### Configuration 

In `master/`, `slave/` and `client/` diretories, some configuration files should be configured correctly.  

**1. master.conf**  

```  
share=/path/to/share
port=6666
```  
**share**: the path of shared file by master and slave nodes. Some information of tests and testing results will be stored in this file. We use *MooseFS* here.  

**port**: the port of the master node.  

**2. slave.conf**  

```
share=/path/to/share
master=ip:6666
adb=/opt/android-sdk/platform-tools/adb
```  
**share**: same with master.conf.  
**master**: the IP address and port of the master node.  
**adb**: the path of adb.  

**3. client.conf**

```
master=ip:6666
share=/path/to/share
```
**master** and **share** are same with slave.conf.  

#### Deploy  
Copy the `master/`, `slave/` and `client/` directories to the appropriate nodes.   

#### Run it  

Start master
  
```
$ ./run.sh
```  

Start slave  

```
$ ./run.sh
```  

### Usage  

**Check the state of slave nodes:**     

```
$ ./client slaves
```  

**Check the state of all of the testing jobs:**     

```
$ ./client jobs
```

**Submit testing job:**     

```
$ ./client subjobs [requirements]
```  
*Requirements* is the json file of testing requirements. For example:  

```
{
	"FrameKind":"monkey",
	"Frame":{
		"AppPath":"TouchMe.apk",
		"PkgName":"com.tc.touchme",
		"Argu":"-v 1"
	},
	"FilterKind":"specify_devices",
	"Filter":{
		"IdList":["*"],
		"Replaceable":true
	}
}
```  

**Check the state of a job**  

```
$ ./client job [ID]
```   
*ID* is the id of the job.   


### Result  
The testing results will be stored in the shared file. For a appointed job, the reuslts of this job will be stored in the directory with its ID.  







  



    


