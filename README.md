# goproxy

This project is a template for a simple golang based reverse proxy
server. It contains all the basics for a useful HTTP reverse proxy,
you just need to add your custom logic.

# Getting Started

goproxy has 3 configuration options:

```
target: The target server (must be a full address)
port: The port you want to run goproxy on
logdir: The path you want goproxy to write its logs to
```

For example:

```
$ goproxy -target http://localhost:4567 -port 8080 -logdir logs
```

# Working with goproxy

You can simply compile and run goproxy, but if you just want a reverse
proxy with no additional options, I recommend using something like
NGINX. If you need custom logic in your proxy, you just need to modify
the `handle` function and insert your code. Here you have full access
to the request and can perform any necessary actions on the request
there. This includes additional processing, routing, and even
termination of unwanted traffic. You are of course welcome to modify
anything in the project, but the intent is that the everything else is
taken care of for you already.

# Pull Requests Welcome

Contributions and discussion is always welcome. If you feel something
is missing or broken, please file an issue or a pull request.
