FROM scratch
COPY hqcli /usr/bin/hqcli
ENTRYPOINT ["/usr/bin/hqcli"]