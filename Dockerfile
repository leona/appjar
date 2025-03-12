FROM golang:1.24-bookworm
RUN apt-get -y update && apt-get install -y nodejs npm libwebkit2gtk-4.0-dev
RUN go install github.com/wailsapp/wails/v2/cmd/wails@latest
WORKDIR /app
