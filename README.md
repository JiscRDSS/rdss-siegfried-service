[![Build Status](https://travis-ci.org/JiscRDSS/rdss-siegfried-service.svg?branch=master)](https://travis-ci.org/JiscRDSS/rdss-siegfried-service) 

# rdss-siegfried-service

This is a prototype.

## Usage

Open a terminal and run the following command and leave it running:

    $ docker run -p "8080" artefactual/rdss-siegfried-service-amd64:v0.1.0

In a different terminal run:

    $ curl -v 127.0.0.1:8080/(echo -n '/siegfried/default.sig' | base64)
