package main

import "fmt"

type WithContactInfo interface {
	getContactInfo() contactInfo
	updateContactInfo(contactInfo)
}

func main() {
	clientContactInfo := contactInfo{"client name", "client lastname", "client-email", "client-cellphone"}
	newContactInfo := contactInfo{"name", "lastname", "email", "cellphone"}
	newAgenda := agenda{"id", make(map[string]day)}
	newSchedule := schedule{"schedule-id", "somewhere", "13:00", scheduleDetails{}, clientContactInfo, false}
	juan := draftman{newContactInfo, newAgenda, map[string]day{}}
	juan.addScheduleReqToDate("24-06-2022", newSchedule)
	fmt.Println(juan.agenda.data["24-06-2022"].schedulesRequests)

}

func printContactInfo(w WithContactInfo) {
	data := w.getContactInfo()
	fmt.Printf("\nfirst name: %v \n last name: %v \n email: %v \n phone: %v", data.firstName, data.lastName, data.email, data.phoneNumber)
}
