FROM ubuntu:bionic

RUN apt-get update -y

RUN apt-get install -y openssh-server python3 python3-pyftpdlib

# Configure SSH Server
# /usr/sbin/sshd -D
RUN mkdir /var/run/sshd
RUN echo 'root:password' | chpasswd
RUN sed -i 's/#PermitRootLogin/PermitRootLogin yes #/' /etc/ssh/sshd_config
RUN sed 's@session\s*required\s*pam_loginuid.so@session optional pam_loginuid.so@g' -i /etc/pam.d/sshd

# HTTP Server
# python3 -m http.server 80

# FTP Server
# python3 -m pyftpdlib -p 21

ENTRYPOINT [ "/bin/bash" ]
