FROM centos:7

ADD bin/download /

RUN chmod 777 /download

ENV PARAMS=""

ENTRYPOINT ["sh","-c","/download $PARAMS"]