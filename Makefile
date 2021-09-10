# This makes the subsequent variables available to child shells
.EXPORT_ALL_VARIABLES:

include .env

# Collect Last Target, convert to variable, and consume the target.
# Allows passing arguments to the target recipes from the make command line.
CMD_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
# Consume them to prevent interpretation as targets
$(eval $(CMD_ARGS):;@:)
# Service for command args
ARGUMENT  := $(word 1,${CMD_ARGS})

##
## Usage:
##  make [target] [ARGUMENT]
##   operates in namespace ${ARGUMENT}
##

help:		## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

commit:		## Short hand for Commit to Prod Remote
	git add .; git commit -m ${ARGUMENT}; git push

fork:		## Short hand for Commit to Fork Remote
	git add . ; git commit -m ${ARGUMENT}; git push fork HEAD:master 

tag:		## Tag a Release
tag: fork
	git tag $${RELEASE_TAG} -am ${ARGUMENT}
	git push fork HEAD:master --tags 

logs:		## Log Pod ${ARGUMENT} by prefix
logs:
	kubectl logs --namespace openfaas-fn $(shell kubectl get pods --namespace openfaas-fn -o=jsonpath='{.items[*].metadata.name}' -l faas_function=${ARGUMENT})

login:  	## ECR Docker Login
	@ aws ecr get-login-password --region $${AWS_REGION} | docker login --username AWS --password-stdin $${AWS_ACCOUNT_ID}.dkr.ecr.$${AWS_REGION}.amazonaws.com
	@ AWS_ACCOUNT_ID=$$(aws sts get-caller-identity --output text --query 'Account'); \
	AWS_IAM_ARN=$$(aws sts get-caller-identity --output text --query 'Arn'); \
	echo "Running as $${AWS_IAM_ARN} in $${AWS_REGION} for $${AWS_ACCOUNT_ID}."

up:		## Run FaaS up
up: login
	# Update micros with new core && code bases
	# ./update-micros.sh telar-core
	# ./update-micros.sh telar-web
	echo "Running FaaS up..."
	faas up --build-arg GO111MODULE=on

version:	## Commit, push and tag new version of telar-web
version:
	echo "On prod/main" && \
	echo "Edit code...."
	echo "Update stack.yml to v2.2.95 ... " && \
	make commit ${ARGUMENT} && \
	echo "Move to fork/gmcd" && \
	git checkout gmcd && \
	git merge main && \
	echo "Update Release Number to v0.1.86 in micros/auth/go.mod and .env " && \ 
	make tag ${ARGUMENT} && \
	git checkout main && \
	git merge gmcd && \
	faas up --build-arg GO111MODULE=on

core:	## Commit, push and tag new version of telar-core
core:
	echo "On prod/main" && \
	echo "Edit code...."
	make commit ${ARGUMENT} && \
	echo "Move to fork/gmcd" && \
	git checkout gmcd && \
	git merge main && \
	echo "Update .env with new RELEASE_TAG"
	make tag ${ARGUMENT} && \
	git checkout main && \
	git merge gmcd && \
	faas up --build-arg GO111MODULE=on
