# Small-scale RPM Build Automation

I run ~~CentOS~~ Rocky Linux on my home server ([PowerEdge R610](https://i.dell.com/sites/csdocuments/Shared-Content_data-Sheets_Documents/en/R610-SpecSheet.pdf))
and regularly found myself creating and building RPMs for software unavailable
in mainstream repositories or rebuilding kernel module RPMs (ZFS, mptsas) when
a kernel upgrade was imminent.

Instead of running a full-blown [Koji](https://docs.pagure.org/koji/) build server
(which wouldn't handle repository updates anyway as far as I know), I created this
simple system which consists of:

* `corebuilder` &ndash; An RPM build environment container image for use with CI/CD (GitLab in my case)
* `repo-manager` &ndash; A repository server written in Go to manage and serve my small RPM repository

## corebuilder

This image is packaged with two scripts, `build-rpms` and `upload-rpms`. `build-rpms`
handles installing build dependencies, moving RPM specs and sources to the appropriate
directories, and the building itself. `upload-rpms` POSTs the built RPMs to the
authenticated repo-manager endpoint, which adds the RPMs to my custom repo and refreshes
the repodata.

Example `.gitlab-ci.yml`

``` yaml
stages:
- build
- deploy

workflow:
  rules:
  - if: '$CI_COMMIT_BRANCH == "core-build"'

build:
  stage: build
  image: reg.laboutpost.net/core/rpms/corebuilder:fresh
  script: build-rpms "${CI_PROJECT_DIR}"
  artifacts:
    paths:
    - build/RPMS/
    - build/SRPMS/

deploy:
  stage: deploy
  image: reg.laboutpost.net/core/rpms/corebuilder:fresh
  script: upload-rpms "${CI_PROJECT_DIR}" "${REPO_MANAGER_SECRET}"
```

## repo-manager

Primarily a static HTTP server (`/repos`), with the exception of one endpoint
(`/addpkgs`). The `/addpkgs` endpoint accepts authenticated mulitpart form data
POST requests containing the built RPMs, and adds them to the repository.

After copying the RPM to disk, `repo-manager` calls [createrepo_c](https://github.com/rpm-software-management/createrepo_c)
to update the repository data.


```
curl -X POST http://repo-manager/addpkgs \
     -F pkg=@package0.rpm \
     -F pkg=@package1.rpm \
     -H "x-auth: $REPO_MANAGER_SECRET"
```