package main

import (
	"testing"
)

const agendaId = "const-agenda_id"
const targetDate = "30-6-2022"
const targetDate2 = "10-10-2022"

func TestDraftman(t *testing.T) {
	draftmanContactInfo := contactInfo{"draftman name", "draftman lastname", "draftman email", "draftman number"}
	agenda := agenda{agendaId, map[string]day{}}
	history := map[string]day{}
	clientContactInfo := contactInfo{"client name", "client lastname", "client email", "client number"}
	targetSchedule := schedule{scheduleId, scheduleDirection, hour, scheduleDetails{}, clientContactInfo, false}

	draftmanTest := draftman{draftmanContactInfo, agenda, history}

	newDay01 := day{targetDate2, map[string]schedule{}, map[string]schedule{}, false, 6}

	//has empty an agenda
	if len(draftmanTest.agenda.data) != 0 {
		t.Error("agenda's data must start empty")
	}

	//date should not exist
	if v, exist := draftmanTest.agenda.data[date]; exist {
		t.Errorf("agenda's data must start empty \n %v", v)
	}

	//adding schedule req to a non existing date
	err := draftmanTest.addScheduleReqToDate(date, targetSchedule)
	if err != nil {
		t.Error(err)
	}

	//now the date must exist
	if _, exist := draftmanTest.agenda.data[date]; !exist {
		t.Errorf("agenda's data must contain date \n %v", date)
	}

	//confirming schedule
	err = draftmanTest.confirmScheduleReq(date, targetSchedule.id)
	if err != nil {
		t.Error(err)
	}
	//schedule must be confirmed
	if _, exist := draftmanTest.agenda.data[date].confirmedSchedules[targetSchedule.id]; !exist {
		t.Errorf("agenda's must contain date \n %v \n as confirmed schedule", date)
	}

	//moving confirmed schedule another to date
	err = draftmanTest.changeScheduleDay(date, targetSchedule.id, targetDate)
	if err != nil {
		t.Error(err)
	}

	//new date must was created
	if _, exist := draftmanTest.agenda.data[targetDate]; !exist {
		t.Errorf("agenda's data must contain date \n %v", targetDate)
	}

	//schedule must be confirmed on new date
	if _, exist := draftmanTest.agenda.data[targetDate].confirmedSchedules[targetSchedule.id]; !exist {
		t.Errorf("agenda's must contain date \n %v \n as confirmed schedule", targetDate)
	}

	//toggles free day
	draftmanTest.agenda.data[targetDate2] = newDay01
	draftmanTest.toggleFreeDay(targetDate2)
	//"10-10-2022" must be a free day
	if !draftmanTest.agenda.data[targetDate2].isFreeDay {
		t.Error("could'nt update free day \n expected true != false")
	}
	draftmanTest.toggleFreeDay(targetDate2)
	//"10-10-2022" should not be a free day any more
	if draftmanTest.agenda.data[targetDate2].isFreeDay {
		t.Error("could'nt update free day \n expected false != true")
	}

	//confirms schedules within a date
	draftmanTest = draftman{draftmanContactInfo, agenda, history}
	err = draftmanTest.addScheduleReqToDate(date, targetSchedule)
	if err != nil {
		t.Error(err)
	}

	draftmanTest.confirmScheduleReq(date, targetSchedule.id)
	//date's schedules req and confirmed len should be 0 and 1
	if len(draftmanTest.agenda.data[date].schedulesRequests) != 0 {
		t.Errorf("expected %v != %v", 0, len(draftmanTest.agenda.data[date].schedulesRequests))
	}
	if len(draftmanTest.agenda.data[date].confirmedSchedules) != 1 {
		t.Errorf("expected %v != %v", 1, len(draftmanTest.agenda.data[date].schedulesRequests))
	}

	//schedule must be confirmed on date
	if _, exist := draftmanTest.agenda.data[targetDate].confirmedSchedules[targetSchedule.id]; !exist {
		t.Errorf("agenda's must contain date \n %v \n as confirmed schedule", targetDate)
	}
}
