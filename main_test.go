package main

import (
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"

	_ "modernc.org/sqlite"
)

func getConnection() *sql.DB {
	db, err := sql.Open("sqlite", "demo.db")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return db
}

func Test_SelectClient_WhenOk(t *testing.T) {
	db := getConnection()
	defer db.Close()
	clientID := 1
	client, err := selectClient(db, clientID)

	require.NoError(t, err)
	require.Equal(t, clientID, client.ID)
	require.NotEmpty(t, client.FIO)
	require.NotEmpty(t, client.Login)
	require.NotEmpty(t, client.Birthday)
	require.NotEmpty(t, client.Email)

}

func Test_SelectClient_WhenNoClient(t *testing.T) {
	db := getConnection()

	clientID := -1

	client, err := selectClient(db, clientID)

	require.Error(t, err)

	require.Equal(t, err, sql.ErrNoRows)

	require.Empty(t, client.FIO)
	require.Empty(t, client.Login)
	require.Empty(t, client.Birthday)
	require.Empty(t, client.Email)

}

func Test_InsertClient_ThenSelectAndCheck(t *testing.T) {
	db := getConnection()

	cl := Client{
		FIO:      "Test",
		Login:    "Test",
		Birthday: "19700101",
		Email:    "mail@mail.com",
	}

	id, err := insertClient(db, cl)

	require.NoError(t, err)
	require.Positive(t, id)
	cl.ID = id

	client, err := selectClient(db, cl.ID)

	require.NoError(t, err)
	require.Equal(t, cl, client)

}

func Test_InsertClient_DeleteClient_ThenCheck(t *testing.T) {
	db := getConnection()
	defer db.Close()
	cl := Client{
		FIO:      "Test",
		Login:    "Test",
		Birthday: "19700101",
		Email:    "mail@mail.com",
	}

	id, err := insertClient(db, cl)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	_, err = selectClient(db, id)
	require.NoError(t, err)

	err = deleteClient(db, id)
	require.NoError(t, err)

	_, err = selectClient(db, id)
	require.NotNil(t, err)
	require.ErrorIs(t, err, sql.ErrNoRows)

}
