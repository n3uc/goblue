#
# NOTE:  This uses git tags to create version labels.  Will not work if git describe does not find a tag
#
IMG_REPO := docker.io/sterlingdevil/goblue
CONENG := podman

# check for .git dir
ifeq ($(shell test -d .git; echo $$?), 0)
  #check for valid repo
  GIT_INFO := $(shell git status 1>/dev/null 2>/dev/null; echo $$?)
  ifneq ($(GIT_INFO),0)
    VER := "GIT-REPO-INVALID"
    CI_COMMIT_SHA := "GIT-REPO-INVALID"
  else	
    # git describe is ok, so lets get our data
    ifeq ($(origin CI_COMMIT_SHA), undefined)
      CI_COMMIT_SHA := $(shell git describe --all --dirty --long)
    endif

    LVER := $(shell git describe --tags)
    SVER := $(shell git describe --tags | rev | cut -c 10- | rev)

    ifeq ($(shell test -z $(LVER); echo $$?), 1)
      LEN := $(shell expr length $(LVER))
      ifeq ($(shell test $(LEN) -gt 10; echo $$?), 0)
   	    VER := $(SVER)
      else
	    VER := $(LVER)
      endif
	else  
      VER := "NO-GIT-TAG"
      CI_COMMIT_SHA := "NO-GIT-TAG"
    endif
  endif
# no .git directory
else
  VER := "NO-GIT-REPO"
  CI_COMMIT_SHA := "NO-GIT-REPO"
endif

build:
	@mkdir -p bin
	go build -ldflags="-X 'main.Version=$(VER)' -X 'main.GitCommit=$(CI_COMMIT_SHA)' -X 'main.BuildTime=$$(date)'" -buildvcs=false -tags netgo -o bin ./...

oci: build
	$(CONENG) build . --tag $(IMG_REPO):latest --tag ${IMG_REPO}:$(VER)

clean:
	@rm -rf bin

push: oci
	$(CONENG) push $(IMG_REPO) 
	$(CONENG) push $(IMG_REPO):$(VER)

helmchart: 
	helm package deploy/helmchart -d bin --app-version $(VER) --version $(VER)

all: clean push helmchart
