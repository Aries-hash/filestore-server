FROM centos:7

ADD bin/upload /

RUN chmod 777 /upload

ENV PARAMS=""

ENTRYPOINT ["sh","-c","/upload $PARAMS"]