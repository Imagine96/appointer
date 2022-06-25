package main

import (
	"errors"
)

type scheduleDetails []string

type contactInfo struct {
	firstName   string
	lastName    string
	email       string
	phoneNumber string
}

type schedule struct {
	id                string
	direction         string
	hour              string
	importantDetails  scheduleDetails
	clientContactInfo contactInfo
	done              bool
}

type day struct {
	date               string
	confirmedSchedules map[string]schedule
	schedulesRequests  map[string]schedule
	isFreeDay          bool
	maxSchedules       int
}

func (d *day) addScheduleReq(s *schedule) error {
	if d.isFreeDay {
		return errors.New("day " + d.date + "is not available for schedule")
	}
	if d.isDayFull() {
		return errors.New("day " + d.date + " is full")
	}
	d.schedulesRequests[s.id] = *s
	return nil
}

func (d *day) rejectScheduleReq(id string) error {
	_, _, err := removeScheduleFromMap(id, d.schedulesRequests)
	if err != nil {
		return err
	}
	return nil
}

func (d *day) confirmScheduleReq(id string) error {
	if d.isDayFull() {
		return errors.New("day " + d.date + " is full")
	}
	targetSchedule, remain, err := removeScheduleFromMap(id, d.schedulesRequests)
	d.schedulesRequests = remain
	if err != nil {
		return err
	}
	d.confirmedSchedules[targetSchedule.id] = *targetSchedule
	if len(d.confirmedSchedules) == d.maxSchedules {
		//redirect to select other day
		d.schedulesRequests = map[string]schedule{}
	}
	return nil
}

func (d *day) onScheduleDone(id string) error {

	if _, exist := d.confirmedSchedules[id]; exist {
		if d.confirmedSchedules[id].done {
			return errors.New("The schedule is over")
		}
		entry := d.confirmedSchedules[id]
		entry.done = true
		d.confirmedSchedules[id] = entry
	}
	return nil
}

func (d day) isDayFull() bool {
	return d.maxSchedules == len(d.confirmedSchedules)
}

func (s schedule) getContactInfo() contactInfo {
	return s.clientContactInfo
}

func (s *schedule) updateContactInfo(newInfo contactInfo) {
	s.clientContactInfo = newInfo
}

func removeScheduleFromMap(targetId string, src map[string]schedule) (*schedule, map[string]schedule, error) {
	if targetSchedule, exist := src[targetId]; exist {
		delete(src, targetId)
		return &targetSchedule, src, nil
	} else {
		return nil, nil, errors.New("could not find schedule " + targetId)
	}
}
