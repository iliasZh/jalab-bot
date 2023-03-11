package db

import (
	"time"
)

type GetJalabStatsQuery struct {
	GroupChatID int64     `db:"group_chat_id"`
	FromDate    time.Time `db:"from_date"`
	ToDate      time.Time `db:"to_date"`
}

func NewGetJalabOfTheMonthQuery(groupChatID int64) GetJalabStatsQuery {
	now := time.Now().UTC()

	return GetJalabStatsQuery{
		GroupChatID: groupChatID,
		FromDate:    firstDayOfMonth(now, currentMonth),
		ToDate:      firstDayOfMonth(now, nextMonth).Add(-24 * time.Hour),
	}
}
func firstDayOfMonth(d time.Time, advance func(time.Time) (int, time.Month)) time.Time {
	y, m := advance(d)
	return time.Date(y, m, 1, 0, 0, 0, 0, time.UTC)
}

func currentMonth(d time.Time) (int, time.Month) {
	return d.Year(), d.Month()
}

func nextMonth(d time.Time) (int, time.Month) {
	const months = 1

	zeroIndexedMonths := d.Month() + months - 1

	addedYears := int(zeroIndexedMonths / 12)
	zeroIndexedMonth := zeroIndexedMonths % 12

	return d.Year() + addedYears, zeroIndexedMonth + 1
}
