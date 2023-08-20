package server

func Init(path string) {
	r := NewRouter()

	r.Run(path)
}
