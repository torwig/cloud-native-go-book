package fanout

func Split(source <-chan int, n int) []<-chan int {
	dests := make([]<-chan int, 0, n)

	for i := 0; i < n; i++ {
		ch := make(chan int)
		dests = append(dests, ch)

		go func() {
			defer close(ch)

			for v := range source {
				ch <- v
			}
		}()
	}

	return dests
}
