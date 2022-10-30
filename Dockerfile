FROM alpine:latest

RUN mkdir /app

COPY workoutApp /app

CMD ["/app/workoutApp"]
