package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
)

var DB *sql.DB

func ConnectDB() {

	err := godotenv.Load(os.Getenv("KO_DATA_PATH") + "/.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	var host string = os.Getenv("DB_HOST")
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatalf("Error converting DB_PORT to integer: %v", err)
	}
	var user string = os.Getenv("DB_USER")
	var password string = os.Getenv("DB_PASSWORD")
	var dbname string = os.Getenv("DB_NAME")

	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname = %s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Printf("failed to connect to database: %v", err)
	}
	log.Printf("Connected to DB - %v - and error is - %v", db, err)
	DB = db
}

func GetComponents() []Component {

	rows, err := DB.Query("SELECT * FROM components")
	if err != nil {
		log.Printf("failed to query components: %v", err)
		return nil
	}
	defer rows.Close()

	var components []Component
	for rows.Next() {
		var c Component
		if err := rows.Scan(&c.ID, &c.Name, &c.Type, &c.Team, &c.Language); err != nil {
			log.Printf("failed to scan component: %v", err)
			return nil
		}
		components = append(components, c)
	}
	return components
}

func GetTeams() []Team {
	rows, err := DB.Query("SELECT * FROM teams")
	if err != nil {
		log.Printf("failed to query teams: %v", err)
		return nil
	}
	defer rows.Close()

	var teams []Team
	for rows.Next() {
		var t Team
		if err := rows.Scan(&t.ID, &t.Name); err != nil {
			log.Printf("failed to scan team: %v", err)
			return nil
		}
		teams = append(teams, t)
	}
	return teams
}

func GetComponentsByTeam(teamID int) []Component {
	rows, err := DB.Query("SELECT * FROM components WHERE team = $1", teamID)
	if err != nil {
		log.Printf("failed to query components: %v", err)
		return nil
	}
	defer rows.Close()

	var components []Component
	for rows.Next() {
		var c Component
		if err := rows.Scan(&c.ID, &c.Name, &c.Type, &c.Team, &c.Language); err != nil {
			log.Printf("failed to scan component: %v", err)
			return nil
		}
		components = append(components, c)
	}
	return components
}

func GetComponentByID(componentID string) (Component, error) {
	var c Component
	err := DB.QueryRow("SELECT * FROM components WHERE id = $1", componentID).Scan(&c.ID, &c.Name, &c.Type, &c.Team, &c.Language)
	if err != nil {
		log.Printf("failed to query component: %v", err)
		return Component{}, err
	}
	return c, nil
}

func GetTeamByID(teamID int) (Team, error) {
	var t Team
	err := DB.QueryRow("SELECT * FROM teams WHERE id = $1", teamID).Scan(&t.ID, &t.Name)
	if err != nil {
		log.Printf("failed to query team: %v", err)
		return Team{}, err
	}
	return t, nil
}

func PostComponents(c Component) error {
	query := "INSERT INTO components (id, name, type, team, language) VALUES ($1, $2, $3, $4, $5)"
	log.Printf("Query is %v", query)
	_, err := DB.Exec(query, c.ID, c.Name, c.Type, c.Team, c.Language)
	if err != nil {
		log.Printf("failed to insert component: %v", err)
		return err
	}
	log.Printf("Coming back")
	return nil
}

func PostTeams(t Team) error {
	_, err := DB.Exec("INSERT INTO teams (id, name) VALUES ($1, $2)", t.ID, t.Name)
	if err != nil {
		log.Printf("failed to insert team: %v", err)
		return err
	}
	return nil
}

func GetImgByComponent(componentID int) ([]byte, error) {
	var img []byte
	err := DB.QueryRow("SELECT image FROM images WHERE id = $1", componentID).Scan(&img)
	if err != nil {
		log.Printf("failed to query image: %v", err)
		return nil, err
	}
	return img, nil
}

func PostImg(i Image) error {
	_, err := DB.Exec("INSERT INTO images (id, image) VALUES ($1,$2)", i.ID, i.Image)
	if err != nil {
		log.Printf("failed to insert image: %v", err)
		return err
	}
	return nil
}
