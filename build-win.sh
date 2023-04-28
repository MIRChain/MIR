#!/bin/bash

### Unfortunately at this time build for windows only at windows - cross compilation not working properly

env CGO_ENABLED=1 go build -v .\cmd\mir