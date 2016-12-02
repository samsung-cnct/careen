# Careen

Careen is a utility to __clone__ repositories and __apply__ patches to them. It does not aim to be a general artifact management system, just to create a consistent set of sources to build from.

Careen is also verb meaning to "turn (a ship) on its side for cleaning, caulking, or repair."

## Usage

```
git clone <careen-repo>
cd careen
make
./careen help
./careen clone configs/docker.yaml
./careen apply configs/docker.yaml
```

Building instructions are contained in each individual project. For an example see the Jenkins job here https://github.com/samsung-cnct/kraken-ci-jobs WIP

## Repository Patch Set Specification

## Options
### Root Options
| Key Name | Required | Type | Description|
| --- | --- | --- | --- |
| version | __Required__ | String | Version of the repository patch set |
| packages | __Required__ | Object Array | Array of package |

### package options
| Key Name | Required | Type | Description|
| --- | --- | --- | --- |
| name | __Required__ | String | Name of package |
| repo | __Required__ | String | URL of the repository |
| revision | __Required__ | String | Commit hash from the repository |
| patches | __Optional__ | Object Array | Array of patch |

### patch options
| Key Name | Required | Type | Description|
| --- | --- | --- | --- |
| name | __Required__ | String | Name of patch |
| filename | __Required__ | String | Filename of patch |
| hash | __Required__ | String | SHA-1 hash of file referred to by filename |
| documentation | __Optional__ | Object Array | Optional array of URLs to PR requests, bug reports, or other documentation |

## Example
```yaml
---
version: 0.0.1
packages:
  - name: docker
    repo: "https://github.com/docker/docker.git"
    revision: "fb99f992c081a1d433c97c99ffb46d12693eeb76"
    patches:
      - name: "Add support for setting sysctls"
        filename: docker-19265.patch
        hash: "71705e0fa7d5dc0d9495ce692e7c9b95a8ddf9ff"
        documentation:
          - "https://github.com/docker/docker/pull/19265"
```
