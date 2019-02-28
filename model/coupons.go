package model

import (
	"fmt"
	"log"
	"strconv"
	"sync"
)

func CreateCouponSchema(event string, coupons []Coupon, c chan error) {

	// check if coupon exists
	data, _, _, err := con.QueryNeoAll(`
		MATCH (n:EVENT)-[r:COUPON]->(:COUPON_SCHEMA)
		WHERE n.name = $event
		RETURN r.coupons
	`, map[string]interface{}{
		"event": event,
	})
	if err != nil {
		c <- err
		return
	}

	str := fmt.Sprintf("%v", data)
	if str == "[[]]" && str == "[]" {
		c <- fmt.Errorf("already exists")
		return
	}

	// create schema; TODO error handling
	mutex := &sync.Mutex{}
	for _, cps := range coupons {

		go func(cp Coupon, mu *sync.Mutex) {

			mu.Lock()
			rs, err := con.ExecNeo(`
				MATCH(n:EVENT) WHERE n.name = $event
				CREATE (:COUPON_SCHEMA {name:$name, description:$desc, day:$day})<-[:COUPON]-(n)
			`, map[string]interface{}{
				"event": event,
				"name":  cp.Name,
				"desc":  cp.Desc,
				"day":   cp.Day,
			})
			mu.Unlock()

			if err != nil {
				log.Println(err)
				c <- err
				return
			}
			log.Println(rs)
		}(cps, mutex)

	}

	c <- nil
	return

}

func MarkPresent(attendance Attendance, c chan MessageReturn) {

	// check if user exists or not
	data, _, _, err := con.QueryNeoAll(`
		MATCH(n:EVENT)-[:ATTENDS]->(b)
		WHERE n.name=$name AND b.email=$rn
		RETURN b.email
	`, map[string]interface{}{
		"name": attendance.EventName,
		"rn":   attendance.Email,
	})
	if err != nil {
		c <- MessageReturn{"Error marking attendance", err}
		return
	}
	if len(data) < 1 {
		c <- MessageReturn{"No participant found", nil}
		return
	}

	// check if already given attendance
	data, _, _, err = con.QueryNeoAll(`
		MATCH(n:EVENT)-[r:PRESENT`+strconv.Itoa(attendance.Day)+`]->(b)
		WHERE n.name=$name AND b.email=$rn
		RETURN b.email
	`, map[string]interface{}{
		"name": attendance.EventName,
		"rn":   attendance.Email,
	})
	if err != nil {
		c <- MessageReturn{"Error marking attendance", err}
		return
	}
	if len(data) > 0 {
		c <- MessageReturn{"Already marked present", nil}
		return
	}

	// check if schema exists
	data, _, _, err = con.QueryNeoAll(
		`MATCH (n:EVENT)-[r:COUPON]->(a:COUPON_SCHEMA)
		 WHERE n.name = $event
		 RETURN a.name, a.description, a.day
	`, map[string]interface{}{
			"event": attendance.EventName,
		})
	if err != nil {
		c <- MessageReturn{"Error occurred while checking if coupon schema exists", err}
		return
	}

	str := fmt.Sprintf("%v", data)

	if str == "[[]]" && str == "[]" {
		// mark present if schema does not exist
		_, err = con.ExecNeo(`
			MATCH(n:EVENT)-[:ATTENDS]->(b)
			WHERE n.name=$name AND b.email=$rn
			CREATE (n)-[:PRESENT`+strconv.Itoa(attendance.Day)+`]->(b) 
		`, map[string]interface{}{
			"name": attendance.EventName,
			"rn":   attendance.Email,
		})
		if err != nil {
			c <- MessageReturn{"Error creating present relation", err}
			return
		}
		c <- MessageReturn{"Done", nil}
		return
	}

	// if schema exists,generate hash, and save for each coupon
	var couponArr []Coupon

	for _, o := range data {
		if o[2].(int64) == int64(attendance.Day) {
			couponArr = append(couponArr, Coupon{
				Name: o[0].(string),
				Desc: o[1].(string),
				Day:  int(o[2].(int64)),
			})
		}
	}
	go couponGen(couponArr, attendance, c)

}

