# 🪶 dummy-endpoint

A minimal Go web service that serves **HTTP** or **HTTPS** depending on environment variables.
Runs as a **non-root** user inside the container for security.

---

## 🧩 Features
- ✅ Optional HTTPS — automatically enabled when both `SSL_CERT_FILE` and `SSL_KEY_FILE` are provided
- ✅ HTTP fallback when TLS is not configured
- ✅ Minimal, dependency-free Go code
- ✅ Non-root container (user `verumex`)
- ✅ Multi-stage build for small, secure images

---

## 📁 Project Structure
```
.
├── Dockerfile
├── go.mod
└── main.go
```

---

## ⚙️ Environment Variables

| Name | Description | Default |
|------|-------------|---------|
| `SSL_CERT_FILE` | Path to TLS certificate (inside container) | *(empty)* |
| `SSL_KEY_FILE`  | Path to TLS private key (inside container) | *(empty)* |

> If both `SSL_CERT_FILE` and `SSL_KEY_FILE` are set, the app will **serve HTTPS only** on port **8443**.
> If not set, it will **serve HTTP only** on port **8080**.

---

## 🏗️ Build the image

```bash
docker build -t dummy-endpoint .
```

---

## 🚀 Run (HTTP only)

```bash
docker run --rm -p 80:8080 dummy-endpoint
```

You should see logs similar to:

```text
TLS envs not set -> starting HTTP ONLY on :8080
```

Then open http://localhost.

---

## 🔐 Run (HTTPS only)

```bash
docker run --rm \
  -p 443:8443 \
  -e SSL_CERT_FILE=/certs/fullchain.pem \
  -e SSL_KEY_FILE=/certs/privkey.pem \
  -v $(pwd)/certs:/certs:ro \
  dummy-endpoint
```

Logs:

```text
Starting HTTPS ONLY on :8443 (cert=/certs/fullchain.pem key=/certs/privkey.pem)
```

Then open https://localhost.

---

## 🧠 Notes

- The container uses non-privileged ports (8080/8443). Map them to 80/443 on the host as shown above.
- SSL certificate files must be **mounted** or **injected** via Kubernetes secrets (do **not** bake them into the image).

---

## 🧪 Test locally (without Docker)

```bash
go run main.go
```

or with HTTPS:

```bash
SSL_CERT_FILE=certs/fullchain.pem SSL_KEY_FILE=certs/privkey.pem go run main.go
```

---

## 🧰 Example Kubernetes Deployment (optional)

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: dummy-endpoint
spec:
  replicas: 1
  selector:
    matchLabels:
      app: dummy-endpoint
  template:
    metadata:
      labels:
        app: dummy-endpoint
    spec:
      containers:
        - name: dummy-endpoint
          image: ghcr.io/karunsiri/dummy-endpoint:latest
          ports:
            - containerPort: 8080
            - containerPort: 8443
          env:
            - name: SSL_CERT_FILE
              value: /certs/fullchain.pem
            - name: SSL_KEY_FILE
              value: /certs/privkey.pem
          volumeMounts:
            - name: certs
              mountPath: /certs
              readOnly: true
      volumes:
        - name: certs
          secret:
            secretName: tls-secret
```
