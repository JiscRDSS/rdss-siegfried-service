[![Build Status](https://travis-ci.org/JiscRDSS/rdss-siegfried-service.svg?branch=master)](https://travis-ci.org/JiscRDSS/rdss-siegfried-service) 

# rdss-siegfried-service

This is a prototype.

## Usage

Open a terminal and run the following command and leave it running:

    $ docker run --rm -p "8080:8080" artefactual/rdss-siegfried-service-amd64:v0.2.0

In a different terminal run:

    $ curl -v 127.0.0.1:8080/$(echo -n '/siegfried/default.sig' | base64)

## Sharing data with the service

Bind mount a volume to make it available in the container as follows:

    $ docker run \
    	--rm -p "8080:8080" \
    	--mount "type=bind,src=/mnt/my-data,dst=/mnt/my-data,readonly" \
    		artefactual/rdss-siegfried-service-amd64:v0.2.0

In a different terminal run:

    $ curl -v 127.0.0.1.8080/$(echo -n '/mnt/my-data/foobar.iso' | base64)

We've shared our local directory `/mnt/my-data` with the container and tried
to identify the `/mnt/my-data/foobar.iso` file.
