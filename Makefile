NAME=interact-server
IMAGE=zxcv9203/piscine-golang-interact
TAG=0.1.4

dpull		:
	docker pull $(IMAGE):$(TAG)

dcrmf		:
	docker container rm -f $(NAME)

dirmf		:
	docker image rm -f $(IMAGE):$(TAG)

dbt			:
	docker build -t $(IMAGE):$(TAG) .

dpush		:
	docker push $(IMAGE):$(TAG)
