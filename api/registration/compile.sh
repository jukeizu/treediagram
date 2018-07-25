#!/bin/bash

protoc registration.proto --go_out=plugins=grpc:.
