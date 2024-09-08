package repository

const (
	selectPersonsQuery = `select * from persons offset $1 limit $2;`
	insertPersonQuery  = `insert into persons(name, age, address, work) values (:name, :age, :address, :work) returning *;`
	selectPersonQuery  = `select * from persons where id=$1 limit 1;`
	updatePersonQuery  = `update persons set name=:name, age=:age, address=:address, work=:work where id=:id returning *;`
	deletePersonQuery  = `delete from persons where id=$1;`
)
