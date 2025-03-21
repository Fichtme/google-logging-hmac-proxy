# Google logging HMAC Proxy

A Go-based HMAC proxy project, containerized using Docker.
This application provides a secured connection that verifies requests using HMAC before forwarding them to the Google
Logging API.

## Environment Variables

Two environment variables are necessary for running this application:
- `HMAC_SECRET`: The secret key used for HMAC verification. If not provided, a default secret "default_secret" will be
  used.
- `GOOGLE_ACCESS_TOKEN`: The access token for Google Logging API. Place your actual access token here.

## Docker Instructions

Build the Docker image: 

```shell
docker build -t hmac-proxy .
```

Run the Docker container:

```shell
docker run -p 8080:8080 -e HMAC_SECRET -e GOOGLE_ACCESS_TOKEN hmac-proxy
```

## License

This is a simple licensed project under the MIT License.

## Contributing

Your contributions are always welcome! Please have a look at the contribution guidelines first. 🎉