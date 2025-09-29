IMAGE=franzego/chat-server
TAG=latest

# Build local image
build:
	docker buildx build -t $(IMAGE):$(TAG) .

# Run locally (maps port 8080)
run:
	docker run --rm -p 8080:8080 $(IMAGE):$(TAG)

# Push to Docker Hub
push:
	docker push $(IMAGE):$(TAG)

# Build for multiple platforms & push directly
release:
	docker buildx build --platform linux/amd64,linux/arm64 -t $(IMAGE):$(TAG) --push .

# Clean up dangling images
clean:
	docker image prune -f
