# Go Shell Executor and Scheduler
## Introduction
This is an application built for schedule and execute shell scripts on Linux envrionments. It was designed to address some limitations of the default cron daeomn, built in to linux OSs
### Features
- Logging capability
- Built in database with quick access to important execution metadata
- In built REST API for,
  - Retrieving execution metadata
  - Overide schedule and remotely execute
## Getting Started
### 1. Clone Repository
```
git clone git@github.com:Janith911/shell-executor.git
```
### 2. Build binary executable
This step builds the binary executable and stores it in ```bin``` directory
```
make build
```
### 3. Create directories
Create directories for Sqlite DB file, Logs and Scripts in a desired location. Please note that these directories are specified in the configuration files later <br />
<br />
```
mkdir db
```
```
mkdir logs
```
```
mkdir scripts
```
### 4. Add configurations in to configuration.json file
This is a .JSON file which includes necessary configuration directives. In this configuration file,<br />
- Specify Sqilte DB path<br />
```
"DbFilePath" : "[ABSOLUTE_FILE_PATH]/db/[DB_FILENAME].db"
```
- Specify log file path<br />
```
"LogFilePath" : "[ABSOLUTE_FILE_PATH]/logs/[LOG_FILE_NAME].log"
```
- Specify API Bind IP and Port<br />
```
"BindIP": "127.0.0.1"
"BindPort": "8080"
```
- Specify script detail including ```Name```, ```ScriptPath``` and ```CronExpression``` for each script<br />
```
"Scripts" : [
        {
            "Name" : "Test01",
            "ScriptPath" : "[ABSOLUTE_FILE_PATH]/scripts/script1.sh",
            "CronExpression" : "* 10 * * *"
        },
        {
            "Name" : "Test02",
            "ScriptPath" : "[ABSOLUTE_FILE_PATH]/scripts/script2.sh",
            "CronExpression" :  "* 11 * * *"
        }
    ],
```
NOTE : A Sample configuration file is included in the repository ```(conf.json)```
### 5. Create environment variable spcifying configuration file path
```
export CONFIG_FILE_PATH=[ABSOLUTE_FILE_PATH]/[CONFIG_FILE_NAME].json

Eg :
export CONFIG_FILE_PATH=/etc/shellexec/conf.json
```
### 6. Start application
After starting the application, Scripts will be executed according to the specified CRON Expression, while listening for API requests on specified Port and IP
```
./shellexec start
```
OUTPUT : 
```
     _          _  _
 ___| |_   ___ | || | ___ __ __ ___  __
(_-<| ' \ / -_)| || |/ -_)\ \ // -_)/ _|
/__/|_||_|\___||_||_|\___|/_\_\\___|\__|
Version :  1.0.0
Author  :  Janith Vinura Bandara
Email   :  janithvinu@gmail.com
2023/09/20 21:53:09 Starting HTTP Endpoint
2023/09/20 21:53:09 Started HTTP Endpoint successfully
2023/09/20 21:53:09 Listening on : http://127.0.0.1:8080
```
### 7. Get execution metadata
```
./shellexec executions
```
OUTPUT : 
```
ID   ScriptName  StartTime                  Status   ExecutionId                    
1    Test01      2023-09-20T14:17:07+05:30  SUCCESS  Manual_2023_09_20_14:17:07     
2    Test02      2023-09-20T14:17:30+05:30  SUCCESS  Manual_2023_09_20_14:17:30     
3    Test02      2023-09-20T14:23:48+05:30  SUCCESS  Manual_2023_09_20_14:23:48     
4    Test02      2023-09-20T14:23:54+05:30  SUCCESS  Manual_2023_09_20_14:23:54     
5    Test01      2023-09-20T15:10:00+05:30  SUCCESS  Scheduled_2023_09_20_15:10:00  
6    Test01      2023-09-20T15:10:01+05:30  SUCCESS  Scheduled_2023_09_20_15:10:01  
7    Test01      2023-09-20T15:10:02+05:30  SUCCESS  Scheduled_2023_09_20_15:10:02  
8    Test01      2023-09-20T15:10:03+05:30  SUCCESS  Scheduled_2023_09_20_15:10:03  
9    Test01      2023-09-20T15:10:04+05:30  SUCCESS  Scheduled_2023_09_20_15:10:04  
10   Test01      2023-09-20T15:10:05+05:30  SUCCESS  Scheduled_2023_09_20_15:10:05
```
### 8. Execute Script manually
```
./shellexec execute [SCRIPT_NAME] [SHELL]

Eg :
./shellexec execute Test01 /bin/bash
```
### 9. List Scripts
```
./shellexec list

Eg :
./shellexec list
```
### 10. Read Scripts
```
./shellexec read [SCRIPT_NAME]

Eg :
./shellexec read Test01
```
