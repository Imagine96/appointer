package main

import (
	"errors"
)

type ScheduleDetails []string

type ContactInfo struct {
	firstName   string
	lastName    string
	email       string
	phoneNumber string
}

type Schedule struct {
	id                string
	direction         string
	hour              string
	importantDetails  ScheduleDetails
	clientContactInfo ContactInfo
	done              bool
}

type Day struct {
	date               string
	confirmedSchedules map[string]Schedule
	schedulesRequests  map[string]Schedule
	isFreeDay          bool
	maxSchedules       int
}

func (d *Day) addScheduleReq(s *Schedule) error {
	if d.isFreeDay {
		return errors.New("day " + d.date + "is not available for schedule")
	}
	if d.isDayFull() {
		return errors.New("day " + d.date + " is full")
	}
	d.schedulesRequests[s.id] = *s
	return nil
}

func (d *Day) rejectScheduleReq(id string) error {
	_, _, err := removeScheduleFromMap(id, d.schedulesRequests)
	if err != nil {
		return err
	}
	return nil
}

func (d *Day) confirmScheduleReq(id string) error {
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
		d.schedulesRequests = map[string]Schedule{}
	}
	return nil
}

func (d *Day) onScheduleDone(id string) error {

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

func (d Day) isDayFull() bool {
	return d.maxSchedules == len(d.confirmedSchedules)
}

func (s Schedule) getContactInfo() ContactInfo {
	return s.clientContactInfo
}

func (s *Schedule) updateContactInfo(newInfo ContactInfo) {
	s.clientContactInfo = newInfo
}

func removeScheduleFromMap(targetId string, src map[string]Schedule) (*Schedule, map[string]Schedule, error) {
	if targetSchedule, exist := src[targetId]; exist {
		delete(src, targetId)
		return &targetSchedule, src, nil
	} else {
		return nil, nil, errors.New("could not find schedule " + targetId)
	}
}
