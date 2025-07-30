package app

func (receiver *Application) processSource() {
	receiver.source.ForwardTo(receiver.sourceChan)
}
