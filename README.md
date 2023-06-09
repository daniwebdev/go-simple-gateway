# GO Simple Gateway

GO Simple Gateway is a lightweight gateway application written in Go that acts as an intermediary between multiple APIs with different domains. It allows you to configure the endpoints and namespaces for each API through a YAML configuration file. This project utilizes the Go Fiber library for building fast and efficient HTTP APIs.

## Features

- Acts as a gateway for multiple APIs with different domains.
- Configuration is stored in a YAML file.
- Implements middleware for handling requests.
- Supports various HTTP methods (GET, POST, PUT, PATCH, DELETE).
- Easy to configure and extend.

## Prerequisites

- Go 1.16 or above
- Go Fiber library (https://github.com/gofiber/fiber)

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/your-username/go-simple-gateway.git

2. Change to the project directory:

    ```bash
    cd go-simple-gateway

3. Build the project:

    ```bash
    go build

4. Run the executable:

    ```bash
    ./go-simple-gateway

    The server will start running on http://localhost:8080.

## Configuration

The configuration for the gateway is stored in a YAML file named config.yml. Update this file according to your API endpoints and namespaces. Here's an example of the configuration file:

```yaml
- endpoint: http://httpbin.org/
  namespace: httpbin
```

In the above example, the gateway will forward requests with the /httpbin namespace to the http://httpbin.org/ endpoint.

## Usage

To access the APIs through the gateway, use the following URL format:

```bash
http://localhost:8080/v1/{namespace}/{path}
```

Replace `{namespace}` with the desired namespace from your configuration file and `{path}` with the specific path for the API endpoint.

## Middleware
The gateway application includes a middleware that can be used for various purposes such as request logging, authentication, rate limiting, etc. Currently, the middleware is empty and can be customized according to your needs. Feel free to add your own logic and functionality to the middleware.

## Contributing
Contributions are welcome! If you find any issues or want to add new features to the project, please open an issue or submit a pull request.

## License
This project is licensed under the MIT License.

## Acknowledgments
- [Go Fiber](https://gofiber.io/) A fast and efficient web framework for Go.

## Contact
If you have any questions or suggestions, feel free to reach out to me at [me@dani.work](mailto://me@dani.work)

---
Feel free to modify the README.md file according to your project's specific details and requirements. Provide instructions for setting up and running the application, describe the available features, and include any necessary guidelines for contributing to the project.