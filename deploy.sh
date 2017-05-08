#!/bin/sh

./build.sh && rsync -av . deploy@ta.do:/opt/ta && ssh 'deploy@ta.do' 'supervisorctl restart ta'
