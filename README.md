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

Further improvements
--------------------

I tried to keep it really simple and minimalistic because it supposed to be a µ-service. In case the project is intended to grow into something monolithic I would, for example, put all the components into a single `server` structure, group handlers into separate packages, etc.

##### What I would do next:

- [ ] Move out sample data from tests' code into separate text files
- [ ] Better error handling, for example, through using JSON Schema for input validation
- [ ] Embed Swagger UI as a static asset
- [ ] Authentication & authorisation (using JWT, for example)
- [ ] [Versioning](https://stripe.com/blog/api-versioning)
- [ ] Status endpoint for health check and monitoring (something like `/healthz`)
- [ ] Pagination (for large responses)
- [ ] Caching
- [ ] [Tune server timeouts](https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779) if necessary
