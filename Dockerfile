FROM scratch
ENTRYPOINT ["/usr/bin/inventory"]
COPY inventory /usr/bin/inventory