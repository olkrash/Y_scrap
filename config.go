package main

type Config struct {
	Host          string `env:"HOST,required"`
	Port          string `env:"PORT" envDefault:"3306"`
	Login         string `env:"LOGIN,required"`
	Password      string `env:"PASSWORD"`
	DB            string `env:"DB,required"`
	PictureFolder string `env:"FOLDER" envDefault:"pictures"`
}
