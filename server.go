package watchman

type server struct {
	commands chan<- string
	events   <-chan string
}
