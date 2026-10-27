module github.com/zaldk/mirea-sem1-tid-private/pr6

go 1.24.0

replace gorm.io/gorm => github.com/go-gorm/gorm v1.31.2

replace gorm.io/driver/sqlite => github.com/go-gorm/sqlite v1.6.0

replace golang.org/x/text => github.com/golang/text v0.43.0

require (
	github.com/go-chi/chi/v5 v5.5.0
	gorm.io/driver/sqlite v1.6.0
	gorm.io/gorm v1.30.0
)

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
	golang.org/x/text v0.20.0 // indirect
)
