#!/bin/bash

set -o allexport
source .env
set +o allexport

_cmd="./$1"
$_cmd