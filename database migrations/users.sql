CREATE TABLE users (
                       id int not null AUTO_INCREMENT,
                       PRIMARY KEY (id),
                       email varchar(255),
                       pass varchar(255),
                       name varchar(255),
                       surname varchar(255),
                       birthdate timestamp
)