package main

import (
	"fmt"
)

type Agenda map[string]Day

type Draftman struct {
	contactInfo ContactInfo
	agenda      Agenda
	history     map[string]Day
}

func (d Draftman) getAgenda() Agenda {
	return d.agenda
}

func (d Draftman) getContactInfo() ContactInfo {
	return d.contactInfo
}

func (d *Draftman) updateContactInfo(newInfo ContactInfo) {
	d.contactInfo = newInfo
}

func (d *Draftman) addScheduleReqToDate(date string, s Schedule) error {
	if err := d.agenda.checkDate(date); err != nil {
		return err
	}

	if v, exist := d.agenda[date]; exist {
		v.addScheduleReq(&s)
		d.agenda[date] = v
	} else {
		d.agenda[date] = Day{date, map[string]Schedule{}, map[string]Schedule{s.id: s}, false, 6}
	}
	return nil
}

func (d *Draftman) addConfirmedScheduleToDate(date string, s Schedule) error {
	if err := d.agenda.checkDate(date); err != nil {
		return err
	}

	if v, exist := d.agenda[date]; exist {
		v.confirmedSchedules[s.id] = s
	} else {
		d.agenda[date] = Day{date, map[string]Schedule{s.id: s}, map[string]Schedule{}, false, 6}
	}
	return nil
}

func (d *Draftman) confirmScheduleReq(date string, id string) error {
	if err := d.agenda.checkDate(date); err != nil {
		return err
	}

	if v, exist := d.agenda[date]; exist {
		if _, exist := v.schedulesRequests[id]; exist {
			v.confirmScheduleReq(id)
			d.agenda[date] = v
			return nil
		}
		return fmt.Errorf("could find not schedule %v on day %v", id, date)
	}
	return fmt.Errorf("could not find any schedule on day %v", date)
}

func (d *Draftman) rejectScheduleReq(date string, id string) error {
	if v, exist := d.agenda[date]; exist {
		if err := v.rejectScheduleReq(id); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("could not find schedule %v on day %v", id, date)
}

func (d *Draftman) toggleFreeDay(date string) {
	if v, exist := d.agenda[date]; exist {
		v.isFreeDay = !v.isFreeDay
		d.agenda[date] = v
		return
	}
	d.agenda[date] = Day{date, map[string]Schedule{}, map[string]Schedule{}, true, 0}
}

func (d *Draftman) onDayEnd(today string) {
	if v, exist := d.agenda[today]; exist {
		d.history[today] = v
		delete(d.agenda, today)
	}
}

func (d *Draftman) setScheduleOver(date string, id string) error {
	if v, exist := d.agenda[date]; exist {
		if s, exist := v.confirmedSchedules[id]; exist {
			s.done = true
			return nil
		}
		return fmt.Errorf("could not schedule %v on day %v", id, date)
	}
	return fmt.Errorf("could not find any schedule on day %v", date)
}

func (d *Draftman) changeScheduleDay(sDate string, sId string, targetDate string) error {

	if v, exist := d.agenda[sDate]; exist {
		if _, exist := v.schedulesRequests[sId]; exist {
			s, src, err := removeScheduleFromMap(sId, v.schedulesRequests)
			if err != nil {
				return err
			}
			v.schedulesRequests = src
			d.addScheduleReqToDate(targetDate, *s)
			return nil
		} else if _, exist := v.confirmedSchedules[sId]; exist {
			s, src, err := removeScheduleFromMap(sId, v.confirmedSchedules)
			if err != nil {
				return err
			}
			v.confirmedSchedules = src
			d.addConfirmedScheduleToDate(targetDate, *s)
		}
		d.agenda[sDate] = v
		return nil
	}
	return fmt.Errorf("could not find any schedule on day %v ", sDate)
}

func (a Agenda) checkDate(date string) error {
	if v, exist := a[date]; exist {
		if v.isDayFull() {
			return fmt.Errorf("day %v is full", date)
		}
	}
	return nil
}
