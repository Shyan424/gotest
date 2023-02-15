package main

import "natstest/router"

// var nsc *nats.Conn
// var jsc nats.JetStreamContext

func main() {
	// ns, err := nats.Connect(nats.DefaultURL, nats.UserInfo("aacc", "passowrd"))
	// whatErr(err)
	// nsc = ns

	// defer ns.Drain()

	// ns.Subscribe("nats.test", func(msg *nats.Msg) {
	// 	log.Printf("resive data %s", string(msg.Data))
	// })

	// js, _ := ns.JetStream()
	// jsc = js

	// Create a Stream
	// js.AddStream(&nats.StreamConfig{
	// 	Name:     "jetStream",
	// 	Subjects: []string{"jetStream.*"},
	// })

	// 會重複消費持久內的資料
	// js.Subscribe("jetStream.test", func(msg *nats.Msg) {
	// 	log.Printf("JetStream Subscribe %s\n", msg.Data)
	// })

	// go func() {
	// 	// 會依照durable的名稱去認定是否消費過，當該資料有ACK時就不會重複消費
	// 	sub, _ := js.PullSubscribe("jetStream.test", "con2")
	// 	for {
	// 		// 用JetStream時不可以用NextMsg()
	// 		msgs, err := sub.Fetch(1)
	// 		whatErr(err)
	// 		for _, m := range msgs {
	// 			log.Printf("PullSubscribe1 Data %s\n", m.Data)
	// 			m.Ack()
	// 		}
	// 	}
	// }()
	// sb, err := ns.SubscribeSync("nats.test")
	// whatErr(err)

	// m, err := sb.NextMsg(1 * time.Second)
	// whatErr(err)
	// log.Println(m.Data)

	// time.Sleep(1 * time.Second)

	// http.ListenAndServe(":8901", nil)

	router.StartServer(8901)
}

// func publishMsg(w http.ResponseWriter, r *http.Request) {
// 	err = ns.Publish("nats.test", []byte(string(r.PostFormValue("msg"))))
// 	whatErr(err)
// }
