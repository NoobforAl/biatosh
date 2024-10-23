# biatosh

biatosh app

## How to Run

To run the `biatosh` application, you can use the provided `Makefile` with the following targets:

- **all**: Build and run the project.
- **build**: Build the Docker image.
- **run**: Run the code outside of Docker.
- **dev**: Run the code outside of Docker with air.
- **dev-add-create-user**: Run the code outside of Docker with air and add create-user.
- **run-docker**: Run the code inside Docker.
- **clean**: Clean up Docker containers and images.
- **help**: Display the help message.

### Usage

To use the `Makefile`, navigate to the project directory and run:

```sh
make [target]
```

Replace `[target]` with one of the targets listed above. For example, to build the Docker image, you would run:

```sh
make build
```

To run the application outside of Docker, you would run:

```sh
make run
```

For more information on each target, you can run:

```sh
make help
```

### Creating Users in Docker

To create new users in the Docker container, you can use the following commands:

```sh
./main create-user --name="ali" --email="ali" --password="ali" --username="ali" --phone="123456789"
./main create-user --name="amin" --email="amin@gmail.com" --password="amin" --username="amin" --phone="987654321"
./main create-user --name="reza" --email="reza@gmail.com" --password="reza" --username="reza" --phone="234567890"
./main create-user --name="behzad" --email="behzad@gmail.com" --password="behzad" --username="behzad" --phone="345678901"
./main create-user --name="omid" --email="omid@gmail.com" --password="omid" --username="omid" --phone="456789012"
./main create-user --name="sepi" --email="sepi@gmail.com" --password="sepi" --username="sepi" --phone="567890123"
./main create-user --name="yasin" --email="yasin@gmail.com" --password="yasin" --username="yasin" --phone="678901234"
./main create-user --name="masih" --email="masih@gmail.com" --password="masih" --username="masih" --phone="789012345"
```

| Username | Password |
|----------|----------|
| ali      | ali      |
| amin     | amin     |
| reza     | reza     |
| behzad   | behzad   |
| omid     | omid     |
| sepi     | sepi     |
| yasin    | yasin    |
| masih    | masih    |
