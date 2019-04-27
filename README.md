Form3 Payments Service
======================

Form3 service that provides a RESTful API to a payment system.

TL;DR
-----

First start:

    $ make compose

This will set up necessary infrastructure (database and networking) and bring up the service.

Use your favourite REST client to query the service on port 80 according to the Swagger specification.

For further control:

    $ make {stop|start}

Teardown:

    $ make decompose
