package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"

	"github.com/Rhymen/go-whatsapp"
)

const saveFile = "./data/wasessio.gob"

func login(wac *whatsapp.Conn) error {
	session, err := readSession()

	// if saved session exists, try to restore
	if err == nil {
		session, err = wac.RestoreWithSession(session)
		// if restore works,
		if err == nil {
			newLog("session restored, saving")
			if err = writeSession(session); err != nil {
				return fmt.Errorf("error saving session: %v", err)
			}
			return nil
		}
	}

	// if restore failed or if saved session doesnt exist, try new login
	newLog("saved session doesn't exist/work, will have to login again")

	qr := make(chan string)
	go func() {
		qcode := <-qr
		newLog("Latest Is ", qcode)
		setState("qr", qcode)
	}()
	session, err = wac.Login(qr)
	if err != nil {
		return fmt.Errorf("error during login: %v", err)
	}
	// if restore worked or even if you did new login, save it
	if err = writeSession(session); err != nil {
		return fmt.Errorf("error saving session: %v", err)
	}
	return nil
}

func readSession() (whatsapp.Session, error) {
	session := whatsapp.Session{}

	file, err := os.Open(saveFile)
	if err != nil {
		return session, err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	if err = decoder.Decode(&session); err != nil {
		return session, err
	}

	return session, nil
}

func writeSession(session whatsapp.Session) error {
	file, err := os.Create(saveFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	if err = encoder.Encode(session); err != nil {
		return err
	}

	return nil
}

func removeSession() {
	err := os.Remove(saveFile)
	if err != nil {
		log.Println("couldnt delete file", err)
	}
}
