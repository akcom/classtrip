package classtrip

import (
	"time"
	"fmt"
	"strconv"
)

func genRow(values []string) (result string) {
	result = "<tr>"
	for i := 0; i < 7; i++ {
		//values starting with an asterisk are emphasized
		if len(values[i])>1 && values[i][0] == '*' {
			result += `<td class="success">` + values[i][1:] + `</td>`
		} else {
			result += "<td>" + values[i] + "</td>"
		}
	}
	result += "</tr>\n"
	return
}

func formatDayNum(dayNum int) (result string) {
	strDay := strconv.Itoa(dayNum)
	result =  `<span class="tbl-num">` + strDay + `</span>`
	return
}

func genHeader(currentDate time.Time) (result string) {
	var (
		nextMonth time.Month
		nextYear int
		prevMonth time.Month
		prevYear int
		currMonth = currentDate.Month()
		currYear = currentDate.Year()
	)
	if currMonth == time.December {
		nextMonth = time.January
		nextYear = currYear + 1 
	} else {
		nextMonth = time.Month(int(currMonth)+1)
		nextYear = currYear
	}
	
	if currMonth == time.January {
		prevMonth = time.December
		prevYear = currYear-1
	} else {
		prevMonth = time.Month(int(currMonth)-1)
		prevYear = currYear
	}
	result = fmt.Sprintf(
		`<h2><a href="/calendar/%d/%d"> < </a> %v %v <a href="/calendar/%d/%d"> > </a></h2>`,
		int(prevMonth), prevYear, currMonth.String(), currYear, int(nextMonth), nextYear)
	return
}


func GenCalendar(month time.Month, year int) (result string) {
	var (
		daysOfWeek = []string{ "Sun", "Mon", "Tue","Wed","Thu","Fri","Sat"}
		days []string = make([]string, 7)
	)
	t := time.Now()
	//create a new Time object on 1st of month, year
	currentDay := time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
	//figure out what day of the week the first falls on
	firstWeekDay := int(currentDay.Weekday())

	result += genHeader(currentDay)
	result += `<div style="min-height: 800px"><table class="table table-bordered" style="height: 80%">`
	//set the first row 
	//each time has a day of th week and a day (#)
	//ex: Tues 1, Wed 2, Thu 3, etc
	for i := 0;i < 7; i++ {
		if i >= firstWeekDay {
			daysOfWeek[i] += " " + formatDayNum(currentDay.Day())
			currentDay = currentDay.Add(24 * time.Hour)
		}
	}
	result += genRow(daysOfWeek)
	//iterate until we've finished the current month
	for currentDay.Month() == month {
		days = []string{"","","","","","",""}
		//increment the day until we hit the end of the month
		for i := 0;i < 7 && currentDay.Month() == month; i++ {
			days[i] += formatDayNum(currentDay.Day())
			currentDay = currentDay.Add(24 * time.Hour)
		}
		result += genRow(days)
	}
	result += `</table></div>`
	return result
}