package main

import (
	"errors"
)

type agenda struct {
	id   string
	data map[string]day
}

type draftman struct {
	contactInfo contactInfo
	agenda      agenda
	history     map[string]day
}

func (d draftman) getAgenda() agenda {
	return d.agenda
}

func (d draftman) getContactInfo() contactInfo {
	return d.contactInfo
}

func (d *draftman) updateContactInfo(newInfo contactInfo) {
	d.contactInfo = newInfo
}

func (d *draftman) addScheduleReqToDate(date string, s schedule) error {
	if err := d.agenda.checkDate(date); err != nil {
		return err
	}

	if v, exist := d.agenda.data[date]; exist {
		v.addScheduleReq(&s)
		d.agenda.data[date] = v
	} else {
		d.agenda.data[date] = day{date, map[string]schedule{}, map[string]schedule{s.id: s}, false, 6}
	}
	return nil
}

func (d *draftman) addConfirmedScheduleToDate(date string, s schedule) error {
	if err := d.agenda.checkDate(date); err != nil {
		return err
	}

	if v, exist := d.agenda.data[date]; exist {
		v.confirmedSchedules[s.id] = s
		return nil
	} else {
		d.agenda.data[date] = day{date, map[string]schedule{s.id: s}, map[string]schedule{}, false, 6}
		return nil
	}

}

func (d *draftman) confirmScheduleReq(date string, id string) error {
	if err := d.agenda.checkDate(date); err != nil {
		return err
	}

	if v, exist := d.agenda.data[date]; exist {
		if _, exist := v.schedulesRequests[id]; exist {
			v.confirmScheduleReq(id)
			d.agenda.data[date] = v
			return nil
		}
		return errors.New("could not schedule " + id + " on day " + date)
	}
	return errors.New("could not find any schedule on day " + date)
}

func (d *draftman) rejectScheduleReq(date string, id string) error {
	if v, exist := d.agenda.data[date]; exist {
		if err := v.rejectScheduleReq(id); err != nil {
			return err
		}
		return nil
	}
	return errors.New("could not find schedule" + id + "on day" + date)
}

func (d *draftman) toggleFreeDay(date string) {
	if v, exist := d.agenda.data[date]; exist {
		v.isFreeDay = !v.isFreeDay
		d.agenda.data[date] = v
		return
	}
	d.agenda.data[date] = day{date, map[string]schedule{}, map[string]schedule{}, true, 0}
}

func (d *draftman) onDayEnd(today string) {
	if v, exist := d.agenda.data[today]; exist {
		d.history[today] = v
		delete(d.agenda.data, today)
	}
}

func (d *draftman) setScheduleOver(date string, id string) error {
	if v, exist := d.agenda.data[date]; exist {
		if s, exist := v.confirmedSchedules[id]; exist {
			s.done = true
			return nil
		}
		return errors.New("could not find schedule " + id + " on day " + date)
	}
	return errors.New("could not find any schedule on day " + date)
}

func (d *draftman) changeScheduleDay(sDate string, sId string, targetDate string) error {

	if v, exist := d.agenda.data[sDate]; exist {
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
		d.agenda.data[sDate] = v
		return nil
	}
	return errors.New("could not find any schedule on day " + sDate)
}

func (a *agenda) checkDate(date string) error {
	if v, exist := a.data[date]; exist {
		if v.isDayFull() {
			return errors.New("day is full")
		}
		return nil
	}
	return nil
}
