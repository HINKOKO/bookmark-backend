# **Bookmarkers - backend repository**

## **Overview**

Hello young gophers ! This repository contains the backend code for the **Bookmarkers** project. <br>
This project purpose is to build a central place for the students of **Systems, Low Level && Algorithms** specialization at Holberton School.

## **Technical Requirements**

At least, the following must be installed on your machine:

- **Golang** programming language latest version - follow installation instructions -> [over here](https://go.dev/doc/install)
  Confirm that installation process is successful by typing in a terminal:
  ```
  go version
  ```
  You should be prompted with something like (os specific according to your machine)
  ```
  go version go1.22.2 linux/amd64
  ```
- **Docker Engine** to be able to setup the **PostgreSQL** Database connected to this project. -> Docker [installation instructions](https://docs.docker.com/engine/install/)
  Similarly, in a terminal, run the command:

```
docker --version
```

You should be prompted with something like

```
Docker version 26.0.0, build 2ae903e
```

With those tools installed and ready to go, you can git clone this repository on your machine !

### **Starting backend**

```
git clone https://github.com/HINKOKO/bookmark-backend.git your-repo
cd your-repo
```

Navigate to the root of the folder you cloned and type

```
docker compose up -d
```

Thanks to the **docker-compose.yaml** file in this repository, this command will download a **postgres image**, build a container
based on this image, start this container on port 5432, with a PostgreSQL database , seeded with the latest migrated version in the **migrations/schema.sql** file

<quote>The **-d** flag starts the container in **detached** mode, avoid polluting your terminal with the container log. <br>
But at any time, if you want some info on what's going on inside this container, you can type

```
docker ps
```

This will list all container up and running, just spot the one named **bookmark-backend-postgres-1** and its **Container ID** (first col of output) , then type:

```
docker log <container_id>
```

This command is very useful, and gives valuable infos when database operations fails :wink:

</quote>

Once this is done, in the root level of the repository, type this command:

```
go run ./cmd/api
```

This will start the application on port 8080 and you should be prompted with these two lines:

```
2024/06/22 20:46:11 Connected to Postgres !
2024/06/22 20:46:11 starting application on port 8080
```

## **Contribute**

Contributions are welcome ! this project, while providing a functionnal MVP in conjunction with [this repository](https://github.com/HINKOKO/bookmarkers-client) is open to improvements and any suggestions ! <br>

### **How to contribute**

**1. Fork the repository**
Click on the "Fork" button at the top right cotner to create a copy of this repository under your Github account

**2. Clone your fork**

- Clone your forked repository to your local machine

```
git clone https://github.com/HINKOKO/bookmark-backend.git your-repo
cd your-repo
```

**3. Create a branch**

- Create a new branch for your changes

```
git checkout -b feature/your-feature-your-name
```

**4. Make changes**
Make your changes to the codebase

**5. Commit your changes**

- Please commit your changes with a meaningful message

```
git add .
git commit -m "Added feature X, improving feature Y..."
```

**6. Push to Github**
Push your changes to your forked repository on Github

```
git push origin feature/your-feature-your-name
```

**7. Create a pull request**
Go to the original repository on Github and create a pull request from your forked repository <br>
Provide a clear description please of your changes, and why they are beneficial

## **Help section**

If you have any questions or need help, feel free to reach me out at pizzoni.anthony@gmail.com

## **License**

Totally Open-source
