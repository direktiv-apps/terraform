#!/bin/sh

docker build -t terraform . && docker run -p 8080:8080 terraform