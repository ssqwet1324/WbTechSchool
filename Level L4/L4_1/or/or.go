package or

// Or - объединение каналов в один
func Or(channels ...<-chan interface{}) <-chan interface{} {
	chDone := make(chan interface{})

	go func() {
		defer close(chDone)

		for _, ch := range channels {
			go func(c <-chan interface{}) {
				<-c
				select {
				case chDone <- struct{}{}:
				default:
				}
			}(ch)
		}
	}()

	return chDone
}
