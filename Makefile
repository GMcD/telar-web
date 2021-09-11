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

stack:		## Update ECR tags in stack.yml
	awk -F "." '/354455067292/ { printf $$1; for(i=2;i<NF;i++) printf FS$$i; print FS$$NF+1 } !/354455067292/ { print }' stack.yml > .stack.yml && mv .stack.yml stack.yml

commit:		## Short hand for Commit to Prod Remote
commit:		##   * All references to version should be equal, and one more than published
commit:		##   *  -> go.mod == v0.1.101
commit:		##   *  -> package.json == v0.1.100
commit: 
	for i in actions admin auth notifications profile setting storage; do pushd micros/$i > /dev/null; go mod tidy; popd > /dev/null; done
	npm --no-git-tag-version version patch
	git add .; git commit -m ${ARGUMENT}; git push

fork:		## Short hand for Commit to Fork Remote
fork: stack
	git add . ; git commit -m ${ARGUMENT}; git push fork HEAD:master 

tag:		## Tag a Release
tag: fork
	git merge main && \
	git tag v$$(cat package.json | jq -j '.version') -am ${ARGUMENT} && \
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
	faas up --build-arg GO111MODULE=on --filter ${ARGUMENT}

telar-web:	## Commit, push and tag new version of telar-web
telar-web:
	echo "Edit code on prod/main" && \
	echo "Bump all go.mods to package.json + 1 " && \ 
	make commit ${ARGUMENT} && \
	echo "Move to fork/gmcd" && \
	git checkout gmcd && \
	git merge main && \
	make tag ${ARGUMENT} && \
	git checkout main && \
	git merge gmcd && \
	faas up --build-arg GO111MODULE=on

telar-core:	## Commit, push and tag new version of telar-core
telar-core:
	echo "Edit code on prod/main" && \
	make commit ${ARGUMENT} && \
	echo "Move to fork/gmcd" && \
	git checkout gmcd && \
	make tag ${ARGUMENT} && \
	git checkout main && \
	git merge gmcd && \
	faas up --build-arg GO111MODULE=on
