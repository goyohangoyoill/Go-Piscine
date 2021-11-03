NAME=interact-server
IMAGE=bigpel66/piscine-golang-interact
TAG=0.0.2

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
