package main

import (
	"testing"
)

const agendaId = "const-agenda_id"
const targetDate = "30-6-2022"
const targetDate2 = "10-10-2022"

func TestDraftman(t *testing.T) {
	draftmanContactInfo := ContactInfo{"draftman name", "draftman lastname", "draftman email", "draftman number"}
	agenda := map[string]Day{}
	history := map[string]Day{}
	clientContactInfo := ContactInfo{"client name", "client lastname", "client email", "client number"}
	targetSchedule := Schedule{scheduleId, scheduleDirection, hour, ScheduleDetails{}, clientContactInfo, false}

	draftmanTest := Draftman{draftmanContactInfo, agenda, history}

	newDay01 := Day{targetDate2, map[string]Schedule{}, map[string]Schedule{}, false, 6}

	//has empty an agenda
	if len(draftmanTest.agenda) != 0 {
		t.Error("agenda's data must start empty")
	}

	//date should not exist
	if v, exist := draftmanTest.agenda[date]; exist {
		t.Errorf("agenda's data must start empty \n %v", v)
	}

	//adding schedule req to a non existing date
	err := draftmanTest.addScheduleReqToDate(date, targetSchedule)
	if err != nil {
		t.Error(err)
	}

	//now the date must exist
	if _, exist := draftmanTest.agenda[date]; !exist {
		t.Errorf("agenda's data must contain date \n %v", date)
	}

	//confirming schedule
	err = draftmanTest.confirmScheduleReq(date, targetSchedule.id)
	if err != nil {
		t.Error(err)
	}
	//schedule must be confirmed
	if _, exist := draftmanTest.agenda[date].confirmedSchedules[targetSchedule.id]; !exist {
		t.Errorf("agenda's must contain date \n %v \n as confirmed schedule", date)
	}

	//moving confirmed schedule another to date
	err = draftmanTest.changeScheduleDay(date, targetSchedule.id, targetDate)
	if err != nil {
		t.Error(err)
	}

	//new date must was created
	if _, exist := draftmanTest.agenda[targetDate]; !exist {
		t.Errorf("agenda's data must contain date \n %v", targetDate)
	}

	//schedule must be confirmed on new date
	if _, exist := draftmanTest.agenda[targetDate].confirmedSchedules[targetSchedule.id]; !exist {
		t.Errorf("agenda's must contain date \n %v \n as confirmed schedule", targetDate)
	}

	//toggles free day
	draftmanTest.agenda[targetDate2] = newDay01
	draftmanTest.toggleFreeDay(targetDate2)
	//"10-10-2022" must be a free day
	if !draftmanTest.agenda[targetDate2].isFreeDay {
		t.Error("could'nt update free day \n expected true != false")
	}
	draftmanTest.toggleFreeDay(targetDate2)
	//"10-10-2022" should not be a free day any more
	if draftmanTest.agenda[targetDate2].isFreeDay {
		t.Error("could'nt update free day \n expected false != true")
	}

	//confirms schedules within a date
	draftmanTest = Draftman{draftmanContactInfo, agenda, history}
	err = draftmanTest.addScheduleReqToDate(date, targetSchedule)
	if err != nil {
		t.Error(err)
	}

	draftmanTest.confirmScheduleReq(date, targetSchedule.id)
	//date's schedules req and confirmed len should be 0 and 1
	if len(draftmanTest.agenda[date].schedulesRequests) != 0 {
		t.Errorf("expected %v != %v", 0, len(draftmanTest.agenda[date].schedulesRequests))
	}
	if len(draftmanTest.agenda[date].confirmedSchedules) != 1 {
		t.Errorf("expected %v != %v", 1, len(draftmanTest.agenda[date].schedulesRequests))
	}

	//schedule must be confirmed on date
	if _, exist := draftmanTest.agenda[targetDate].confirmedSchedules[targetSchedule.id]; !exist {
		t.Errorf("agenda's must contain date \n %v \n as confirmed schedule", targetDate)
	}
}
