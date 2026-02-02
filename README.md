# Cartero

✅ A fast and lightweight SMTP outgoing catch-all built with Golang, designed to intercept and display all outgoing emails in a user-friendly web interface.

## Features

- **Catches outgoing emails** via SMTP on port 1025.
- **Parses and Displays emails** in a list format on `localhost:10122`.
- **Click on an email** to view its parsed HTML and raw email data.
- **Extracts links** from emails and displays them in plain text.

## How It Works

Cartero listens for outgoing emails on SMTP port 1025, captures all email traffic, and presents a list of those emails through a local web interface. You can click on any email to see both the raw email data and its parsed HTML. Any links present in the email content are automatically extracted and displayed in plain text.

## Setup

1. Clone the repository:

    ```
    git clone https://github.com/tego101/cartero.git
    cd cartero
    ```

2. Install dependencies:

    ```
    go mod tidy
    ```

3. Run the application:

    ```
    go run main.go
    ```

4. Visit `http://localhost:10122` in your browser to view the captured emails.

## Docker Setup

### Using Docker Compose (Recommended)

1. Clone the repository:

    ```
    git clone https://github.com/tego101/cartero.git
    cd cartero
    ```

2. Run with Docker Compose:

    ```
    docker-compose up -d
    ```

3. Visit `http://localhost:10122` in your browser to view the captured emails.

### Using Docker directly

1. Build the Docker image:

    ```
    docker build -t cartero .
    ```

2. Run the container:

    ```
    docker run -d -p 1025:1025 -p 10122:10122 --name cartero cartero
    ```

3. Visit `http://localhost:10122` in your browser to view the captured emails.

#### Docker Volumes

To persist the email database, you can mount a volume:

```
docker run -d -p 1025:1025 -p 10122:10122 -v $(pwd)/data:/app/data --name cartero cartero
```

#### Environment Variables in Docker

You can override the default environment variables:

```
docker run -d -p 1025:1025 -p 10122:10122 \
  -e SMTP_HOST=0.0.0.0 \
  -e SMTP_PORT=1025 \
  -e WEB_HOST=0.0.0.0 \
  -e WEB_PORT=10122 \
  --name cartero cartero
```

## Configuration (env variables)

- `SMTP_HOST`: The host on which Cartero listens for outgoing emails. Default is `localhost`.
- `SMTP_PORT`: The port on which Cartero listens for outgoing emails. Default is `1025`.
- `WEB_HOST`: The host on which the web interface is served. Default is `localhost`.
- `WEB_PORT`: The port on which the web interface is served. Default is `10122`.

### Sample ENV file.

```
# .env
SMTP_PORT=1025
SMTP_HOST=0.0.0.0
WEB_PORT=10122
WEB_HOST=0.0.0.0
```

## Dependencies

- Go 1.18+
- Templ
- TailwindCSS
- Alpine.js

## Screen Shots

### Inbox

![Inbox View](/screen-shots/inbox-view.png)

### E-Mail View

![E-Mail View](/screen-shots/email-view.png)

## License

MIT License. See the [LICENSE](LICENSE) file for more details.

---
