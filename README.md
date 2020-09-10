# gtitle
**Fast Golang tool to get Pages title**




# Install
```
$ go get -u github.com/yghonem14/gtitle
```

## Basic Usage
gtitle accepts only Stdin Inputs:

```
$ cat yahoo.txt | gtitle
https://election2020.yahoo.com | Yahoo
```


## Concurrency

You can set the concurrency value with the `-c` flag:

```
$ cat yahoo.txt | gtitle -c 35
```

## Timeout

You can set the timeout by using the `-t`:

```
$ cat yahoo.txt | gtitle -t 3
```

## Follow Redirects

You can set Follow Redirects or not by using the `-r`:

```
$ cat yahoo.txt | gtitle -r
```
