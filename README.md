# Gossip DB

From smartphones, IoTs for homes and smart cities, we generate enormous amounts of data while constantly accessing services on the internet. These cloud-enabled services often consider connected devices as edges and need them to be in sync. Therefore, we need ways to support distributed systems that are highly fault-tolerant and available while abiding limited connectivity and resources.

We introduce a distributed database that is adaptable to systems with low network bandwidth and node stability. By relaxing the consistency requirement of a database we support systems that look for eventual consistency. In this session, you will learn the basics of distributed systems, its design concerning the CAP theory, and a brave usage of recent reconciliation protocols that revolutionizes the fundamentals of traditional databases. With this work, we hope to inspire and enable similar approaches that effectively decentralize workloads and aggregate the results.



## Authors

@[Bowen Song](https://github.com/Bowenislandsong) *bowenson@usc.edu*

@[Yichen Ma](https://github.com/Never-M) *yichenm2@uci.edu*

@[Fuyao Wang](https://github.com/wfystx) *fuyao@bu.edu*



## Get Started

- Firstly, clone the project

  ```bash
  $ git clone https://github.com/Never-M/MyGossip.git
  $ cd Networking-inside-of-Kubernetes
  ```

- Make sure you have docker installed. Then use the script to start a basic 3-node cluster.

  ```bash
  $ ./start.sh
  Creating node3 ... done
  Creating node2 ... done
  Creating node1 ... done
  ```

- Now you can see that the nodes are running:

  ```bash
  $ docker ps -a
  CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS              PORTS               NAMES
  2ef8ad29f27f        mygossip_node3      "sh"                35 seconds ago      Up 34 seconds       8001-8002/tcp       node3
  f62d1cacb4b8        mygossip_node2      "sh"                35 seconds ago      Up 34 seconds       8001-8002/tcp       node2
  1f56ee2657be        mygossip_node1      "sh"                35 seconds ago      Up 34 seconds       8001-8002/tcp       node1
  ```

- Set up three terminals to get into the nodes seperately:

  ```bash
  $ docker exec -it node1 bash
  ```

- Now you can have fun with **Gossiper**, enter the node's name and its IP to start, and type `help` for instructions:

  ```bash
  $ go run main.go
  >> Enter node name
  node1
  >> Enter node ip
  172.28.1.1
  
  [node1]: help
  >> Please use commands below:
  >> exit: Shut down and exit current node
  >> add: Add a peer to current node
  >> remove: remove a peer from current node
  >> show: Print out peers of current node
  >> put: Put a key & value pair to the database
  >> get: get the value of a specific key
  >> delete: delete a key from the database
  ```

![](https://ibb.co/ZMRVMyF)

