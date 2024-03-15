package service

// Initialization Process
func init() {
	//execution
	Connect()

	//if not acquired
	if StreamConf.Country != "" {
		StreamConf = GetStreamConf()
	}

	//start streaming
	if false /*StreamConf.Status*/ {

		//set current stream_id
		SetCurrentSub()
		FmtPrintln("blue", "set current stream_id")

		//open websocket
		go Gvg()
		FmtPrintln("blue", "open websocket")

		//send stream_id to start streaming
		Buffer <- GetBuffer()
		<-ReqFlg
	}
}
