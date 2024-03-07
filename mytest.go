package main

import "WarpGPT/test"

func main() {
	// fmt.Println("Hello, World!")
	// token, err := funcaptcha.GetOpenAIArkoseToken(0, "")
	// if err != nil {
	// 	fmt.Println("r", err)
	// } else {
	// 	println("token:", token)
	// }
	// token, err := funcaptcha.GetOpenAIArkoseToken(0, "")

	//funcaptcha.Decrypt()
	//funcaptcha.GetOpenAIArkoseToken()
	//reg := test.NewRegClient("https://ipv6.lookup.test-ipv6.com/ip/?asn=1&testdomain=test-ipv6.com&testname=test_asn6", "adsd")
	//reg := test.NewRegClient("https://myip.ipip.net", "adsd")
	// reg := test.NewRegClient("https://mandrillapp.com/track/click/31165340/auth0.openai.com?p=eyJzIjoicktfQUxkRXFMWmVOQ245UHdYVGtUN0cyR1k0IiwidiI6MSwicCI6IntcInVcIjozMTE2NTM0MCxcInZcIjoxLFwidXJsXCI6XCJodHRwczpcXFwvXFxcL2F1dGgwLm9wZW5haS5jb21cXFwvdVxcXC9lbWFpbC12ZXJpZmljYXRpb24_dGlja2V0PUREeDE3OEJ0OG1qSlVRY0thRlloSHJGeERGSzRicGhXI1wiLFwiaWRcIjpcIjc4NGMwMzE0Zjc5MjQxNTJhODg5ODJlOTEzZGFmYjEzXCIsXCJ1cmxfaWRzXCI6W1wiMWM3OTUyMjNiMmQ0YmUwMjBmZDJhNTBmMmM5YzQxZjEwMThlNDU0Y1wiXX0ifQ", "adsd")
	// reg.Start()

	// auth := tools.NewAuthenticator("mygod@addmao.com", "4MjcTD.EtQ8.@sy", "")
	// auth.Begin()

	// mail := test.NewImapMail("hts999@gmail.com", "adjgnsmwkegbukty")
	// err := mail.Login()
	// if err != nil {
	// 	println("login Fail :", err)
	// }
	// msg, err := mail.GetNewMail("b0228@addmao.com", 10)
	// if err != nil {
	// 	println("login Fail :", err)
	// } else {
	// 	println("login success :", msg)
	// }

	//test.MyProxy()

	// s := test.Socket{Addr: "0.0.0.0:6050"}
	// s.Start()

	test.StartHttpProxy("0.0.0.0:6050")
}
