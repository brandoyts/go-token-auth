# Environment Configuration

This project uses environment variables to configure authentication, database, and server settings. Below is a breakdown of each `.env` variable and its purpose.

## 🔐 Authentication

| Variable            | Description                                                | Example              |
| ------------------- | ---------------------------------------------------------- | -------------------- |
| `JWT_SECRET`        | Secret key used to sign JWT tokens. Keep this secure.      | `dal4$h@049h&j93...` |
| `ACCESS_TOKEN_TTL`  | Lifetime of the access token. Accepts Go duration format.  | `"1m"` (1 minute)    |
| `REFRESH_TOKEN_TTL` | Lifetime of the refresh token. Accepts Go duration format. | `"2m"` (2 minutes)   |

> **Note**: Duration values must follow [Go’s duration syntax](https://pkg.go.dev/time#ParseDuration) (e.g., `"1h"`, `"30m"`, `"15s"`).

---

## 🗄️ MongoDB

| Variable              | Description                  | Example                     |
| --------------------- | ---------------------------- | --------------------------- |
| `MONGO_URI`           | Full MongoDB connection URI  | `mongodb://localhost:27017` |
| `MONGO_DATABASE_NAME` | Name of the MongoDB database | `default`                   |
| `MONGO_USERNAME`      | MongoDB username             | `root`                      |
| `MONGO_PASSWORD`      | MongoDB password             | `password`                  |

> If `MONGO_URI` is provided, it overrides username/password configuration.

---

## 🧠 Redis

| Variable        | Description                                | Example          |
| --------------- | ------------------------------------------ | ---------------- |
| `REDIS_ADDRESS` | Redis server address                       | `localhost:6379` |
| `REDIS_SSL`     | Whether to use TLS/SSL (`true` or `false`) | `false`          |

---

## 🚀 Server

| Variable | Description              | Example |
| -------- | ------------------------ | ------- |
| `PORT`   | Port for the HTTP server | `8090`  |

---

## 🧪 Example `.env` File

```env
JWT_SECRET=dal4$h@049h&j93oqju4o5xl-qrk4z%$jz^86_-^v4+1w6o7ko
ACCESS_TOKEN_TTL="1m"
REFRESH_TOKEN_TTL="2m"

MONGO_URI=
MONGO_DATABASE_NAME=default
MONGO_USERNAME=root
MONGO_PASSWORD=password

REDIS_ADDRESS=
REDIS_SSL=false

PORT=8090
```
