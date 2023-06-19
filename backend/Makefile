# Image version
V=1.0
# Name of the google cloud project
PROJECT=spotify-party-390221

build:
	@go build -o ./bin/spotify-party

run: build
	@./bin/spotify-party

image:
	@docker build . -t spotify-party:$(V)

prod: image
	@docker tag spotify-party:$(V) gcr.io/$(PROJECT)/spotify-party:$(V)
	@docker push gcr.io/$(PROJECT)/spotify-party:$(V)

clean:
	@rm -rf ./bin


