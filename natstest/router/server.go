package router

import (
	"encoding/json"
	"fmt"
	"io"
	"natstest/mq"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
)

func StartServer(port int) {
	LogConfig()
	mq.JeStreamConnect()
	defer mq.JeStreamDisConnect()

	mq.SubMsg()

	http.HandleFunc("/js", jsHandle)
	http.HandleFunc("/many", sendManyToMq)
	http.HandleFunc("/req", natsResponse)
	http.HandleFunc("/noack", noack)

	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func LogConfig() {
	// 設定顯示等級
	// logrus.SetLevel(logrus.InfoLevel)
	// 設定
	logrus.SetReportCaller(true)
	// logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.JSONFormatter{})
}

func jsHandle(w http.ResponseWriter, r *http.Request) {
	msg := jsonDecode(r.Body)

	pubMsg := mq.PublishMsg{Subject: mq.DEFAULT_SUBJECT, Message: msg}
	err := pubMsg.Publish()
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Something Wrong: %s", err)
		return
	}

	fmt.Fprintf(w, "Message %s is send", msg)
}

// 一次丟多一點
func sendManyToMq(w http.ResponseWriter, r *http.Request) {
	msg := jsonDecode(r.Body)
	defer r.Body.Close()

	for i := 0; i < 100; i++ {
		pm := mq.PublishMsg{
			Subject: mq.DEFAULT_SUBJECT,
			Message: msg["msg"] + strconv.Itoa(i),
		}

		err := pm.Publish()
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "Something Wrong: %s", err)
			return
		}

		fmt.Fprintf(w, "Message %s is send", msg)
	}
}

func natsResponse(w http.ResponseWriter, r *http.Request) {
	jReqMap := jsonDecode(r.Body)

	req := &mq.NatsRequest{
		Subject: "re",
		Message: jReqMap["msg"],
	}

	res := req.SendRequest()

	fmt.Fprintf(w, res.Message)
}

func noack(w http.ResponseWriter, r *http.Request) {
	jMsg := jsonDecode(r.Body)

	pub := mq.PublishMsg{
		Subject: "jetStream.noack",
		Message: jMsg["msg"],
	}

	if err := pub.Publish(); err == nil {
		fmt.Fprint(w, "okㄅ")
	}

}

func jsonDecode(body io.Reader) map[string]string {
	bodyMap := map[string]string{}
	whatErr(json.NewDecoder(body).Decode(&bodyMap))

	return bodyMap
}

func whatErr(err error) {
	if err != nil {
		logrus.Error(err)
	}
}
