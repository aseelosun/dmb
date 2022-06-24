package db

import (
	"database/sql"
	entity "doMassageBot/internal/entity"
	freeTime "doMassageBot/internal/time"
	"log"
	"time"
)

func GetUserStatus(db *sql.DB, id int) (int, error) {
	var (
		status int
	)
	sqlStatement := `SELECT status FROM "doMassageBot".users WHERE userid = $1;`
	err := db.QueryRow(sqlStatement, id).Scan(&status)
	if err != nil {
		return status, err
	}
	return status, nil
}
func UpdateUsername(db *sql.DB, id int, username string) error {
	sqlStatement := `UPDATE "doMassageBot".users SET username= $2 WHERE userid = $1;`
	_, err := db.Exec(sqlStatement, id, username)
	if err != nil {
		return err
	}
	return nil
}
func UpdateFullname(db *sql.DB, id int, fullname string) error {
	sqlStatement := `UPDATE "doMassageBot".users SET fullname= $2 WHERE userid = $1;`
	_, err := db.Exec(sqlStatement, id, fullname)
	if err != nil {
		return err
	}
	return nil
}
func UpdateEmail(db *sql.DB, id int, email string) error {
	sqlStatement := `UPDATE "doMassageBot".users SET email= $2 WHERE userid = $1;`
	_, err := db.Exec(sqlStatement, id, email)
	if err != nil {
		return err
	}
	return nil
}
func UpdatePhoneNum(db *sql.DB, id int, phoneNum string) error {
	sqlStatement := `UPDATE "doMassageBot".users SET phoneNumber = $2 WHERE userid = $1;`
	_, err := db.Exec(sqlStatement, id, phoneNum)
	if err != nil {
		return err
	}
	return nil
}

func UpdateUserStatus(db *sql.DB, id int, status int) error {
	sqlStatement := `UPDATE "doMassageBot".users SET status = $2 WHERE userid = $1;`
	_, err := db.Exec(sqlStatement, id, status)
	if err != nil {
		return err
	}
	return nil
}

func UpdateUserIdStatus(db *sql.DB, id int, status int) error {
	sqlStatement := `UPDATE "doMassageBot".users SET status = $2 WHERE userId = $1;`
	_, err := db.Exec(sqlStatement, id, status)
	if err != nil {
		return err
	}
	return nil
}

