# Take pictures with a Raspberry Pi when someone tweet about us

Whenever someone tweet on our selected key the camera takes a pictures and
upload it on AWS S3.

## Run it

Prepare a `config.json` file (use the dist as example)

Just run it!

```
go run main.go
```

Now tweet something:

```
Hey #cloudconf2015 where is my picture?
```

