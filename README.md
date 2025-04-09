# MCP Link - Convert Any OpenAPI V3 API to MCP Server

[![Join our Discord](https://img.shields.io/discord/1234567890?color=7289da&label=Discord&logo=discord&logoColor=white)](https://discord.gg/qkzfbqdSa9)

## 🧩 Architecture

![MCP Link](assets/diagrams.png)

## 🤔 Why MCP Link?

There is a notable gap in the current AI Agent ecosystem:

- Most MCP Servers are simple wrappers around Web APIs
- Functionality interfaces may not be complete, depending on developer implementation
- Manual creation of MCP interfaces is time-consuming and error-prone
- Lack of standardized conversion processes

MCP Link solves these issues through automation and standardization, allowing any API to easily join the AI-driven application ecosystem.

## 🌟 Key Features

- **Automatic Conversion**: Generate complete MCP Servers based on OpenAPI Schema
- **Seamless Integration**: Make existing RESTful APIs immediately compatible with AI Agent calling standards
- **Complete Functionality**: Ensure all API endpoints and features are correctly mapped
- **Zero Code Modification**: Obtain MCP compatibility without modifying the original API implementation
- **Open Standard**: Follow the MCP specification to ensure compatibility with various AI Agent frameworks

## 🌐 Online Version

Try our hosted version at [mcp-link.vercel.app](https://mcp-link.vercel.app) to quickly convert and test your APIs without installation.

## 🔧 Installation and Setup

### Prerequisites

Before installing MCP Link, make sure you have the following prerequisites:

1. **Go Installation**: MCP Link requires Go 1.18 or later.
   ```bash
   # Check if Go is installed and which version
   go version
   
   # If not installed, download and install from https://golang.org/dl/
   # Then verify installation
   go version
   ```

2. **MongoDB**: MCP Link uses MongoDB for persistent storage.
   ```bash
   # For local development, you can use Docker
   docker run --name mongodb -d -p 47017:27017 mongo:latest
   
   # Or install MongoDB directly on your system:
   # https://docs.mongodb.com/manual/installation/
   ```

### Clone Repository

```bash
# Clone repository
git clone https://github.com/automation-ai-labs/mcp-link.git
cd mcp-link
```

### Install Dependencies

```bash
# Download all required Go modules
go mod download

# Verify dependencies
go mod verify
```

## ⚙️ Configuration

MCP Link can be configured using environment variables or command-line flags.

### Environment Variables

You can create a `.env` file in the project root directory with the following variables:

```
MONGODB_URI=mongodb://localhost:47017
MONGODB_DATABASE=ominmcp
BASE_URL=http://localhost:8080
API_SERVER_HOST=localhost
API_SERVER_PORT=8080
API_SERVER_ENABLE_CORS=true
```

### Command-line Flags

When running the server, you can provide configuration using command-line flags:

```bash
go run main.go serve --port 8080 --host localhost --mongodb-uri mongodb://localhost:47017 --mongodb-database ominmcp --base-url http://localhost:8080
```

## 🚀 Running the Application

### Development Mode

For development purposes, you can run the application directly using Go:

```bash
# Run with default settings
go run main.go serve

# Or specify custom port and host
go run main.go serve --port 8080 --host 0.0.0.0
```

### Building for Production

To build the application for production deployment:

```bash
# Build the binary
go build -o mcp-link

# Make it executable (Linux/macOS)
chmod +x mcp-link

# Run the compiled binary
./mcp-link serve
```

## 📦 Docker Deployment

MCP Link can be easily deployed using Docker:

### Using Docker Compose

1. Edit the `docker-compose.yml` file to configure your environment.
2. Run Docker Compose:

```bash
docker-compose up -d
```

### Building Docker Image Manually

```bash
# Build the Docker image
docker build -t mcp-link .

# Run the container
docker run -d -p 8080:8080 \
  -e MONGODB_URI=mongodb://mongo:27017 \
  -e MONGODB_DATABASE=ominmcp \
  -e BASE_URL=http://localhost:8080 \
  --name mcp-link mcp-link
```

## 🔄 Using MCP Link

### Parameter Description

When using the SSE endpoint, the following parameters are available:

- `s=` - URL of the OpenAPI specification file
- `u=` - Base URL of the target API
- `h=` - Authentication header format, in the format of `header-name:value-prefix`
- `f=` - Path filter expressions to include or exclude API endpoints. Syntax:
  - `+/path/**` - Include all endpoints under /path/
  - `-/path/**` - Exclude all endpoints under /path/
  - `+/users/*:GET` - Include only GET endpoints for /users/{id}
  - Multiple filters can be separated by semicolons: `+/**:GET;-/internal/**`
  - Wildcards: `*` matches any single path segment, `**` matches zero or more segments

### Examples

| API | MCP Link URL | Authentication Method |
|-----|-------------|---------|
| ![Brave](https://img.logo.dev/brave.com) | https://mcp-link.vercel.app/links/brave | API Key |
| ![DuckDuckGo](https://img.logo.dev/duckduckgo.com) | https://mcp-link.vercel.app/links/duckduckgo | None |
| ![Figma](https://img.logo.dev/figma.com) | https://mcp-link.vercel.app/links/figma | API Token |
| ![GitHub](https://img.logo.dev/github.com) | https://mcp-link.vercel.app/links/github | Bearer Token |
| ![Home Assistant](https://img.logo.dev/home-assistant.io) | https://mcp-link.vercel.app/links/homeassistant | Bearer Token |
| ![Notion](https://img.logo.dev/notion.so) | https://mcp-link.vercel.app/links/notion | Bearer Token |
| ![Slack](https://img.logo.dev/slack.com) | https://mcp-link.vercel.app/links/slack | Bearer Token |
| ![Stripe](https://img.logo.dev/stripe.com) | https://mcp-link.vercel.app/links/stripe | Bearer Token |
| ![TMDB](https://img.logo.dev/themoviedb.org) | https://mcp-link.vercel.app/links/tmdb | Bearer Token |
| ![YouTube](https://img.logo.dev/youtube.com) | https://mcp-link.vercel.app/links/youtube | Bearer Token |

### Usage in AI Agents

```json
{
  "mcpServers": {
    "@service-name": {
      "url": "http://localhost:8080/sse?s=[OpenAPI-Spec-URL]&u=[API-Base-URL]&h=[Auth-Header]:[Value-Prefix]"
    }
  }
}
```

These URLs allow any API with an OpenAPI specification to be immediately converted into an MCP-compatible interface accessible to AI Agents.

## 💾 Using Persistent Configuration

MCP Link supports persistent storage of SSE configurations through MongoDB, allowing you to create a configuration once and reference it by ID without passing complete configuration parameters each time.

### Creating a Configuration

First, create a configuration using the API:

```bash
curl -X POST "http://localhost:8080/api/v1/config" \
  -H "Content-Type: application/json" \
  -d '{
    "schemaURL": "https://petstore3.swagger.io/api/v3/openapi.json",
    "baseURL": "https://petstore3.swagger.io",
    "headers": {
      "Authorization": "Bearer your-api-key"
    },
    "filters": [
      "+/pet/**:GET POST PUT",
      "+/store/**:GET",
      "-/user/**"
    ]
  }'
```

Upon successful creation, you'll receive a configuration ID and SSE URL:

```json
{
  "id": "645f8a1b2c3d4e5f6a7b8c9d",
  "sseUrl": "http://localhost:8080/sse/config?configId=645f8a1b2c3d4e5f6a7b8c9d",
  "message": "Configuration created successfully",
  "status": true
}
```

### Using Configuration by ID

You can access the SSE service using the configuration ID in either of two ways:

1. Using the dedicated configuration endpoint:
   ```
   http://localhost:8080/sse/config?configId=645f8a1b2c3d4e5f6a7b8c9d
   ```

2. Using the compatible standard SSE endpoint:
   ```
   http://localhost:8080/sse?configId=645f8a1b2c3d4e5f6a7b8c9d
   ```

### Usage in AI Agents

Use the SSE URL in your AI agent configuration:

```json
{
  "mcpServers": {
    "@petstore": {
      "url": "http://localhost:8080/sse?configId=645f8a1b2c3d4e5f6a7b8c9d"
    }
  }
}
```

### Managing Configurations

- Get configuration details: `GET /api/v1/config/{id}`
- Update configuration: `PUT /api/v1/config/{id}`
- Delete configuration: `DELETE /api/v1/config/{id}`

## 📋 Future Development

- **MCP Protocol OAuthflow**: Implement OAuth authentication flow support for MCP Protocol
- **Resources Support**: Add capability to handle resource-based API interactions
- **MIME Types**: Enhance support for various MIME types in API requests and responses

## 🔍 API Server Configuration

MCP Link now includes API server configuration functionality, allowing you to manage multiple API servers and their configurations.

### Creating an API Server Configuration

```bash
curl -X POST "http://localhost:8080/api/v1/api-server/config" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Pet Store API",
    "description": "Swagger Petstore API",
    "schemaUrl": "https://petstore3.swagger.io/api/v3/openapi.json",
    "baseUrl": "https://petstore3.swagger.io/"
  }'
```

### Managing API Server Configurations

- List all configurations: `GET /api/v1/api-server/config?all=true`
- Get configuration by ID: `GET /api/v1/api-server/config/{id}`
- Update configuration: `PUT /api/v1/api-server/config/{id}`
- Delete configuration: `DELETE /api/v1/api-server/config/{id}`

## 🛠️ Health Check Endpoint

MCP Link provides a health check endpoint to verify the application status:

```
GET /health
```

A successful response indicates that the application is running and connected to MongoDB:

```json
{
  "status": "ok",
  "message": "Service is healthy"
}
```

## 📄 License

[Apache 2.0](LICENSE)
