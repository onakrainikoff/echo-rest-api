-- +migrate Up
CREATE TABLE category(
  id		 SERIAL,
  name   VARCHAR(100) NOT NULL,
  constraint category_pk primary key(id)
);

CREATE TABLE product(
  id		 SERIAL,
  category INTEGER NOT NULL,
  name   VARCHAR(200) NOT NULL,
  description TEXT NOT NULL,
  price numeric(10,2) NOT NULL,
  constraint product_pk primary key(id),
  constraint product_to_category foreign key (category) references category(id) ON DELETE CASCADE
);

-- +migrate Down
DROP TABLE product;
DROP TABLE category;