// generate coupons and save
func couponGen(coupons []Coupon, attendance Attendance, c chan MessageReturn) {

	var (
		str string
		cps []string
	)
	//SALT, _ := strconv.Atoi(os.Getenv("SALT"))

	for _, coupon := range coupons {
		str = attendance.EventName + strconv.Itoa(attendance.Day) + coupon.Name + attendance.Email
		// bytes, err := bcrypt.GenerateFromPassword([]byte(str), SALT)
		// fmt.Println(str)
		// if err != nil {
		// 	c <- MessageReturn{"Error while hashing", err}
		// 	return
		// }
		cps = append(cps, str)
	}
	// create coupon relation
	_, err := con.ExecNeo(`
			MATCH(n:EVENT)-[:ATTENDS]->(b)
			WHERE n.name=$name AND b.email=$rn
			CREATE (n)-[:PRESENT`+strconv.Itoa(attendance.Day)+`{coupons:$cps}]->(b) 
		`, map[string]interface{}{
		"name": attendance.EventName,
		"rn":   attendance.Email,
		"cps":  cps,
	})
	if err != nil {
		c <- MessageReturn{"Error creating coupon relation", err}
		return
	}
	c <- MessageReturn{"Successfully marked present for the day", nil}
	return
}

// redeem a coupon

func RedeemCoupon(attendance Attendance, couponName string, c chan MessageReturn) {

	// build coupon
	//SALT, _ := strconv.Atoi(os.Getenv("SALT"))
	str := attendance.EventName + strconv.Itoa(attendance.Day) + couponName + attendance.Email

	fmt.Println(str)
	// bytes, err := bcrypt.GenerateFromPassword([]byte(str), SALT)
	// if err != nil {
	// 	c <- MessageReturn{"Error while hashing", err}
	// 	return
	// }
	//coupon := string(bytes)

	// check if coupon exists
	data, _, _, err := con.QueryNeoAll(`
	MATCH (n:EVENT)-[r:PRESENT`+strconv.Itoa(attendance.Day)+`]->(a)
	WHERE a.email=$email
	RETURN [x IN r.coupons WHERE x = $coupon];
	`, map[string]interface{}{
		"email":  attendance.Email,
		"coupon": str,
	})

	if err != nil {
		c <- MessageReturn{"Error checking if coupon exists", err}
		return
	}

	str = fmt.Sprintf("%v", data)

	if str == "[[[]]]" || str == "[]" {
		c <- MessageReturn{"No match found for this coupon", nil}
		return
	}

	// check if empty coupon node
	// cp := ViewCoupon(attendance)
	// if len(cp) < 1 || cp[0] == "" {
	// 	ce := make(chan MessageReturn)
	// 	go DeleteCoupons(attendance, ce)

	// 	msg := <-ce
	// 	if err = msg.Err; err != nil {
	// 		c <- msg
	// 		return
	// 	}

	// 	c <- MessageReturn{"No more coupons exist for this user", nil}
	// 	return

	// }

	// remove from array
	_, err = con.ExecNeo(`
		MATCH (n:EVENT)-[c:PRESENT`+strconv.Itoa(attendance.Day)+`]->(a)
		WHERE a.email=$email AND n.name=$eventName
		SET c.coupons=[x IN c.coupons WHERE x <> $coupon];
		`, map[string]interface{}{
		"eventName": attendance.EventName,
		"email":     attendance.Email,
		"coupon":    str,
	})

	if err != nil {
		c <- MessageReturn{"Some error occurred", err}
		return
	}

	// check if empty node
	// cp = ViewCoupon(attendance)
	// if cp[0] == "" {
	// 	ce := make(chan MessageReturn)
	// 	go DeleteCoupons(attendance, ce)

	// 	msg := <-ce
	// 	if err := msg.Err; err != nil {
	// 		c <- msg
	// 		return
	// 	}

	// 	c <- MessageReturn{"No more coupons exist for this user", nil}
	// 	return

	// }

	c <- MessageReturn{"Successfully posted coupon", nil}
	return
}