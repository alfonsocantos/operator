FROM rastasheep/ubuntu-sshd
COPY certs/id_rsa.pub /tmp/
RUN cat /tmp/id_rsa.pub >> /root/.ssh/authorized_keys && \
    rm /tmp/id_rsa.pub && \
    chmod 600 /root/.ssh/authorized_keys && \
    chmod 700 /root/.ssh
