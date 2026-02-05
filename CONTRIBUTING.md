# Local Development Setup

This application depends on a local version of the `go-auth` library. Follow these steps to set up the development environment:

## 1. Clone Both Repositories

```bash
# Create a parent directory for both repositories
mkdir wristband-dev
cd wristband-dev

# Clone the go-auth library
git clone https://github.com/wristband-dev/go-auth.git

# Clone this demo application
git clone https://github.com/wristband-dev/go-demo-app.git
```


## 2. Replace dependency with local

At the bottom of the `go.mod` file, add the following line to use local version of `go-auth`.

```
replace github.com/wristband-dev/go-auth => ../go-auth
```

_Adjust the path accordingly._
