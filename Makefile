
.PHONY: all
all:
	cd rutte
	make build
	# cd ../arib
	# npm install
	# npm run tsc
	# cd ../opstelten
	# npm install
	# npm run tsc

## invoke this to create a hotfix release (e.g. 1.6 -> 1.7)
.PHONY: create-tag-hotfix
create-tag-hotfix: assert-no-pending-tags
	git fetch --tags
	@export VERSION=$$(git tag --sort=v:refname | tail -1) && \
	export MAJOR=$$(echo $$VERSION | awk '{ print substr($$0,3,1) }') && \
	export MINOR=$$(echo $$VERSION | awk '{ print substr($00,5,1) }') && \
	export NEW_VERSION=v0.$${MAJOR}.$$((MINOR+1)) && \
	echo "$$VERSION -> $$NEW_VERSION" && \
	git tag 0.$$NEW_VERSION && \
	git push --tags

## invoke this to create a major release (e.g. 1.6 -> 2.0)
.PHONY: create-tag-release
create-tag-release: assert-no-pending-tags
	git fetch --tags
	@export VERSION=$$(git tag --sort=v:refname | tail -1) && \
	export MAJOR=$$(echo $$VERSION | awk '{ print substr($$0,3,1) }') && \
	export MINOR=$$(echo $$VERSION | awk '{ print substr($00,5,1) }') && \
	export NEW_VERSION=v0.$$((MAJOR+1)).0 && \
	echo "$$VERSION -> $$NEW_VERSION" && \
	git tag $$NEW_VERSION && \
	git push --tags

assert-no-pending-tags:
	@if git push --tags --dry-run 2>&1 | grep -F 'new tag'; then echo "Some unpushed new tags already exist"; exit 1; fi
