# Latest golang image on apline linux
FROM golang:1.22

# Work directory
WORKDIR /user-authentication

# Installing dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copying all the files
COPY . .

# Starting our application
CMD ["go", "run", "main.go"]

# Exposing server port
