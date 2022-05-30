FROM alpine

COPY autoscaler-testcase ./autoscaler-testcase

RUN ["ls"]

RUN ["chmod +x ./autoscaler-testcase"]

ENTRYPOINT ["./autoscaler-testcase"]