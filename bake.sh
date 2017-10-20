# Run golang container
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

# Run go in workspace
task:go() {
  task:wrap go $@
}

# Execute go run ... in workspace
task:run() {
  task:go run $@
}

# Build current subproject
task:build() {
  local FILE=$1
  shift 1

  task:go build -o build/$FILE $@
}

# Create new sub-project in project's root
task:new() {
  local NAME=$1

  if [ -d "$NAME" ]
  then
    echo "Project '$NAME' already exists"
    exit 1
  fi

  mkdir $NAME
  echo "# $NAME" > $NAME/readme.md
}

# Initialize subproject
task:init() {
  cd $CWD
  # Install dependency management tool
  task:go get -u github.com/golang/dep/cmd/dep
}

# Run go's dep utility
task:dep() {
  task:wrap dep $@
}
