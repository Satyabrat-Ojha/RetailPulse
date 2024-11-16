<!-- cd deployment -->

docker build -t retail-pulse-app .
docker run --rm -p 8080:8080 retail-pulse-app
