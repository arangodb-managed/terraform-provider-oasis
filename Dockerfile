FROM alpine
ARG GOARCH=amd64
ARG ROLEID

RUN apk add -u ca-certificates
ADD ./bin/linux/${GOARCH}/datad /app/
RUN mkdir -p /etc/datad
RUN echo ${ROLEID} > /etc/datad/role-id

WORKDIR /app/
ENTRYPOINT [ "/app/datad" ]
