export DOCKER_BUILDKIT := 1
images: dns rest
all: images compose

# TODO: In final release it would be better to default to amd64
#PLATFORM=linux/amd64
PLATFORM=linux/arm
dns:
	@docker build -t dns:`echo ${PLATFORM} | sed 's@/@_@'` --target dns --platform ${PLATFORM} .
rest:
	@docker build -t rest:`echo ${PLATFORM} | sed 's@/@_@'` --target rest --platform ${PLATFORM} .
compose:
	@TAG=`echo ${PLATFORM} | sed 's@/@_@'` docker-compose up
