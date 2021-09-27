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

stack:		## Update ECR tags in stack.yml on main
	git checkout main && \
	if (( $$(git status --porcelain | wc -l) > 0 )); then \
	    printf "$${GREEN}Module $${RED}ts-serverless$${GREEN} has changes, run $${CYAN}make commit <message>$${GREEN} first.$${NC}\n"; \
	    exit 1; \
	fi && \
	awk -F "." '/354455067292/ { printf $$1; for(i=2;i<NF;i++) printf FS$$i; print FS$$NF+1 } !/354455067292/ { print }' stack.yml > .stack.yml && mv .stack.yml stack.yml

bump:		## Update go mod version numbers on main
bump: stack
	npm --no-git-tag-version version patch && \
	for mod in $$(find ./micros -name \*.mod); do \
		awk -F "0.1." '/telar-web v/ { printf $$1; for(i=2;i<NF;i++) printf FS$$i; print FS$$NF+1 } !/telar-web v/ { print }' $$mod > $${mod}.tmp && mv $${mod}.tmp $$mod; \
	done
	GOPRIVATE=github.com/GMcD for micro in $$(ls -d micros/*/); do pushd ./$${micro}; go mod tidy; popd; done && \
	git add . ; git commit -m Version-$$(cat package.json | jq -j '.version'); git push

commit:		## Short hand for Commit to main
commit: 
	git add .; git commit -m ${ARGUMENT}; git push

fork:		## Short hand for Commit main to Fork Remote
fork: bump
	git checkout gmcd && \
	git merge main && \
	git add . ; git commit -m ${ARGUMENT}; git push fork HEAD:master 

tag:		## Tag a Release
tag:		##   * Checkout fork/gmcd, update, push, tag, and release. Checkout prod/main.
tag: fork 
	git tag v$$(cat package.json | jq -j '.version') -am ${ARGUMENT} && \
	git push fork HEAD:master --tags && \
	git checkout main

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
	GOPRIVATE=github.com/GMcD faas up --build-arg GO111MODULE=on # --filter ${ARGUMENT}

## Deprecated
##   Workflow is now simply `make commit ...`, `make tag ...` and `make up`

# telar-web:	# Commit, push and tag new version of telar-web
# telar-web:
# 	echo "Edit code on prod/main" && \
# 	echo "Bump all go.mods to package.json + 1 " && \ 
# 	make commit ${ARGUMENT} && \
# 	echo "Move to fork/gmcd" && \
# 	git checkout gmcd && \
# 	git merge main && \
# 	make tag ${ARGUMENT} && \
# 	git checkout main && \
# 	git merge gmcd && \
# 	faas up --no-cache --build-arg GO111MODULE=on

# telar-core:	# Commit, push and tag new version of telar-core
# telar-core:
# 	echo "Edit code on prod/main" && \
# 	make commit ${ARGUMENT} && \
# 	echo "Move to fork/gmcd" && \
# 	git checkout gmcd && \
# 	make tag ${ARGUMENT} && \
# 	git checkout main && \
# 	git merge gmcd && \
# 	faas up --build-arg GO111MODULE=on
