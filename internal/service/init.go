package service

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/harveytvt/movie-reservation-system/internal/biz/genre"
	"github.com/harveytvt/movie-reservation-system/internal/biz/movie"
	"github.com/harveytvt/movie-reservation-system/internal/biz/user"
	"github.com/harveytvt/movie-reservation-system/internal/config"
	"github.com/harveytvt/movie-reservation-system/internal/data/mysql/repository"
	"github.com/harveytvt/movie-reservation-system/internal/data/r2"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
)

var (
	r2Client r2.Client
	bunDB    bun.DB

	userRepo *repository.UserRepository
	userBiz  user.Biz

	movieRepo *repository.MovieRepository
	movieBiz  movie.Biz

	genreRepo      *repository.GenreRepository
	movieGenreRepo *repository.MovieGenreRepository
	genreBiz       genre.Biz
)

func init() {
	cfConf := config.Get().Cloudflare
	r2Client = r2.NewClient(cfConf.AccountID, cfConf.AccessKeyID, cfConf.AccessKeySecret)

	sqldb, err := sql.Open("mysql", config.Get().MysqlDSN)
	if err != nil {
		panic(err)
	}
	bunDB := bun.NewDB(sqldb, mysqldialect.New())

	userRepo = repository.NewUserRepository(bunDB)
	userBiz = user.NewBiz(userRepo)

	movieRepo = repository.NewMovieRepository(bunDB)
	genreRepo = repository.NewGenreRepository(bunDB)
	movieGenreRepo = repository.NewMovieGenreRepository(bunDB)
	movieBiz = movie.NewBiz(movieRepo, genreRepo, movieGenreRepo)

	genreBiz = genre.NewBiz(genreRepo, movieGenreRepo)
}
