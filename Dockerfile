FROM centos:7

RUN yum install -y openldap-clients

CMD ["/bin/bash"]
