# Bundle

# Development!

This repo is being developed.

### Build 

Not automated yet.   

go version 1.16  

Manual Build:
```
sudo dnf install gpgme-devel libassuan-devel btrfs-progs-devel device-mapper-devel
cd cmd/oc-bundle && go build .
```

### Test

No unit tests... yet.  
  
  
Function test:  
1. Download pull secret and place at ~/.docker/config.josn   
2. Build and test:  
```
mkdir -p cmd/oc-bundle/test/
cp data/bundle-config.yaml cmd/oc-bundle/test/
cp data/.metadata.json cmd/oc-bundle/test/src/publish
make build
go build .
./bin/oc-bundle create full --dir=cmd/oc-bundle/test  --log-level=debug
```   




## Overview
Bundle is an OpenShift Client (oc) plugin that manages OpenShift installation, operator, and associated container image bundles.   

Bundle managment is a two part process:  
  
Part 1: Bundle Creation (Internet Connected)  
Part 2: Bundle Publishing (Disconnected)  

## Usage
```
Command: oc bundle   
  
Requires: bundle-config.yaml in target directory  
  
Sub-commands:   
create  
  Options:  
    full  
    diff  
  Flags:  
    --directory (string | oc bundle managed directory)  
    --bundle-name (string | name of bundle archive | optional)  
publish  
  Flags:  
    --from-bundle (string | archive name)  
    --install (optional)  
    --no-mirror (optional)  
    --tls-verify (boolean | optional)  
    --to-directory (string | oc bundle managed directory)  
    --to-mirror (string | target registry url)  
```  

## Bundle Spec

Note: The `bundle-config.yaml` needs to be placed in the directory specified by the --dir flag.
Note: The `bundle-config.yaml` is only used during bundle creation.
```
bundleSpec: {{ version of bundle-config | String | Required }}
targetOCPRelease: {{ version of OCP to be bundled | String | Optional }}
operators: {{ Optional }}
  - channel: {{ update stream | string | if operators, then required}}
    full: {{ bolean }}
    operatorList: {{ list of strings | if full catalog, then required}}
additionalImages: {{ list of additional images to be bundled }}
blockedImages: {{ list of blocked images / no download }}
pullSecret: {{ cloud.redhat.com pull secret }}
bundleSize: {{ Number in GB to limit bundle file size to }}
```

