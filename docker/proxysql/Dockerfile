FROM proxysql/proxysql:2.4.7

# Update packages
RUN apt update -y --fix-missing

# Install common packages
RUN apt install -y mariadb-client

# Set custom config
ADD proxysql.cnf /etc/
