
DROP USER 'schsvr'@'%';

CREATE USER 'schsvr'@'%' IDENTIFIED BY '4DY1ObpCRA6wTUYd9IMNRqfVD62y9N7s';

CREATE DATABASE schbsddb;

GRANT ALL ON schbsddb.* TO 'schsvr'@'%';

