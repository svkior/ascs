FROM ubuntu:14.04
MAINTAINER Sergey V. Kior <svkior@gmail.com>
RUN apt-get -qq update
RUN apt-get -qqy upgrade
COPY ./work /work
CMD ["/work/main_cli",  "-c", "/work/lnx_switch24_docker.json"]

