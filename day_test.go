package main

import (
	"reflect"
	"testing"
)

const maxSchedules = 6
const scheduleId = "const-schedule-id"
const scheduleDirection = "const-schedule-dir"
const timeFormat = "YY-MM-YY"
const date = "30-06-2022"
const hour = "14:00"

func TestDay(t *testing.T) {

	clientContactInfo := contactInfo{"client name", "client lastname", "client email", "client number"}
	targetSchedule := schedule{scheduleId, scheduleDirection, hour, scheduleDetails{}, clientContactInfo, false}
	newSchedule01 := schedule{"schedule-1", scheduleDirection, hour, scheduleDetails{}, clientContactInfo, false}
	newSchedule02 := schedule{"schedule-2", scheduleDirection, hour, scheduleDetails{}, clientContactInfo, false}
	newSchedule03 := schedule{"schedule-3", scheduleDirection, hour, scheduleDetails{}, clientContactInfo, false}
	newSchedule04 := schedule{"schedule-4", scheduleDirection, hour, scheduleDetails{}, clientContactInfo, false}
	newSchedule05 := schedule{"schedule-5", scheduleDirection, hour, scheduleDetails{}, clientContactInfo, false}
	newSchedule06 := schedule{"schedule-6", scheduleDirection, hour, scheduleDetails{}, clientContactInfo, false}
	schedules := []schedule{newSchedule01, newSchedule02, newSchedule03, newSchedule04, newSchedule05, newSchedule06}

	//in free day
	testDay := day{date, map[string]schedule{}, map[string]schedule{}, true, maxSchedules}

	//should not append requests on free days
	err := testDay.addScheduleReq(&targetSchedule)
	if err == nil {
		t.Error("should not append requests on free days")
	}

	//testDay is not full
	if testDay.isDayFull() {
		t.Errorf("testDay is not full \n expected %v != %v", false, testDay.isDayFull())
	}

	//in working day

	testDay = day{date, map[string]schedule{}, map[string]schedule{}, false, maxSchedules}

	//testDay is not full
	if testDay.isDayFull() {
		t.Errorf("testDay is not full \n expected %v != %v", false, testDay.isDayFull())
	}

	//push new schedule requests
	err = testDay.addScheduleReq(&targetSchedule)
	if err != nil {
		t.Error(err)
	}

	//checking the just pushed schedule
	if _, exist := testDay.schedulesRequests[targetSchedule.id]; !exist {
		t.Errorf("could not find the schedule %v on scheduleRequests", targetSchedule.id)
	}

	//now lengths should be 1 and 0 for scheduleReq and confirmedSchedules
	if len(testDay.schedulesRequests) != 1 {
		t.Errorf("bad schedulesRequests length \n expected %v != 1", len(testDay.schedulesRequests))
	}
	if len(testDay.confirmedSchedules) != 0 {
		t.Errorf("bad confirmedSchedules length \n expected %v != 0", len(testDay.schedulesRequests))
	}

	//updates contact info in schedule requests
	newTargetContactInfo := contactInfo{"new client name", "new client lastname", "new client email", "new client number"}
	if v, exist := testDay.schedulesRequests[scheduleId]; exist {
		v.updateContactInfo(newTargetContactInfo)
		testDay.schedulesRequests[scheduleId] = v
	}
	if !reflect.DeepEqual(testDay.schedulesRequests[scheduleId].getContactInfo(), newTargetContactInfo) {
		t.Errorf("contact info did not update, \n expected %v != %v", newTargetContactInfo, testDay.schedulesRequests[scheduleId].getContactInfo())
	}

	//upgrading "const-schedule-id" to confirmedSchedules
	err = testDay.confirmScheduleReq(scheduleId)
	if err != nil {
		t.Error(err)
	}

	//now lengths should be 0 and 1 for scheduleReq and confirmedSchedules
	if len(testDay.schedulesRequests) != 0 {
		t.Errorf("bad schedulesRequests length \n expected %v != 0", len(testDay.schedulesRequests))
	}
	if len(testDay.confirmedSchedules) != 1 {
		t.Errorf("bad confirmedSchedules length \n expected %v != 1", len(testDay.schedulesRequests))
	}

	//checking the just confirmed schedule
	if _, exist := testDay.confirmedSchedules[targetSchedule.id]; !exist {
		t.Errorf("could not find schedule %v on confirmedSchedule", targetSchedule.id)
	}

	//adding more schedule
	for _, s := range schedules {
		err := testDay.addScheduleReq(&s)
		if err != nil {
			t.Errorf("could not find add %v req", s.id)
		}
	}

	//making 3 pass
	if err := testDay.confirmScheduleReq("schedule-1"); err != nil {
		t.Error(err)
	}
	if err := testDay.confirmScheduleReq("schedule-2"); err != nil {
		t.Error(err)
	}
	if err := testDay.confirmScheduleReq("schedule-3"); err != nil {
		t.Error(err)
	}
	if err := testDay.confirmScheduleReq("schedule-4"); err != nil {
		t.Error(err)
	}

	//rejecting one

	if err := testDay.rejectScheduleReq("schedule-6"); err != nil {
		t.Error(err)
	}

	//now lengths should be 1 and 5 for scheduleReq and confirmedSchedules
	if len(testDay.schedulesRequests) != 1 {
		t.Errorf("bad schedulesRequests length \n expected %v != 2", len(testDay.schedulesRequests))
	}
	if len(testDay.confirmedSchedules) != 5 {
		t.Errorf("bad confirmedSchedules length \n expected %v != 5", len(testDay.schedulesRequests))
	}

	//day is not full
	if testDay.isDayFull() {
		t.Errorf("day should not be full \n expected %v != %v", false, testDay.isDayFull())
	}

	//filling day
	if err := testDay.confirmScheduleReq("schedule-5"); err != nil {
		t.Error(err)
	}

	//day is now full
	if !testDay.isDayFull() {
		t.Errorf("day should not be full \n expected %v != %v", true, testDay.isDayFull())
	}

	//now lengths should be 0 and 6 for scheduleReq and confirmedSchedules
	if len(testDay.schedulesRequests) != 0 {
		t.Errorf("bad schedulesRequests length \n expected %v != 0", len(testDay.schedulesRequests))
	}
	if len(testDay.confirmedSchedules) != maxSchedules {
		t.Errorf("bad confirmedSchedules length \n expected %v != %v", len(testDay.schedulesRequests), maxSchedules)
	}

	//should fail if try to push a new schedule req
	if err := testDay.addScheduleReq(&newSchedule06); err == nil {
		t.Errorf("should fail if try to push a new schedule req \n expected 0 != %v", len(testDay.schedulesRequests))
	}

	//updates contact info in confirmed schedules
	if v, exist := testDay.confirmedSchedules[scheduleId]; exist {
		v.updateContactInfo(clientContactInfo)
		testDay.confirmedSchedules[scheduleId] = v
	}
	if !reflect.DeepEqual(testDay.confirmedSchedules[scheduleId].getContactInfo(), clientContactInfo) {
		t.Errorf("contact info did not update, \n expected %v != %v", newTargetContactInfo, testDay.confirmedSchedules[scheduleId].getContactInfo())
	}
}
