FROM ubuntu

COPY containerized_server ./containerized_server

CMD ["./containerized_server"]
