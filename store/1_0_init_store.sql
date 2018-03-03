-- +migrate Up notransaction
CREATE TABLE category2(
  id		 SERIAL,
  name   VARCHAR(100) NOT NULL,
  constraint category_pk2 primary key(id)
);

CREATE TABLE product2(
  id		 SERIAL,
  category INTEGER NOT NULL,
  name   VARCHAR(200) NOT NULL,
  description TEXT NOT NULL,
  price numeric(10,2) NOT NULL,
  constraint product_pk2 primary key(id),
  constraint product_to_category2 foreign key (category) references category2(id) ON DELETE CASCADE
);

-- +migrate Down notransaction
DROP TABLE product2;
DROP TABLE category2;