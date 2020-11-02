"# RabbitMQConsumer" 

func Initialize(server *Server) {

	conn, err := amqp.Dial(os.Getenv("RABITTMQ_CRED"))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	defer ch.Close()

	msgs, err := ch.Consume(os.Getenv("QUEUE_NAME"), "", true, false, false, false, nil)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	var M model.Recieved
	var H []model.Hotel
	var R []model.Room
	var RP []model.Rateplan
	forever := make(chan bool)
	//dd := make(map[string]interface{})
	go func() {
		for d := range msgs {
			//fmt.Printf("Recieved Message: %s\n", d.Body)
			err = json.Unmarshal(d.Body, &M)
			if err != nil {
				fmt.Println(err)
				panic(err)
			}
			for _, j := range M.Offers {
				err = j.Hotel.Validate()
				if err != nil {
					fmt.Println(err)
					panic(err)
				} else {
					H = append(H, j.Hotel)
				}
				err = j.Room.Validate()
				if err != nil {
					fmt.Println(err)
					panic(err)
				} else {
					R = append(R, j.Room)
				}
				err = j.Rateplan.Validate()
				if err != nil {
					fmt.Println(err)
					panic(err)
				} else {
					RP = append(RP, j.Rateplan)
				}

			}
			for _, j := range H {
				j.Prepare()
				j.SaveHotel(server.DB)
			}
			for _, j := range R {
				j.Prepare()
				j.SaveRoom(server.DB)
			}
			for _, j := range RP {
				j.Prepare()
				j.SaveRateplan(server.DB)
			}
			fmt.Println("Hotel Details", H)
			fmt.Println("Room Details", R)
			fmt.Println("RatePlan Details", RP)

		}
	}()

	fmt.Println("Successfully connected to Rabbit")
	fmt.Println("[*] - waiting for messages")
	<-forever
}
