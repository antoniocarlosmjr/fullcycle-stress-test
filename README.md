# fullcycle-stress-test

# How to run

1) Construct the docker image
```bash
docker build -t stress_test .
```

2) After the image is built, we can execute the image docker run, for example:

```bash
docker run load_tester --url=http://google.com --requests=1000 --concurrency=10
```

2.1) The parameters are:
- url: The URL to be tested
- requests: The number of requests to be made
- concurrency: The number of requests to be made concurrently


3) The output will be something like:
```bash
Report:
Total requests: 1000
Time taken: 1.000s
Status code distribution:
  [200] 1000 requests
```
