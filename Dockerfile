#
# Cleanup Dockerfile (Light)
#
FROM alpine:3.11.5 AS download

# The VERSION build argument specifies the Cleanup release
# version to be downloaded from GitHub.
ARG VERSION

RUN apk add --no-cache \
    curl \
    tar

RUN curl -LO https://github.com/dominikbraun/cleanup/releases/download/${VERSION}/cleanup-linux-amd64.tar.gz && \
    tar -xzvf cleanup-linux-amd64.tar.gz -C /bin && \
    rm -f cleanup-linux-amd64.tar.gz

FROM alpine:3.11.5 AS final

LABEL org.label-schema.schema-version="1.0"
LABEL org.label-schema.name="Cleanup"
LABEL org.label-schema.description="Remove gone Git branches with ease."
LABEL org.label-schema.url="https://github.com/dominikbraun/cleanup"
LABEL org.label-schema.vcs-url="https://github.com/dominikbraun/cleanup"
LABEL org.label-schema.version=${VERSION}
# ToDo: Add default command
# LABEL org.label-schema.docker.cmd="docker container run -v $(pwd)/my-blog:/project dominikbraun/espresso"

COPY --from=download ["/bin/cleanup", "/bin/cleanup"]

RUN mkdir /project

ENTRYPOINT ["/bin/cleanup"]
CMD ["branches", "/project"]