func CheckIfUserExists(db *sql.DB, userId int) (bool, error) {
	sqlStmt := `SELECT userId FROM "doMassageBot".users WHERE userId = $1 AND status = 4;`
	err := db.QueryRow(sqlStmt, userId).Scan(&userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func InsertIntoUsers(db *sql.DB, userId int, fullname string, username string, email string, phoneNumber string, status int) (int, error) {
	var id int
	sqlStmt := `INSERT INTO "doMassageBot".users (userId, fullName, username, email,phoneNumber, status) VALUES ($1, $2, $3, $4, $5, $6) returning id;`
	err := db.QueryRow(sqlStmt, userId, fullname, username, email, phoneNumber, status).Scan(&id)
	if err != nil {
		return id, err
	}
	return id, nil
}

func GetMyProfile(db *sql.DB, userId int) (entity.MyProfile, error) {
	var (
		obj      entity.MyProfile
		name     string
		email    string
		phoneNum string
	)
	rows, err := db.Query(`select fullname, email, phoneNumber from "doMassageBot".users where userId = $1;`, userId)
	if err != nil {
		return obj, err
	}

	for rows.Next() {
		err := rows.Scan(&name, &email, &phoneNum)
		if err != nil {
			return obj, err
		}
		obj.Name = name
		obj.Email = email
		obj.PhoneNum = phoneNum

	}
	return obj, nil

}

func UpdateScheduleStatus(db *sql.DB, id int, status int) error {
	sqlStatement := `UPDATE "doMassageBot".massageSchedule SET status = $2 WHERE id = $1;`
	_, err := db.Exec(sqlStatement, id, status)
	if err != nil {
		return err
	}
	return nil
}
func UpdateScheduleTime(db *sql.DB, id int, mtime string) error {
	sqlStatement := `UPDATE "doMassageBot".massageSchedule SET mTime= $2 WHERE id = $1;`
	_, err := db.Exec(sqlStatement, id, mtime)
	if err != nil {
		return err
	}
	return nil
}

func RefreshUserList(db *sql.DB) error {
	sqlStatement := `DELETE FROM "doMassageBot".users u WHERE u.id <> (SELECT min(v.id) FROM "doMassageBot".users v WHERE  u.userId = v.userId);`
	_, err := db.Exec(sqlStatement)
	if err != nil {
		return err
	}
	return nil
}

func RefreshMassageSchedule(db *sql.DB) error {
	sqlStatement := `DELETE FROM "doMassageBot".massageSchedule WHERE status != 2;`
	_, err := db.Exec(sqlStatement)
	if err != nil {
		return err
	}
	return nil
}
func GetCurrentSchedule(db *sql.DB, mId, status int) (entity.MySchedule, error) {
	var (
		mDate time.Time
		mTime time.Time
		mType string
		obj   entity.MySchedule
	)
	rows, err := db.Query(`Select m.mDate, m.mTime, mt.mType from "doMassageBot".massageSchedule as m join "doMassageBot".massageType as mt on m.mid = mt.id join "doMassageBot".users u on u.id = m.uId where m.status = $2 and m.isCanceled <> true and m.id = $1 and m.mdate >= CURRENT_DATE  order by m.mtime`, mId, status)
	if err != nil {
		return obj, err
	}
	for rows.Next() {
		err := rows.Scan(&mDate, &mTime, &mType)
		if err != nil {
			return obj, err
		}
		obj.MType = mType
		obj.MDate = mDate.Format("2006-01-02") //mDate[:len(mDate)-10]
		obj.MTime = mTime.Format("15:04")      //mTime[11 : len(mTime)-4]

	}
	return obj, nil

}

func GetCanceledSchedule(db *sql.DB, mId int) (entity.MySchedule, error) {
	var (
		mDate time.Time
		mTime time.Time
		mType string
		obj   entity.MySchedule
	)
	rows, err := db.Query(`Select m.mDate, m.mTime, mt.mType from "doMassageBot".massageSchedule as m join "doMassageBot".massageType as mt on m.mid = mt.id join "doMassageBot".users u on u.id = m.uId where m.status = 2 and m.isCanceled <> false and m.id = $1 and m.mdate >= CURRENT_DATE order by m.mtime`, mId)
	if err != nil {
		return obj, err
	}
	for rows.Next() {
		err := rows.Scan(&mDate, &mTime, &mType)
		if err != nil {
			return obj, err
		}
		obj.MType = mType
		obj.MDate = mDate.Format("2006-01-02")
		obj.MTime = mTime.Format("15:04")

	}
	return obj, nil

}

func GetMySchedule(db *sql.DB, userId int) ([]entity.MySchedule, error) {
	var (
		obj      entity.MySchedule
		mDate    time.Time
		mTime    time.Time
		mType    string
		id       int
		objArray []entity.MySchedule
	)
	rows, err := db.Query(`Select m.id, m.mDate, m.mTime, mt.mType from "doMassageBot".massageSchedule as m join "doMassageBot".massageType as mt on m.mid = mt.id join "doMassageBot".users u on u.id = m.uId where u.userId = $1 and m.status = 2 and m.isCanceled <> true and m.mdate >= CURRENT_DATE  order by m.mtime`, userId)

	if err != nil {
		if err != sql.ErrNoRows {
			return objArray, nil
		}
		return objArray, err
	}

	for rows.Next() {
		err := rows.Scan(&id, &mDate, &mTime, &mType)
		if err != nil {
			if err != sql.ErrNoRows {
				return objArray, nil
			}
			return objArray, err
		}
		obj.Id = id
		obj.MType = mType
		obj.MDate = mDate.Format("2006-01-02")
		obj.MTime = mTime.Format("15:04")

		objArray = append(objArray, obj)

	}
	return objArray, nil

}
func GetAllScheduleForToday(db *sql.DB) ([]entity.AllSchedule, error) {
	var (
		obj      entity.AllSchedule
		fullname string
		email    string
		phoneNum string
		mDate    time.Time
		mTime    time.Time
		mType    string
		objArray []entity.AllSchedule
	)
	rows, err := db.Query(`select u.fullname, u.email, u.phoneNumber,m.mtype,ms.mdate,ms.mtime from "doMassageBot".massageSchedule as ms join "doMassageBot".massageType as m on ms.mid = m.id join "doMassageBot".users u on ms.uId = u.id where ms.mDate = CURRENT_DATE and ms.iscanceled = false order by ms.mtime;`)

	if err != nil {
		return objArray, err
	}

	for rows.Next() {
		err := rows.Scan(&fullname, &email, &phoneNum, &mType, &mDate, &mTime)
		if err != nil {
			return objArray, err
		}
		obj.Name = fullname
		obj.Email = email
		obj.PhoneNum = phoneNum
		obj.MType = mType
		obj.MDate = mDate.Format("2006-01-02")
		obj.MTime = mTime.Format("15:04")
		objArray = append(objArray, obj)

	}
	return objArray, nil

}

func GetAllScheduleForTomorrow(db *sql.DB, Date string) ([]entity.AllSchedule, error) {
	var (
		obj      entity.AllSchedule
		fullname string
		email    string
		phoneNum string
		mDate    time.Time
		mTime    time.Time
		mType    string
		objArray []entity.AllSchedule
	)

	if Date == "Понедельник" {
		rows, err := db.Query(`select u.fullname, u.email,u.phoneNumber,m.mtype,ms.mdate,ms.mtime from "doMassageBot".massageSchedule as ms join "doMassageBot".massageType as m on ms.mid = m.id join "doMassageBot".users u on ms.uId = u.id where ms.mDate = CURRENT_DATE + interval '3 day' and ms.iscanceled = false order by ms.mtime;`)

		if err != nil {
			return objArray, err
		}

		for rows.Next() {
			err := rows.Scan(&fullname, &email, &phoneNum, &mType, &mDate, &mTime)
			if err != nil {
				return objArray, err
			}
			obj.Name = fullname
			obj.Email = email
			obj.PhoneNum = phoneNum
			obj.MType = mType
			obj.MDate = mDate.Format("2006-01-02")
			obj.MTime = mTime.Format("15:04")

			objArray = append(objArray, obj)

		}
	} else {
		rows, err := db.Query(`select u.fullname, u.email,u.phoneNumber,m.mtype,ms.mdate,ms.mtime from "doMassageBot".massageSchedule as ms join "doMassageBot".massageType as m on ms.mid = m.id join "doMassageBot".users u on ms.uId = u.id where ms.mDate = CURRENT_DATE + interval '1 day' and ms.iscanceled = false order by ms.mtime;`)

		if err != nil {
			return objArray, err
		}

		for rows.Next() {
			err := rows.Scan(&fullname, &email, &phoneNum, &mType, &mDate, &mTime)
			if err != nil {
				return objArray, err
			}
			obj.Name = fullname
			obj.Email = email
			obj.PhoneNum = phoneNum
			obj.MType = mType
			obj.MDate = mDate.Format("2006-01-02")
			obj.MTime = mTime.Format("15:04")

			objArray = append(objArray, obj)

		}
	}
	//rows, err := db.Query(`select u.fullname, u.email,u.phoneNumber,m.mtype,ms.mdate,ms.mtime from "doMassageBot".massageSchedule as ms join "doMassageBot".massageType as m on ms.mid = m.id join "doMassageBot".users u on ms.uId = u.id where ms.mDate = CURRENT_DATE + interval '1 day' and ms.iscanceled = false order by ms.mtime;`)

	//if err != nil {
	//	return objArray, err
	//}
	//
	//for rows.Next() {
	//	err := rows.Scan(&fullname, &email, &phoneNum, &mType, &mDate, &mTime)
	//	if err != nil {
	//		return objArray, err
	//	}
	//	obj.Name = fullname
	//	obj.Email = email
	//	obj.PhoneNum = phoneNum
	//	obj.MType = mType
	//	obj.MDate = mDate[:len(mDate)-10]
	//	obj.MTime = mTime[11 : len(mTime)-4]
	//
	//	objArray = append(objArray, obj)
	//
	//}
	return objArray, nil

}

func GetAllScheduleForMonday(db *sql.DB) ([]entity.AllSchedule, error) {
	var (
		obj      entity.AllSchedule
		fullname string
		email    string
		phoneNum string
		mDate    time.Time
		mTime    time.Time
		mType    string
		objArray []entity.AllSchedule
	)
	rows, err := db.Query(`select u.fullname, u.email, u.phoneNumber,m.mtype,ms.mdate,ms.mtime from "doMassageBot".massageSchedule as ms join "doMassageBot".massageType as m on ms.mid = m.id join "doMassageBot".users u on ms.uId = u.id where ms.mDate = CURRENT_DATE + interval '3 day' and ms.iscanceled = false order by ms.mtime;`)

	if err != nil {
		return objArray, err
	}

	for rows.Next() {
		err := rows.Scan(&fullname, &email, &phoneNum, &mType, &mDate, &mTime)
		if err != nil {
			return objArray, err
		}
		obj.Name = fullname
		obj.Email = email
		obj.PhoneNum = phoneNum
		obj.MType = mType
		obj.MDate = mDate.Format("2006-01-02")
		obj.MTime = mTime.Format("15:04")
		objArray = append(objArray, obj)

	}
	return objArray, nil

}

func GenerateTimeCollarToday(db *sql.DB) ([]string, error) {
	var (
		hoursArray []string
	)

	freeTimes1 := freeTime.GetFreeTime(freeTime.NewFreeTime(9, 0), freeTime.NewFreeTime(12, 30), freeTime.NewFreeTime(time.Now().Hour(), time.Now().Minute()), 30)
	freeTimes2 := freeTime.GetFreeTime(freeTime.NewFreeTime(14, 0), freeTime.NewFreeTime(17, 0), freeTime.NewFreeTime(time.Now().Hour(), time.Now().Minute()), 30)
	takenTime, err := GetTakenTimeCollarToday(db)
	for _, hours := range freeTimes1 {
		if err != nil {
			return hoursArray, err
		}
		if !stringInSlice(hours.ToString(), takenTime) {
			hoursArray = append(hoursArray, hours.ToString())
		}
	}
	for _, hours := range freeTimes2 {
		if err != nil {
			return hoursArray, err
		}
		if !stringInSlice(hours.ToString(), takenTime) {
			hoursArray = append(hoursArray, hours.ToString())
		}
	}

	return hoursArray, nil
}
func GenerateTimeCollarMonday(db *sql.DB) ([]string, error) {
	var (
		hoursArray []string
	)

	freeTimes1 := freeTime.GetFreeTime(freeTime.NewFreeTime(9, 0), freeTime.NewFreeTime(12, 30), freeTime.NewFreeTime(6, 0), 30)
	freeTimes2 := freeTime.GetFreeTime(freeTime.NewFreeTime(14, 0), freeTime.NewFreeTime(17, 0), freeTime.NewFreeTime(6, 0), 30)

	takenTime, err := GetTakenTimeCollarMonday(db)

	for i, _ := range freeTimes1 {
		if err != nil {
			return hoursArray, err
		}
		if !stringInSlice(freeTimes1[i].ToString(), takenTime) {
			hoursArray = append(hoursArray, freeTimes1[i].ToString())
		}
	}

	for _, hours := range freeTimes2 {
		if err != nil {
			return hoursArray, err
		}
		if !stringInSlice(hours.ToString(), takenTime) {
			hoursArray = append(hoursArray, hours.ToString())
		}
	}

	return hoursArray, nil
}

func GenerateTimeCollarTomorrow(db *sql.DB) ([]string, error) {
	var (
		hoursArray []string
	)

	freeTimes1 := freeTime.GetFreeTime(freeTime.NewFreeTime(9, 0), freeTime.NewFreeTime(12, 30), freeTime.NewFreeTime(6, 0), 30)
	freeTimes2 := freeTime.GetFreeTime(freeTime.NewFreeTime(14, 0), freeTime.NewFreeTime(17, 0), freeTime.NewFreeTime(6, 0), 30)
	takenTime, err := GetTakenTimeCollarTomorrow(db)

	for _, hours := range freeTimes1 {
		if err != nil {
			return hoursArray, err
		}
		if !stringInSlice(hours.ToString(), takenTime) {
			hoursArray = append(hoursArray, hours.ToString())
		}
	}
	for _, hours := range freeTimes2 {
		if err != nil {
			return hoursArray, err
		}
		if !stringInSlice(hours.ToString(), takenTime) {
			hoursArray = append(hoursArray, hours.ToString())
		}
	}

	return hoursArray, nil
}
func GenerateTimeMedicalToday(db *sql.DB) ([]string, error) {
	var (
		hoursArray []string
	)

	freeTimes1 := freeTime.GetFreeTime(freeTime.NewFreeTime(8, 0), freeTime.NewFreeTime(8, 30), freeTime.NewFreeTime(time.Now().Hour(), time.Now().Minute()), 30)
	freeTimes2 := freeTime.GetFreeTime(freeTime.NewFreeTime(13, 0), freeTime.NewFreeTime(13, 30), freeTime.NewFreeTime(time.Now().Hour(), time.Now().Minute()), 30)
	freeTimes3 := freeTime.GetFreeTime(freeTime.NewFreeTime(18, 0), freeTime.NewFreeTime(18, 0), freeTime.NewFreeTime(time.Now().Hour(), time.Now().Minute()), 30)

	takenTime, err := GetTakenTimeMedicalToday(db)
	for _, hours := range freeTimes1 {
		if err != nil {
			return hoursArray, err
		}
		if !stringInSlice(hours.ToString(), takenTime) {
			hoursArray = append(hoursArray, hours.ToString())
		}
	}

	for _, hours := range freeTimes2 {
		if err != nil {
			return hoursArray, err
		}
		if !stringInSlice(hours.ToString(), takenTime) {
			hoursArray = append(hoursArray, hours.ToString())
		}
	}
	for _, hours := range freeTimes3 {
		if err != nil {
			return hoursArray, err
		}
		if !stringInSlice(hours.ToString(), takenTime) {
			hoursArray = append(hoursArray, hours.ToString())
		}
	}

	return hoursArray, nil

}
func GenerateTimeMedicalMonday(db *sql.DB) ([]string, error) {
	var (
		hoursArray []string
	)

	freeTimes1 := freeTime.GetFreeTime(freeTime.NewFreeTime(8, 0), freeTime.NewFreeTime(8, 30), freeTime.NewFreeTime(6, 0), 30)
	freeTimes2 := freeTime.GetFreeTime(freeTime.NewFreeTime(13, 0), freeTime.NewFreeTime(13, 30), freeTime.NewFreeTime(6, 0), 30)
	freeTimes3 := freeTime.GetFreeTime(freeTime.NewFreeTime(18, 0), freeTime.NewFreeTime(18, 0), freeTime.NewFreeTime(6, 0), 30)

	for _, hours := range freeTimes1 {
		takenTime, err := GetTakenTimeMedicalMonday(db)
		if err != nil {
			return hoursArray, err
		}
		if !stringInSlice(hours.ToString(), takenTime) {
			hoursArray = append(hoursArray, hours.ToString())
		}
	}

	for _, hours := range freeTimes2 {
		takenTime, err := GetTakenTimeMedicalMonday(db)
		if err != nil {
			return hoursArray, err
		}
		if !stringInSlice(hours.ToString(), takenTime) {
			hoursArray = append(hoursArray, hours.ToString())
		}
	}
	for _, hours := range freeTimes3 {
		takenTime, err := GetTakenTimeMedicalMonday(db)
		if err != nil {
			return hoursArray, err
		}
		if !stringInSlice(hours.ToString(), takenTime) {
			hoursArray = append(hoursArray, hours.ToString())
		}
	}

	return hoursArray, nil

}

func GenerateTimeMedicalTomorrow(db *sql.DB) ([]string, error) {
	var (
		hoursArray []string
	)
	freeTimes1 := freeTime.GetFreeTime(freeTime.NewFreeTime(8, 0), freeTime.NewFreeTime(8, 30), freeTime.NewFreeTime(6, 0), 30)
	freeTimes2 := freeTime.GetFreeTime(freeTime.NewFreeTime(13, 0), freeTime.NewFreeTime(13, 30), freeTime.NewFreeTime(6, 0), 30)
	freeTimes3 := freeTime.GetFreeTime(freeTime.NewFreeTime(18, 0), freeTime.NewFreeTime(18, 0), freeTime.NewFreeTime(6, 0), 30)

	for _, hours := range freeTimes1 {
		println("This:", hours.ToString())
		takenTime, err := GetTakenTimeMedicalTomorrow(db)
		if err != nil {
			return hoursArray, err
		}
		if !stringInSlice(hours.ToString(), takenTime) {
			hoursArray = append(hoursArray, hours.ToString())
		}
	}

	for _, hours := range freeTimes2 {
		println("This:", hours.ToString())
		takenTime, err := GetTakenTimeMedicalTomorrow(db)
		if err != nil {
			return hoursArray, err
		}
		if !stringInSlice(hours.ToString(), takenTime) {
			hoursArray = append(hoursArray, hours.ToString())
		}
	}
	for _, hours := range freeTimes3 {
		println("This:", hours.ToString())
		takenTime, err := GetTakenTimeMedicalTomorrow(db)
		if err != nil {
			return hoursArray, err
		}
		if !stringInSlice(hours.ToString(), takenTime) {
			hoursArray = append(hoursArray, hours.ToString())
		}
	}

	return hoursArray, nil

}

func GetTakenTimeCollarToday(db *sql.DB) ([]string, error) {
	var (
		hours      time.Time
		hoursArray []string
	)
	rows, err := db.Query(`select ms.mtime from "doMassageBot".massageSchedule as ms join "doMassageBot".massageType as m on ms.mid = m.id join "doMassageBot".users u on ms.uId = u.id  where ms.status = 2 and ms.mDate = CURRENT_DATE and m.id = 1 and ms.isCanceled = false order by ms.mtime;`)

	if err != nil {
		return hoursArray, err
	}
	for rows.Next() {
		err := rows.Scan(&hours)
		if err != nil {
			log.Fatal("Failed to execute query: ", err)
		}
		hoursArray = append(hoursArray, hours.Format("15:04"))

	}
	return hoursArray, nil
}
func GetTakenTimeCollarMonday(db *sql.DB) ([]string, error) {
	var (
		hours      time.Time
		hoursArray []string
	)
	rows, err := db.Query(`select ms.mtime from "doMassageBot".massageSchedule as ms join "doMassageBot".massageType as m on ms.mid = m.id join "doMassageBot".users u on ms.uId = u.id  where ms.status = 2 and ms.mDate = CURRENT_DATE + interval '3 day' and m.id = 1 and ms.isCanceled = false order by ms.mtime;`)

	if err != nil {
		return hoursArray, err
	}
	for rows.Next() {
		err := rows.Scan(&hours)
		if err != nil {
			log.Fatal("Failed to execute query: ", err)
		}
		hoursArray = append(hoursArray, hours.Format("15:04"))

	}
	return hoursArray, nil
}
func GetTakenTimeCollarTomorrow(db *sql.DB) ([]string, error) {
	var (
		hours      time.Time
		hoursArray []string
	)
	rows, err := db.Query(`select ms.mtime from "doMassageBot".massageSchedule as ms join "doMassageBot".massageType as m on ms.mid = m.id join "doMassageBot".users u on ms.uId = u.id  where ms.status = 2 and ms.mDate = CURRENT_DATE  + interval '1 day' and m.id = 1 and ms.isCanceled = false order by ms.mtime;`)

	if err != nil {
		return hoursArray, err
	}
	for rows.Next() {
		err := rows.Scan(&hours)
		if err != nil {
			return hoursArray, err
		}
		hoursArray = append(hoursArray, hours.Format("15:04"))

	}
	return hoursArray, nil

}

func GetTakenTimeMedicalToday(db *sql.DB) ([]string, error) {
	var (
		hours      time.Time
		hoursArray []string
	)
	rows, err := db.Query(`select ms.mtime from "doMassageBot".massageSchedule as ms join "doMassageBot".massageType as m on ms.mid = m.id join "doMassageBot".users u on ms.uId = u.id  where ms.status = 2 and ms.mDate = CURRENT_DATE and m.id = 2 and ms.isCanceled = false order by ms.mtime`)

	if err != nil {
		return hoursArray, err
	}
	for rows.Next() {
		err := rows.Scan(&hours)
		if err != nil {
			log.Fatal("Failed to execute query: ", err)
		}
		hoursArray = append(hoursArray, hours.Format("15:04"))

	}
	return hoursArray, nil
}

func GetTakenTimeMedicalMonday(db *sql.DB) ([]string, error) {
	var (
		hours      time.Time
		hoursArray []string
	)
	rows, err := db.Query(`select ms.mtime from "doMassageBot".massageSchedule as ms join "doMassageBot".massageType as m on ms.mid = m.id join "doMassageBot".users u on ms.uId = u.id  where ms.status = 2 and ms.mDate = CURRENT_DATE + interval '3 day' and m.id = 2 and ms.isCanceled = false order by ms.mtime`)

	if err != nil {
		return hoursArray, err
	}
	for rows.Next() {
		err := rows.Scan(&hours)
		if err != nil {
			log.Fatal("Failed to execute query: ", err)
		}
		hoursArray = append(hoursArray, hours.Format("15:04"))

	}
	return hoursArray, nil
}

func GetTakenTimeMedicalTomorrow(db *sql.DB) ([]string, error) {
	var (
		hours      time.Time
		hoursArray []string
	)
	rows, err := db.Query(`select ms.mtime from "doMassageBot".massageSchedule as ms join "doMassageBot".massageType as m on ms.mid = m.id join "doMassageBot".users u on ms.uId = u.id  where ms.status = 2 and ms.mDate = CURRENT_DATE  + interval '1 day' and m.id = 2 and ms.isCanceled = false order by ms.mtime;`)

	if err != nil {
		return hoursArray, err
	}
	for rows.Next() {
		err := rows.Scan(&hours)
		if err != nil {
			log.Fatal("Failed to execute query: ", err)
		}
		hoursArray = append(hoursArray, hours.Format("15:04"))

	}
	return hoursArray, nil
}
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

//
//func GetScheduleDay(db *sql.DB, id int) (bool, error) {
//	var (
//		status bool
//	)
//	sqlStatement := `SELECT (CASE WHEN mDate = CURRENT_DATE THEN 1 ELSE 0 END) FROM "doMassageBot".massageSchedule WHERE id = $1;`
//	err := db.QueryRow(sqlStatement, id).Scan(&status)
//	if err != nil {
//		return status, err
//	}
//	return status, nil
//}

func GetScheduleStatus(db *sql.DB, id int) (int, error) {
	var (
		status int
	)
	sqlStatement := `SELECT status FROM "doMassageBot".massageSchedule WHERE id = $1;`
	err := db.QueryRow(sqlStatement, id).Scan(&status)
	if err != nil {
		return status, err
	}
	return status, nil
}

func InsertIntoSchedule(db *sql.DB, mType string, mDate string, mTime string, userId int, status int) (int, error) {
	var id int
	if mDate == "Сегодня" {
		sqlStmt := `INSERT INTO "doMassageBot".massageSchedule(mid, mDate, mTime, uId,status) VALUES((SELECT id from "doMassageBot".massageType WHERE mType=$1),CURRENT_DATE,$2, (SELECT id from "doMassageBot".users WHERE userId= $3),$4) returning id;`
		err := db.QueryRow(sqlStmt, mType, mTime, userId, status).Scan(&id)
		if err != nil {
			return id, err
		}
	} else if mDate == "Понедельник" {
		sqlStmt := `INSERT INTO "doMassageBot".massageSchedule(mid, mDate, mTime, uId,status) VALUES((SELECT id from "doMassageBot".massageType WHERE mType=$1),CURRENT_DATE + interval '3 day',$2, (SELECT id from "doMassageBot".users WHERE userId= $3),$4) returning id;`
		err := db.QueryRow(sqlStmt, mType, mTime, userId, status).Scan(&id)
		if err != nil {
			return id, err
		}
	} else {
		sqlStmt := `INSERT INTO "doMassageBot".massageSchedule(mid, mDate, mTime, uId,status) VALUES((SELECT id from "doMassageBot".massageType WHERE mType=$1),CURRENT_DATE + interval '1 day',$2, (SELECT id from "doMassageBot".users WHERE userId= $3),$4) returning id;`
		err := db.QueryRow(sqlStmt, mType, mTime, userId, status).Scan(&id)
		if err != nil {
			return id, err
		}
	}
	return id, nil
}

func UpdateScheduleMType(db *sql.DB, mType string, id int, status int) error {
	sqlStatement := `UPDATE "doMassageBot".massageSchedule SET mId = (SELECT id from "doMassageBot".massageType WHERE mType= $1), status = $2 WHERE id = $3;`
	_, err := db.Exec(sqlStatement, mType, status, id)
	if err != nil {
		return err
	}
	return nil
}

func CancelEntry(db *sql.DB, id int) (bool, error) {
	sqlStatement := `UPDATE "doMassageBot".massageSchedule SET isCanceled = true WHERE id = $1;`
	_, err := db.Exec(sqlStatement, id)
	if err != nil {
		return false, err

	}
	return true, nil
}
