package main

type config struct {
	Tensorflow tensorflowConfig
	Mongo      mongoConfig
	Queue      queueConfig
}

type queueConfig struct {
	Host  string
	Topic string
}

type tensorflowConfig struct {
	ModelDir string
	Host     string
}

type mongoConfig struct {
	Host     string
	Database string
	Table    string
}
