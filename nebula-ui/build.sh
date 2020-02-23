#!/bin/bash
cd $(dirname $0)

yarn
yarn run build
