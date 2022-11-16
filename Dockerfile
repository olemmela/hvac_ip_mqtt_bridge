# docker build --pull --label gsasha/hvac_ip_mqtt_bridge:latest -t gsasha/hvac_ip_mqtt_bridge:latest .
# docker build --pull --label gsasha/hvac_ip_mqtt_bridge:latest_arm -t gsasha/hvac_ip_mqtt_bridge:latest_arm .
# docker build --label gsasha/hvac_ip_mqtt_bridge:latest .
# docker push gsasha/hvac_ip_mqtt_bridge:latest
# docker push gsasha/hvac_ip_mqtt_bridge:latest_arm
FROM docker.io/library/golang:alpine AS build

LABEL maintainer="Sasha Gontmakher <gsasha@gmail.com>"

WORKDIR /data

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o bridge

FROM gcr.io/distroless/static
WORKDIR /app
USER nonroot:nonroot

COPY --from=build --chown=nonroot:nonroot /data/bridge .
COPY ac14k_m.pem .

EXPOSE 8080

CMD ["./bridge", "--config_file=/config/config.yaml"]

