CREATE TABLE "bank"(
  "id" INTEGER PRIMARY KEY AUTOINCREMENT,
  "item_id" INTEGER,
  "nom" VARCHAR(64) NULL,
  "category" INTEGER,
  "count" INTEGER
);


CREATE TABLE "favori"(
  "id" INTEGER NULL
);
