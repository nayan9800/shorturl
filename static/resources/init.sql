CREATE TABLE IF NOT EXISTS users (
 	user_id serial PRIMARY KEY,
  	name VARCHAR NOT NULL,
  	password VARCHAR UNIQUE NOT NULL,
	email VARCHAR UNIQUE NOT NULL
  );
    
CREATE Table IF NOT EXISTS urlservice (
 	urlservice_id serial PRIMARY KEY NOT NULL,
  	url VARCHAR UNIQUE NOT NULL,
  	code VARCHAR UNIQUE NOT NULL,
  	user_id INT,
  	FOREIGN KEY(user_id) 
  	REFERENCES users(user_id)
);









  
  	
  	