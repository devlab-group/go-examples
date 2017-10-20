task:wrap() {
  cd $CWD
  local PORT=${PORT:-4000}

  docker run --rm \
    -v $PWD:/app/src/project \
    -v $PWD/.go:/go \
    -w /app/src/project \
    -p $PORT:8080 \
    -e GOPATH=/go:/app \
    -ti golang $@
}

task:go() {
  task:wrap go $@
}

task:dep() {
  task:wrap dep $@
}

task:build() {
  local FILE=$1
  shift 1

  task:go build -o build/$FILE $@
}

task:init() {
  # Install dependency management tool
  task:go get -u github.com/golang/dep/cmd/dep
}
