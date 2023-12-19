FROM golang:1.21 AS build
WORKDIR /app
COPY . /app
RUN make NAME=mackerel-sql-metric-collector

FROM gcr.io/distroless/static-debian11:nonroot
COPY --from=build --chown=nonroot:nonroot /app/bin/mackerel-sql-metric-collector /
WORKDIR /
ENTRYPOINT ["/mackerel-sql-metric-collector"]
