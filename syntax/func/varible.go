package main

func YourName(name string, aliases ...string) {
	if len(aliases) > 0 {
		println(aliases[0])
	}
}

func YourNameInvoke() {
	YourName("smf")
	YourName("smf", "sun")
	YourName("smf", "sun", "fff")
}
