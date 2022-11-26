```
  ______  __    __ _______   ______  _______   ______  
 /      \|  \  |  \       \ /      \|       \ /      \ 
|  ▓▓▓▓▓▓\ ▓▓  | ▓▓ ▓▓▓▓▓▓▓\  ▓▓▓▓▓▓\ ▓▓▓▓▓▓▓\  ▓▓▓▓▓▓\
| ▓▓__| ▓▓ ▓▓  | ▓▓ ▓▓__| ▓▓ ▓▓  | ▓▓ ▓▓__| ▓▓ ▓▓__| ▓▓
| ▓▓    ▓▓ ▓▓  | ▓▓ ▓▓    ▓▓ ▓▓  | ▓▓ ▓▓    ▓▓ ▓▓    ▓▓
| ▓▓▓▓▓▓▓▓ ▓▓  | ▓▓ ▓▓▓▓▓▓▓\ ▓▓  | ▓▓ ▓▓▓▓▓▓▓\ ▓▓▓▓▓▓▓▓
| ▓▓  | ▓▓ ▓▓__/ ▓▓ ▓▓  | ▓▓ ▓▓__/ ▓▓ ▓▓  | ▓▓ ▓▓  | ▓▓
| ▓▓  | ▓▓\▓▓    ▓▓ ▓▓  | ▓▓\▓▓    ▓▓ ▓▓  | ▓▓ ▓▓  | ▓▓
 \▓▓   \▓▓ \▓▓▓▓▓▓ \▓▓   \▓▓ \▓▓▓▓▓▓ \▓▓   \▓▓\▓▓   \▓▓
```

# What is Aurora ? #

Aurora is a small PaaS `Platform as a Service` that can be run on a small single node with 1GB of RAM all the way to a multi-node setup.\
The goal is to make it as easy as possible to simply focus on writing your software no matter who you are and what your resources might be
without having to pay for a PaaS or have a difficult time dealing with steps outside your normal workflow.

## NB ##

 This is not really a production level PaaS and there are definitely a lot of wrong assumptions but, it could be at some point suitable in the future.

### Installation and Setup ###

#### Pre-requisites ####

- Linux system or WSL-2 with golang 1.19.X to compile the code
- Debian Linux or any linux environment with docker and docker swarm installed and systemctl
- Open 80 and 8080 open, the rest can stay closed and preferably internal ethernet as well for node to node communication

#### STEPS ####

On a computer with the desired golang compiler for the target architecture running linux , you can compile the golang source code by running `./full_build.sh`
This might be x86_64, RISC-V, or ARM-64.

This will result into a tar-ball with the target architecture for your servers that you can transfer using

```bash
    scp aurora-service.tar.gz hostname@IP-ADDRESS:~
    tar -xvf aurora-service.tar.gz 
    sudo chmod +x aurora-service  
    sudo mv aurora-service /usr/local/bin/
```

Create a new group and user

```bash
   sudo groupadd --systemctl aurora
   sudo useradd -s /sbin/nologin --system -g aurora aurora
```

Create a folder under `/etc~` to hold configurations and files for the `aurora-service daemon`

```bash
    sudo mkdir -p /etc/aurora/configurations/
```

Change the folder and file ownership

```bash
    sudo chown -R aurora:aurora /etc/aurora/configurations 
    sudo chown -R 755 /etc/aurora/configurations
```

Create a systemd file

```bash
    sudo nano /etc/systemd/system/aurora-service.service
```

Add the following lines

```bash
    [Unit]
    Description=Aurora Service [Daemon]
    After=network-online.target
    Wants=network-online.target
    [Service]
    Type=simple
    User=aurora
    Group=aurora
    ExecStart=/usr/local/bin/aurora-service --mode=run
    ExecReload=/bin/kill -HUP $MAINPID
    KillSignal=SIGINT
    TimeoutStopSec=5
    Restart=on-failure
    SyslogIdentifier=aurora-service
    [Install]
    WantedBy=multi-user.target 
```

Save and close the file when you are finished

Next, reload the system daemon with the following command

```bash
    sudo systemctl daemon-reload
    sudo systemctl start aurora-service
    sudo systemctl enable aurora-service
```

And check to see if the daemon is running

```bash
    sudo systemctl status aurora-service
```

Then the last thing is to add aurora user to the docker user group

```bash
    sudo usermod -aG docker aurora
```

Run the command to create a system administration use and re-own the created database

```bash
    sudo aurora-service --mode=setup --email=<email-address> --password=**********
    sudo chown aurora:aurora -R /etc/aurora/configurations/store.db
    sudo systemctl restart aurora-service
```

Everything should be running ok now

### What's Next ? ###

To use the `PaaS` you can check out [Aurora CLI](https://github.com/space-fold-technologies/aurora-client) which can be used on your developer work station. \

To add more nodes to your cluster, you can use the [Aurora Agent](https://github.com/space-fold-technologies/aurora-agent)

### Features ###
  
- [x] Docker Swarm Integration
- [x] Traefik Reverse Proxy Support
- [ ] Kubernetes Integration
- [ ] Podman Integration
- [x] Agent Node Integration
- [ ] Agent Node Stats
- [ ] Management Dashboard
- [ ] Scheduled container checks
- [ ] Integrated Services / Catridges[`PostgreSQL`,`MariaDB`,`Redis`,`RabbitMQ`]


