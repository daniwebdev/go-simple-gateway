package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/daniwebdev/go-simple-gateway/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gopkg.in/yaml.v2"
)



type Config struct {
	Endpoint  string            `yaml:"endpoint"`
	Namespace string            `yaml:"namespace"`
	Headers   map[string]string `yaml:"header_auth"`
}

type GatewayConfig struct {
	Configs []Config `yaml:"configs"`
}

func main() {
	configFile := "config.yml"
	
	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}

	// Load configuration from YAML file
	config, err := loadConfig(configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	app := fiber.New()

	// Middleware
	// app.Use(logger.New())
	app.Use(middleware.LoggerMiddleware())
	app.Use(recover.New())

	// Routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Ok")
	})
	app.All("/v1/*", func(c *fiber.Ctx) error {
		path := c.Params("*")

		config, ok := checkNamespace(path, config.Configs)

		if !ok {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Namespace not found",
			})
		}

		// Remove the namespace segment from the path
		segments := strings.Split(path, "/")[1:]
		path = "/" + strings.Join(segments, "/")

		// Create the target URL by appending the remaining path to the configured endpoint
		targetURL := strings.TrimSuffix(config.Endpoint, "/") + path

		log.Println(targetURL)

		// Create a new HTTP request to the target URL
		req, err := http.NewRequest(c.Method(), targetURL, nil)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create request",
			})
		}

		// Set query parameters from the request
		queryParams := c.Request().URI().QueryArgs()
		if queryParams.Len() > 0 {
			q := req.URL.Query()
			queryParams.VisitAll(func(key, value []byte) {
				q.Add(string(key), string(value))
			})
			req.URL.RawQuery = q.Encode()
		}

		// Set body parameter from the request (if applicable)
		if c.Method() == fiber.MethodPost ||
			c.Method() == fiber.MethodPut ||
			c.Method() == fiber.MethodPatch ||
			c.Method() == fiber.MethodDelete {
			bodyBytes := c.Body()
			req.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
			req.ContentLength = int64(len(bodyBytes))
		}

		// Pass all headers from the incoming request to the new request
		c.Request().Header.VisitAll(func(key, value []byte) {
			req.Header.Set(string(key), string(value))
		})

		// Set headers from the configuration, allowing client headers to overwrite
		for key, value := range config.Headers {
			if req.Header.Get(key) == "" {
				req.Header.Set(key, value)
			}
		}

		// Perform the HTTP request
		client := http.DefaultClient
		resp, err := client.Do(req)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to perform request",
			})
		}
		defer resp.Body.Close()

		// Read the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to read response body",
			})
		}

		// Set the response body as the Fiber response
		c.Status(resp.StatusCode).Send(body)

		return nil
	})

	// Start server
	log.Fatal(app.Listen(":8080"))
}


func loadConfig(filename string) (GatewayConfig, error) {
	var config GatewayConfig
	file, err := os.ReadFile(filename)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func checkNamespace(path string, configs []Config) (Config, bool) {
	segments := strings.Split(path, "/")
	if len(segments) > 0 {
		namespace := segments[0]
		for _, config := range configs {
			if config.Namespace == namespace {
				return config, true
			}
		}
	}
	return Config{}, false
}

