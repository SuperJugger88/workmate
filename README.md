### Building and running your application

When you're ready, start your application in docker by running:
`make dc`.

Or you can run it alternatively in local k8s cluster (ingress controller is needed to be installed):  
`make k8s`.

Your services will be available at following routes: 

```
http://localhost/tasks/*     - Application API Endpoints
http://localhost/swagger     - Swagger Specs    
http://localhost/admin/*     - Caddy Admin Endpoints    
```

### Deploying your application to the cloud

First, build your image, e.g.: `docker build -t workmate-app .`.
If your cloud uses a different CPU architecture than your development
machine (e.g., you are on a Mac M1 and your cloud provider is amd64),
you'll want to build the image for that platform, e.g.:
`docker build --platform=linux/amd64 -t workmate-app .`.

Then, push it to your registry, e.g. `docker push myregistry.com/myapp`.

Consult Docker's [getting started](https://docs.docker.com/go/get-started-sharing/)
docs for more detail on building and pushing.

### References
* [Docker's Go guide](https://docs.docker.com/language/golang/)