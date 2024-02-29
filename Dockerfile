FROM --platform=$BUILDPLATFORM golang:1.21 AS build
ARG GIT_REVISION
ARG TARGETOS
ARG TARGETARCH
WORKDIR /app
COPY . /app
# To avoid downloading Go modules in the build of each platform, we share the module cache and lock during builds.
RUN --mount=type=cache,sharing=locked,target=/go/pkg/mod/ \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} make GIT_REVISION=${GIT_REVISION} NAME=mackerel-sql-metric-collector

FROM gcr.io/distroless/static-debian11:nonroot
COPY --from=build --chown=nonroot:nonroot /app/bin/mackerel-sql-metric-collector /
WORKDIR /
ENTRYPOINT ["/mackerel-sql-metric-collector"]
