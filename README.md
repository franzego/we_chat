# 🗨️ Go WebSocket Chat Server

A simple **real-time chat application** written in Go, powered by **WebSockets** and **Redis** for message broadcasting.  
Includes a lightweight **vanilla JavaScript frontend** for testing(couldn't make it work..bad at frontend :(    ).

---

## 🚀 Features

- Real-time chat via WebSockets
- Multi-user support (messages broadcast to all clients)
- Graceful shutdown handling (Ctrl+C, SIGTERM)
- Configurable via **environment variables** (`.env` for local, Vercel/Render for prod)
- Minimal frontend (`index.html`) for quick testing

---

## 📂 Project Structure

```
.
├── main.go              # Entry point, starts HTTP/WebSocket server
├── manager.go           # Connection manager (handles clients & broadcasts)
├── index.html           # Simple frontend (vanilla JS chat UI)
├── go.mod / go.sum      # Go module files
├── .env   
├── client.go            # Manages the various clients, the read and write messages
├── event.go             # Manages the events
├── redisclient.go       # for managing the redis connection
├── DockerFile
├── docker-compose.yaml
└── README.md            # Documentation
```

---

## ⚙️ Setup

### 1. Clone repo

```bash
git clone https://github.com/your-username/go-websocket-chat.git
cd go-websocket-chat
```

### 2. Install dependencies

```bash
go mod tidy
```

### 3. Create `.env`

For **local development**, create a `.env` file:

```
PORT=8080
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
```

> ⚠️ Never commit `.env` to GitHub — it’s for local use only.

---

## ▶️ Running Locally

1. Start **Redis** (Docker example):

   ```bash
   docker run -d -p 6379:6379 redis
   ```

2. Start the Go server:

   ```bash
   go run main.go
   ```

3. Open the frontend:

   - Open `index.html` directly in a browser, or  
   - Serve it with a simple HTTP server:

     ```bash
     python3 -m http.server 5500
     ```

   Then visit [http://localhost:5500](http://localhost:5500).

---

## 🌐 WebSocket Endpoint

- Local: `ws://localhost:8080/ws`  
- Production: `wss://<your-deployment-url>/ws`

---

## 🖥️ Frontend Preview
Use wscat in various terminals to simulate different users
- wscat -c ws://localhost:8080/ws    to connect to the server

In Alice terminal 
- {"type":"join-message","payload":{"username":"alice"}}
In Bob terminal:
- {"type":"join-message","payload":{"username":"bob"}}
Test Group Chat
- {"type":"send-message","payload":{"sender":"alice","message":"hello everyone"}}
Test DM
- {"type":"send-message","payload":{"sender":"alice","recipient":"bob","message":"hi bob"}}
Test Offine DM
- Disconnect Bob (Ctrl+C his terminal).
Alice sends:
{"type":"send-message","payload":{"sender":"alice","recipient":"bob","message":"are you there?"}}
Reconnect Bob:
wscat -c ws://localhost:8080/ws
{"type":"join-message","payload":{"username":"bob"}}


---

## ☁️ Deployment

### On Vercel / Render

1. Push code to GitHub.
2. Import the repo into Vercel or Render.
3. Set **Environment Variables** in dashboard:
   - `PORT` → `8080`
   - `REDIS_ADDR` → your Redis host (e.g., `redis-1234.c10.us-east-1-3.ec2.cloud.redislabs.com:12345`)
   - `REDIS_PASSWORD` → Redis password if set
4. Deploy 🚀

The server will automatically read environment variables using `os.Getenv`.

---

## 🛠️ Development Notes

- Local `.env` is optional — system env vars take precedence.
- `godotenv` is used for local dev convenience.
- For production, environment variables must be configured in the hosting platform.

---

## 📜 License

MIT License © 2025 Your Name
